package socketPack

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// 연결이 오래 지속되다보면, 소켓연결이 끊길 수 있음 그래서 서버에서 핑신호를 보내면, 클라이언트 들이
// 퐁신호를 보내서 응답하면, 연결이 살아있다는 것이다. 이를 확인해서 연결이 살아있는지 확인
var (
	pongWait = 3 * time.Second
	// pong신호를 기다리는 시간보다 길면 pongWait후에 연결이 끊김
	// 왜냐하면, write함수와 read함수가 동시에 돌아가는데, pingInterval값이 더 크면
	// write함수에서 ping신호를 보내기 전에 read함수에서 pong신호를 읽으려고 함
	// 그래서 에러가 나서 연결이 끊기고 뒤늦게 ping신호를 보내보지만, 이미 연결이 끊겼음.
	pingInterval = 2 * time.Second
)

type ClientList map[*Client]bool

type Client struct {
	conn       *websocket.Conn
	manager    *Manager
	egress     chan Event // 동시에 쓰는 것을 피하기 위한 채널
	clientType string     // 어떤 종류의 클라인지 확인
}

func NewCilent(conn *websocket.Conn, manager *Manager, clientType string) *Client {
	return &Client{
		conn:       conn,
		manager:    manager,
		egress:     make(chan Event),
		clientType: clientType,
	}
}

func (c *Client) writeMessages() {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		c.manager.removeClient(c)
	}()

	for {
		select {
		case message, ok := <-c.egress:
			if !ok { // egress로 메세지가 이상하게 왔다는 것을 의미
				if err := c.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed: ", err)
				}
				return
			}
			// ok이므로 메세지 쓰기
			// Event객체인 message를 다시 json으로 변환
			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return
			}
			// json으로 변환한 data를 클라들한테 보내기
			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println(err)
			}
		case <-ticker.C:
			log.Println("ping")
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("writemsg: ", err)
				return
			}
		}
	}
}

func (c *Client) readMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	// 악의적인 사용을 방지하기 위해 512byte이상의 메세지가 들어오면, 연결이 끊기게 설정
	c.conn.SetReadLimit(512)

	// pong신호를 보내기 위해 설정하는 작업
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		return
	}

	// 이 함수에서는 클라이언트에서 pong신호를 보내면, Deadline시간을 초기화 하는 듯하다.
	// 고릴라 패키지 설명이 있는 공식문서 같은게 없어서 예측만 가능
	c.conn.SetPongHandler(c.pongHandler)

	// 중간에 채널도 없어서 계속 동작을 수행할 거 같은데 왜 잘 돌아가지?
	for {
		_, payload, err := c.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}
		var request Event
		// payload로 읽어온 이벤트를 request에 옮겨 담음
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("error marshalling message: %v", err)
			break
		}
		// 읽어온 이벤트를 routeEvent로 넘김
		log.Println("Event type:", request.Type)
		if err := c.manager.routeEvent(request, c); err != nil {
			log.Println("Error handling Message: ", err)
		}
	}
}

func (c *Client) pongHandler(pongMsg string) error {
	log.Println("Pong")
	return c.conn.SetReadDeadline(time.Now().Add(pongWait))
}
