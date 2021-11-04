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

		// requested path contains a shorturl
		shorturl := r.URL.Path[1:]

		if r.Method == "GET" {
			shorturlGET(w, r, shorturl)
		} else {
			weberr.InternalServerError(w, r, "unexpected method '"+r.Method+"' for '/"+shorturl+"'")
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
	log.Printf("request POST '/' from %s with '%s'", r.RemoteAddr, r.FormValue("longurl"))

	// TODO: fetch/store the shortened url
	// tmp: mock shorturl result
	shorturl, err := "mock_url", error(nil)

	if err != nil {
		// error with the input long URL
		data := pages.HomeData{
			HasError:     true,
			HasResult:    false,
			HasForm:      true,
			BaseURL:      "",
			ShortURL:     "",
			ErrorMessage: "the input long URL resulted in error",
		}
		data.Render(w, r)

		log.Printf("response to POST '/' from %s - bad long URL", r.RemoteAddr)
		return
	}

	// no error, reply with the shortened URL
	data := pages.HomeData{
		HasError:     false,
		HasResult:    true,
		HasForm:      false,
		BaseURL:      "http://" + app.Config.AppURL,
		ShortURL:     shorturl,
		ErrorMessage: "",
	}
	data.Render(w, r)

	log.Printf("response to POST '/' from %s - shorturl '%s'", r.RemoteAddr, shorturl)
}

func shorturlGET(w http.ResponseWriter, r *http.Request, shorturl string) {
	log.Printf("request GET '/%s' from %s", shorturl, r.RemoteAddr)

	// TODO: fetch and validate the shorturl
	// tmp: mock shorturl result
	longurl, err := "http://www.google.com", error(nil)

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

		log.Printf("response to GET '/%s' from %s - bad short URL", shorturl, r.RemoteAddr)
		return
	}

	// no error, simply redirect
	log.Printf("response to GET '/%s' from %s - redirect to '%s'", shorturl, r.RemoteAddr, longurl)
	http.Redirect(w, r, longurl, http.StatusFound)
}
