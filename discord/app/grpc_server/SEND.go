package grpc_server

import "context"

func (s *Server) SendDmText(ctx context.Context, req *SendDmTextRequest) (*Empty, error) {
	s.ds.SendDmText(req.Text, req.AuthorID)
	return &Empty{}, nil
}

func (s *Server) Send(ctx context.Context, req *SendRequest) (*TextResponse, error) {
	id := s.ds.Send(req.Chatid, req.Text)
	return &TextResponse{Text: id}, nil
}

func (s *Server) SendChannelDelSecond(ctx context.Context, req *SendChannelDelSecondRequest) (*Empty, error) {
	s.ds.SendChannelDelSecond(req.Chatid, req.Text, int(req.Second))
	return &Empty{}, nil
}

func (s *Server) SendEmbedTime(ctx context.Context, req *SendEmbedTimeRequest) (*TextResponse, error) {
	id := s.ds.SendEmbedTime(req.Chatid, req.Text)
	return &TextResponse{Text: id}, nil
}

func (s *Server) SendComplexContent(ctx context.Context, req *SendComplexContentRequest) (*TextResponse, error) {
	id := s.ds.SendComplexContent(req.Chatid, req.Text)
	return &TextResponse{Text: id}, nil
}

func (s *Server) SendComplex(ctx context.Context, req *SendComplexRequest) (*TextResponse, error) {
	id := s.ds.SendComplex(req.Chatid, req.MapEmbeds)
	return &TextResponse{Text: id}, nil
}

func (s *Server) SendEmbedText(ctx context.Context, req *SendEmbedTextRequest) (*TextResponse, error) {
	m := s.ds.SendEmbedText(req.Chatid, req.Title, req.Text)
	if m == nil {
		return &TextResponse{}, nil
	}
	return &TextResponse{Text: m.ID}, nil
}

func (s *Server) SendHelp(ctx context.Context, req *SendHelpRequest) (*TextResponse, error) {
	id := s.ds.SendHelp(req.Chatid, req.Title, req.Description, req.OldMidHelps, req.IfUser)
	return &TextResponse{Text: id}, nil
}

func (s *Server) SendWebhook(ctx context.Context, req *SendWebhookRequest) (*TextResponse, error) {
	id := s.ds.SendWebhook(req.Text, req.Username, req.Chatid, req.Avatar)
	return &TextResponse{Text: id}, nil
}
func (s *Server) SendOrEditEmbedImage(ctx context.Context, req *SendEmbedImageRequest) (*ErrorResponse, error) {
	err := s.ds.SendOrEditEmbedImage(req.GetChatid(), req.GetTitle(), req.GetImageurl())
	if err != nil {
		return &ErrorResponse{ErrorMessage: err.Error()}, err
	}
	return &ErrorResponse{}, nil
}
func (s *Server) SendOrEditEmbedImageFileName(ctx context.Context, req *SendEmbedImageFileNameRequest) (*ErrorResponse, error) {
	err := s.ds.SendOrEditEmbedImageFileName(req.GetChatId(), req.GetTitle(), req.GetFileNameScoreboard())
	if err != nil {
		return &ErrorResponse{ErrorMessage: err.Error()}, err
	}
	return &ErrorResponse{}, nil
}
