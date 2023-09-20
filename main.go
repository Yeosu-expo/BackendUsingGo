package main

import (
	"kiosk/kioskPack"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func setupAPI(router *mux.Router) {
	manager := NewManager()

	router.Handle("/", http.FileServer(http.Dir("./frontend")))
	router.HandleFunc("/ws", manager.serveWS)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/admin", kioskPack.OpenAdminHtml).Methods("GET")
	router.HandleFunc("/client", kioskPack.OpenClientHtml).Methods("GET")
	router.HandleFunc("/admin", kioskPack.PostAndStoreJson).Methods("POST")

	setupAPI(router)

	http.ListenAndServe("localhost:8080", router)
}
