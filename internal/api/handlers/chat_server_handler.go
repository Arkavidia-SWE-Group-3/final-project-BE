package handlers

import (
	"fmt"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type ChatServerHandler struct {
	Clients map[string]*websocket.Conn
	Rooms   map[string][]*websocket.Conn
	mu      sync.Mutex
}

type MessageObject struct {
	Data           string `json:"message"`
	From           string `json:"sender"`
	ProfilePicture string `json:"profile_picture"`
}

func NewChatServerHandler() *ChatServerHandler {
	return &ChatServerHandler{
		Clients: make(map[string]*websocket.Conn),
		Rooms:   make(map[string][]*websocket.Conn),
	}
}

func (h *ChatServerHandler) SetupRoutes(app *fiber.App) {
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:room", websocket.New(func(c *websocket.Conn) {
		room := c.Params("room")
		h.mu.Lock()
		h.Clients[c.RemoteAddr().String()] = c
		h.Rooms[room] = append(h.Rooms[room], c)
		h.mu.Unlock()

		defer func() {
			h.mu.Lock()
			delete(h.Clients, c.RemoteAddr().String())
			for i, client := range h.Rooms[room] {
				if client == c {
					h.Rooms[room] = append(h.Rooms[room][:i], h.Rooms[room][i+1:]...)
					break
				}
			}
			h.mu.Unlock()
		}()

		for {
			var msg MessageObject
			if err := c.ReadJSON(&msg); err != nil {
				fmt.Println("Error reading JSON:", err)
				break
			}
			h.broadcastMessage(room, msg)
		}
	}))
}

func (h *ChatServerHandler) broadcastMessage(room string, msg MessageObject) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, client := range h.Rooms[room] {
		if err := client.WriteJSON(msg); err != nil {
			fmt.Println("Error sending message:", err)
		}
	}
}
