package utils

import "net/http"

func IsHTMXRequest(req *http.Request) bool {
	return req.Header.Get("HX-Request") == "true"
}
