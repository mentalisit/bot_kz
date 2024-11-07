package grpc_server

import (
	"context"
	"discord/models"
)

func (s *Server) SendBridgeArrayMessages(ctx context.Context, req *SendBridgeArrayMessagesRequest) (*SendBridgeArrayMessagesResponse, error) {
	in := models.BridgeSendToMessenger{
		Text:      req.Text,
		Sender:    req.Username,
		ChannelId: req.ChannelID,
		Avatar:    req.Avatar,
	}
	if len(req.Extra) > 0 {
		for _, info := range req.Extra {
			in.Extra = append(in.Extra, models.FileInfo{
				Name:   info.Name,
				Data:   info.Data,
				URL:    info.Url,
				Size:   info.Size,
				FileID: info.FileID,
			})
		}

	}
	if req.Reply != nil {
		in.Reply = &models.BridgeMessageReply{
			TimeMessage: req.Reply.TimeMessage,
			Text:        req.Reply.Text,
			Avatar:      req.Reply.Avatar,
			UserName:    req.Reply.UserName,
		}
	}

	messageIds := s.ds.SendBridgeFuncRest(in)
	mids := make([]*MessageIds, len(messageIds))
	for _, id := range messageIds {
		mids = append(mids, &MessageIds{
			MessageId: id.MessageId,
			ChatId:    id.ChatId,
		})
	}
	return &SendBridgeArrayMessagesResponse{MessageIds: mids}, nil
}
