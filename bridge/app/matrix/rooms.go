package matrix

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"bridge/models"
)

func (m *Matrix) JoinRoom(roomID string) {
	m.JoinRoomAs(roomID, m.Config.Matrix.Username)
}

func (m *Matrix) JoinRoomAs(roomID, userID string) {
	// 1. Check if user is already in the room via cache (Thread-safe)
	if m.IsUserInRoom(roomID, userID) {
		return
	}

	// 2. If it's a ghost, try to join.
	// For ghosts, we often need an invite if the room is private.
	if userID != m.Config.Matrix.Username {
		m.InviteUser(roomID, userID)
	}

	// 3. Try to join
	_, err := m.apiCall("POST", "/_matrix/client/v3/rooms/"+roomID+"/join", userID, nil)
	if err == nil {
		m.updateRoomMembersCache(roomID, userID)
		log.Printf("[Matrix] %s joined %s", userID, roomID)
	} else {
		// If join fails, it might be because the user is already there or
		// the invite didn't work. We still mark as 'checked' to reduce spam.
		if strings.Contains(err.Error(), "M_FORBIDDEN") {
			// Possibly already in or unauthorized, we'll try to send anyway
			m.updateRoomMembersCache(roomID, userID)
		}
		log.Printf("[Matrix] Join attempt for %s in %s: %v", userID, roomID, err)
	}
}

func (m *Matrix) IsUserInRoom(roomID, userID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if members, ok := m.RoomMembers[roomID]; ok {
		return members[userID]
	}
	return false
}

func (m *Matrix) updateRoomMembersCache(roomID, userID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.RoomMembers[roomID] == nil {
		m.RoomMembers[roomID] = make(map[string]bool)
	}
	m.RoomMembers[roomID][userID] = true
}

func (m *Matrix) SendText(roomID, text string) string {
	return m.SendTextAs(roomID, m.Config.Matrix.Username, text, "")
}

func (m *Matrix) SendTextAs(roomID, userID, text string, replyID string) string {
	txnID := fmt.Sprintf("%d", time.Now().UnixNano())
	path := fmt.Sprintf("/_matrix/client/v3/rooms/%s/send/m.room.message/%s", roomID, txnID)

	payload := map[string]interface{}{
		"msgtype": "m.text",
		"body":    text,
	}

	if replyID != "" {
		payload["m.relates_to"] = map[string]interface{}{
			"m.in_reply_to": map[string]interface{}{
				"event_id": replyID,
			},
		}
	}

	resp, err := m.apiCall("PUT", path, userID, payload)
	if err != nil {
		log.Printf("Error sending message from %s to %s: %v", userID, roomID, err)
		return ""
	}

	var res struct {
		EventID string `json:"event_id"`
	}
	json.Unmarshal(resp, &res)
	return res.EventID
}

func (m *Matrix) SendMediaAs(roomID, userID string, file models.FileInfo, replyID string) string {
	var mxcURL string
	var contentType string
	var err error

	if len(file.Data) > 0 {
		contentType = http.DetectContentType(file.Data)
		mxcURL, err = m.UploadMedia(file.Data, contentType, file.Name)
	} else if file.URL != "" {
		mxcURL, contentType, err = m.UploadMediaFromURL(file.URL, file.Name)
	} else {
		return ""
	}

	if err != nil {
		log.Printf("[Matrix] Failed to upload media: %v", err)
		return ""
	}

	txnID := fmt.Sprintf("%d", time.Now().UnixNano())
	eventType := "m.room.message"
	msgType := "m.file"

	if strings.HasPrefix(contentType, "image/") {
		msgType = "m.image"
		if contentType == "image/webp" {
			eventType = "m.sticker"
		}
	} else if strings.HasPrefix(contentType, "video/") {
		msgType = "m.video"
	} else if strings.HasPrefix(contentType, "audio/") {
		msgType = "m.audio"
	}

	path := fmt.Sprintf("/_matrix/client/v3/rooms/%s/send/%s/%s", roomID, eventType, txnID)

	payload := map[string]interface{}{
		"body": file.Name,
		"url":  mxcURL,
		"info": map[string]interface{}{
			"size":     file.Size,
			"mimetype": contentType,
		},
	}

	if eventType == "m.room.message" {
		payload["msgtype"] = msgType
	}

	if replyID != "" && eventType == "m.room.message" {
		payload["m.relates_to"] = map[string]interface{}{
			"m.in_reply_to": map[string]interface{}{
				"event_id": replyID,
			},
		}
	}

	resp, err := m.apiCall("PUT", path, userID, payload)
	if err != nil {
		log.Printf("[Matrix] Error sending media from %s to %s: %v", userID, roomID, err)
		return ""
	}

	var res struct {
		EventID string `json:"event_id"`
	}
	json.Unmarshal(resp, &res)
	return res.EventID
}

func (m *Matrix) GetRoomName(roomID string) string {
	resp, err := m.apiCall("GET", "/_matrix/client/v3/rooms/"+roomID+"/state/m.room.name", m.Config.Matrix.Username, nil)
	if err != nil {
		return ""
	}
	var nr struct {
		Name string `json:"name"`
	}
	json.Unmarshal(resp, &nr)
	return nr.Name
}

