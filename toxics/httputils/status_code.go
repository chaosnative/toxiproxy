package httputils

import (
	"net/http"
)

var statusText = map[int]string{
	200: "200 OK",
	201: "201 Created",
	202: "202 Accepted",
	204: "204 No Content",
	300: "300 Multiple Choices",
	301: "301 Moved Permanently",
	302: "302 Found",
	304: "304 Not Modified",
	305: "305 Use Proxy",
	307: "307 Temporary Redirect",
	400: "400 Bad Request",
	401: "401 Unauthorized",
	403: "403 Forbidden",
	404: "404 Not Found",
	500: "500 Internal Server Error",
	501: "501 Not Implemented",
	502: "502 Bad Gateway",
	503: "503 Service Unavailable",
	504: "504 Gateway Timeout",
}

// SetHttpStatusCode sets the status code of the response.
func SetHttpStatusCode(r *http.Response, statusCode int) {
	if _, exists := statusText[statusCode]; statusCode >= 200 && statusCode < 600 && exists {
		r.StatusCode = statusCode
		r.Status = statusText[statusCode]
	}
	// if the status code is not recognized, do not change it
}

func SetErrorResponseBody(r *http.Response, statusCode int) {
	if _, exists := StatusBodyTemplate[statusCode]; statusCode >= 200 && statusCode < 600 && exists {
		EditResponseBody(r, StatusBodyTemplate[statusCode], "", "text/html")
	}
}
