package demoweb

import (
	"crypto/sha512"
	"encoding/base64"
	"net/http"
	"strings"
)

func getRealIP(req *http.Request) string {
	ts := []string{"X-Forwarded-For", "x-forwarded-for", "X-FORWARDED-FOR"}
	for _, k := range ts {
		if v := req.Header.Get(k); v != "" {
			return v
		}
	}
	return ""
}

func hashSha512(s string) string {
	sum := sha512.Sum512([]byte(s))
	return base64.StdEncoding.EncodeToString(sum[:])
}

func getUserID(req *http.Request) string {
	ip := getRealIP(req)
	if ip == "" {
		ip = strings.Split(req.RemoteAddr, ":")[0]
	}
	return ip + "_" + hashSha512(req.UserAgent())[:5]
}

func urlToName(s string) string {
	ss := strings.Split(s, "_")
	suffix := ss[len(ss)-1]
	switch suffix {
	case "1":
		return "etcd1"
	case "2":
		return "etcd2"
	case "3":
		return "etcd3"
	default:
		return "unknown"
	}
}

func boldHTMLMsg(msg string) string {
	return "<br><b>########## " + msg + "</b><br>"
}
