package matrix

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"
)

// MatrixCache структура для хранения всех кэшей в БД
type MatrixCache struct {
	ProfileCache map[string]GhostProfile    `json:"profile_cache"`
	AvatarCache  map[string]string          `json:"avatar_cache"`
	RoomMembers  map[string]map[string]bool `json:"room_members"`
}

// LoadCacheFromDB загружает кэш из базы данных при старте
func (m *Matrix) LoadCacheFromDB() {
	if m.db == nil {
		log.Println("[Matrix] Database not configured, skipping cache load")
		return
	}

	var cacheData sql.NullString
	err := m.db.QueryRow(`SELECT cache_data FROM rs_bot.matrix_cache WHERE id = 1`).Scan(&cacheData)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("[Matrix] No cache found in database, starting with empty cache")
			return
		}
		log.Printf("[Matrix] Error loading cache from database: %v", err)
		return
	}

	if !cacheData.Valid || cacheData.String == "" {
		log.Println("[Matrix] Cache data is empty, starting with empty cache")
		return
	}

	var cache MatrixCache
	if err := json.Unmarshal([]byte(cacheData.String), &cache); err != nil {
		log.Printf("[Matrix] Error parsing cache data: %v", err)
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if cache.ProfileCache != nil {
		m.ProfileCache = cache.ProfileCache
	}
	if cache.AvatarCache != nil {
		m.AvatarCache = cache.AvatarCache
	}
	if cache.RoomMembers != nil {
		m.RoomMembers = cache.RoomMembers
	}

	log.Printf("[Matrix] Cache loaded: %d profiles, %d avatars, %d rooms",
		len(m.ProfileCache), len(m.AvatarCache), len(m.RoomMembers))
}

// SaveCacheToDB сохраняет кэш в базу данных
func (m *Matrix) SaveCacheToDB() {
	if m.db == nil {
		return
	}

	m.mu.RLock()
	cache := MatrixCache{
		ProfileCache: m.ProfileCache,
		AvatarCache:  m.AvatarCache,
		RoomMembers:  m.RoomMembers,
	}
	m.mu.RUnlock()

	cacheJSON, err := json.Marshal(cache)
	if err != nil {
		log.Printf("[Matrix] Error marshaling cache: %v", err)
		return
	}

	// Upsert: вставить или обновить запись с id=1
	_, err = m.db.Exec(`
		INSERT INTO rs_bot.matrix_cache (id, cache_data, updated_at)
		VALUES (1, $1, NOW())
		ON CONFLICT (id) DO UPDATE SET cache_data = $1, updated_at = NOW()
	`, string(cacheJSON))

	if err != nil {
		log.Printf("[Matrix] Error saving cache to database: %v", err)
		return
	}

	log.Printf("[Matrix] Cache saved: %d profiles, %d avatars, %d rooms",
		len(cache.ProfileCache), len(cache.AvatarCache), len(cache.RoomMembers))
}

// startCacheSaver запускает горутину для периодического сохранения кэша
func (m *Matrix) startCacheSaver(stopChan <-chan struct{}) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.SaveCacheToDB()
		case <-stopChan:
			// Сохраняем кэш перед завершением
			m.SaveCacheToDB()
			log.Println("[Matrix] Cache saver stopped, final save completed")
			return
		}
	}
}

// createCacheTable создаёт таблицу для кэша, если она не существует
func (m *Matrix) createCacheTable() {
	if m.db == nil {
		return
	}

	_, err := m.db.Exec(`
		CREATE TABLE IF NOT EXISTS rs_bot.matrix_cache (
			id INTEGER PRIMARY KEY,
			cache_data JSONB NOT NULL,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Printf("[Matrix] Error creating cache table: %v", err)
	}
}
