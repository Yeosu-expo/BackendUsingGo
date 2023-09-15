package main

import (
	"net/http"

	"github.com/Yeosu-expo/backendUsingGo/kioskPack"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/admin", kioskPack.OpenHtml).Methods("GET")
	router.HandleFunc("/admin", kioskPack.GetAndStoreJson).Methods("POST")

	http.ListenAndServe("localhost:8080", router)
}
