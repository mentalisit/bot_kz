package bridge

import (
	"context"
	"fmt"
	"whatsapp/models"

	"github.com/mentalisit/logger"
	"google.golang.org/grpc"
)

type Client struct {
	conn   *grpc.ClientConn
	client BridgeServiceClient
	log    *logger.Logger
}

func NewClient(log *logger.Logger) *Client {
	conn, err := grpc.Dial("bridge:50051", grpc.WithInsecure())
	if err != nil {
		log.ErrorErr(err)
		return nil
	}
	fmt.Println("connect to bridge grpc ok")
	return &Client{
		conn:   conn,
		client: NewBridgeServiceClient(conn),
		log:    log,
	}
}

// Close закрывает соединение
func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) SendToBridge(i models.ToBridgeMessage) error {
	in := ToBridgeMessage{
		Text:        i.Text,
		Sender:      i.Sender,
		Tip:         i.Tip,
		ChatId:      i.ChatId,
		MesId:       i.MesId,
		GuildId:     i.GuildId,
		TimeMessage: i.TimestampUnix,
		Avatar:      i.Avatar,
		ReplyMap:    i.ReplyMap,
	}
	if len(i.Extra) > 0 {
		for _, info := range i.Extra {
			in.Extra = append(in.Extra, &FileInfo{
				Name:   info.Name,
				Data:   info.Data,
				Url:    info.URL,
				Size:   info.Size,
				FileId: info.FileID,
			})
		}
	}
	if i.Reply != nil && i.Reply.UserName != "" {
		in.Reply = &BridgeMessageReply{
			TimeMessage: i.Reply.TimeMessage,
			Text:        i.Reply.Text,
			Avatar:      i.Reply.Avatar,
			UserName:    i.Reply.UserName,
		}
	}

	// 4. Обработка конфигурации (Config)
	if i.Config != nil && i.Config.HostRelay != "" {
		// Инициализируем Protobuf-структуру конфигурации
		conf := &Bridge2Config{
			Id:                int32(i.Config.Id), // Убедимся, что Id — это int32
			NameRelay:         i.Config.NameRelay,
			HostRelay:         i.Config.HostRelay,
			Role:              i.Config.Role,
			ForbiddenPrefixes: i.Config.ForbiddenPrefixes,
		}

		// Инициализация поля Channel (map<string, Bridge2ConfigsList>)
		// Мы предполагаем, что Channel в Bridge2Config — это map[string]*Bridge2Config_Bridge2ConfigsList
		conf.Channel = make(map[string]*Bridge2Config_Bridge2ConfigsList)

		// --- Заполнение Channel ---
		if len(i.Config.Channel) > 0 {
			for key, configs := range i.Config.Channel {
				// 4a. Создаем вспомогательный список Protobuf
				listProto := &Bridge2Config_Bridge2ConfigsList{}

				for _, cfg := range configs {
					// 4b. Создаем Protobuf-структуру Bridge2Configs и добавляем ее в список
					listProto.Configs = append(listProto.Configs, &Bridge2Configs{
						ChannelId:       cfg.ChannelId,
						GuildId:         cfg.GuildId,
						CorpChannelName: cfg.CorpChannelName,
						AliasName:       cfg.AliasName,
						MappingRoles:    cfg.MappingRoles,
					})
				}

				// 4c. Добавляем список во вспомогательную карту Protobuf
				conf.Channel[key] = listProto
			}
		}

		in.Config = conf
	}

	// 6. Вызов клиента
	_, err := c.client.InboxBridge(context.Background(), &in) // Передаем указатель 'in'
	return err
}
