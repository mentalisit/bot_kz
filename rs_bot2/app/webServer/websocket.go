package webServer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rs/models"
	"rs/storage/postgresV2"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// rooms maps room names to a map of clients in that room
	rooms       map[string]map[*Client]bool
	broadcast   chan models.ChatMessage
	register    chan *Client
	unregister  chan *Client
	mu          sync.Mutex
	db          *postgresV2.Db
	onBroadcast func(models.ChatMessage, []string)
	// typingUsers tracks who is currently typing in each room
	typingUsers map[string]map[string]*TypingUser // room -> uuid -> TypingUser
}

type TypingUser struct {
	UUID      string
	Nickname  string
	Avatar    string
	Timestamp int64
}

func NewHub(db *postgresV2.Db) *Hub {
	return &Hub{
		broadcast:   make(chan models.ChatMessage),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		rooms:       make(map[string]map[*Client]bool),
		db:          db,
		typingUsers: make(map[string]map[string]*TypingUser),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			internalRoom := client.InternalRoom()
			h.mu.Lock()
			if h.rooms[internalRoom] == nil {
				h.rooms[internalRoom] = make(map[*Client]bool)
			}
			h.rooms[internalRoom][client] = true
			h.mu.Unlock()

			// Send history to the new client
			go h.sendHistory(client)
			// Update online users list for everyone in the room
			go h.broadcastUserList(internalRoom, client.Room)

		case client := <-h.unregister:
			internalRoom := client.InternalRoom()
			h.mu.Lock()
			if roomClients, ok := h.rooms[internalRoom]; ok {
				if _, ok := roomClients[client]; ok {
					delete(roomClients, client)
					close(client.send)
					if len(roomClients) == 0 {
						delete(h.rooms, internalRoom)
					}
				}
			}
			h.mu.Unlock()
			// Update online users list
			go h.broadcastUserList(internalRoom, client.Room)

		case msg := <-h.broadcast:
			h.mu.Lock()
			msgBytes, _ := json.Marshal(msg)
			var onlineUUIDs []string
			seenIds := make(map[string]bool)

			if strings.HasPrefix(msg.Room, "dm_") {
				parts := strings.Split(msg.Room, "_")
				if len(parts) == 3 {
					u1, u2 := parts[1], parts[2]
					for _, clients := range h.rooms {
						for client := range clients {
							if client.UserInfo.UUID == u1 || client.UserInfo.UUID == u2 {
								select {
								case client.send <- msgBytes:
									if !seenIds[client.UserInfo.UUID] {
										onlineUUIDs = append(onlineUUIDs, client.UserInfo.UUID)
										seenIds[client.UserInfo.UUID] = true
									}
								default:
									close(client.send)
									delete(clients, client)
								}
							}
						}
					}
				}
			} else {
				internalRoom := msg.GID + ":" + msg.Room
				clients := h.rooms[internalRoom]
				for client := range clients {
					select {
					case client.send <- msgBytes:
						if !seenIds[client.UserInfo.UUID] {
							onlineUUIDs = append(onlineUUIDs, client.UserInfo.UUID)
							seenIds[client.UserInfo.UUID] = true
						}
					default:
						close(client.send)
						delete(clients, client)
					}
				}
			}
			h.mu.Unlock()

			if h.onBroadcast != nil {
				go h.onBroadcast(msg, onlineUUIDs)
			}
		}
	}
}

func (h *Hub) sendHistory(client *Client) {
	messages, err := h.db.GetWsMessages(client.GID, client.Room, 50)
	if err != nil {
		log.Printf("error getting history: %v", err)
		return
	}

	// Batch load reactions for all messages
	if len(messages) > 0 {
		var msgIDs []int
		for _, m := range messages {
			if m.ID > 0 {
				msgIDs = append(msgIDs, m.ID)
			}
		}
		if reactionsMap, err := h.db.GetReactionsForMessages(msgIDs); err == nil && reactionsMap != nil {
			for i := range messages {
				if r, ok := reactionsMap[messages[i].ID]; ok {
					messages[i].Reactions = r
				}
			}
		}
	}

	historyMsg := models.ChatMessage{
		Type:    "history",
		Room:    client.Room,
		History: messages,
	}
	msgBytes, _ := json.Marshal(historyMsg)
	client.send <- msgBytes
}

func (h *Hub) broadcastUserList(internalRoom string, room string) {
	h.mu.Lock()
	clients := h.rooms[internalRoom]
	var users []models.UserInfo
	for client := range clients {
		users = append(users, client.UserInfo)
	}
	h.mu.Unlock()

	userListMsg := models.ChatMessage{
		Type:  "user_list",
		Room:  room,
		Users: users,
	}
	msgBytes, _ := json.Marshal(userListMsg)

	h.mu.Lock()
	for client := range clients {
		select {
		case client.send <- msgBytes:
		default:
		}
	}
	h.mu.Unlock()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	Room     string
	GID      string
	UserInfo models.UserInfo
}

