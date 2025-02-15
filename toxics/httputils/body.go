package httputils

import (
	"io"
	"net/http"
	"strings"
)

var StatusBodyTemplate = map[int]string{
	// status200 default nginx status page
	200: `<html><head><title>200 Status OK</title></head>
		<body><center><h1>200 Status OK</h1></body></html>`,
	// status201 default nginx error page
	201: `<html><head><title>201 Created</title></head>
		<body><center><h1>301 Moved Permanently</h1></body></html>`,
	// status202 default nginx status page
	202: `<html><head><title>202 Accepted</title></head>
		<body><center><h1>202 Accepted</h1></body></html>`,
	// status204 default nginx status page
	204: `<html><head><title>204 No Content</title></head>
		<body><center><h1>204 No Content</h1></body></html>`,
	// status301 default nginx error page
	301: `<html><head><title>301 Moved Permanently</title></head>
		<body><center><h1>301 Moved Permanently</h1></body></html>`,
	// status302 default nginx error page
	302: `<html><head><title>302 Found</title></head>
		<body><center><h1>302 Found</h1></body></html>`,
	// status304 default nginx error page
	304: `<html><head><title>304 Not Modified</title></head>
		<body><center><h1>304 Not Modified</h1></body></html>`,
	// status305 default nginx error page
	305: `<html><head><title>305 Use Proxy</title></head>
		<body><center><h1>305 Use Proxy</h1></body></html>`,
	// status305 default nginx error page
	307: `<html><head><title>307 Temporary Redirect</title></head>
		<body><center><h1>307 Temporary Redirect</h1></body></html>`,
	// status400 default nginx error page
	400: `<html><head><title>400 Bad Request</title></head>
		<body><center><h1>400 Bad Request</h1><hr></body></html>`,
	// status401 default nginx error page
	401: `<html><head><title>401 Unauthorized</title></head>
	<body><center><h1>401 Unauthorized</h1><hr></body></html>`,
	// status403 default nginx error page
	403: `<html><head><title>403 Forbidden</title></head>
		<body><center><h1>403 Forbidden</h1><hr></body></html>`,
	// status404 default nginx error page
	404: `<html><head><title>404 Not Found</title></head>
		<body><center><h1>404 Not Found</h1><hr></body></html>`,
	// status500 default nginx error page
	500: `<html><head><title>500 Internal Server Error</title></head>
		<body><center><h1>500 Internal Server Error</h1></body></html>`,
	// status501 default nginx error page
	501: `<html><head><title>501 Not Implemented</title></head>
		<body><center><h1>501 Not Implemented</h1></body></html>`,
	// status502 default nginx error page
	502: `<html><head><title>502 Bad Gateway</title></head>
		<body><center><h1>502 Bad Gateway</h1></body></html>`,
	// status503 default nginx error page
	503: `<html><head><title>503 Service Unavailable</title></head>
		<body><center><h1>503 Service Unavailable</h1></body></html>`,
	// status504 default nginx error page
	504: `<html><head><title>504 Gateway Timeout</title></head>
		<body><center><h1>504 Gateway Timeout</h1></body></html>`,
}

func EditResponseBody(r *http.Response, body, encoding, contentType string) {
	compressedBody := encodeBody(r, []byte(body), encoding)

	r.ContentLength = int64(len(body))
	r.Header.Set("Content-Type", getContentType(contentType))
	r.Body = io.NopCloser(strings.NewReader(string(compressedBody)))
}

func getContentType(contentType string) string {
	if contentType == "" {
		contentType = "text/plain"
	}
	return contentType
}

func encodeBody(r *http.Response, body []byte, encoding string) []byte {
	r.Header.Set("Content-Encoding", "")
	switch encoding {
	case "gzip":
		compressedBody, err := Gzip(body)
		if err == nil {
			r.Header.Set("Content-Encoding", "gzip")
			return compressedBody
		}
	case "deflate":
		compressedBody, err := Deflate(body)
		if err == nil {
			r.Header.Set("Content-Encoding", "deflate")
			return compressedBody
		}
	}
	return body
}
