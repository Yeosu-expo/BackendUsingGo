package main

import (
	"kiosk/kioskPack"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/admin", kioskPack.OpenAdminHtml).Methods("GET")
	router.HandleFunc("/client", kioskPack.OpenClientHtml).Methods("GET")
	router.HandleFunc("/admin", kioskPack.PostAndStoreJson).Methods("POST")
	router.HandleFunc("/chat", kioskPack.OpenChatHtml).Methods("GET")

	http.ListenAndServe("localhost:8080", router)
}
