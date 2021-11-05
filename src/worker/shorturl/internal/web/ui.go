package web

import (
	"log"
	"net/http"

	app "shorturl/internal/app"
	weberr "shorturl/internal/web/err"
	pages "shorturl/internal/web/pages"
)

func UIHandler(w http.ResponseWriter, r *http.Request) {
	// common handler for all requests

	if len(r.URL.Path) == 1 {

		// default requested path is "/"
		if r.Method == "GET" {
			homeGET(w, r)
		} else if r.Method == "POST" {
			homePOST(w, r)
		} else {
			weberr.InternalServerError(w, r, "unexpected method '"+r.Method+"' for '/'")
		}
		return

	} else if len(r.URL.Path) > 1 {

		// requested path contains a short URL
		surl := r.URL.Path[1:]

		if r.Method == "GET" {
			shorturlGET(w, r, surl)
		} else {
			weberr.InternalServerError(w, r, "unexpected method '"+r.Method+"' for '/"+surl+"'")
		}
		return

	}

	// shouldn't be here
	weberr.InternalServerError(w, r, "reached handler end without returning first")
}

func homeGET(w http.ResponseWriter, r *http.Request) {
	log.Printf("request GET '/' from %s", r.RemoteAddr)

	data := pages.HomeData{
		HasError:     false,
		HasResult:    false,
		HasForm:      true,
		BaseURL:      "",
		ShortURL:     "",
		ErrorMessage: "",
	}
	data.Render(w, r)

	log.Printf("response to GET '/' from %s", r.RemoteAddr)
}

func homePOST(w http.ResponseWriter, r *http.Request) {
	lurl := r.FormValue("longurl")
	log.Printf("request POST '/' from %s with '%s'", r.RemoteAddr, lurl)

	surl, err := app.DBInstance.StoreLongURL(lurl)

	if err != nil {
		// error with the input long URL
		data := pages.HomeData{
			HasError:     true,
			HasResult:    false,
			HasForm:      true,
			BaseURL:      "",
			ShortURL:     "",
			ErrorMessage: "cannot store the short URL",
		}
		data.Render(w, r)

		log.Printf("response to POST '/' from %s - short URL not stored", r.RemoteAddr)
		return
	}

	// no error, reply with the short URL
	data := pages.HomeData{
		HasError:     false,
		HasResult:    true,
		HasForm:      false,
		BaseURL:      "http://" + app.Config.AppURL,
		ShortURL:     surl,
		ErrorMessage: "",
	}
	data.Render(w, r)

	log.Printf("response to POST '/' from %s - created short URL '%s'", r.RemoteAddr, surl)
}

func shorturlGET(w http.ResponseWriter, r *http.Request, surl string) {
	log.Printf("request GET '/%s' from %s", surl, r.RemoteAddr)

	lurl, err := app.DBInstance.GetLongURL(surl)

	if err != nil {
		// error with the input long URL
		data := pages.HomeData{
			HasError:     true,
			HasResult:    false,
			HasForm:      true,
			BaseURL:      "",
			ShortURL:     "",
			ErrorMessage: "the input short URL resulted in error",
		}
		data.Render(w, r)

		log.Printf("response to GET '/%s' from %s - bad short URL", surl, r.RemoteAddr)
		return
	}

	// no error, simply redirect
	log.Printf("response to GET '/%s' from %s - redirect to '%s'", surl, r.RemoteAddr, lurl)
	http.Redirect(w, r, lurl, http.StatusFound)
}
