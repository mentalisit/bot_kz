package grpc_server

import "context"

func (s *Server) DeleteMessageSecond(ctx context.Context, req *DeleteMessageSecondRequest) (*Empty, error) {
	s.ds.DeleteMesageSecond(req.Chatid, req.Mesid, int(req.Second))
	return &Empty{}, nil
}
func (s *Server) DeleteMessage(ctx context.Context, req *DeleteMessageRequest) (*Empty, error) {
	s.ds.DeleteMessage(req.Chatid, req.Mesid)
	return &Empty{}, nil
}
