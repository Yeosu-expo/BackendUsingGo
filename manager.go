package main

import (
	"errors"
	"fmt"
	"kiosk/kioskPack"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     checkOrigin, // 정해진 도메인에서의 접근이 아니면 서버에 연결을 못하게 함
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")

	switch origin {
	case "http://localhost:8080": // 이 주소로 들어오지 않으면 연결이 종료됨
		return true
	default:
		return false
	}
}

var ErrEventNotSupported = errors.New("this event type is not supported")

type Manager struct {
	clients ClientList
	sync.RWMutex
	// 들어오는 이벤트 종류마다 적절한 조치를 취하기 위해 맵으로 매칭해주기 위한 필드
	handlers map[string]EventHandler
}

func NewManager() *Manager {
	m := &Manager{
		clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
	}
	m.setupEventHandlers()
	return m
}

func (m *Manager) setupEventHandlers() {
	// EventSendMessage라는 이벤트에 맞는 handler정의 및 연결
	m.handlers[EventSendMessage] = func(event Event, c *Client) error {
		fmt.Println(event)
		return nil
	}
}

// 이벤트를 다른 클라이언트한테 보냄
func (m *Manager) routeEvent(event Event, c *Client) error {
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return ErrEventNotSupported
	}
}

func (m *Manager) serveWS(w http.ResponseWriter, r *http.Request) {
	log.Println("New Connection")

	conn, err := upgrader.Upgrade(w, r, nil)
	kioskPack.CheckErr(err)

	client := NewCilent(conn, m)
	m.addClient(client)

	go client.readMessages()
	go client.writeMessages()
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.clients[client] = true
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		client.conn.Close()
		delete(m.clients, client)
	}
}
