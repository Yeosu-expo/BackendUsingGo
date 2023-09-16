package kioskPack

import (
	"database/sql"
	"html/template"
	"net/http"
)

func OpenAdminHtml(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./template/admin.html"))
	tmpl.Execute(w, nil)
}

func OpenClientHtml(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./template/client.html"))

	db, err := sql.Open("mysql", "root:9250@tcp(127.0.0.1:3306)/kiosk")
	defer db.Close()
	CheckErr(err)

	rows, err := db.Query("SELECT * FROM menu")
	defer rows.Close()
	CheckErr(err)

	menusData := new(MenusData)
	for rows.Next() {
		var id int
		var name string
		var category sql.NullString
		var price sql.NullString

		err = rows.Scan(&id, &name, &category, &price)
		CheckErr(err)
		var categoryy string
		var pricee string
		if category.Valid {
			categoryy = category.String
		} else {
			categoryy = "null"
		}
		if price.Valid {
			pricee = price.String
		} else {
			pricee = "null"
		}
		tmpMenu := new(Menu)
		tmpMenu.Name = name
		tmpMenu.Category = categoryy
		tmpMenu.Price = pricee
		menusData.Menus = append(menusData.Menus, *tmpMenu)
	}

	tmpl.Execute(w, menusData)
}
