package webServer

import (
	"fmt"
	"log"
	"net/http"
	"rs/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ==================== User Chats ====================

// GetUserChats возвращает список чатов пользователя
func (s *Server) GetUserChats(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	chats, err := s.db.GetUserChats(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Error getting user chats: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user chats"})
		return
	}

	c.JSON(http.StatusOK, chats)
}

// ==================== Chat Roles ====================

// GetChatRoles возвращает роли чата
func (s *Server) GetChatRoles(c *gin.Context) {
	chatIDStr := c.Param("chatId")
	userIDStr := c.Query("user_id")

	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat_id"})
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	roles, err := s.db.GetChatRoles(c.Request.Context(), chatID, userID)
	if err != nil {
		log.Printf("Error getting chat roles: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get chat roles"})
		return
	}

	c.JSON(http.StatusOK, roles)
}

// CreateRole создает новую роль
func (s *Server) CreateRole(c *gin.Context) {
	chatIDStr := c.Param("chatId")
	userIDStr := c.Query("user_id")

	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat_id"})
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	// Проверяем права администратора
	isAdmin, err := s.db.IsChatAdmin(c.Request.Context(), chatID, userID)
	if err != nil {
		log.Printf("Error checking admin status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check permissions"})
		return
	}

	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin rights required"})
		return
	}

	var req models.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role name is required"})
		return
	}

	if len(req.Name) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role name too long"})
		return
	}

	role := &models.Role{
		ChatID:    chatID,
		Name:      req.Name,
		CreatedBy: userID,
	}

	if err := s.db.CreateRole(c.Request.Context(), role); err != nil {
		log.Printf("Error creating role: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role created successfully", "success": true})
}

// UpdateRole обновляет имя роли
func (s *Server) UpdateRole(c *gin.Context) {
	chatIDStr := c.Param("chatId")
	roleIDStr := c.Param("roleId")
	userIDStr := c.Query("user_id")

	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat_id"})
		return
	}

	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role_id"})
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	isAdmin, err := s.db.IsChatAdmin(c.Request.Context(), chatID, userID)
	if err != nil {
		log.Printf("Error checking admin status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check permissions"})
		return
	}

	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin rights required"})
		return
	}

	var req models.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role name is required"})
		return
	}

	if len(req.Name) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role name too long"})
		return
	}

	if err := s.db.UpdateRoleName(c.Request.Context(), roleID, chatID, req.Name); err != nil {
		log.Printf("Error updating role: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role updated successfully", "success": true})
}

// DeleteRole удаляет роль
func (s *Server) DeleteRole(c *gin.Context) {
	chatIDStr := c.Param("chatId")
	roleIDStr := c.Param("roleId")
	userIDStr := c.Query("user_id")

	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat_id"})
		return
	}

	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role_id"})
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	isAdmin, err := s.db.IsChatAdmin(c.Request.Context(), chatID, userID)
	if err != nil {
		log.Printf("Error checking admin status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check permissions"})
		return
	}

	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin rights required"})
		return
	}

	if err := s.db.DeleteRole(c.Request.Context(), roleID, chatID); err != nil {
		log.Printf("Error deleting role: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully", "success": true})
}

// ==================== Role Membership ====================

// JoinRole добавляет пользователя в роль
func (s *Server) JoinRole(c *gin.Context) {
	s.handleRoleMembership(c, true)
}

// LeaveRole удаляет пользователя из роли
func (s *Server) LeaveRole(c *gin.Context) {
	s.handleRoleMembership(c, false)
}

func (s *Server) handleRoleMembership(c *gin.Context, join bool) {
	chatIDStr := c.Param("chatId")
	roleIDStr := c.Param("roleId")
	userIDStr := c.Query("user_id")

	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat_id"})
		return
	}

	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role_id"})
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	if join {
		if err := s.db.JoinRole(userID, roleID, chatID); err != nil {
			log.Printf("Error joining role: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to join role"})
			return
		}
	} else {
		if err := s.db.LeaveRole(userID, roleID); err != nil {
			log.Printf("Error leaving role: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to leave role"})
			return
		}
	}

	action := "joined"
	if !join {
		action = "left"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully %s role", action),
		"success": true,
	})
}

// ==================== Chat Users ====================

