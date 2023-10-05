package main

import (
	"kiosk/kioskPack"
	sp "kiosk/socketPack"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func setupAPI(manager *sp.Manager, router *mux.Router) {
	router.HandleFunc("/", kioskPack.OpenChatHtml).Methods("GET")
	router.HandleFunc("/ws", manager.ServeWS)
}

func setupOrderAPIforClient(manager *sp.Manager, router *mux.Router) {
	router.HandleFunc("/client", kioskPack.OpenClientHtml).Methods("GET")
	router.HandleFunc("/client/ws", manager.ServeWS)
}

func setupOrderAPIforAdmin(manager *sp.Manager, router *mux.Router) {
	router.HandleFunc("/admin", kioskPack.OpenAdminHtml).Methods("GET")
	router.HandleFunc("/admin/ws", manager.ServeWSforAdmin)
}

func main() {
	router := mux.NewRouter()

	//router.HandleFunc("/admin", kioskPack.OpenAdminHtml).Methods("GET")
	//router.HandleFunc("/client", kioskPack.OpenClientHtml).Methods("GET")
	router.HandleFunc("/admin", kioskPack.PostAndStoreJson).Methods("POST")

	manager := sp.NewManager()
	setupAPI(manager, router)
	setupOrderAPIforClient(manager, router)
	setupOrderAPIforAdmin(manager, router)

	http.ListenAndServe(":8080", router)
}
