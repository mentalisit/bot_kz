package webServer

import (
	"encoding/json"
	"net/http"
	"rs/models"
	"strings"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/gin-gonic/gin"
)

const (
	vapidPublicKey  = "BKt7_2x5v93ObZPuoLVxsWHNqM5H1IUoAxAhweiwcEmMXDTjEaJ8yZnkdN6bLwHYy674vMWZsaCQ9Pc_PTb-p_8"
	vapidPrivateKey = "zUxnsPX14uDHOUe6AUud6jmNt6QnvOKhcVP-gArVyf0"
)

func (s *Server) postPushSubscribe(c *gin.Context) {
	var req struct {
		UUID         string `json:"uuid"`
		Subscription string `json:"subscription"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.UUID == "" || req.Subscription == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid and subscription are required"})
		return
	}

	// Extract endpoint for indexing/uniqueness
	var sub webpush.Subscription
	if err := json.Unmarshal([]byte(req.Subscription), &sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription format"})
		return
	}

	err := s.db.SavePushSubscription(req.UUID, sub.Endpoint, req.Subscription)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (s *Server) sendPushNotification(msg models.ChatMessage, ignoreUUIDs []string) {
	// 1. Get all subscriptions
	subs, err := s.db.GetAllPushSubscriptions()
	if err != nil {
		s.log.ErrorErr(err)
		return
	}

	// Create a map for faster lookup of users to ignore
	ignoreMap := make(map[string]bool)
	for _, id := range ignoreUUIDs {
		ignoreMap[id] = true
	}

	// 2. Collect unique UUIDs to load settings and nicknames
	uuidSet := make(map[string]bool)
	for _, ps := range subs {
		if !ignoreMap[ps.UUID] {
			uuidSet[ps.UUID] = true
		}
	}

	// Convert to slice
	uuids := make([]string, 0, len(uuidSet))
	for uuid := range uuidSet {
		uuids = append(uuids, uuid)
	}

	// 3. Load nicknames for all users
	nicknames, err := s.db.GetUserNicknamesByUUIDs(uuids)
	if err != nil {
		s.log.ErrorErr(err)
		nicknames = make(map[string]string)
	}

	// 4. Load settings for all users (we'll load individually for now, can be optimized later)
	userSettings := make(map[string]map[string]any)
	for _, uuid := range uuids {
		settings, err := s.db.GetChatUserSettings(uuid, msg.GID)
		if err != nil {
			s.log.ErrorErr(err)
			continue
		}
		userSettings[uuid] = settings
	}

	// 5. Prepare payload
	// Truncate text to avoid "payload has exceeded maximum length" error (Web Push limit is ~4KB)
	displayBody := msg.Text
	if msg.Type == "image" {
		displayBody = "[Изображение]"
	} else if len(displayBody) > 200 {
		displayBody = displayBody[:197] + "..."
	}

	payload, _ := json.Marshal(gin.H{
		"title": "Eco Night Chat",
		"body":  msg.Sender + ": " + displayBody,
		"icon":  msg.Avatar,
		"data": gin.H{
			"room": msg.Room,
		},
	})

	// 6. Send to those who pass settings filter
	for _, ps := range subs {
		if ignoreMap[ps.UUID] {
			continue // Skip users who are already online in the room
		}

		// Check user settings
		settings := userSettings[ps.UUID]
		if settings == nil {
			settings = make(map[string]any)
		}

		// Get global mode (default: 'all')
		globalMode := "all"
		if mode, ok := settings["global_mode"].(string); ok && mode != "" {
			globalMode = mode
		}

		// Skip if globally muted or no_push
		if globalMode == "muted" || globalMode == "no_push" {
			continue
		}

		// Check channel-specific settings
		channelNotifSettings := make(map[string]any)
		if chSettings, ok := settings["channel_notif_settings"].(map[string]any); ok {
			channelNotifSettings = chSettings
		}

		channelMode := ""
		if chMode, ok := channelNotifSettings[msg.Room].(string); ok {
			channelMode = chMode
		}

		// Determine effective mode: channel setting overrides global
		effectiveMode := globalMode
		if channelMode != "" && channelMode != "all" {
			effectiveMode = channelMode
		}

		// Skip if channel is muted or no_push
		if effectiveMode == "muted" || effectiveMode == "no_push" {
			continue
		}

		// Check mentions mode
		if effectiveMode == "mentions" {
			nickname := nicknames[ps.UUID]
			if nickname == "" {
				continue // Can't check mentions without nickname
			}
			// Check if user is mentioned in the text
			if !strings.Contains(msg.Text, "@"+nickname) {
				continue // Not mentioned, skip push
			}
		}

		// Send push notification
		var sub webpush.Subscription
		json.Unmarshal([]byte(ps.Subscription), &sub)

		resp, err := webpush.SendNotification(payload, &sub, &webpush.Options{
			Subscriber:      "mailto:admin@mentalisit.tsl.rocks",
			VAPIDPublicKey:  vapidPublicKey,
			VAPIDPrivateKey: vapidPrivateKey,
			TTL:             3600,
		})
		if err != nil {
			s.log.Error("Push failed for " + ps.Endpoint + ": " + err.Error())
			continue
		}
		resp.Body.Close()
	}
}
