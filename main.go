package main

import (
	"kiosk/kioskPack"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func setupAPI(manager *Manager, router *mux.Router) {
	router.Handle("/", http.FileServer(http.Dir("./chatFront")))
	router.HandleFunc("/ws", manager.serveWS)
}

func setupOrderAPIforClient(manager *Manager, router *mux.Router) {
	router.HandleFunc("/client", kioskPack.OpenClientHtml).Methods("GET")
	router.HandleFunc("/client/ws", manager.serveWS)
}

func setupOrderAPIforAdmin(manager *Manager, router *mux.Router) {
	router.HandleFunc("/admin", kioskPack.OpenAdminHtml).Methods("GET")
	router.HandleFunc("/admin/ws", manager.serveWSforAdmin)
}

func main() {
	router := mux.NewRouter()

	//router.HandleFunc("/admin", kioskPack.OpenAdminHtml).Methods("GET")
	//router.HandleFunc("/client", kioskPack.OpenClientHtml).Methods("GET")
	router.HandleFunc("/admin", kioskPack.PostAndStoreJson).Methods("POST")

	manager := NewManager()
	setupAPI(manager, router)
	setupOrderAPIforClient(manager, router)
	setupOrderAPIforAdmin(manager, router)

	http.ListenAndServe(":8080", router)
}
