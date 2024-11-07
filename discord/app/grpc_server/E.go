package grpc_server

import "context"

func (s *Server) EditComplexButton(ctx context.Context, req *EditComplexButtonRequest) (*ErrorResponse, error) {
	err := s.ds.EditComplexButton(req.Dsmesid, req.Dschatid, req.MapEmbed)
	if err != nil {
		return &ErrorResponse{ErrorMessage: err.Error()}, nil
	}
	return &ErrorResponse{}, nil
}
func (s *Server) EditWebhook(ctx context.Context, req *EditWebhookRequest) (*Empty, error) {
	s.ds.EditWebhook(req.Text, req.Username, req.ChatID, req.MID, req.AvatarURL)
	return &Empty{}, nil
}
