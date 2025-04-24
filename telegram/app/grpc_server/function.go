package grpc_server

import (
	"context"
	"strconv"
	"telegram/models"
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
	err := s.tg.SendPic(in.GetChatid(), in.GetText(), in.GetImageBytes())
	if err != nil {
		return &ErrorResponse{
			ErrorMessage: err.Error(),
		}, err
	}
	return &ErrorResponse{}, nil
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
	send, err := s.tg.SendChannelDelSecond(in.GetChatID(), in.GetText(), int(in.GetSecond()))
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
	}
	if req.GetExtra() != nil && len(req.Extra) > 0 {
		in.Extra = make([]models.FileInfo, 0, len(req.Extra))
		for _, info := range req.GetExtra() {
			in.Extra = append(in.Extra, models.FileInfo{
				Name:   info.Name,
				Data:   info.Data,
				URL:    info.Url,
				Size:   info.Size,
				FileID: info.FileID,
			})
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
	mids := make([]*MessageIds, len(messageIds))
	for _, id := range messageIds {
		mids = append(mids, &MessageIds{
			MessageId: id.MessageId,
			ChatId:    id.ChatId,
		})
	}
	return &SendBridgeArrayMessagesResponse{MessageIds: mids}, nil
}
