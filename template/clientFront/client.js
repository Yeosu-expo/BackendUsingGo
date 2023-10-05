class SendOrderEvent{
    constructor(user, name, category, price){
        this.user=user;
        this.name=name;
        this.category=category;
        this.price=price;
    }
}

class Event{
    constructor(type, payload){
        this.type=type;
        this.payload=payload;
    }
}

function toggleMenu(category) {
    var container = document.getElementById(category + "-item");
    if (container.style.display == "none") {
        container.style.display = "block";
    } else {
        container.style.display = "none";
    }
}

function setMenu(name, category, price) {
    document.getElementById("order-name").value = name;
    document.getElementById("order-category").value = category;
    document.getElementById("order-price").value = price;
}

function sendEvent(eventName, payload) {
    const event = new Event(eventName, payload);
    conn.send(JSON.stringify(event));
}

function sendOrder() {
    //event.preventDefault();
    user = document.getElementById("order-user").value;
    name = document.getElementById("order-name").value;
    category = document.getElementById("order-category").value;
    price = document.getElementById("order-price").value;

    if(name!=null){
        let outgoingEvent = new SendOrderEvent(user, name, category, price);
        sendEvent("send_order", outgoingEvent);
    }
    return false;
}

window.onload = function() {
    document.getElementById("order-send").onsubmit = sendOrder;

    if(window["WebSocket"]) {
        console.log("supports websocket");
        conn =  new WebSocket("ws://" + document.location.host + "/client/ws");
    } else {
        alert("Not supporting websocket")
    }
}