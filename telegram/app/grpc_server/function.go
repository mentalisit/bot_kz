package grpc_server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path"
	"strconv"
	"strings"
	"telegram/models"
	"time"
)

func (s *Server) DeleteMessage(ctx context.Context, in *DeleteMessageRequest) (*ErrorResponse, error) {
	mId, _ := strconv.Atoi(in.GetMesid())

	err := s.tg.DelMessage(in.GetChatid(), mId)
	if err != nil {
		return &ErrorResponse{
			ErrorMessage: err.Error(),
		}, err
	}
	return &ErrorResponse{}, nil
}
func (s *Server) DeleteMessageSecond(ctx context.Context, in *DeleteMessageSecondRequest) (*ErrorResponse, error) {
	err := s.tg.DelMessageSecond(in.GetChatid(), in.GetMesid(), int(in.GetSecond()))
	if err != nil {
		return &ErrorResponse{
			ErrorMessage: err.Error(),
		}, err
	}
	return &ErrorResponse{}, nil
}
func (s *Server) EditMessage(ctx context.Context, in *EditMessageRequest) (*ErrorResponse, error) {
	mId, _ := strconv.Atoi(in.GetMID())
	err := s.tg.EditText(in.GetChatID(), mId, in.GetTextEdit(), in.GetParseMode())
	if err != nil {
		return &ErrorResponse{
			ErrorMessage: err.Error(),
		}, err
	}
	return &ErrorResponse{}, nil
}
func (s *Server) EditMessageTextKey(ctx context.Context, in *EditMessageTextKeyRequest) (*ErrorResponse, error) {
	err := s.tg.EditMessageTextKey(in.GetChatId(), int(in.GetEditMesId()), in.GetTextEdit(), in.GetLvlkz())
	if err != nil {
		return &ErrorResponse{
			ErrorMessage: err.Error(),
		}, err
	}
	return &ErrorResponse{}, nil
}
func (s *Server) SendPic(ctx context.Context, in *SendPicRequest) (*ErrorResponse, error) {
	id, err := s.tg.SendPic(in.GetChatid(), in.GetText(), in.GetImageBytes())
	if err != nil {
		return &ErrorResponse{
			ErrorMessage: err.Error(),
		}, err
	}
	return &ErrorResponse{Mesid: id}, nil
}
func (s *Server) SendPicScoreboard(ctx context.Context, in *ScoreboardRequest) (*ScoreboardResponse, error) {
	mid, err := s.tg.SendPicScoreboard(in.ChaatId, in.Text, in.FileNameScoreboard)
	if err != nil {
		return &ScoreboardResponse{ErrorMessage: err.Error()}, err
	}
	return &ScoreboardResponse{Mid: mid}, nil
}
func (s *Server) CheckAdmin(ctx context.Context, in *CheckAdminRequest) (*FlagResponse, error) {
	admin := s.tg.CheckAdminTg(in.GetChatid(), in.GetName())
	return &FlagResponse{Flag: admin}, nil
}
func (s *Server) SendChannelDelSecond(ctx context.Context, in *SendMessageRequest) (*FlagResponse, error) {
	send, err := s.tg.SendChannelDelSecond(in.GetChatID(), in.GetText(), in.GetParseMode(), int(in.GetSecond()))
	return &FlagResponse{Flag: send}, err
}
func (s *Server) GetAvatarUrl(ctx context.Context, in *GetAvatarUrlRequest) (*TextResponse, error) {
	url := s.tg.GetAvatarUrl(in.GetUserid())
	return &TextResponse{Text: url}, nil
}
func (s *Server) Send(ctx context.Context, in *SendMessageRequest) (*TextResponse, error) {
	id, err := s.tg.SendChannel(in.GetChatID(), in.GetText(), in.GetParseMode())
	if err != nil && err.Error() == "Forbidden: bot can't initiate conversation with a user" {
		return &TextResponse{Text: "Forbidden"}, nil
	}
	return &TextResponse{Text: id}, err
}
func (s *Server) SendPoll(ctx context.Context, in *SendPollRequest) (*TextResponse, error) {
	id := s.tg.SendPoll(models.Request{
		Data:    in.GetData(),
		Options: in.GetOptions(),
	})
	return &TextResponse{Text: id}, nil
}
func (s *Server) SendHelp(ctx context.Context, in *SendHelpRequest) (*TextResponse, error) {
	id, err := s.tg.SendHelp(in.GetChatId(), in.GetText(), in.GetOldMidHelps(), in.GetIfUser())
	return &TextResponse{Text: id}, err
}
func (s *Server) SendEmbedText(ctx context.Context, in *SendEmbedRequest) (*IntResponse, error) {
	id, err := s.tg.SendEmbed(in.GetLevel(), in.GetChatId(), in.GetText())
	return &IntResponse{Result: int32(id)}, err
}
func (s *Server) SendEmbedTime(ctx context.Context, in *SendMessageRequest) (*IntResponse, error) {
	id, err := s.tg.SendEmbedTime(in.GetChatID(), in.GetText())
	return &IntResponse{Result: int32(id)}, err
}
func (s *Server) SendChannelTyping(ctx context.Context, in *SendChannelTypingRequest) (*Empty, error) {
	_ = s.tg.ChatTyping(in.GetChannelID())
	return &Empty{}, nil
}
func (s *Server) SendBridgeArrayMessages(ctx context.Context, req *SendBridgeArrayMessagesRequest) (*SendBridgeArrayMessagesResponse, error) {
	in := models.BridgeSendToMessenger{
		Text:      req.GetText(),
		Sender:    req.GetUsername(),
		ChannelId: req.GetChannelID(),
		Avatar:    req.GetAvatar(),
		ReplyMap:  req.ReplyMap,
	}
	if req.GetExtra() != nil && len(req.Extra) > 0 {
		for _, info := range req.GetExtra() {
			fi := models.FileInfo{
				Name:   info.Name,
				Data:   info.Data,
				URL:    info.Url,
				Size:   info.Size,
				FileID: info.FileID,
			}
			if fi.URL != "" && len(fi.Data) == 0 {
				err := downloadFile(&fi)
				if err != nil {
					s.log.ErrorErr(err)
				}
			}
			in.Extra = append(in.Extra, fi)
		}

	}
	if req.GetReply() != nil {
		in.Reply = &models.BridgeMessageReply{
			TimeMessage: req.Reply.GetTimeMessage(),
			Text:        req.Reply.GetText(),
			Avatar:      req.Reply.GetAvatar(),
			UserName:    req.Reply.GetUserName(),
		}
	}

	messageIds := s.tg.SendBridgeFuncRest(in)
	var mids []*MessageIds
	for _, id := range messageIds {
		mids = append(mids, &MessageIds{
			MessageId: id.MessageId,
			ChatId:    id.ChatId,
		})
	}
	return &SendBridgeArrayMessagesResponse{MessageIds: mids}, nil
}

