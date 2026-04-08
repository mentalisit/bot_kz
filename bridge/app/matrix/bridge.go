package matrix

import (
	"bridge/models"
	"fmt"
	"log"
	"strings"
	"sync"
)

// SendInMessage is called by the bridge logic to send messages to Matrix
func (m *Matrix) SendInMessage(in models.ToBridgeMessage) []models.MessageIds {
	var mids []models.MessageIds

	// 1. Find or create the target room in the space
	roomID, err := m.GetRoomIDByNameInSpace(in.Config.NameRelay, "@mentalisit:mentalisit.myds.me", "мосты")
	if err != nil {
		log.Printf("[Matrix] Failed to get room ID: %v", err)
		return nil
	}

	// 2. Prepare ghost user ID (localpart)
	ghostLocalpart := fmt.Sprintf("%s_%s", in.Tip, in.SenderId)
	ghostLocalpart = strings.ReplaceAll(ghostLocalpart, "@", "_")
	ghostLocalpart = strings.ReplaceAll(ghostLocalpart, ":", "_")

	sender := fmt.Sprintf("(%s) %s", strings.ToUpper(in.Tip), in.Sender)

	// 3. Get or Create ghost user, update profile
	//log.Printf("[Matrix] Avatar URL for %s: %s", sender, in.Avatar)
	ghostMXID := m.GetOrCreateGhost(ghostLocalpart, sender, in.Avatar)

	// 4. Ensure ghost is in the room
	m.JoinRoomAs(roomID, ghostMXID)

	// 5. Check for reply
	replyID := ""
	if in.ReplyMap != nil {
		replyID = in.ReplyMap[roomID]
	}

	// 6. Send message as ghost
	if in.Text != "" {
		eventID := m.SendTextAs(roomID, ghostMXID, in.Text, replyID)
		if eventID != "" {
			mids = append(mids, models.MessageIds{MessageId: eventID, ChatId: roomID})
		}
	}

	// 7. Send attachments
	if len(in.Extra) > 0 {
		for _, file := range in.Extra {
			eventID := m.SendMediaAs(roomID, ghostMXID, file, replyID)
			if eventID != "" {
				mids = append(mids, models.MessageIds{MessageId: eventID, ChatId: roomID})
			}
		}
	}
	return mids
}

func (m *Matrix) SendBridgeArrayMessage(resultChannel chan<- models.MessageIds, wg *sync.WaitGroup, in models.ToBridgeMessage) {
	defer wg.Done()

	ids := m.SendInMessage(in)
	for _, id := range ids {
		resultChannel <- id
	}
}
