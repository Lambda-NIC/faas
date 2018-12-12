// Copyright (c) Sean Choi 2018. All rights reserved.

package handlers

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

func generateResponse(w http.ResponseWriter, r *http.Request, isHealth bool) (int, error) {
	var body string
	if isHealth {
		body = "Hello world"
	} else {
		body = "OK"
	}
	res := &http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          ioutil.NopCloser(bytes.NewBufferString(body)),
		ContentLength: int64(len(body)),
		Request:       r,
		Header:        make(http.Header, 0),
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	copyHeaders(w.Header(), &res.Header)

	// Write status code
	w.WriteHeader(res.StatusCode)

	if res.Body != nil {
		// Copy the body over
		io.CopyBuffer(w, res.Body, nil)
	}

	return res.StatusCode, nil
}
