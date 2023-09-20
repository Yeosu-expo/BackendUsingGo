package main

import "encoding/json"

// 클라이언트에서 쓰는 이벤트, 클라이언트에서 서버에서 메세지를 전달 받는 이벤트
// 등에 따라 다른 조취를 취하기 위해 이벤트 클래스로 이벤트를 분류에 각각 다른 방법으로
// 대처하기 위해 사용되는 클래스
type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

const (
	// 이벤트 타입을 상수로 정의하는 곳
	EventSendMessage = "send_message"
)

type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}
