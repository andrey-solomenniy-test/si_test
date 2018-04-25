package main

import (
	"math/rand"
	"net/http"
	"strings"
)

var iplist = []string{
	"192.169.140.100:80",
	"191.37.87.223:39880",
	"176.58.125.65:80",
	"118.139.178.67:23540",
	"111.118.169.39:43156",
	"46.105.57.149:3537",
}

func getIPPart(s string) string {
	return strings.Split(s, ":")[0]
}

func getIPFromRequest(r *http.Request) string {
	// return getIPPart(r.RemoteAddr)
	return getIPPart(iplist[rand.Intn(len(iplist))])
}
