package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", root)
	http.HandleFunc("/ip", ip)
	http.HandleFunc("/json", jsonPrint)
	http.HandleFunc("/json/indent", jsonIndentPrint)
	fmt.Println("listening...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}

func root(res http.ResponseWriter, req *http.Request) {
	ip(res, req)
	res.Write([]byte("\n"))
}
func ip(res http.ResponseWriter, req *http.Request) {
	ips := req.Header["X-Forwarded-For"]
	length := len(ips)
	fmt.Fprintf(res, "%s", ips[length-1])
}

func jsonPrint(res http.ResponseWriter, req *http.Request) {
	j, _ := json.Marshal(req)
	fmt.Fprintln(res, string(j))
}

func jsonIndentPrint(res http.ResponseWriter, req *http.Request) {
	j, _ := json.MarshalIndent(req, "", " ")
	fmt.Fprintln(res, string(j))
}
