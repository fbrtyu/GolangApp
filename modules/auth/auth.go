package auth

import (
	"fmt"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.ParseForm()

	fmt.Println(r.FormValue("login"))
	fmt.Println(r.FormValue("password"))

	w.Write([]byte("OK"))
}