func (m *Matrix) SetupBridgeSpace(spaceName, adminUserID string) {
	createPayload := map[string]interface{}{
		"name": spaceName,
		"creation_content": map[string]interface{}{
			"type": "m.space",
		},
		"preset": "public_chat",
	}

	resp, err := m.apiCall("POST", "/_matrix/client/v3/createRoom", m.Config.Matrix.Username, createPayload)
	if err != nil {
		log.Printf("Failed to create space %s: %v", spaceName, err)
		return
	}

	var createResp struct {
		RoomID string `json:"room_id"`
	}
	json.Unmarshal(resp, &createResp)
	spaceID := createResp.RoomID
	log.Printf("Space '%s' created with ID: %s", spaceName, spaceID)

	m.InviteUser(spaceID, adminUserID)
	m.SetAdmin(spaceID, adminUserID)
}

func (m *Matrix) InviteUser(roomID, userID string) {
	payload := map[string]string{"user_id": userID}
	_, err := m.apiCall("POST", "/_matrix/client/v3/rooms/"+roomID+"/invite", m.Config.Matrix.Username, payload)
	if err != nil {
		// Only log if it's NOT a 403 (already in room/invited)
		if !strings.Contains(err.Error(), "403") {
			log.Printf("[Matrix] Invite failed for %s to %s: %v", userID, roomID, err)
		}
	} else {
		log.Printf("[Matrix] Invited %s to %s", userID, roomID)
	}
}

func (m *Matrix) SetAdmin(roomID, userID string) {
	resp, err := m.apiCall("GET", "/_matrix/client/v3/rooms/"+roomID+"/state/m.room.power_levels", m.Config.Matrix.Username, nil)
	if err != nil {
		log.Printf("Failed to get power levels: %v", err)
		return
	}

	var pl map[string]interface{}
	json.Unmarshal(resp, &pl)

	users, ok := pl["users"].(map[string]interface{})
	if !ok {
		users = make(map[string]interface{})
	}
	users[userID] = 100
	pl["users"] = users

	_, err = m.apiCall("PUT", "/_matrix/client/v3/rooms/"+roomID+"/state/m.room.power_levels", m.Config.Matrix.Username, pl)
	if err != nil {
		log.Printf("Failed to set power levels for %s: %v", userID, err)
	} else {
		log.Printf("User %s is now admin in %s", userID, roomID)
	}
}

func (m *Matrix) GetRoomIDByNameInSpace(roomName, adminUserID, spaceName string) (string, error) {
	spaceID := m.FindJoinedRoomByName(spaceName)
	if spaceID == "" {
		return "", fmt.Errorf("space '%s' not found", spaceName)
	}

	roomID := m.FindJoinedRoomByName(roomName)
	if roomID != "" {
		// Pre-fill cache if we found the room
		m.LoadRoomMembers(roomID)
		return roomID, nil
	}

	log.Printf("Room '%s' not found, creating...", roomName)
	createPayload := map[string]interface{}{
		"name":   roomName,
		"preset": "private_chat",
		"invite": []string{adminUserID},
		"initial_state": []map[string]interface{}{
			{
				"type":      "m.room.parent",
				"state_key": spaceID,
				"content": map[string]interface{}{
					"via":       []string{m.getHomeserverDomain()},
					"canonical": true,
				},
			},
		},
	}

	resp, err := m.apiCall("POST", "/_matrix/client/v3/createRoom", m.Config.Matrix.Username, createPayload)
	if err != nil {
		return "", fmt.Errorf("failed to create room %s: %v", roomName, err)
	}

	var createResp struct {
		RoomID string `json:"room_id"`
	}
	json.Unmarshal(resp, &createResp)
	newRoomID := createResp.RoomID
	m.LinkRoomToSpace(spaceID, newRoomID)

	// Load members to fill cache
	m.LoadRoomMembers(newRoomID)

	return newRoomID, nil
}

func (m *Matrix) LoadRoomMembers(roomID string) {
	resp, err := m.apiCall("GET", "/_matrix/client/v3/rooms/"+roomID+"/joined_members", m.Config.Matrix.Username, nil)
	if err != nil {
		return
	}

	var members struct {
		Joined map[string]interface{} `json:"joined"`
	}
	json.Unmarshal(resp, &members)

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.RoomMembers[roomID] == nil {
		m.RoomMembers[roomID] = make(map[string]bool)
	}
	for userID := range members.Joined {
		m.RoomMembers[roomID][userID] = true
	}
	//log.Printf("[Matrix] Loaded %d members for room %s", len(members.Joined), roomID)
}

func (m *Matrix) FindJoinedRoomByName(name string) string {
	resp, err := m.apiCall("GET", "/_matrix/client/v3/joined_rooms", m.Config.Matrix.Username, nil)
	if err != nil {
		return ""
	}

	var joined struct {
		Rooms []string `json:"joined_rooms"`
	}
	json.Unmarshal(resp, &joined)

	for _, roomID := range joined.Rooms {
		stateResp, err := m.apiCall("GET", "/_matrix/client/v3/rooms/"+roomID+"/state/m.room.name", m.Config.Matrix.Username, nil)
		if err == nil {
			var nr struct {
				Name string `json:"name"`
			}
			json.Unmarshal(stateResp, &nr)
			if nr.Name == name {
				return roomID
			}
		}
	}
	return ""
}

func (m *Matrix) LinkRoomToSpace(spaceID, roomID string) {
	payload := map[string]interface{}{
		"via": []string{m.getHomeserverDomain()},
	}
	_, err := m.apiCall("PUT", "/_matrix/client/v3/rooms/"+spaceID+"/state/m.space.child/"+roomID, m.Config.Matrix.Username, payload)
	if err != nil {
		log.Printf("Failed to link room %s to space %s: %v", roomID, spaceID, err)
	}
}