// GetChatUsers возвращает пользователей чата
func (s *Server) GetChatUsers(c *gin.Context) {
	chatIDStr := c.Param("chatId")

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat_id"})
		return
	}

	users, err := s.db.GetChatUsers(c.Request.Context(), chatID)
	if err != nil {
		log.Printf("Error getting chat users: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get chat users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// ==================== User Roles (Admin) ====================

// SetUserRole назначает/снимает роль пользователя (для админов)
func (s *Server) SetUserRole(c *gin.Context) {
	chatIDStr := c.Param("chatId")
	userIDStr := c.Param("userId")
	roleIDStr := c.Param("roleId")
	adminIDStr := c.Query("admin_id")

	if adminIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "admin_id is required"})
		return
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat_id"})
		return
	}

	targetUserID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role_id"})
		return
	}

	adminID, err := strconv.ParseInt(adminIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid admin_id"})
		return
	}

	isAdmin, err := s.db.IsChatAdmin(c.Request.Context(), chatID, adminID)
	if err != nil {
		log.Printf("Error checking admin status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check permissions"})
		return
	}

	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin rights required"})
		return
	}

	assign := c.Request.Method == "POST"

	if err := s.db.SetUserRole(c.Request.Context(), targetUserID, roleID, chatID, assign); err != nil {
		log.Printf("Error setting user role: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set user role"})
		return
	}

	action := "assigned"
	if !assign {
		action = "removed"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Role %s successfully", action),
		"success": true,
	})
}

// ==================== Permissions ====================

// GetUserPermissions возвращает права пользователя в чате
func (s *Server) GetUserPermissions(c *gin.Context) {
	chatIDStr := c.Param("chatId")
	userIDStr := c.Query("user_id")

	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat_id"})
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	isAdmin, err := s.db.IsChatAdmin(c.Request.Context(), chatID, userID)
	if err != nil {
		log.Printf("Error checking permissions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check permissions"})
		return
	}
	if !isAdmin {
		isAdmin, _ = s.cl.Tg.CheckAdminTg(fmt.Sprintf("%d", chatID), fmt.Sprintf("UserID %d", userID))
	}

	c.JSON(http.StatusOK, gin.H{
		"is_admin": isAdmin,
		"chat_id":  chatID,
		"user_id":  userID,
	})
}

// ==================== Role Members ====================

// GetRoleMembers возвращает участников конкретной роли
func (s *Server) GetRoleMembers(c *gin.Context) {
	chatIDStr := c.Param("chatId")
	roleIDStr := c.Param("roleId")

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat_id"})
		return
	}

	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role_id"})
		return
	}

	users, err := s.db.GetChatUsers(c.Request.Context(), chatID)
	if err != nil {
		log.Printf("Error getting chat users: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get chat users"})
		return
	}

	var roleName string
	err = s.db.GetRoleName(c.Request.Context(), roleID, &roleName)
	if err != nil {
		log.Printf("Error getting role name: %v", err)
	}

	var roleUsers []models.ChatUser
	if roleName == "all" {
		roleUsers = users
	} else {
		for _, user := range users {
			if user.Roles[roleID] {
				roleUsers = append(roleUsers, user)
			}
		}
	}

	c.JSON(http.StatusOK, roleUsers)
}

// ==================== Corp Members ====================

// GetCorpMembers возвращает участников корпорации из всех источников
func (s *Server) GetCorpMembers(c *gin.Context) {
	chatIDStr := c.Param("chatId")

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat_id"})
		return
	}

	var allMembers []models.CompendiumCorpMember

	if members, err := s.db.GetCorpMembersMyCompendium(c.Request.Context(), chatID); err == nil {
		allMembers = append(allMembers, members...)
	} else {
		log.Printf("Error getting corp members from my_compendium: %v", err)
	}

	if members, err := s.db.GetCorpMembersHSCompendium(c.Request.Context(), chatID); err == nil {
		allMembers = append(allMembers, members...)
	} else {
		log.Printf("Error getting corp members from hs_compendium: %v", err)
	}

	c.JSON(http.StatusOK, allMembers)
}

// RemoveCorpMember удаляет участника из корпорации
func (s *Server) RemoveCorpMember(c *gin.Context) {
	chatIDStr := c.Param("chatId")
	userIDStr := c.Param("userId")

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat_id"})
		return
	}

	tableSource := c.Query("tableSource")
	if tableSource == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tableSource parameter is required"})
		return
	}

	fmt.Printf("RemoveCorpMember %s %d %s\n", tableSource, chatID, userIDStr)

	var removeErr error
	switch tableSource {
	case "my_compendium":
		removeErr = s.db.RemoveCorpMemberMyCompendium(c.Request.Context(), chatID, userIDStr)
	case "hs_compendium":
		removeErr = s.db.RemoveCorpMemberHSCompendium(c.Request.Context(), chatID, userIDStr)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tableSource parameter"})
		return
	}

	if removeErr != nil {
		log.Printf("Error removing corp member from %s: %v", tableSource, removeErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove corp member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
