package matrix

import (
	"bridge/models"
	"fmt"
	"log"
	"strings"
	"time"
)

func (m *Matrix) handleEvent(ev *Event) {
	// Handle Invites
	if ev.Type == "m.room.member" {
		membership, _ := ev.Content["membership"].(string)
		if membership == "invite" && ev.StateKey != nil && *ev.StateKey == m.Config.Matrix.Username {
			log.Printf("[Matrix] Bot invited to room %s, joining...", ev.RoomID)
			m.JoinRoom(ev.RoomID)
		}
	}

	// Handle Messages (can be extended to forward to other bridges)
	if (ev.Type == "m.room.message" || ev.Type == "m.sticker") && ev.Sender != m.Config.Matrix.Username && !m.IsGhost(ev.Sender) {
		body, _ := ev.Content["body"].(string)

		if m.OnMessage != nil {
			msg := models.ToBridgeMessage{
				Text:          body,
				Sender:        strings.Split(strings.TrimPrefix(ev.Sender, "@"), ":")[0],
				SenderId:      ev.Sender,
				Tip:           "matrix",
				ChatId:        ev.RoomID,
				MesId:         ev.EventID,
				TimestampUnix: time.Now().Unix(),
			}

			// Handle reply
			if relatesTo, ok := ev.Content["m.relates_to"].(map[string]interface{}); ok {
				if inReplyTo, ok := relatesTo["m.in_reply_to"].(map[string]interface{}); ok {
					if originalEventID, ok := inReplyTo["event_id"].(string); ok {
						msg.ReplyMap = map[string]string{ev.RoomID: originalEventID}
					}
				}
			}

			// Handle Media
			msgType, _ := ev.Content["msgtype"].(string)
			mxcURL, _ := ev.Content["url"].(string)
			if mxcURL != "" {
				info, _ := ev.Content["info"].(map[string]interface{})
				size := int64(0)
				if s, ok := info["size"].(float64); ok {
					size = int64(s)
				}

				msg.Extra = append(msg.Extra, models.FileInfo{
					Name: body,
					URL:  m.GetMediaURL(mxcURL),
					Size: size,
				})

				if msgType != "m.text" {
					msg.Text = "" // Don't use filename as text if it's media
				}
			}

			//log.Printf("[Matrix] Forwarding message from %s in %s (Type: %s)", ev.Sender, ev.RoomID, ev.Type)
			m.OnMessage(msg)
		}
	}
}

func (m *Matrix) GetMediaURL(mxc string) string {
	if !strings.HasPrefix(mxc, "mxc://") {
		return mxc
	}
	parts := strings.SplitN(mxc[6:], "/", 2)
	if len(parts) != 2 {
		return mxc
	}
	return fmt.Sprintf("%s/_matrix/media/v3/download/%s/%s",
		m.Config.Matrix.HomeserverURL, parts[0], parts[1])
}
