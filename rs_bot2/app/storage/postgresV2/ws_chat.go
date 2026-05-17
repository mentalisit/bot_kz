package postgresV2

import (
	"database/sql"
	"rs/models"
	"time"
)

type WsMessage struct {
	ID            int    `json:"id"`
	Room          string `json:"room"`
	SenderName    string `json:"sender"`
	AvatarURL     string `json:"avatar"`
	ContentType   string `json:"type"` // "text", "image", or "voice"
	Content       string `json:"text"` // for images and voice too
	Timestamp     string `json:"timestamp"`
	ReplyToID     *int   `json:"reply_to_id,omitempty"`
	VoiceDuration int    `json:"voice_duration,omitempty"`
}

func (d *Db) SaveWsMessage(msg models.ChatMessage) (int, error) {

	var replyToID *int
	if msg.ReplyTo != nil {
		replyToID = &msg.ReplyTo.ID
	}

	query := `
		INSERT INTO my_chat.messages (gid, room, sender_uuid, sender_name, avatar_url, content_type, content, reply_to_id, voice_duration)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`
	var id int
	err := d.db.QueryRow(query, msg.GID, msg.Room, msg.SenderUUID, msg.Sender, msg.Avatar, msg.Type, msg.Text, replyToID, msg.VoiceDuration).Scan(&id)
	if err != nil {
		d.log.ErrorErr(err)
	}
	return id, err
}

func (d *Db) GetWsMessages(gid, room string, limit int) ([]models.ChatMessage, error) {

	query := `
		SELECT m.id, m.sender_uuid, m.sender_name, m.avatar_url, m.content_type, m.content, m.timestamp,
		       m.reply_to_id, m.voice_duration,
		       m.edited,
		       r.sender_name as reply_sender, r.content as reply_text
		FROM my_chat.messages m
		LEFT JOIN my_chat.messages r ON m.reply_to_id = r.id
		WHERE m.gid = $1 AND m.room = $2
		ORDER BY m.timestamp DESC
		LIMIT $3
	`
	rows, err := d.db.Query(query, gid, room, limit)
	if err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}
	defer rows.Close()

	var messages []models.ChatMessage
	for rows.Next() {
		var msg models.ChatMessage
		var ts time.Time
		var replyID *int
		var replySender, replyText *string
		var voiceDuration sql.NullInt32
		err := rows.Scan(&msg.ID, &msg.SenderUUID, &msg.Sender, &msg.Avatar, &msg.Type, &msg.Text, &ts,
			&replyID, &voiceDuration, &msg.Edited, &replySender, &replyText)
		if err != nil {
			return nil, err
		}
		msg.Timestamp = ts.Unix()
		msg.Room = room
		if voiceDuration.Valid {
			msg.VoiceDuration = int(voiceDuration.Int32)
		} else {
			msg.VoiceDuration = 0
		}

		// Build reply reference if exists
		if replyID != nil && replySender != nil {
			msg.ReplyTo = &models.ReplyRef{
				ID:     *replyID,
				Sender: *replySender,
				Text:   "",
			}
			if replyText != nil {
				msg.ReplyTo.Text = *replyText
			}
		}

		messages = append(messages, msg)
	}

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (d *Db) EditWsMessage(id int, text string, uuid string) error {

	// Only author can edit
	query := `UPDATE my_chat.messages SET content = $1, edited = TRUE WHERE id = $2 AND sender_uuid = $3`
	_, err := d.db.Exec(query, text, id, uuid)
	return err
}

func (d *Db) DeleteWsMessage(id int, uuid string) error {

	// Only author can delete
	query := `DELETE FROM my_chat.messages WHERE id = $1 AND sender_uuid = $2`
	_, err := d.db.Exec(query, id, uuid)
	return err
}

func (d *Db) MarkMessageAsRead(msgID int, uuid string) error {

	query := `INSERT INTO my_chat.message_reads (msg_id, user_uuid) VALUES ($1, $2) ON CONFLICT (msg_id, user_uuid) DO NOTHING`
	_, err := d.db.Exec(query, msgID, uuid)
	return err
}

