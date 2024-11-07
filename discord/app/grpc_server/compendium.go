package grpc_server

import (
	"context"
)

func (s *Server) SendPic(ctx context.Context, req *SendPicRequest) (*ErrorResponse, error) {
	err := s.ds.SendPic(req.Chatid, req.Text, req.ImageBytes)
	if err != nil {
		return &ErrorResponse{ErrorMessage: err.Error()}, err
	}
	return &ErrorResponse{}, nil
}

func (s *Server) SendPoll(ctx context.Context, req *SendPollRequest) (*TextResponse, error) {
	id := s.ds.SendPoll(req.Data, req.Options)
	return &TextResponse{Text: id}, nil
}

func (s *Server) CheckRole(ctx context.Context, req *CheckRoleRequest) (*FlagResponse, error) {
	flag := s.ds.CheckRole(req.Guild, req.Memberid, req.Roleid)
	return &FlagResponse{Flag: flag}, nil
}
func (s *Server) GetAvatarUrl(ctx context.Context, req *GetAvatarUrlRequest) (*TextResponse, error) {
	url := s.ds.GetAvatarUrl(req.Userid)
	return &TextResponse{Text: url}, nil
}
func (s *Server) GetRoles(ctx context.Context, req *GuildRequest) (*GetRolesResponse, error) {
	roles := s.ds.GetRoles(req.GetGuild())
	var roles2 []*CorpRole
	for _, role := range roles {
		roles2 = append(roles2, &CorpRole{
			Id:   role.ID,
			Name: role.Name,
		})
	}
	return &GetRolesResponse{Roles: roles2}, nil
}
func (s *Server) GetMembersRoles(ctx context.Context, req *GuildRequest) (*MembersRolesResponse, error) {
	mm := s.ds.GetMembersRoles(req.Guild)
	var roles2 []*MembersRoles
	for _, roles := range mm {
		roles2 = append(roles2, &MembersRoles{
			Userid:  roles.Userid,
			RolesId: roles.RolesId,
		})
	}
	return &MembersRolesResponse{Memberroles: roles2}, nil
}
