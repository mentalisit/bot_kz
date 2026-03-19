package matrix

import (
	"bridge/config"
	"bridge/models"
	"database/sql"
	"log"
	"strings"
	"sync"
)

// GhostProfile хранит закэшированную информацию о профиле ghost-пользователя
type GhostProfile struct {
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

type Matrix struct {
	Config       *config.ConfigBot
	RoomMembers  map[string]map[string]bool // roomID -> userID -> true
	AvatarCache  map[string]string          // httpURL -> mxcURI
	ProfileCache map[string]GhostProfile    // userID -> GhostProfile
	mu           sync.RWMutex
	OnMessage    func(models.ToBridgeMessage)
	db           *sql.DB
	stopChan     chan struct{}
}

// RunMatrixBridge initializes the bridge and ensures the bot user is registered
func RunMatrixBridge(cfg *config.ConfigBot, db *sql.DB) *Matrix {
	m := &Matrix{
		Config:       cfg,
		RoomMembers:  make(map[string]map[string]bool),
		AvatarCache:  make(map[string]string),
		ProfileCache: make(map[string]GhostProfile),
		db:           db,
		stopChan:     make(chan struct{}),
	}

	// 1. Create cache table if not exists
	m.createCacheTable()

	// 2. Load cache from database
	m.LoadCacheFromDB()

	// 3. Start periodic cache saver (every hour)
	go m.startCacheSaver(m.stopChan)

	// 4. Ensure the bot user is registered on the homeserver
	m.Register(cfg.Matrix.Username)

	// 5. Start AppService HTTP server in background
	go m.startAppServiceServer()

	log.Printf("Matrix AppService mode enabled for %s", cfg.Matrix.Username)
	return m
}

func (m *Matrix) IsGhost(userID string) bool {
	if !strings.HasPrefix(userID, "@") {
		return false
	}
	// Localpart is between @ and :
	localpart := userID[1:]
	if strings.Contains(localpart, ":") {
		localpart = strings.Split(localpart, ":")[0]
	}

	// Check for known bridge prefixes (tg_, ds_, wa_)
	prefixes := []string{"tg_", "ds_", "wa_"}
	for _, p := range prefixes {
		if strings.HasPrefix(localpart, p) {
			return true
		}
	}
	return false
}

func (m *Matrix) Shutdown() {
	log.Println("Matrix AppService shutting down")
	// Stop cache saver and trigger final save
	if m.stopChan != nil {
		close(m.stopChan)
	}
}
