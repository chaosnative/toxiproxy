package toxics_test

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"testing"

	"github.com/chaosnative/toxiproxy/v2/toxics"
	"github.com/chaosnative/toxiproxy/v2/toxics/httputils"
)

// status500 default nginx error page.
var status500 = httputils.StatusBodyTemplate[500]

func echoHelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func TestToxicModifiesHTTPStatusCode(t *testing.T) {
	http.HandleFunc("/status-code", echoHelloWorld)

	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal("Failed to create TCP server", err)
	}

	go http.Serve(ln, nil)
	defer ln.Close()

	proxy := NewTestProxy("test", ln.Addr().String())
	proxy.Start()
	defer proxy.Stop()

	resp, err := http.Get("http://" + proxy.Listen)
	if err != nil {
		t.Error("Failed to connect to proxy", err)
	}
	defer resp.Body.Close()

	AssertStatusCodeNotEqual(t, resp.StatusCode, 500)

	proxy.Toxics.AddToxicJson(ToxicToJson(
		t, "", "status_code", "downstream",
		&toxics.StatusCodeToxic{StatusCode: 500, ModifyResponseBody: 0},
	))

	resp, err = http.Get("http://" + proxy.Listen)
	if err != nil {
		t.Error("Failed to connect to proxy", err)
	}
	defer resp.Body.Close()

	AssertStatusCodeEqual(t, resp.StatusCode, 500)
}

func TestToxicStatusCodeEmptyBodyModifyBodyTrue(t *testing.T) {
	http.HandleFunc("/status-empty-body", echoHelloWorld)

	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal("Failed to create TCP server", err)
	}

	go http.Serve(ln, nil)
	defer ln.Close()

	proxy := NewTestProxy("test", ln.Addr().String())
	proxy.Start()
	defer proxy.Stop()

	resp, err := http.Get("http://" + proxy.Listen)
	if err != nil {
		t.Error("Failed to connect to proxy", err)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	AssertStatusCodeNotEqual(t, resp.StatusCode, 500)
	AssertBodyNotEqual(t, body, []byte(status500))

	proxy.Toxics.AddToxicJson(ToxicToJson(
		t, "", "status_code", "downstream",
		&toxics.StatusCodeToxic{
			StatusCode: 500, ResponseBody: "",
			ModifyResponseBody: 1},
	))

	resp, err = http.Get("http://" + proxy.Listen)
	if err != nil {
		t.Error("Failed to connect to proxy", err)
	}

	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	AssertStatusCodeEqual(t, resp.StatusCode, 500)
	AssertBodyEqual(t, body, []byte(status500))
}

func TestToxicStatusCodeWithBodyModifyBodyTrue(t *testing.T) {
	http.HandleFunc("/status-body", echoHelloWorld)

	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal("Failed to create TCP server", err)
	}

	go http.Serve(ln, nil)
	defer ln.Close()

	proxy := NewTestProxy("test", ln.Addr().String())
	proxy.Start()
	defer proxy.Stop()

	resp, err := http.Get("http://" + proxy.Listen)
	if err != nil {
		t.Error("Failed to connect to proxy", err)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	AssertStatusCodeNotEqual(t, resp.StatusCode, 500)
	AssertBodyNotEqual(t, body, []byte(status500))

	proxy.Toxics.AddToxicJson(ToxicToJson(
		t, "", "status_code", "downstream",
		&toxics.StatusCodeToxic{
			StatusCode: 500, ResponseBody: "Hello World",
			ModifyResponseBody: 1},
	))

	resp, err = http.Get("http://" + proxy.Listen)
	if err != nil {
		t.Error("Failed to connect to proxy", err)
	}

	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	AssertStatusCodeEqual(t, resp.StatusCode, 500)
	AssertBodyEqual(t, body, []byte("Hello World"))
}

func TestToxicStatusCodeBodyModifyBodyFalse(t *testing.T) {
	http.HandleFunc("/status-body-false", echoHelloWorld)

	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal("Failed to create TCP server", err)
	}

	go http.Serve(ln, nil)
	defer ln.Close()

	proxy := NewTestProxy("test", ln.Addr().String())
	proxy.Start()
	defer proxy.Stop()

	resp, err := http.Get("http://" + proxy.Listen)
	if err != nil {
		t.Error("Failed to connect to proxy", err)
	}

	defer resp.Body.Close()

	prevBody, _ := io.ReadAll(resp.Body)

	AssertStatusCodeNotEqual(t, resp.StatusCode, 500)
	AssertBodyNotEqual(t, prevBody, []byte(status500))

	proxy.Toxics.AddToxicJson(ToxicToJson(
		t, "", "status_code", "downstream",
		&toxics.StatusCodeToxic{StatusCode: 500, ModifyResponseBody: 0},
	))

	resp, err = http.Get("http://" + proxy.Listen)
	if err != nil {
		t.Error("Failed to connect to proxy", err)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	AssertStatusCodeEqual(t, resp.StatusCode, 500)
	AssertBodyEqual(t, body, prevBody)
}

func TestUnsupportedStatusCode(t *testing.T) {
	http.HandleFunc("/status-unsupported", echoHelloWorld)

	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal("Failed to create TCP server", err)
	}

	go http.Serve(ln, nil)
	defer ln.Close()

	proxy := NewTestProxy("test", ln.Addr().String())
	proxy.Start()
	defer proxy.Stop()

	resp, err := http.Get("http://" + proxy.Listen)
	if err != nil {
		t.Error("Failed to connect to proxy", err)
	}

	defer resp.Body.Close()

	statusCode := resp.StatusCode
	initialBody, _ := io.ReadAll(resp.Body)

	proxy.Toxics.AddToxicJson(ToxicToJson(
		t, "", "status_code", "downstream",
		&toxics.StatusCodeToxic{StatusCode: 1000, ModifyResponseBody: 1},
	))

	resp, err = http.Get("http://" + proxy.Listen)
	if err != nil {
		t.Error("Failed to connect to proxy", err)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	AssertStatusCodeEqual(t, resp.StatusCode, statusCode)
	AssertBodyEqual(t, body, initialBody)
}

func AssertStatusCodeEqual(t *testing.T, respStatusCode, expectedStatusCode int) {
	if respStatusCode != expectedStatusCode {
		t.Errorf("Response status code {%v} not equal to expected status code {%v}.",
			respStatusCode, expectedStatusCode)
	}
}

func AssertStatusCodeNotEqual(t *testing.T, respStatusCode, expectedStatusCode int) {
	if respStatusCode == expectedStatusCode {
		t.Errorf("Response status code {%v} equal to expected status code {%v}.",
			respStatusCode, expectedStatusCode)
	}
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
