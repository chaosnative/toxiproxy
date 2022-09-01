package httputils_test

import (
	"io"
	"testing"

	"github.com/Shopify/toxiproxy/v2/toxics/httputils"
)

var status500 = httputils.StatusBodyTemplate[500]

func TestSetHttpStatusCodeWithCorrectCode(t *testing.T) {
	resp := createHttpResponse("")
	code := 500
	status := "500 Internal Server Error"
	defer resp.Body.Close()

	AssertStatusCodeNotEqual(t, resp.StatusCode, code)
	AssertBodyNotEqual(t, []byte(resp.Status), []byte(status))

	httputils.SetHttpStatusCode(resp, code)

	AssertStatusCodeEqual(t, resp.StatusCode, code)
	AssertBodyEqual(t, []byte(resp.Status), []byte(status))
}

func TestSetHttpStatusCodeWithIncorrectCode(t *testing.T) {
	resp := createHttpResponse("")
	prevCode := resp.StatusCode
	prevStatus := resp.Status
	defer resp.Body.Close()

	httputils.SetHttpStatusCode(resp, 615)

	AssertStatusCodeEqual(t, resp.StatusCode, prevCode)
	AssertBodyEqual(t, []byte(resp.Status), []byte(prevStatus))
}

func TestSetResponseBodyWithBody(t *testing.T) {
	resp := createHttpResponse("Everything Okay")
	checkBody := "Everything not okay"
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	AssertBodyNotEqual(t, body, []byte(checkBody))

	httputils.EditResponseBody(resp, checkBody, "", "text/plain")

	body, _ = io.ReadAll(resp.Body)
	AssertBodyEqual(t, body, []byte(checkBody))
}

func TestSetResponseBodyWithStatusCodeBody(t *testing.T) {
	resp := createHttpResponse("Everything Okay")
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	AssertBodyNotEqual(t, body, []byte(status500))

	httputils.SetErrorResponseBody(resp, 500)

	body, _ = io.ReadAll(resp.Body)
	AssertBodyEqual(t, body, []byte(status500))
}

func AssertStatusCodeEqual(t *testing.T, statusCode, expectedStatusCode int) {
	if statusCode != expectedStatusCode {
		t.Errorf("Response status code {%v} not equal to expected status code {%v}.",
			statusCode, expectedStatusCode)
	}
}

func AssertStatusCodeNotEqual(t *testing.T, respStatusCode, expectedStatusCode int) {
	if respStatusCode == expectedStatusCode {
		t.Errorf("Response status code {%v} equal to expected status code {%v}.",
			respStatusCode, expectedStatusCode)
	}
}
