package pb

import (
	"compendium/models"
	"compendium/storage"
	"compendium/storage/postgres"
	"compendium/storage/postgres/multi"
	"compendium/storage/postgres/postgresV2"
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/google/uuid"
	"github.com/mentalisit/logger"
	"google.golang.org/grpc"
)

type Server struct {
	log *logger.Logger
	db  *postgres.Db
	In  chan models.IncomingMessage
	LogicServiceServer
	Multi *multi.Db
	DBv2  *postgresv2.Db
	st    *storage.Storage
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
		DBv2:  st.V2,
		st:    st,
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
		GuildId:     req.GuildId,
		GuildName:   req.GuildName,
		GuildAvatar: req.GuildAvatar,
		Type:        req.Type,
		Language:    req.Language,
	}
	if req.GuildId == "" {
		req.GuildId = "DM"
		req.GuildName = "DM"
		in.ChannelId = in.DmChat
		in.Type = in.Type[:2]
	}

	multiAccount, _ := s.db.Multi.FindMultiAccountByUserId(in.NameId)
	if multiAccount != nil && multiAccount.TelegramID != "" && multiAccount.DiscordID != "" {
		if in.Avatar != "" {
			if multiAccount.AvatarURL != in.Avatar {
				multiAccount.AvatarURL = in.Avatar
				_, _ = s.db.Multi.UpdateMultiAccountAvatarUrl(*multiAccount)
			}
		}
		in.MultiAccount = multiAccount
	}

	//v2
	in.MAcc, _ = s.DBv2.FindMultiAccountByUserId(req.NameId)
	fmt.Printf("FindMultiAccountByUserId %s\n", req.NameId)
	if in.MAcc == nil || in.MAcc.Nickname == "" {
		oldAcc, _ := s.st.DB.Multi.FindMultiAccountByUserId(req.NameId)
		user, _ := s.db.UsersGetByUserId(req.NameId)
		if oldAcc == nil && user != nil {
			oldAcc = &models.MultiAccount{
				UUID:      uuid.New(),
				Nickname:  user.GameName,
				AvatarURL: user.AvatarURL,
				Alts:      user.Alts,
			}
			if oldAcc.Nickname == "" {
				oldAcc.Nickname = user.Username
			}
			//if req.Type
			switch req.Type {
			case "ds":
				oldAcc.DiscordID = user.ID
				oldAcc.DiscordUsername = user.Username
			case "tg":
				oldAcc.TelegramID = user.ID
				oldAcc.TelegramUsername = user.Username
			case "wa":
				oldAcc.WhatsappID = user.ID
				oldAcc.WhatsappUsername = user.Username
			}
		}
		if oldAcc != nil {
			//копируем
			in.MAcc, _ = s.DBv2.CreateMultiAccountFull(*oldAcc)
		} else {
			// Создаем новый аккаунт
			in.MAcc, _ = s.DBv2.CreateMultiAccountWithPlatform(req.NameId, req.Name, req.Type, req.Name)
		}
	}
	if in.MAcc.AvatarURL != req.Avatar {
		in.MAcc.AvatarURL = req.Avatar
		in.MAcc, _ = s.DBv2.UpdateMultiAccountAvatarUrl(*in.MAcc)

	}
	fmt.Printf("in.ma %+v\n", in.MAcc)
	guild2, err := s.DBv2.GuildGetChatId(req.GuildId)
	if err != nil && guild2 == nil {
		g := models.MultiAccountGuildV2{
			GuildName: req.GuildName,
			Channels:  make(map[string][]string),
			AvatarUrl: req.GuildAvatar,
		}
		g.Channels[req.Type] = append(g.Channels[req.Type], req.GuildId)
		guild2, err = s.DBv2.GuildInsert(g)
		if err != nil {
			s.log.ErrorErr(err)
		}
	} else if guild2 != nil && guild2.AvatarUrl != req.GuildAvatar {
		guild2.AvatarUrl = req.GuildAvatar
		err = s.DBv2.GuildUpdateAvatar(*guild2)
		if err != nil {
			s.log.ErrorErr(err)
		}
	}
	in.MGuild = guild2

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
				Ts:    ints.Ts,
				Level: int32(ints.Level),
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

//type db interface {
//	CorpMembersRead(guildid string) ([]models.CorpMember, error)
//	CorpMembersApiRead(guildid, userid string) ([]models.CorpMember, error)
//	GuildRolesRead(guildid string) ([]models.CorpRole, error)
//	GuildRolesExistSubscribe(guildid, RoleName, userid string) bool
//	ListUserGetUserIdAndGuildId(token string) (userid string, guildid string, err error)
//	UsersGetByUserId(userid string) (*models.User, error)
//	//GuildGet(guildid string) (*models.Guild, error)
//	TechGet(username, userid, guildid string) ([]byte, error)
//	TechUpdate(username, userid, guildid string, tech []byte) error
//	CodeGet(code string) (*models.Code, error)
//}
