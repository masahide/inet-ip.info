package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/masahide/go-libGeoIP"
	"gopkg.in/yaml.v2"
)

type HostInfo struct {
	IP              string
	Hostname        string
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

type TemplateParams struct {
	HostInfo
	Json string
}

var geoIp *libgeo.GeoIP
var tpl *template.Template

func init() {
	var err error
	geoIp, err = libgeo.Initialize(MustAsset("data/GeoIP.dat"))
	if err != nil {
		panic(err)
	}
	tpl = template.Must(template.New("index").Parse(string(MustAsset("data/index.html"))))
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

func checkUa(ua string) bool {
	lua := strings.ToLower(ua)
	if strings.Contains(lua, "curl") {
		return true
	}
	return strings.Contains(lua, "wget")
}

func getIp(req *http.Request) string {
	ips := req.Header["X-Forwarded-For"]
	return ips[len(ips)-1]
}

func getInfo(req *http.Request) HostInfo {
	ip := getIp(req)
	loc := geoIp.GetLocationByIP(ip)
	hostname, _ := net.LookupAddr(ip)
	info := HostInfo{
		IP:          ip,
		Hostname:    fmt.Sprintf("%s", hostname),
		CountryCode: loc.CountryCode,
		CountryName: loc.CountryName,
	}

	info.RequestURI = req.RequestURI

	info.Via, _ = req.Header["Via"]
	info.Accept, _ = req.Header["Accept"]
	info.UserAgent, _ = req.Header["User-Agent"]

	info.XForwardedFor, _ = req.Header["X-Forwarded-For"]
	info.XForwardedPort, _ = req.Header["X-Forwarded-Port"]
	info.XForwardedProto, _ = req.Header["X-Forwarded-Proto"]

	info.AcceptEncoding, _ = req.Header["Accept-Encoding"]
	info.AcceptLanguage, _ = req.Header["Accept-Language"]

	return info
}

func getJson(req *http.Request) string {
	j, err := json.Marshal(getInfo(req))
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(j)
}
func getJsonIndent(req *http.Request) []byte {
	j, err := json.MarshalIndent(getInfo(req), "", " ")
	if err != nil {
		log.Println(err)
		return []byte{}
	}
	return j
}

func getYaml(req *http.Request) []byte {
	y, err := yaml.Marshal(getInfo(req))
	if err != nil {
		log.Println(err)
		return []byte{}
	}
	return y
}
func getToml(req *http.Request) []byte {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(getInfo(req)); err != nil {
		log.Println(err)
		return []byte{}
	}
	return buf.Bytes()
}

func jsonPrint(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, getJson(req))
	go PageView(req)
}
func jsonIndentPrint(res http.ResponseWriter, req *http.Request) {
	res.Write(getJsonIndent(req))
	go PageView(req)
}
func yamlPrint(res http.ResponseWriter, req *http.Request) {
	res.Write(getYaml(req))
	go PageView(req)
}
func tomlPrint(res http.ResponseWriter, req *http.Request) {
	res.Write(getToml(req))
	go PageView(req)
}

func root(res http.ResponseWriter, req *http.Request) {
	ua, _ := req.Header["User-Agent"]
	defer func() { go PageView(req) }()
	if checkUa(fmt.Sprintln(ua)) {
		fmt.Fprintln(res, getIp(req))
		return
	}
	p := TemplateParams{
		HostInfo: getInfo(req),
		Json:     getJson(req),
	}
	tpl.Execute(res, p)
}
func ip(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte(getIp(req)))
	go PageView(req)
}
