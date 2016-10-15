package main

import (
	"fmt"
	"github.com/FrostyFeet/shortener/models"
	"html"
	"net/http"
	"regexp"
)

func main() {
	http.HandleFunc("/", showUrl)
	http.HandleFunc("/api/create", writeUrl)
	http.ListenAndServe(":3100", nil)
}

func writeUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}
	url := r.URL.Query().Get("url")
	hash := r.URL.Query().Get("hash")
	api := r.URL.Query().Get("api")
	fmt.Println(url, hash)
	if hashCheck(hash) && api == "PDQRFYiWuumLTzCl6t8FzMy1d55IICih" {
		fmt.Println(models.PutUrl(url, hash))
		fmt.Println("URL inserted:", url, "with hash:", hash)
	} else {
		fmt.Println("The URL is wrong")
	}
}

func hashCheck(hash string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9]*$")
	if !re.MatchString(hash) {
		fmt.Println("Non alphanumeric characters detected")
		return false
	}
	if len(hash) > 40 || len(hash) == 0 {
		fmt.Println("Hash too long")
		return false
	} else {
		return true
	}
}

func showUrl(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	id := r.URL.Path[1:len(r.URL.Path)]
	re := regexp.MustCompile("^[a-zA-Z0-9]*$")
	if !re.MatchString(id) {
		fmt.Println("Non alphanumeric characters detected")
		http.Error(w, http.StatusText(404), 404)
		return
	}
	if len(r.URL.Path) > 40 {
		fmt.Println("Hash too long")
		http.Error(w, http.StatusText(404), 404)
		return
	}
	fmt.Println("OK: ", id)
	if id == "" {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	url, err := models.GetUrl(id)
	if err == models.ErrNoUrl {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	fmt.Println("Query:", id)
	http.Redirect(w, r, url, 301)
}

func showApi(w http.ResponseWriter, r *http.Request) {
	fmt.Println("API call:", html.EscapeString(r.URL.Path))
	fmt.Fprintf(w, "%s \n", "lalala")
}
