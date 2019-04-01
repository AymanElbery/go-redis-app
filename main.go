package main

import (
	"fmt"
	"go-redis-app/models"
	"net/http"
	"strconv"
)

func main() {

	// Use the showAlbum handler for all requests with a URL path beginning
	// '/album'.
	http.HandleFunc("/album", showAlbum)
	http.HandleFunc("/like", addLike)
	http.HandleFunc("/popular", listPopular)
	http.ListenAndServe(":3000", nil)
}

// showAlbum Handler ...
func showAlbum(w http.ResponseWriter, r *http.Request) {

	// Unless the request is using the GET method, return a 405 'Method Not
	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		http.Error(w, http.StatusText(405), 405)
		return
	}

	// Retrieve the id from the request URL query string ...
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Validate that the id is a valid integer ...
	if _, err := strconv.Atoi(id); err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Call the FindAlbum() function passing in the user-provided id ...
	album, err := models.FindAlbum(id)
	if err == models.ErrNoAlbum {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Write the album details as plain text to the client.
	fmt.Fprintf(w, "%s by %s: £%.2f [%d likes] \n", album.Title, album.Artist, album.Price, album.Likes)
}

// addLike Handler ...
func addLike(w http.ResponseWriter, r *http.Request) {

	// Unless the request is using the POST method, return a 405 'Method Not
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, http.StatusText(405), 405)
		return
	}

	// Retreive the id from the POST request body ...
	id := r.PostFormValue("id")
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Validate that the id is a valid integer ...
	if _, err := strconv.Atoi(id); err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Call the IncrementLikes() function passing in the user-provided id ...
	err := models.IncrementLikes(id)
	if err == models.ErrNoAlbum {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Redirect the client to the GET /ablum route to see impact of increasing ...
	http.Redirect(w, r, "/album?id="+id, 303)
}

// listPopular Handler ...
func listPopular(w http.ResponseWriter, r *http.Request) {

	// Unless the request is using the GET method, return a 405 'Method Not ...
	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		http.Error(w, http.StatusText(405), 405)
		return
	}

	// Call the FindTopThree() function ...
	abs, err := models.FindTopThree()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Loop through the 3 albums, writing the details as a plain text list ...
	for i, ab := range abs {
		fmt.Fprintf(w, "%d) %s by %s: £%.2f [%d likes] \n", i+1, ab.Title, ab.Artist, ab.Price, ab.Likes)
	}
}
