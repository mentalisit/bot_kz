package postgres

import (
	"database/sql"
	"discord/config"
	"encoding/json"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/mentalisit/restapi/models"
)

type Attachment struct {
	ID        int64
	FileID    string
	LocalPath string
	RemoteURL string
	SHA256    string
	MessageDB int64
	FileName  string
}

func (d *Db) ReadTop5Level(corpname string) []string {

	query := `
        SELECT lvlkz, COUNT(*) AS lvlkz_count
        FROM kzbot.sborkz
        WHERE corpname=$1
          AND date::timestamp >= CURRENT_DATE - INTERVAL '5 days'
        GROUP BY lvlkz
        ORDER BY lvlkz_count DESC
        LIMIT 5;
    `

	// Выполнение запроса
	rows, err := d.db.Query(query, corpname)
	defer rows.Close()
	if err != nil {
		d.log.ErrorErr(err)
	}

	var levels []string

	// Итерация по результатам запроса
	for rows.Next() {
		var lvlkz string
		var lvlkzCount int
		if err = rows.Scan(&lvlkz, &lvlkzCount); err != nil {
			d.log.ErrorErr(err)
		}
		levels = append(levels, lvlkz)
	}
	if err = rows.Err(); err != nil {
		d.log.ErrorErr(err)
	}
	return levels
}

