package toxics

import (
	"bufio"
	"bytes"
	"io"
	"net/http"

	"github.com/chaosnative/toxiproxy/v2/stream"
	"github.com/chaosnative/toxiproxy/v2/toxics/httputils"
)

type ModifyBodyToxic struct {
	Body             string `json:"body"`
	ContentEnconding string `json:"content_encoding"`
	ContentType      string `json:"content_type"`
}

func (t *ModifyBodyToxic) ModifyResponseBody(resp *http.Response) {
	if t.Body == "" {
		return
	}
	httputils.EditResponseBody(resp, t.Body, t.ContentEnconding, t.ContentType)
}

func (t *ModifyBodyToxic) Pipe(stub *ToxicStub) {
	buffer := bytes.NewBuffer(make([]byte, 0, 32*1024))
	writer := stream.NewChanWriter(stub.Output)
	reader := stream.NewChanReader(stub.Input)
	reader.SetInterrupt(stub.Interrupt)
	for {
		tee := io.TeeReader(reader, buffer)
		resp, err := http.ReadResponse(bufio.NewReader(tee), nil)

		if err == stream.ErrInterrupted {
			buffer.WriteTo(writer)
			return
		} else if err == io.EOF {
			stub.Close()
			return
		}
		if err != nil {
			buffer.WriteTo(writer)
		} else {
			t.ModifyResponseBody(resp)
			resp.Write(writer)
		}
		buffer.Reset()
	}
}

func init() {
	Register("modify_body", new(ModifyBodyToxic))
}