func (d *Db) GetMessageReaders(msgID int) ([]models.ReaderInfo, error) {

	query := `
		SELECT m.uuid, m.nickname, m.avatarurl, EXTRACT(EPOCH FROM r.read_at)::bigint
		FROM my_chat.message_reads r
		JOIN my_compendium.multi_accounts m ON r.user_uuid = m.uuid
		WHERE r.msg_id = $1
		ORDER BY r.read_at DESC
	`
	rows, err := d.db.Query(query, msgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var readers []models.ReaderInfo
	for rows.Next() {
		var r models.ReaderInfo
		err := rows.Scan(&r.UUID, &r.Nickname, &r.Avatar, &r.ReadAt)
		if err == nil {
			readers = append(readers, r)
		}
	}
	return readers, nil
}

func (d *Db) SavePushSubscription(uuid, endpoint, subscription string) error {

	query := `
		INSERT INTO my_chat.push_subscriptions (uuid, endpoint, subscription)
		VALUES ($1, $2, $3)
		ON CONFLICT (endpoint) DO UPDATE SET
			uuid = EXCLUDED.uuid,
			subscription = EXCLUDED.subscription,
			updated_at = NOW()
	`
	_, err := d.db.Exec(query, uuid, endpoint, subscription)
	if err != nil {
		d.log.ErrorErr(err)
	}
	return err
}

func (d *Db) GetAllPushSubscriptions() ([]models.PushSubscription, error) {

	query := `SELECT id, uuid, endpoint, subscription, updated_at FROM my_chat.push_subscriptions`
	rows, err := d.db.Query(query)
	if err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}
	defer rows.Close()

	var subs []models.PushSubscription
	for rows.Next() {
		var s models.PushSubscription
		err := rows.Scan(&s.ID, &s.UUID, &s.Endpoint, &s.Subscription, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}
	return subs, nil
}

func (d *Db) SaveChatChannel(ch models.ChatChannel) error {

	query := `
		INSERT INTO my_chat.channels (id, gid, name, creator_uuid)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			creator_uuid = EXCLUDED.creator_uuid
	`
	_, err := d.db.Exec(query, ch.ID, ch.GID, ch.Name, ch.CreatorUUID)
	if err != nil {
		d.log.ErrorErr(err)
	}
	return err
}

func (d *Db) GetChatChannels(gid string) ([]models.ChatChannel, error) {

	query := `SELECT id, name, creator_uuid, created_at FROM my_chat.channels WHERE gid = $1 ORDER BY created_at ASC`
	rows, err := d.db.Query(query, gid)
	if err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}
	defer rows.Close()

	var channels []models.ChatChannel
	for rows.Next() {
		var ch models.ChatChannel
		err := rows.Scan(&ch.ID, &ch.Name, &ch.CreatorUUID, &ch.CreatedAt)
		if err != nil {
			return nil, err
		}
		channels = append(channels, ch)
	}
	return channels, nil
}

// ToggleReaction adds or removes a reaction. Returns true if added, false if removed.
func (d *Db) ToggleReaction(msgID int, userUUID, emoji string) (bool, error) {

	// Check if reaction already exists
	var count int
	err := d.db.QueryRow(
		`SELECT COUNT(*) FROM my_chat.reactions WHERE msg_id = $1 AND user_uuid = $2 AND emoji = $3`,
		msgID, userUUID, emoji).Scan(&count)
	if err != nil {
		return false, err
	}

	if count > 0 {
		// Remove
		_, err = d.db.Exec(
			`DELETE FROM my_chat.reactions WHERE msg_id = $1 AND user_uuid = $2 AND emoji = $3`,
			msgID, userUUID, emoji)
		return false, err
	}

	// Add
	_, err = d.db.Exec(
		`INSERT INTO my_chat.reactions (msg_id, user_uuid, emoji) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`,
		msgID, userUUID, emoji)
	return true, err
}

// GetReactionsForMessage returns map[emoji][]user_uuid for a single message.
func (d *Db) GetReactionsForMessage(msgID int) (map[string][]string, error) {

	query := `SELECT emoji, user_uuid FROM my_chat.reactions WHERE msg_id = $1 ORDER BY created_at ASC`
	rows, err := d.db.Query(query, msgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reactions := make(map[string][]string)
	for rows.Next() {
		var emoji, userUUID string
		if err := rows.Scan(&emoji, &userUUID); err == nil {
			reactions[emoji] = append(reactions[emoji], userUUID)
		}
	}
	return reactions, nil
}

// GetReactionsForMessages returns map[msgID]map[emoji][]user_uuid for batch loading.
func (d *Db) GetReactionsForMessages(msgIDs []int) (map[int]map[string][]string, error) {
	if len(msgIDs) == 0 {
		return nil, nil
	}

	query := `SELECT msg_id, emoji, user_uuid FROM my_chat.reactions WHERE msg_id = ANY($1) ORDER BY created_at ASC`
	rows, err := d.db.Query(query, msgIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int]map[string][]string)
	for rows.Next() {
		var msgID int
		var emoji, userUUID string
		if err := rows.Scan(&msgID, &emoji, &userUUID); err == nil {
			if result[msgID] == nil {
				result[msgID] = make(map[string][]string)
			}
			result[msgID][emoji] = append(result[msgID][emoji], userUUID)
		}
	}
	return result, nil
}
