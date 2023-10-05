/* function login() {
    let formData = {
        "username": document.getElementById("username").value,
        "password": document.getElementById("password").value
    }
    fetch("login", {
        method: 'post',
        body: JSON.stringify(formData),
        mode: 'cors',
    }).then((response) => {
        if (response.ok) {
            return response.json();
        } else {
            throw 'unauthorized';
        }
    }).then((data) => {
        connectWebsocket(data.otp);
    }).catch((e) => { alert(e); });
    return false;

function connectWebsocket(otp) {
    if (window[WebSocket]) {
        console.log("supports websockets")
        conn = new WebSocket("ws://" + document.location.host + "/ws?otp=" + otp)
        conn.onopen = function (evt) {
            document.getElementById("connection-header").innerHTML = "Connect to Websocket: true";
        }
        conn.onclose = function (evt) {
            document.getElementById("connection-header").innerHTML = "Connected to Websocket: false";
        
        conn.onmessage = function (evt) {
            console.log(evt)
            const eventData = JSON.parse(evt.data);
            const event = Object.assign(new Event, eventData);
            routeEvent(event);
        }
    } else {
        alert("Not supporting websockets");
    }
} */

class SendMessageEvent {
    constructor(message, from) {
        this.message = message;
        this.from = from;
    }
}

class NewMessageEvent {
    constructor(message, from, sent) {
        this.message = message;
        this.from = from;
        this.sent = sent;
    }
}

var selectedRoom = "yet";

class Event {
    constructor(type, payload) {
        this.type = type;
        this.payload = payload;
    }
}

function routeEvent(event) {
    if (event.type === undefined) {
        alert("no 'type' field in event");
    }
    switch (event.type) {
        case "new_message":
            const messageEvent = Object.assign(new NewMessageEvent, event.payload);
            appendChatMessage(messageEvent);
            break;
        default:
            alert("unsupported message type");
            break;
    }
}

function appendChatMessage(messageEvent) {
    var date = new Date(messageEvent.sent);
    const formattedMsg = `${date.toLocaleDateString()}: ${messageEvent.message}`;
    textarea = document.getElementById("chatMessages");
    textarea.innerHTML = textarea.innerHTML + "\n" + formattedMsg;
    textarea.scrollTop = textarea.scrollHeight;
}

// JSON으로 패키징해서 다양한 정보를 다루기 쉽게 할 수 있음
function sendEvent(eventName, payload) {
    const event = new Event(eventName, payload);
    conn.send(JSON.stringify(event));
}

function changeChatRoom() {
    var newRoom = document.getElementById("chatroom");
    if (newRoom != null && newRoom.value != selectedRoom) {
        console.log(newRoom);
    }
    return false;
}

function sendMessage() {
    var newMessage = document.getElementById("message");
    if (newMessage != null) {
        let outgoingEvent = new SendMessageEvent(newMessage.value, "DMin");
        sendEvent("send_message", outgoingEvent);
    }
    return false;
}

window.onload = function () {
    document.getElementById("chatroom-selection").onsubmit = changeChatRoom;
    document.getElementById("chatroom-message").onsubmit = sendMessage;
    if (window["WebSocket"]) {
        console.log("supports websockets");
        // Connect to websocket
        conn = new WebSocket("ws://" + document.location.host + "/ws");
        
        conn.onmessage = function (evt) {
            // 받은 메세지 중에 data만 뽑음
            const eventData = JSON.parse(evt.data);
            // 뽑은 data를 Event객체로 가공
            const event = Object.assign(new Event, eventData);
            routeEvent(event);
        }
    
    } else {
        alert("Not supporting websockets");
    }
}