func downloadFile(fi *models.FileInfo) error {
	resp, err := http.Get(fi.URL)
	if err != nil {
		return fmt.Errorf("HTTP GET failed for URL %s: %w", fi.URL, err)
	}
	defer resp.Body.Close()

	// 2. Проверка HTTP-статуса
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status code %d for URL %s", resp.StatusCode, fi.URL)
	}

	// 3. Определение размера файла
	// Пытаемся получить размер из заголовка Content-Length
	if lengthStr := resp.Header.Get("Content-Length"); lengthStr != "" {
		fi.Size, _ = strconv.ParseInt(lengthStr, 10, 64)
	}

	// 4. Определение имени файла
	// Пытаемся получить имя из заголовка Content-Disposition
	if disposition := resp.Header.Get("Content-Disposition"); disposition != "" {
		if _, params, dispErr := parseContentDisposition(disposition); dispErr == nil {
			if filename, ok := params["filename"]; ok {
				fi.Name = filename
			}
		}
	}

	// Если имя не найдено, берем его из URL
	if fi.Name == "" {
		fi.Name = path.Base(fi.URL)
		// Убедимся, что имя не является пустым или слишком общим (например, URL=/)
		if fi.Name == "." || fi.Name == "/" || fi.Name == "" {
			// Запасной вариант: уникальное имя с расширением по типу
			ext := ""
			if contentType := resp.Header.Get("Content-Type"); contentType != "" {
				// Простая попытка получить расширение из MIME-типа
				ext = strings.Split(contentType, "/")[1]
			}
			fi.Name = fmt.Sprintf("downloaded_file_%d.%s", time.Now().UnixNano(), ext)
		}
	}

	// 5. Чтение данных
	fi.Data, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Если Content-Length был пуст, но мы прочитали данные, обновляем size
	if fi.Size == 0 {
		fi.Size = int64(len(fi.Data))
	}
	return nil
}

// parseContentDisposition - вспомогательная функция для парсинга Content-Disposition
// (Простая версия, без использования mime/multipart)
func parseContentDisposition(s string) (string, map[string]string, error) {
	parts := strings.Split(s, ";")
	if len(parts) == 0 {
		return "", nil, fmt.Errorf("invalid Content-Disposition header")
	}

	cdType := strings.TrimSpace(parts[0])
	params := make(map[string]string)

	for _, part := range parts[1:] {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) == 2 {
			key := strings.ToLower(strings.TrimSpace(kv[0]))
			value := strings.Trim(kv[1], "\" ")
			params[key] = value
		}
	}
	return cdType, params, nil
}
