package pb

import (
	"compendium/models"
	"compendium/storage"
	"compendium/storage/postgres/multi"
	"context"
	"errors"
	"fmt"
	"github.com/mentalisit/logger"
	"google.golang.org/grpc"
	"net"
	"strings"
)

type Server struct {
	log *logger.Logger
	db  db
	In  chan models.IncomingMessage
	LogicServiceServer
	Multi *multi.Db
}

func GrpcMain(log *logger.Logger, st *storage.Storage) *Server {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.ErrorErr(err)
	}

	s := grpc.NewServer()
	serv := &Server{
		log:   log,
		db:    st.DB,
		In:    make(chan models.IncomingMessage, 10),
		Multi: st.Multi,
	}

	RegisterLogicServiceServer(s, serv)

	//fmt.Println("Server is running on port :50051")
	fmt.Printf("gRPC server is starting on %s\n", lis.Addr().String())
	go func() {
		if err = s.Serve(lis); err != nil {
			log.ErrorErr(err)
		}
	}()
	return serv
}

func (s *Server) InboxMessage(ctx context.Context, req *IncomingMessage) (*Empty, error) {
	in := models.IncomingMessage{
		Text:        req.Text,
		DmChat:      req.DmChat,
		Name:        req.Name,
		MentionName: req.MentionName,
		NameId:      req.NameId,
		NickName:    req.NickName,
		Avatar:      req.Avatar,
		ChannelId:   req.ChannelId,
		//GuildId:     req.GuildId,
		//GuildName:   req.GuildName,
		//GuildAvatar: req.GuildAvatar,
		Type:     req.Type,
		Language: req.Language,
	}
	if req.GuildId == "" {
		req.GuildId = "DM"
		req.GuildName = "DM"
		in.ChannelId = in.DmChat
		in.Type = in.Type[:2]
	}
	guild, _ := s.Multi.GuildGet(req.GuildId)
	if guild == nil {
		err := s.Multi.GuildInsert(models.MultiAccountGuild{
			GuildName: req.GuildName,
			Channels:  []string{req.GuildId},
			AvatarUrl: req.GuildAvatar,
		})
		if err != nil {
			s.log.ErrorErr(err)
		}
		guild, _ = s.Multi.GuildGet(req.GuildId)
	} else if guild.AvatarUrl != req.GuildAvatar {
		guild.AvatarUrl = req.GuildAvatar
		err := s.Multi.GuildUpdateAvatar(*guild)
		if err != nil {
			s.log.ErrorErr(err)
		}
	}
	in.MultiGuild = guild

	s.In <- in
	return &Empty{}, nil
}
func (s *Server) CorpMembersApiRead(ctx context.Context, req *ReqCorpMembersApiRead) (*ResCorpMembersApiRead, error) {

	read, err := s.db.CorpMembersApiRead(req.GuildId, req.Userid)
	if len(read) == 0 || read == nil || err != nil {
		return &ResCorpMembersApiRead{}, err
	}
	var memb []models.CorpMember
	for _, member := range read {
		if strings.Contains(member.UserId, req.Userid) {
			memb = append(memb, member)
		}
	}
	if len(memb) == 0 {
		return &ResCorpMembersApiRead{}, errors.New("member null")
	}
	var cMember []*CorpMember

	for _, m := range memb {
		cm := &CorpMember{
			Name:        m.Name,
			UserId:      m.UserId,
			GuildId:     m.GuildId,
			Avatar:      m.Avatar,
			AvatarUrl:   m.AvatarUrl,
			LocalTime:   m.LocalTime,
			LocalTime24: m.LocalTime24,
			TimeZone:    m.TimeZone,
			ZoneOffset:  int32(m.ZoneOffset),
			AfkFor:      m.AfkFor,
			AfkWhen:     int32(m.AfkWhen),
		}
		cm.Tech = make(map[int32]*TechLevels)

		for i, ints := range m.Tech {

			cm.Tech[int32(i)].Tech = append(cm.Tech[int32(i)].Tech, &TechLevel{
				Ts:    int64(ints[0]),
				Level: int32(ints[1]),
			})
		}

		cMember = append(cMember, cm)
	}

	return &ResCorpMembersApiRead{Array: cMember}, nil
}
func (s *Server) ApiGetUserAlts(ctx context.Context, req *ReqApiGetUserAlts) (*ResApiGetUserAlts, error) {
	read, _ := s.db.UsersGetByUserId(req.GetUserid())
	if read == nil {
		return &ResApiGetUserAlts{}, errors.New("error userid not found")
	}

	return &ResApiGetUserAlts{Alts: read.Alts}, nil
}

type db interface {
	CorpMembersRead(guildid string) ([]models.CorpMember, error)
	CorpMembersApiRead(guildid, userid string) ([]models.CorpMember, error)
	GuildRolesRead(guildid string) ([]models.CorpRole, error)
	GuildRolesExistSubscribe(guildid, RoleName, userid string) bool
	ListUserGetUserIdAndGuildId(token string) (userid string, guildid string, err error)
	UsersGetByUserId(userid string) (*models.User, error)
	GuildGet(guildid string) (*models.Guild, error)
	TechGet(username, userid, guildid string) ([]byte, error)
	TechUpdate(username, userid, guildid string, tech []byte) error
	CodeGet(code string) (*models.Code, error)
}
