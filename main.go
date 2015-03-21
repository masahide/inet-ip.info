package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/masahide/go-libGeoIP"
	"gopkg.in/yaml.v2"
)

type HostInfo struct {
	IP              string
	CountryCode     string
	CountryName     string
	Accept          []string
	AcceptEncoding  []string
	AcceptLanguage  []string
	UserAgent       []string
	Via             []string
	XForwardedFor   []string
	XForwardedPort  []string
	XForwardedProto []string
	RequestURI      string
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
	http.HandleFunc("/yml", yamlPrint)
	http.HandleFunc("/yaml", yamlPrint)
	http.HandleFunc("/toml", tomlPrint)
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
	loc := geoIp.GetLocationByIP(ip)
	info := HostInfo{
		IP:          ip,
		CountryCode: loc.CountryCode,
		CountryName: loc.CountryName,
	}

	info.RequestURI = req.RequestURI

	info.Via, _ = req.Header["Via"]
	info.Accept, _ = req.Header["Accept"]
	info.UserAgent, _ = req.Header["UserAgent"]

	info.XForwardedFor, _ = req.Header["X-Forwarded-For"]
	info.XForwardedPort, _ = req.Header["X-Forwarded-Port"]
	info.XForwardedProto, _ = req.Header["X-Forwarded-Proto"]

	info.AcceptEncoding, _ = req.Header["Accept-Encoding"]
	info.AcceptLanguage, _ = req.Header["AcceptLanguage"]

	return info
}

func root(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, getIp(req))
}
func ip(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte(getIp(req)))
}

func jsonPrint(res http.ResponseWriter, req *http.Request) {
	j, err := json.Marshal(getInfo(req))
	if err != nil {
		log.Println(err)
		return
	}
	res.Write(j)
}

func jsonIndentPrint(res http.ResponseWriter, req *http.Request) {
	j, err := json.MarshalIndent(getInfo(req), "", " ")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Fprintln(res, string(j))
}

func yamlPrint(res http.ResponseWriter, req *http.Request) {
	y, err := yaml.Marshal(getInfo(req))
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Fprintln(res, string(y))
}
func tomlPrint(res http.ResponseWriter, req *http.Request) {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(getInfo(req)); err != nil {
		log.Println(err)
		return
	}
	fmt.Fprintln(res, buf.String())
}
