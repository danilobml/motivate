package handlers

import "net/http"

func getHealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
