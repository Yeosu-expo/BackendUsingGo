package main

import (
	"kiosk/kioskPack"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

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

	http.ListenAndServe(":8080", router)
}
