package routes

import (
	"fmt"
	"net/http"
	"path/filepath"
)

func StartServer(port uint32) {
	http.HandleFunc("/", routeDefault)
	http.HandleFunc("/poker", routePoker)
	http.HandleFunc("/poker/", routePoker)
	http.HandleFunc("/api/poker/create", routeApiPokerCreate)
	http.HandleFunc("/api/poker/", routeApiPokerGame)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

var mimeTypes = map[string]string{
	".json": "application/json",
	".mjs":  "application/javascript",
	".js":   "application/javascript",
	".html": "text/html",
	".css":  "text/css",
	".png":  "image/png",
	".svg":  "image/svg+xml",
}

func routeDefault(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/poker", http.StatusPermanentRedirect)
	}
	path := "web" + r.URL.Path
	ext := filepath.Ext(path)
	if mimeType, ok := mimeTypes[ext]; ok {
		w.Header().Set("Content-Type", mimeType)
	}
	http.ServeFile(w, r, path)
}

func routePoker(w http.ResponseWriter, r *http.Request) {

}

func routeApiPokerCreate(w http.ResponseWriter, r *http.Request) {

}

func routeApiPokerGame(w http.ResponseWriter, r *http.Request) {

}
