package pages

import (
	"html/template"
	"net/http"

	weberr "shorturl/internal/web/err"
)

type HomeData struct {
	HasError     bool
	HasResult    bool
	HasForm      bool
	BaseURL      string
	ShortURL     string
	ErrorMessage string
}

func (data HomeData) Render(w http.ResponseWriter, r *http.Request) {
	home, err := template.ParseFiles("home.html")
	if err != nil {
		weberr.InternalServerError(w, r, "cannot open home.html")
		return
	}

	err = home.Execute(w, data)
	if err != nil {
		weberr.InternalServerError(w, r, "error while rendering home.html")
		return
	}
}
