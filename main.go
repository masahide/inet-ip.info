package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", ip)
	http.HandleFunc("/json", jsonPrint)
	fmt.Println("listening...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}

func ip(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "%s", req.Header["X-Forwarded-For"])
}

func jsonPrint(res http.ResponseWriter, req *http.Request) {
	j, _ := json.MarshalIndent(req, "", " ")
	fmt.Fprintln(res, string(j))
}
