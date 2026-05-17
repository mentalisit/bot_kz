package postgres

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"telegram/config"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"github.com/google/uuid"
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

// CalculateFileSHA256 вычисляет SHA256 хеш файла
func CalculateFileSHA256(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// DownloadFile скачивает файл по URL и сохраняет по указанному пути
func DownloadFile(url string, savePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = os.MkdirAll(filepath.Dir(savePath), 0755)
	if err != nil {
		return err
	}

	out, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// SaveAttachment сохраняет вложение с проверкой дубликатов по SHA256
func (d *Db) SaveAttachment(attachment *Attachment) error {
	// Считаем SHA256
	sha, err := CalculateFileSHA256(attachment.LocalPath)
	if err != nil {
		return err
	}
	attachment.SHA256 = sha

	// Ищем дубликат
	var existingID int64
	err = d.db.QueryRow(`
		SELECT id
		FROM my_chat.other_message_attachments
		WHERE sha256 = $1
		LIMIT 1
	`, sha).Scan(&existingID)

	if err == nil {
		// Файл уже есть - сохраняем ссылку без дублирования файла
		_, err = d.db.Exec(`
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

	// Новый файл - сохраняем полную информацию
	_, err = d.db.Exec(`
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

func (d *Db) SaveTelegramMessage(communityID uuid.UUID, msg *tgbotapi.Message) error {
	raw, err := json.Marshal(msg)
	if err != nil {
		d.log.ErrorErr(err)
		return err
	}

	var replyToID any

	// ищем reply в нашей БД
	if msg.ReplyToMessage != nil {

		err = d.db.QueryRow(`
			SELECT id
			FROM my_chat.other_messages
			WHERE platform = 1
			  AND channel_id = $1
			  AND message_id = $2
		`,
			fmt.Sprint(msg.Chat.ID),
			fmt.Sprint(msg.ReplyToMessage.MessageID),
		).Scan(&replyToID)

		if err != nil {
			replyToID = nil
		}
	}

	var username string
	var displayName string

	if msg.From != nil {
		username = msg.From.UserName
		displayName = msg.From.FirstName

		if msg.From.LastName != "" {
			displayName += " " + msg.From.LastName
		}
	}

	var messageText string

	switch {
	case msg.Text != "":
		messageText = msg.Text

	case msg.Caption != "":
		messageText = msg.Caption
	}

	guildID := fmt.Sprint(msg.Chat.ID)
	guildName := msg.Chat.Title

	ThreadID := msg.MessageThreadID
	if !msg.IsTopicMessage && ThreadID != 0 {
		ThreadID = 0
	}
	channelId := strconv.FormatInt(msg.Chat.ID, 10) + fmt.Sprintf("/%d", ThreadID)
	channelName := ""

	// Проверяем если это создание топика - сохраняем в кэш
	if msg.IsTopicMessage && msg.ReplyToMessage != nil && msg.ReplyToMessage.ForumTopicCreated != nil {
		channelName = msg.ReplyToMessage.ForumTopicCreated.Name
		_ = d.SaveTopicCache(msg.Chat.ID, msg.MessageThreadID, channelName)
	} else if msg.IsTopicMessage && ThreadID != 0 {
		// Пробуем получить имя из кэша
		cachedName, err := d.GetTopicNameFromCache(msg.Chat.ID, ThreadID)
		if err == nil {
			channelName = cachedName
		}
	}

	// для Telegram guild обычно нет
	// можно использовать linked chat/forum/etc если хочешь

	_, err = d.db.Exec(`
		INSERT INTO my_chat.other_messages (community_id, platform,	message_id,
			guild_id, guild_name, channel_id, channel_name,
		    user_id, username, display_name,
			message_text, reply_to_db_id, raw_json)
		VALUES ($1,	1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (platform, channel_id, message_id)
		DO UPDATE SET
			message_text = EXCLUDED.message_text,
			edited_at = NOW(),
			raw_json = EXCLUDED.raw_json`,
		communityID, fmt.Sprint(msg.MessageID),
		guildID, guildName, channelId, channelName,
		fmt.Sprint(msg.From.ID), username, displayName,
		messageText, replyToID,
		raw,
	)

	if err != nil {
		d.log.ErrorErr(err)
		return err
	}

	return d.saveTelegramAttachments(msg)
}

// SaveTopicCache сохраняет имя топика в кэш
func (d *Db) SaveTopicCache(chatID int64, threadID int, topicName string) error {
	_, err := d.db.Exec(`
		INSERT INTO telegram.topic_cache (chat_id, thread_id, topic_name, updated_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (chat_id, thread_id)
		DO UPDATE SET 
			topic_name = EXCLUDED.topic_name,
			updated_at = NOW()
	`, chatID, threadID, topicName)
	return err
}

// GetTopicNameFromCache получает имя топика из кэша
func (d *Db) GetTopicNameFromCache(chatID int64, threadID int) (string, error) {
	var topicName string
	err := d.db.QueryRow(`
		SELECT topic_name 
		FROM telegram.topic_cache 
		WHERE chat_id = $1 AND thread_id = $2
	`, chatID, threadID).Scan(&topicName)
	return topicName, err
}

func (d *Db) saveTelegramAttachments(msg *tgbotapi.Message) error {
	var dbMessageID int64

	tokenTg := config.Instance.Token.TokenTelegram

	err := d.db.QueryRow(`
		SELECT id
		FROM my_chat.other_messages
		WHERE platform = 1 AND guild_id = $1 AND message_id = $2`,
		fmt.Sprint(msg.Chat.ID), fmt.Sprint(msg.MessageID),
	).Scan(&dbMessageID)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			d.log.ErrorErr(err)
			return err
		}
		return nil
	}

	// Получаем file_path через getFile API
	getFileURL := func(fileID string) (string, error) {
		url := fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", tokenTg, fileID)
		resp, err := http.Get(url)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		var result struct {
			Ok     bool `json:"ok"`
			Result struct {
				FilePath string `json:"file_path"`
			} `json:"result"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return "", err
		}
		if !result.Ok {
			return "", fmt.Errorf("telegram API error")
		}
		return result.Result.FilePath, nil
	}

	saveAttachment := func(fileID, fileName, fileType, mimeType string, fileSize int, width, height, duration int) error {
		filePath, err := getFileURL(fileID)
		if err != nil {
			return fmt.Errorf("getFile failed: %w", err)
		}

		downloadURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", tokenTg, filePath)
		localPath := filepath.Join("docker", "web", "attachments", filePath)

		if err := DownloadFile(downloadURL, localPath); err != nil {
			return fmt.Errorf("download failed: %w", err)
		}

		att := &Attachment{
			MessageDB: dbMessageID,
			FileID:    fileID,
			FileName:  fileName,
			LocalPath: localPath,
			RemoteURL: downloadURL,
		}

		if err := d.SaveAttachment(att); err != nil {
			return fmt.Errorf("save attachment failed: %w", err)
		}

		// Обновляем дополнительные поля в БД
		_, err = d.db.Exec(`
			UPDATE my_chat.other_message_attachments
			SET file_type = $1,
			    mime_type = $2,
			    file_size = $3,
			    width = $4,
			    height = $5,
			    duration = $6
			WHERE message_db_id = $7 AND attachment_id = $8
		`, fileType, mimeType, fileSize, width, height, duration, dbMessageID, fileID)

		return err
	}

	// PHOTO
	if len(msg.Photo) > 0 {
		photo := msg.Photo[len(msg.Photo)-1]
		if err := saveAttachment(photo.FileID, "", "image", "", photo.FileSize, photo.Width, photo.Height, 0); err != nil {
			return err
		}
	}

	// DOCUMENT
	if msg.Document != nil {
		if err := saveAttachment(msg.Document.FileID, msg.Document.FileName, "file", msg.Document.MimeType, int(msg.Document.FileSize), 0, 0, 0); err != nil {
			return err
		}
	}

	// VIDEO
	if msg.Video != nil {
		if err := saveAttachment(msg.Video.FileID, "", "video", msg.Video.MimeType, int(msg.Video.FileSize), msg.Video.Width, msg.Video.Height, msg.Video.Duration); err != nil {
			return err
		}
	}

	// VOICE
	if msg.Voice != nil {
		if err := saveAttachment(msg.Voice.FileID, "", "audio", msg.Voice.MimeType, int(msg.Voice.FileSize), 0, 0, msg.Voice.Duration); err != nil {
			return err
		}
	}

	return nil
}
