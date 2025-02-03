package main

//
//import (
//	"fmt"
//	"github.com/celestix/gotgproto"
//	"github.com/celestix/gotgproto/dispatcher/handlers"
//	"github.com/celestix/gotgproto/dispatcher/handlers/filters"
//	"github.com/celestix/gotgproto/ext"
//	"github.com/celestix/gotgproto/functions"
//	"github.com/celestix/gotgproto/sessionMaker"
//	"github.com/glebarez/sqlite"
//	"github.com/go-faster/errors"
//	"log"
//)
//
//func main() {
//	client, err := gotgproto.NewClient(
//		// Get AppID from https://my.telegram.org/apps
//		26854446,
//		// Get ApiHash from https://my.telegram.org/apps
//		"4f08c2aba6a0753c78dae9cb1f8681d0",
//		// ClientType, as we defined above
//		gotgproto.ClientTypeBot("7076265870:AAH1PJ3tYBHNi4flBnpU7NSoWKAfHTWLKSY"), //.ClientTypePhone("+380989033544"),
//		// Optional parameters of client
//		&gotgproto.ClientOpts{
//			Session: sessionMaker.SqlSession(sqlite.Open("hs_bot")),
//		},
//	)
//
//	if err != nil {
//		log.Fatalln("failed to start client:", err)
//	}
//	dispatcher := client.Dispatcher
//
//	//dispatcher.AddHandler(handlers.NewAnyUpdate(anyUpdate))
//	dispatcher.AddHandlerToGroup(handlers.NewMessage(filters.Message.All, update), 1)
//	dispatcher.AddHandlerToGroup(handlers.NewMessage(filters.Message.Media, download), 1)
//
//	fmt.Printf("client (@%s) has been started...\n", client.Self.Username)
//
//	client.Idle()
//}
//
//func update(ctx *ext.Context, update *ext.Update) error {
//	eUser := update.EffectiveUser()
//	eMessage := update.EffectiveMessage
//	eChannel := update.GetChannel()
//
//	if eChannel != nil {
//		fmt.Printf("GuildName %s\n", eChannel.Title)
//		fmt.Printf("GuildId '%d'\n", eChannel.ID)
//		if eChannel.GetPhoto() != nil {
//			photo, ok := eChannel.GetPhoto().AsNotEmpty()
//			if ok {
//
//				fmt.Printf("ok %+v photo %v\n", ok, photo.PhotoID)
//			}
//		}
//	}
//
//	ChatId := getChatId(eMessage)
//	//Send(ctx, ChatId, "textmessage"+eMessage.Text)
//
//	fmt.Printf("text %s\n", eMessage.Text)
//	fmt.Printf("DmChat %d\n", eUser.ID)
//	fmt.Printf("Name %s\n", eUser.Username)
//	fmt.Printf("MentionName @%s\n", eUser.Username)
//	fmt.Printf("NameId %d\n", eUser.ID)
//	fmt.Printf("ChatId %s\n", ChatId)
//
//	fmt.Printf("anyUpdate %+v\n", update)
//
//	type IncomingMessage struct {
//		NickName    string
//		Avatar      string
//		GuildAvatar string
//		Type        string
//		Language    string
//	}
//
//	//if update.EffectiveMessage != nil {
//	//	fmt.Printf("EffectiveMessage %+v\n", update.EffectiveMessage)
//	//}
//	//if update.EffectiveChat() != nil {
//	//	fmt.Printf("EffectiveChat: %+v\n", update.EffectiveChat())
//	//}
//	////if update.GetChannel() != nil {
//	////	fmt.Printf("Channel: %+v\n", update.GetChannel())
//	////}
//	////if update.Entities != nil {
//	////	fmt.Printf("Entities: %+v\n", update.Entities)
//	////}
//	//if update.EffectiveUser() != nil {
//	//	fmt.Printf("EffectiveUser: %+v\n", update.EffectiveUser())
//	//}
//	//if update.GetChat() != nil {
//	//	fmt.Printf("Chat: %+v\n", update.GetChat())
//	//}
//	//if update.GetUserChat() != nil {
//	//	fmt.Printf("UserChat: %+v\n", update.GetUserChat())
//	//}
//	//if update.CallbackQuery != nil {
//	//	fmt.Printf("CallbackQuery: %+v\n", update.CallbackQuery)
//	//}
//	//if update.ChannelParticipant != nil {
//	//	fmt.Printf("ChannelParticipant: %+v\n", update.ChannelParticipant)
//	//}
//	//if update.ChatParticipant != nil {
//	//	fmt.Printf("ChatParticipant: %+v\n", update.ChatParticipant)
//	//}
//	//if update.ChatJoinRequest != nil {
//	//	fmt.Printf("ChatJoinRequest: %+v\n", update.ChatJoinRequest)
//	//}
//	//if update.InlineQuery != nil {
//	//	fmt.Printf("InlineQuery: %+v\n", update.InlineQuery)
//	//}
//	//if update.UpdateClass != nil {
//	//	fmt.Printf("UpdateClass: %+v\n", update.UpdateClass)
//	//}
//	//if update.Args() != nil {
//	//	fmt.Printf("Args: %+v\n", update.Args())
//	//}
//
//	return nil //err
//}
//func download(ctx *ext.Context, update *ext.Update) error {
//	filename, err := functions.GetMediaFileNameWithId(update.EffectiveMessage.Media)
//	if err != nil {
//		return errors.Wrap(err, "failed to get media file name")
//	}
//
//	_, err = ctx.DownloadMedia(
//		update.EffectiveMessage.Media,
//		ext.DownloadOutputPath(filename),
//		nil,
//	)
//	if err != nil {
//		return errors.Wrap(err, "failed to download media")
//	}
//
//	msg := fmt.Sprintf(`File "%s" downloaded`, filename)
//	_, err = ctx.Reply(update, msg, nil)
//	if err != nil {
//		return errors.Wrap(err, "failed to reply")
//	}
//
//	fmt.Println(msg)
//
//	return nil
//}
