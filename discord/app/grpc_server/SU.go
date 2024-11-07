package grpc_server

import "context"

func (s *Server) Unsubscribe(ctx context.Context, req *SubscrRequest) (*IntResponse, error) {
	code := s.ds.Unsubscribe(req.Nameid, req.ArgRoles, req.Guildid)
	return &IntResponse{Result: int32(code)}, nil
}
func (s *Server) Subscribe(ctx context.Context, req *SubscrRequest) (*IntResponse, error) {
	code := s.ds.Subscribe(req.Nameid, req.ArgRoles, req.Guildid)
	return &IntResponse{Result: int32(code)}, nil
}
