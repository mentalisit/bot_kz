package roles

import (
	"fmt"
	"strings"
	"sync"
	"telegram/telegram/types"
	"time"
)

type Manager struct {
	roles map[string]*types.Role
	mux   sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		roles: make(map[string]*types.Role),
	}
}

// Создание роли с привязкой к чату
func (m *Manager) CreateRole(name, description string, createdBy int64, chatID int64, chatTitle string) *types.Role {
	m.mux.Lock()
	defer m.mux.Unlock()

	role := &types.Role{
		ID:          fmt.Sprintf("role_%d_%d", chatID, time.Now().UnixNano()),
		Name:        strings.TrimSpace(name),
		Description: strings.TrimSpace(description),
		CreatedAt:   time.Now(),
		CreatedBy:   createdBy,
		ChatID:      chatID,
		ChatTitle:   chatTitle,
		Subscribers: []int64{},
	}

	m.roles[role.ID] = role
	return role
}

// Получение ролей для конкретного чата
func (m *Manager) GetChatRoles(chatID int64) []*types.Role {
	m.mux.RLock()
	defer m.mux.RUnlock()

	result := make([]*types.Role, 0)
	for _, role := range m.roles {
		if role.ChatID == chatID {
			result = append(result, role)
		}
	}
	return result
}

// Получение всех ролей пользователя (во всех чатах)
func (m *Manager) GetUserRoles(userID int64) []*types.Role {
	m.mux.RLock()
	defer m.mux.RUnlock()

	result := make([]*types.Role, 0)
	for _, role := range m.roles {
		for _, sub := range role.Subscribers {
			if sub == userID {
				result = append(result, role)
				break
			}
		}
	}
	return result
}

// Получение ролей пользователя в конкретном чате
func (m *Manager) GetUserChatRoles(userID int64, chatID int64) []*types.Role {
	m.mux.RLock()
	defer m.mux.RUnlock()

	result := make([]*types.Role, 0)
	for _, role := range m.roles {
		if role.ChatID != chatID {
			continue
		}
		for _, sub := range role.Subscribers {
			if sub == userID {
				result = append(result, role)
				break
			}
		}
	}
	return result
}

// Получение всех ролей (для админки)
func (m *Manager) GetAllRoles() []*types.Role {
	m.mux.RLock()
	defer m.mux.RUnlock()

	result := make([]*types.Role, 0, len(m.roles))
	for _, role := range m.roles {
		result = append(result, role)
	}
	return result
}

// Остальные методы остаются, но добавляем проверку chatID
func (m *Manager) DeleteRole(roleID string, userID int64) bool {
	m.mux.Lock()
	defer m.mux.Unlock()

	role, exists := m.roles[roleID]
	if !exists || role.CreatedBy != userID {
		return false
	}

	delete(m.roles, roleID)
	return true
}

func (m *Manager) SubscribeToRole(userID int64, roleID string) bool {
	m.mux.Lock()
	defer m.mux.Unlock()

	role, exists := m.roles[roleID]
	if !exists {
		return false
	}

	// Проверяем, не подписан ли уже
	for _, sub := range role.Subscribers {
		if sub == userID {
			return true // Уже подписан
		}
	}

	role.Subscribers = append(role.Subscribers, userID)
	return true
}

func (m *Manager) UnsubscribeFromRole(userID int64, roleID string) bool {
	m.mux.Lock()
	defer m.mux.Unlock()

	role, exists := m.roles[roleID]
	if !exists {
		return false
	}

	for i, sub := range role.Subscribers {
		if sub == userID {
			role.Subscribers = append(role.Subscribers[:i], role.Subscribers[i+1:]...)
			return true
		}
	}

	return false
}

func (m *Manager) IsUserSubscribed(userID int64, roleID string) bool {
	m.mux.RLock()
	defer m.mux.RUnlock()

	role, exists := m.roles[roleID]
	if !exists {
		return false
	}

	for _, sub := range role.Subscribers {
		if sub == userID {
			return true
		}
	}
	return false
}

func (m *Manager) GetRole(roleID string) *types.Role {
	m.mux.RLock()
	defer m.mux.RUnlock()
	return m.roles[roleID]
}

// Новый метод: проверяет, есть ли у пользователя права на управление ролями в чате
func (m *Manager) CanManageRolesInChat(userID int64, chatID int64) bool {
	// Здесь можно добавить логику проверки прав пользователя в чате
	// Например, проверка что пользователь администратор и т.д.
	// Пока возвращаем true для всех
	return true
}
