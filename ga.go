package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/pborman/uuid"
)

const (
	endpoint = "http://www.google-analytics.com/collect"
	v        = "1"
)

var (
	defaultVlas   = url.Values{}
	namespaceUUID = uuid.UUID{}
)

func init() {
	namespaceUUID = uuid.Parse(namespace)
}

type Measurement struct {
	V   string
	Tid string
	Cid string
	Uid string
	Sc  string
	Uip string
	Ua  string
	Dr  string
	Sr  string
	Vp  string
	De  string
	Sd  string
	Ul  string

	T string // pageview,event,social,timing,exception

	Dl string

	Dh     string
	Dp     string
	Dt     string
	Cd     string
	LinkId string

	An   string
	Aid  string
	Av   string
	Aiid string

	// Event
	Ec string
	Ea string
	El string
	Ev int

	// social
	Sn string
	Sa string
	St string

	// timing
	Utc string
	Utv string
	Utt int
	Utl string
	Plt string
	Dns int
	Pdt int
	Rtt int
	Tcp int
	Srt int

	// exception
	Exd string
	Exf bool
}

func mkCid(ip string, req *http.Request) string {
	b := []byte(ip + fmt.Sprintf("%s", req.UserAgent()))
	return uuid.NewSHA1(namespaceUUID, b).String()
}

func PageView(req *http.Request) error {
	ip := req.Header.Get("X-Forwarded-For")

	vals := url.Values{}
	vals.Add("v", v)
	vals.Add("tid", tid)
	vals.Add("cid", mkCid(ip, req))
	vals.Add("uip", ip)
	vals.Add("ua", req.UserAgent())

	vals.Add("t", "pageview")
	vals.Add("dr", req.Referer())
	vals.Add("dh", "inet-ip.info")
	vals.Add("dp", req.RequestURI)

	//vals.Add("t", "event")
	//vals.Add("ec", "cli")
	//vals.Add("ea", "get")
	//vals.Add("el", "json")
	//vals.Add("ev", "200")
	//log.Printf("vals:%# v", vals)
	_, err := http.PostForm(endpoint, vals)
	return err

}
