package kioskPack

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Menu struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Price string `json:"price"`
}

func GetAndStoreJson(w http.ResponseWriter, r *http.Request) {
	menu := new(Menu)
	err := json.NewDecoder(r.Body).Decode(menu)
	if err != nil {
		log.Println(err, "line 51")
		http.Error(w, "JSON 데이터를 파싱하는 데 실패했습니다", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("mysql", "root:9250@tcp(127.0.0.1:3306)/kiosk")
	defer db.Close()
	CheckErr(err, w)

	insertQuery := "INSERT INTO menu (name, price) VALUES(?,?)"
	price, _ := strconv.Atoi(menu.Price)
	_, err = db.Exec(insertQuery, menu.Name, price)
	CheckErr(err, w)
}

func OpenHtml(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./template/admin.html"))
	tmpl.Execute(w, nil)
}

func CheckErr(err error, w http.ResponseWriter) bool {
	if err != nil {
		log.Println(err, "not25line")
		return true
	}
	return false
}
