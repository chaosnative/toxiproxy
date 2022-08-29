package httputils_test

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Shopify/toxiproxy/v2/toxics/httputils"
)

func createHttpResponse(body string) *http.Response {
	resp := http.Response{
		StatusCode: 200,
		Status:     "200 OK",
		Header: http.Header{
			"Foo": []string{"Bar"},
		},
		Body: io.NopCloser(strings.NewReader(body)),
	}
	return &resp
}

func TestEditResponseBody(t *testing.T) {
	resp := createHttpResponse("World Hello")

	body, _ := io.ReadAll(resp.Body)

	AssertBodyEqual(t, body, []byte("World Hello"))

	httputils.EditResponseBody(resp, "Hello World")

	body, _ = io.ReadAll(resp.Body)

	AssertBodyEqual(t, body, []byte("Hello World"))
}

func AssertBodyEqual(t *testing.T, respBody, expectedBody []byte) {
	if !bytes.Equal(respBody, expectedBody) {
		t.Errorf("Response body {%v} not equal to expected body {%v}.",
			string(respBody), string(expectedBody))
	}
}

func AssertBodyNotEqual(t *testing.T, respBody, expectedBody []byte) {
	if bytes.Equal(respBody, expectedBody) {
		t.Errorf("Response body {%v} equal to expected body {%v}.",
			string(respBody), string(expectedBody))
	}
}
