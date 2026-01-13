package server

import (
	"compendium/models"
	"compendium/server/pb"
	"compendium/storage"
	"compendium/storage/postgres"
	postgresv2 "compendium/storage/postgres/postgresV2"

	"github.com/mentalisit/logger"
)

type Server struct {
	log  *logger.Logger
	db   *postgres.Db
	dbV2 *postgresv2.Db
	In   chan models.IncomingMessage
}

func NewServer(log *logger.Logger, st *storage.Storage) *Server {
	g := pb.GrpcMain(log, st)
	s := &Server{
		log:  log,
		db:   st.DB,
		dbV2: st.V2,
		In:   g.In,
	}

	go s.RunServerRestApi()
	return s
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