func (d *Db) ReadMesIdDS(mesid string) (string, error) {

	sel := "SELECT lvlkz FROM kzbot.sborkz WHERE dsmesid = $1 AND active = 0"
	results, err := d.db.Query(sel, mesid)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
	}
	a := []string{}
	var dsmesid string
	for results.Next() {
		var t models.Sborkz
		err = results.Scan(&t.Lvlkz)
		a = append(a, t.Lvlkz)
	}
	a = RemoveDuplicates(a)
	if len(a) > 0 {
		dsmesid = a[0]
		return dsmesid, nil
	} else {
		return "", err
	}
}
func RemoveDuplicates[T comparable](a []T) []T {
	result := make([]T, 0, len(a))
	temp := map[T]struct{}{}
	for _, item := range a {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func (d *Db) NumActiveEvent(CorpName string) (event1 int) { //запрос номера ивента

	sel := "SELECT numevent FROM kzbot.rsevent WHERE corpname=$1 AND activeevent=1 ORDER BY numevent DESC LIMIT 1"
	row := d.db.QueryRow(sel, CorpName)
	err := row.Scan(&event1)
	if err != nil {
		if err == sql.ErrNoRows {
			event1 = 0
		} else {
			d.log.ErrorErr(err)
		}
	}
	return event1
}

func (d *Db) saveDiscordAttachments(m *discordgo.MessageCreate) error {
	var dbMessageID int64

	err := d.db.QueryRow(`
		SELECT id
		FROM my_chat.other_messages
		WHERE platform = 2 AND channel_id = $1 AND message_id = $2`,
		m.ChannelID, m.ID,
	).Scan(&dbMessageID)

	if err != nil {
		if err != sql.ErrNoRows {
			d.log.ErrorErr(err)
			return err
		}
		return nil
	}

	// Сохраняем вложения Discord
	for _, attachment := range m.Attachments {
		att := &Attachment{
			MessageDB: dbMessageID,
			FileID:    attachment.ID,
			FileName:  attachment.Filename,
			LocalPath: "", // Будет загружено при необходимости
			RemoteURL: attachment.URL,
		}

		if err := d.SaveAttachment(att); err != nil {
			d.log.ErrorErr(err)
			continue
		}

		// Обновляем дополнительные поля в БД
		_, err = d.db.Exec(`
			UPDATE my_chat.other_message_attachments
			SET file_type = $1,
			    mime_type = $2,
			    file_size = $3,
			    width = $4,
			    height = $5
			WHERE message_db_id = $6 AND attachment_id = $7
		`, "file", attachment.ContentType, int(attachment.Size), attachment.Width, attachment.Height, dbMessageID, attachment.ID)

		if err != nil {
			d.log.ErrorErr(err)
		}
	}

	return nil
}

// SaveAttachment сохраняет вложение в БД
func (d *Db) SaveAttachment(attachment *Attachment) error {
	_, err := d.db.Exec(`
		INSERT INTO my_chat.other_message_attachments (
			message_db_id,
			attachment_id,
			file_name,
			local_path,
			remote_url,
			sha256
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`,
		attachment.MessageDB,
		attachment.FileID,
		attachment.FileName,
		attachment.LocalPath,
		attachment.RemoteURL,
		attachment.SHA256,
	)

	return err
}

// GetGildUUIDMyCompendium получает UUID гильдии из my_compendium
func (d *Db) GetGildUUIDMyCompendium(guildID string) (gid *uuid.UUID, err error) {
	guildQuery := `SELECT gid FROM my_compendium.guilds WHERE EXISTS (
		SELECT 1 FROM jsonb_object_keys(channels) AS k
		CROSS JOIN jsonb_array_elements_text(channels->k) AS v
		WHERE v = $1
	)`
	err = d.db.QueryRow(guildQuery, guildID).Scan(&gid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("guild not found for guild %s", guildID)
		}
		return nil, fmt.Errorf("failed to find guild: %w", err)
	}
	return gid, nil
}

func (d *Db) GuildSave(g config.MultiAccountGuildV2) (*config.MultiAccountGuildV2, error) {
	if d.db == nil {
		d.log.Error("pizdec")
	}
	if g.GId == uuid.Nil {
		g.GId = uuid.New()
	}

	query := `
		INSERT INTO my_compendium.guilds (gid, guildname, channels, avatarurl)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (gid)
		DO UPDATE SET
			guildname = EXCLUDED.guildname,
			channels = EXCLUDED.channels,
			avatarurl = EXCLUDED.avatarurl
		RETURNING
			gid,
			guildname,
			channels,
			avatarurl
	`

	var res config.MultiAccountGuildV2

	marshal, err := json.Marshal(g.Channels)
	if err != nil {
		return nil, err
	}
	var ch []byte
	err = d.db.QueryRow(query,
		g.GId, g.GuildName, marshal, g.AvatarUrl).Scan(
		&res.GId, &res.GuildName, &ch, &res.AvatarUrl)

	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(ch, &res.Channels)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// SaveDiscordMessageWithNames сохраняет сообщение Discord в БД с переданными именами
func (d *Db) SaveDiscordMessageWithNames(communityID uuid.UUID, m *discordgo.MessageCreate, channelName, guildName string) error {
	raw, err := json.Marshal(m.Message)
	if err != nil {
		d.log.ErrorErr(err)
		return err
	}

	var replyToID any

	// ищем reply в нашей БД
	if m.Message.ReferencedMessage != nil {
		err = d.db.QueryRow(`
			SELECT id
			FROM my_chat.other_messages
			WHERE platform = 2
			  AND channel_id = $1
			  AND message_id = $2
		`,
			m.ChannelID,
			m.Message.ReferencedMessage.ID,
		).Scan(&replyToID)

		if err != nil {
			replyToID = nil
		}
	}

	var username string
	var displayName string

	if m.Author != nil {
		username = m.Author.Username
		displayName = m.Author.Username

		if m.Member != nil && m.Member.Nick != "" {
			displayName = m.Member.Nick
		}
	}

	var messageText string = m.Content

	guildID := m.GuildID

	_, err = d.db.Exec(`
		INSERT INTO my_chat.other_messages (community_id, platform,	message_id,
			guild_id, guild_name, channel_id, channel_name,
		    user_id, username, display_name,
			message_text, reply_to_db_id, raw_json)
		VALUES ($1,	2, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (platform, channel_id, message_id)
		DO UPDATE SET
			message_text = EXCLUDED.message_text,
			edited_at = NOW(),
			raw_json = EXCLUDED.raw_json`,
		communityID, m.ID,
		guildID, guildName, m.ChannelID, channelName,
		m.Author.ID, username, displayName,
		messageText, replyToID,
		raw,
	)

	if err != nil {
		d.log.ErrorErr(err)
		return err
	}

	return d.saveDiscordAttachments(m)
}
