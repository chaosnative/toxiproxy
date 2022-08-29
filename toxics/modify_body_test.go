package toxics_test

import (
	"io"
	"net"
	"net/http"
	"testing"

	"github.com/Shopify/toxiproxy/v2/toxics"
)

func TestToxicModifiesHTTPResponseBody(t *testing.T) {
	http.HandleFunc("/", echoHelloWorld)

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

	AssertBodyNotEqual(t, body, []byte(status500))

	proxy.Toxics.AddToxicJson(ToxicToJson(
		t, "", "modify_body", "downstream",
		&toxics.ModifyBodyToxic{Body: status500},
	))

	resp, err = http.Get("http://" + proxy.Listen)
	if err != nil {
		t.Error("Failed to connect to proxy", err)
	}

	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)

	AssertBodyEqual(t, body, []byte(status500))
}
