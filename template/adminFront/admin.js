class Event{
    constructor(type, payload){
        this.type=type;
        this.payload=payload;
    }
}

class NewOrderEvent {
    constructor(user, name, category, price, sent){
        this.user=user;
        this.name=name;
        this.category=category;
        this.price=price;
        this.sent=sent;
    }
}

function getJson() {
    var nameValue = document.getElementsByName("name")[0].value;
    var categoryValue = document.getElementsByName("category")[0].value;
    var priceValue = document.getElementsByName("price")[0].value;

    if (nameValue === "" || priceValue === "") {
        alert("모든 필드를 입력하세요.");
        return;
    }

    var jsonData = {
        "name": nameValue,
        "category": categoryValue,
        "price": priceValue
    };

    $.ajax({
        type: "POST",
        url: "http://localhost:8080/admin",
        data: JSON.stringify(jsonData),
        contentType: "application/json",
        dataType: "json",
    });

    event.preventDefault();
}

function clearInputFields() {
    document.getElementsByName("name")[0].value = "";
    document.getElementsByName("category")[0].value = "";
    document.getElementsByName("price")[0].value = "";
}

function appendOrderList(orderEvent) {
    var date = new Date(orderEvent.sent);
    console.log("!");
    const formattedOrder = ` 주문좌석: ${orderEvent.user}, 음식: ${orderEvent.name}, 분류: ${orderEvent.category}, 가격: ${orderEvent.price}, 주문시간: ${date.toLocaleString()}`;
    textarea = document.getElementById("orderArea");
    textarea.innerHTML = textarea.innerHTML + "\n" + formattedOrder;
    textarea.scrollTop = textarea.scrollHeight;
}

function getEvent(event) {
    // 받은 이벤트가 처리할 수 있는 이벤트타입인지 체크
    switch (event.type) {
        case "new_order":
            // 이벤트 payload를 NewOrderEvent객체로 변환
            const orderEvent = Object.assign(new NewOrderEvent, event.payload);
            // orderEvent를 admin페이지 textarea에 표시
            appendOrderList(orderEvent);
            break;
        default:
            alert("unsupported order type");
            break;
    }
}

window.onload = function() {
    if(window["WebSocket"]) {
        console.log("supports websocket")
        conn = new WebSocket("ws://"+document.location.host+"/admin/ws");
    }

    conn.onmessage = function(evt) {
        // 받은 이벤트 중에 data를 뽑음
        const eventData = JSON.parse(evt.data);
        // 뽑은 data를 event에 Event객체로 변환하여 저장
        const event = Object.assign(new Event, eventData);
        getEvent(event);
    }
}