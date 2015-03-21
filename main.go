package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/masahide/go-libGeoIP"
)

type HostInfo struct {
	IP       string
	Location *libgeo.Location
	Request  *http.Request
}

var geoIp *libgeo.GeoIP

func init() {
	var err error
	geoIp, err = libgeo.Initialize(MustAsset("data/GeoIP.dat"))
	if err != nil {
		panic(err)
	}
}

//go:generate go-bindata data/
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

func getIp(req *http.Request) string {
	ips := req.Header["X-Forwarded-For"]
	return ips[len(ips)-1]
}

func getInfo(req *http.Request) HostInfo {
	ip := getIp(req)
	return HostInfo{
		IP:       ip,
		Location: geoIp.GetLocationByIP(ip),
		Request:  req,
	}
}

func root(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte(getIp(req) + "\n"))
}
func ip(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte(getIp(req)))
}

func jsonPrint(res http.ResponseWriter, req *http.Request) {
	j, _ := json.Marshal(getInfo(req))
	fmt.Fprintln(res, string(j))
}

func jsonIndentPrint(res http.ResponseWriter, req *http.Request) {
	j, _ := json.MarshalIndent(getInfo(req), "", " ")
	fmt.Fprintln(res, string(j))
}
