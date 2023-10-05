package socketPack

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

const (
	// 이벤트 타입을 상수로 정의하는 곳
	EventSendMessage = "send_message"
	EventNewMessage  = "new_message"
	EventSendOrder   = "send_order"
	EventNewOrder    = "new_order"
)

// 클라이언트에서 쓰는 이벤트, 클라이언트에서 서버에서 메세지를 전달 받는 이벤트
// 등에 따라 다른 조취를 취하기 위해 이벤트 클래스로 이벤트를 분류에 각각 다른 방법으로
// 대처하기 위해 사용되는 클래스
type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type NewMessageEvent struct {
	SendMessageEvent
	Sent time.Time `json:"sent"`
}

type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}

type SendOrderEvent struct {
	User     string `json:"user"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Price    string `json:"price"`
}

type NewOrderEvent struct {
	SendOrderEvent
	Sent time.Time `json:"sent"`
}

type EventHandler func(event Event, c *Client) error

func SendMessageHandler(event Event, c *Client) error {
	// 클라에서 보내준 메세지를 받음
	var chatevent SendMessageEvent
	if err := json.Unmarshal(event.Payload, &chatevent); err != nil {
		// log도 찍고 error도 반환
		return fmt.Errorf("bad payload in request: %v", err)
	}

	// 메세지를 다른 클라에 뿌리기 전에 재 가공
	var broadMessage NewMessageEvent

	broadMessage.Sent = time.Now()
	broadMessage.Message = chatevent.Message
	broadMessage.From = chatevent.From

	// 가공된 메세지를 json으로 보내기 위해 변환
	data, err := json.Marshal(broadMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	// 최종으로 Event클래스에 값이 담긴 payload와 Event 타입을 지정함
	var outgoingEvent Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = EventNewMessage

	// 연결된 각 클라 모두에게 채널로 Event를 전송
	for client := range c.manager.clients {
		client.egress <- outgoingEvent
	}

	return nil
}

func SendOrderHandler(event Event, c *Client) error {
	var order SendOrderEvent
	if err := json.Unmarshal(event.Payload, &order); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	var newOrder NewOrderEvent
	newOrder.User = order.User
	newOrder.Name = order.Name
	newOrder.Category = order.Category
	newOrder.Price = order.Price
	newOrder.Sent = time.Now()

	data, err := json.Marshal(newOrder)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast order: %v", err)
	}

	var outgoingEvent Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = EventNewOrder

	// 클라에서 보내줘야하는 대상은 Admin밖에 없음
	for client := range c.manager.clients {
		if client.clientType == "Admin" {
			log.Println("sent well")
			client.egress <- outgoingEvent
		}
	}

	return nil
}
