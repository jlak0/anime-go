package api

import (
	"net/http"
)

func Serve() {
	fs := http.FileServer(http.Dir("./dist"))
	http.Handle("/", fs)
	http.HandleFunc("/hello", jsonHandler)
	http.HandleFunc("/group", groupHandler)
	http.HandleFunc("/animes", animesHandler)

	err := http.ListenAndServe(":8099", nil)

	if err != nil {
		panic(err)
	}
}
