package httputils_test

import (
	"io"
	"testing"

	"github.com/Shopify/toxiproxy/v2/toxics/httputils"
)

var status500 = httputils.StatusBodyTemplate[500]

func TestSetHttpStatusCodeWithCorrectCode(t *testing.T) {
	resp := createHttpResponse("")

	AssertStatusCodeEqual(t, resp.StatusCode, 200)
	AssertBodyEqual(t, []byte(resp.Status), []byte("200 OK"))

	httputils.SetHttpStatusCode(resp, 500)

	AssertStatusCodeEqual(t, resp.StatusCode, 500)
	AssertBodyEqual(t, []byte(resp.Status), []byte("500 Internal Server Error"))
}

func TestSetHttpStatusCodeWithIncorrectCode(t *testing.T) {
	resp := createHttpResponse("")

	AssertStatusCodeEqual(t, resp.StatusCode, 200)
	AssertBodyEqual(t, []byte(resp.Status), []byte("200 OK"))

	httputils.SetHttpStatusCode(resp, 615)

	AssertStatusCodeEqual(t, resp.StatusCode, 200)
	AssertBodyEqual(t, []byte(resp.Status), []byte("200 OK"))
}

func TestSetResponseBodyWithBody(t *testing.T) {
	resp := createHttpResponse("Everything Okay")

	body, _ := io.ReadAll(resp.Body)
	AssertBodyEqual(t, body, []byte("Everything Okay"))
	httputils.SetResponseBody(resp, 200, "Everything not okay")
	body, _ = io.ReadAll(resp.Body)
	AssertBodyEqual(t, body, []byte("Everything not okay"))
}

func TestSetResponseBodyWithStatusCodeBody(t *testing.T) {
	resp := createHttpResponse("Everything Okay")

	body, _ := io.ReadAll(resp.Body)
	AssertBodyEqual(t, body, []byte("Everything Okay"))
	httputils.SetResponseBody(resp, 500, "")
	body, _ = io.ReadAll(resp.Body)
	AssertBodyEqual(t, body, []byte(status500))
}

func AssertStatusCodeEqual(t *testing.T, statusCode, expectedStatusCode int) {
	if statusCode != expectedStatusCode {
		t.Errorf("Response status code {%v} not equal to expected status code {%v}.",
			statusCode, expectedStatusCode)
	}
}
