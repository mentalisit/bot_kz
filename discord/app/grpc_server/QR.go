package grpc_server

import "context"

func (s *Server) QueueSend(ctx context.Context, req *QueueSendRequest) (*Empty, error) {
	s.ds.QueueSend(req.Text)
	return &Empty{}, nil
}

func (s *Server) ReplaceTextMessage(ctx context.Context, req *ReplaceTextMessageRequest) (*TextResponse, error) {
	newtext := s.ds.ReplaceTextMessage(req.Text, req.Guildid)
	return &TextResponse{Text: newtext}, nil
}

func (s *Server) RoleToIdPing(ctx context.Context, req *RoleToIdPingRequest) (*TextResponse, error) {
	ping, err := s.ds.RoleToIdPing(req.RolePing, req.Guildid)
	if err != nil {
		return nil, err
	}
	return &TextResponse{Text: ping}, nil
}