func (c *Client) InternalRoom() string {
	if strings.HasPrefix(c.Room, "dm_") {
		return "00000000-0000-0000-0000-000000000000:" + c.Room
	}
	return c.GID + ":" + c.Room
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var msg models.ChatMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("error unmarshaling message: %v", err)
			continue
		}

		requestedRoom := msg.Room

		// Fill in metadata from client
		msg.Room = c.Room
		msg.GID = c.GID
		msg.Sender = c.UserInfo.Nickname
		msg.SenderUUID = c.UserInfo.UUID
		msg.Avatar = c.UserInfo.Avatar
		msg.Timestamp = time.Now().Unix()

		// Allow room override only for forwarding action.
		if msg.Action == "forward" && requestedRoom != "" {
			if strings.HasPrefix(requestedRoom, "dm_") {
				parts := strings.Split(requestedRoom, "_")
				if len(parts) != 3 || (parts[1] != c.UserInfo.UUID && parts[2] != c.UserInfo.UUID) {
					log.Printf("forward denied for %s to room %s", c.UserInfo.UUID, requestedRoom)
					continue
				}
			}
			msg.Room = requestedRoom
		}

		log.Printf("Incoming message in room %s from %s: %s", msg.Room, msg.Sender, msg.Type)

		if msg.Type == "edit_msg" {
			if err := c.hub.db.EditWsMessage(msg.ID, msg.Text, msg.SenderUUID); err == nil {
				msg.Edited = true
				c.hub.broadcast <- msg
			}
		} else if msg.Type == "delete_msg" {
			if err := c.hub.db.DeleteWsMessage(msg.ID, msg.SenderUUID); err == nil {
				c.hub.broadcast <- msg
			}
		} else if msg.Type == "read_msg" {
			c.hub.db.MarkMessageAsRead(msg.ID, msg.SenderUUID)
		} else if msg.Type == "req_status" {
			readers, err := c.hub.db.GetMessageReaders(msg.ID)
			if err == nil {
				msg.Readers = readers
				msgBytes, _ := json.Marshal(msg)
				c.send <- msgBytes
			}
		} else if msg.Type == "reaction" {
			// Toggle reaction in DB
			if msg.ID > 0 && msg.Emoji != "" {
				c.hub.db.ToggleReaction(msg.ID, msg.SenderUUID, msg.Emoji)
				// Get updated reactions and broadcast
				if reactions, err := c.hub.db.GetReactionsForMessage(msg.ID); err == nil {
					update := models.ChatMessage{
						Type:      "reaction_update",
						ID:        msg.ID,
						Room:      msg.Room,
						GID:       msg.GID,
						Reactions: reactions,
					}
					c.hub.broadcast <- update
				}
			}
		} else if msg.Type == "typing" {
			// Typing indicator - update hub and broadcast
			c.hub.updateTyping(msg.Room, c.UserInfo, msg.Typing)
		} else if msg.Type == "voice" {
			// Voice message
			if id, err := c.hub.db.SaveWsMessage(msg); err == nil {
				msg.ID = id
			} else {
				log.Printf("error saving voice message: %v", err)
			}
			c.hub.broadcast <- msg
		} else {
			// standard message (text, image)
			if id, err := c.hub.db.SaveWsMessage(msg); err == nil {
				msg.ID = id
			} else {
				log.Printf("error saving message: %v", err)
			}
			c.hub.broadcast <- msg
		}
	}
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

// ServeWs handles websocket requests from the peer.
func (s *Server) ServeWs(c *gin.Context) {
	room := c.DefaultQuery("room", "general")
	uuidStr := c.Query("uuid")
	gidStr := c.Query("gid")
	if gidStr == "" {
		gidStr = "00000000-0000-0000-0000-000000000000"
	}

	if uuidStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "uuid required"})
		return
	}

	// Fetch user info from DB
	user, err := s.db.FindMultiAccountByUUId(uuidStr)
	if err != nil || user == nil {
		if room == "bot_relay" {
			// Allow guest access for broadcast room
			user = &models.MultiAccount{
				Nickname:  "Guest",
				AvatarURL: "",
			}
			user.UUID, _ = uuid.Parse(uuidStr)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid uuid or user not found"})
			return
		}
	}

	// DM Room Security: dm_uuid1_uuid2
	if strings.HasPrefix(room, "dm_") {
		parts := strings.Split(room, "_")
		if len(parts) == 3 {
			currentUUID := user.UUID.String()
			if parts[1] != currentUUID && parts[2] != currentUUID {
				c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this private room"})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid DM room format"})
			return
		}
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("upgrade error: %v", err)
		return
	}

	client := &Client{
		hub:  s.hub,
		conn: conn,
		send: make(chan []byte, 256),
		Room: room,
		GID:  gidStr,
		UserInfo: models.UserInfo{
			UUID:     user.UUID.String(),
			Nickname: user.Nickname,
			Avatar:   user.AvatarURL,
		},
	}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

// updateTyping tracks users who are typing in a room
func (h *Hub) updateTyping(room string, user models.UserInfo, isTyping bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.typingUsers[room] == nil {
		h.typingUsers[room] = make(map[string]*TypingUser)
	}

	if isTyping {
		h.typingUsers[room][user.UUID] = &TypingUser{
			UUID:      user.UUID,
			Nickname:  user.Nickname,
			Avatar:    user.Avatar,
			Timestamp: time.Now().Unix(),
		}
	} else {
		delete(h.typingUsers[room], user.UUID)
	}

	h.broadcastTypingUpdate(room)
}

// broadcastTypingUpdate sends the current typing users list to all clients in the room
func (h *Hub) broadcastTypingUpdate(room string) {
	internalRoom := fmt.Sprintf("%s:%s", room, "")

	typingList := make([]models.UserInfo, 0)
	for _, user := range h.typingUsers[room] {
		// Only include users who typed in the last 10 seconds
		if time.Now().Unix()-user.Timestamp < 10 {
			typingList = append(typingList, models.UserInfo{
				UUID:     user.UUID,
				Nickname: user.Nickname,
				Avatar:   user.Avatar,
			})
		}
	}

	msg := models.ChatMessage{
		Type:  "typing_update",
		Room:  room,
		Users: typingList,
	}

	msgBytes, _ := json.Marshal(msg)

	if clients, ok := h.rooms[internalRoom]; ok {
		for client := range clients {
			select {
			case client.send <- msgBytes:
			default:
				close(client.send)
				delete(clients, client)
			}
		}
	}
}
