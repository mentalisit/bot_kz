package grpc_server

import "context"

func (s *Server) ChannelTyping(ctx context.Context, req *ChannelTypingRequest) (*Empty, error) {
	s.ds.ChannelTyping(req.ChannelID)
	return &Empty{}, nil
}

func (s *Server) CleanRsBotOtherMessage(ctx context.Context, req *Empty) (*Empty, error) {
	s.ds.CleanRsBotOtherMessage()
	return &Empty{}, nil
}

func (s *Server) CheckAdmin(ctx context.Context, req *CheckAdminRequest) (*FlagResponse, error) {
	admin := s.ds.CheckAdmin(req.Nameid, req.Chatid)
	return &FlagResponse{Flag: admin}, nil
}

func (s *Server) CleanChat(ctx context.Context, req *CleanChatRequest) (*Empty, error) {
	s.ds.CleanChat(req.Chatid, req.Mesid, req.Text)
	return &Empty{}, nil
}

func (s *Server) CleanOldMessageChannel(ctx context.Context, req *CleanOldMessageChannelRequest) (*Empty, error) {
	s.ds.CleanOldMessageChannel(req.ChatId, req.Lim)
	return &Empty{}, nil
}
