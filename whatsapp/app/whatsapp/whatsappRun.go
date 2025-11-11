package wa

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
	"whatsapp/config"
	"whatsapp/models"
	"whatsapp/storage"
	"whatsapp/whatsapp/restapi"

	"github.com/mdp/qrterminal"
	"github.com/mentalisit/logger"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
	_ "modernc.org/sqlite" // needed for sqlite
)

type Whatsapp struct {
	log                    *logger.Logger
	bridgeConfig           []models.Bridge2Config
	bridgeConfigUpdateTime int64
	Storage                *storage.Storage
	wc                     *whatsmeow.Client
	users                  map[string]types.ContactInfo
	userAvatars            map[string]string
	contacts               map[types.JID]types.ContactInfo
	mapMentionsNames       map[string]*mentionsNames
	joinedGroups           []*types.GroupInfo
	startedAt              time.Time
	*sync.RWMutex
	cfg *config.ConfigBot
	api *restapi.Recover
}
type mentionsNames struct {
	DefaultServer string
	HiddenServer  string
}

func NewWhatsapp(log *logger.Logger, cfg *config.ConfigBot, st *storage.Storage) *Whatsapp {
	b := &Whatsapp{
		log:                    log,
		bridgeConfig:           nil,
		bridgeConfigUpdateTime: 0,
		Storage:                st,
		users:                  make(map[string]types.ContactInfo),
		userAvatars:            make(map[string]string),
		mapMentionsNames:       make(map[string]*mentionsNames),
		cfg:                    cfg,
		api:                    restapi.NewRecover(log),
	}
	b.RWMutex = new(sync.RWMutex)
	err := b.Connect()
	if err != nil {
		log.ErrorErr(err)
		panic(err)
	}

	go b.DeleteMessageTimer()
	return b
}

// Connect to WhatsApp. Required implementation of the Bridger interface
func (b *Whatsapp) Connect() error {
	device, err := b.getDevice()
	if err != nil {
		return err
	}

	number := b.cfg.Whatsapp.Number
	if number == "" {
		return errors.New("whatsapp's telephone number need to be configured")
	}

	fmt.Println("Connecting to WhatsApp..")

	b.wc = whatsmeow.NewClient(device, waLog.Stdout("Client", "INFO", true))
	b.wc.AddEventHandler(b.eventHandler)

	firstLogin := false
	var qrChan <-chan whatsmeow.QRChannelItem
	if b.wc.Store.ID == nil {
		firstLogin = true
		qrChan, err = b.wc.GetQRChannel(context.Background())
		if err != nil && !errors.Is(err, whatsmeow.ErrQRStoreContainsID) {
			return errors.New("failed to to get QR channel:" + err.Error())
		}
	}

	err = b.wc.Connect()
	if err != nil {
		return errors.New("failed to connect to WhatsApp: " + err.Error())
	}

	if b.wc.Store.ID == nil {
		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				b.log.InfoStruct("QR channel result: %s", evt.Event)
			}
		}
	}

	// disconnect and reconnect on our first login/pairing
	// for some reason the GetJoinedGroups in JoinChannel doesn't work on first login
	if firstLogin {
		b.wc.Disconnect()
		time.Sleep(time.Second)

		err = b.wc.Connect()
		if err != nil {
			return errors.New("failed to connect to WhatsApp: " + err.Error())
		}
	}

	fmt.Println("WhatsApp connection successful")

	// Fix: Add context.Background() as first parameter
	b.contacts, err = b.wc.Store.Contacts.GetAllContacts(context.Background())
	if err != nil {
		return errors.New("failed to get contacts: " + err.Error())
	}

	b.joinedGroups, err = b.wc.GetJoinedGroups(context.Background())
	if err != nil {
		fmt.Println(err)
		return errors.New("failed to get list of joined groups: " + err.Error())
	}

	b.startedAt = time.Now()

	// map all the users
	for id, contact := range b.contacts {
		if !isGroupJid(id.String()) && id.String() != "status@broadcast" {
			// it is user
			b.users[id.String()] = contact
			if b.mapMentionsNames[contact.PushName] == nil {
				b.mapMentionsNames[contact.PushName] = &mentionsNames{}
			}
			if id.Server == "s.whatsapp.net" {
				b.mapMentionsNames[contact.PushName].DefaultServer = id.String()
			} else if id.Server == "lid" {
				b.mapMentionsNames[contact.PushName].HiddenServer = id.String()
			}
		}
	}

	// get user avatar asynchronously
	fmt.Printf("Getting user avatars")

	for jid := range b.users {
		info, err := b.GetProfilePicThumb(jid)
		if err != nil {
			//fmt.Printf("Could not get profile photo of %s: %v\n", jid, err)
			fmt.Printf(".")
		} else {
			b.Lock()
			if info != nil {
				b.userAvatars[jid] = info.URL
			}
			b.Unlock()
		}
	}

	fmt.Println("Finished getting avatars")

	return nil
}

func (b *Whatsapp) Disconnect() error {
	b.wc.Disconnect()
	return nil
}

type Message struct {
	Text      string    `json:"text"`
	Channel   string    `json:"channel"`
	Username  string    `json:"username"`
	UserID    string    `json:"userid"` // userid on the bridge
	Avatar    string    `json:"avatar"`
	Account   string    `json:"account"`
	Event     string    `json:"event"`
	Protocol  string    `json:"protocol"`
	Gateway   string    `json:"gateway"`
	ParentID  string    `json:"parent_id"`
	Timestamp time.Time `json:"timestamp"`
	ID        string    `json:"id"`
	Extra     []models.FileInfo
}

//	func (t *Telegram) Close() {
//		t.t.StopReceivingUpdates()
//		t.api.Close()
//	}
func (b *Whatsapp) DeleteMessageTimer() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mes := b.Storage.Db.TimerReadMessage("wa")
			if len(mes) > 0 {
				for _, m := range mes {
					if m.MesId != "" {
						_ = b.DeleteMessage(m.ChatId, m.MesId)
						b.Storage.Db.TimerDeleteMessage(m)
					}
				}
			}
		}
	}
}
