// Copyright (c) Sean Choi 2018. All rights reserved.

package handlers

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

func sendReceiveLambdaNic(addrStr string, port int, data string) string {
	remoteUDPAddr := fmt.Sprintf("%s:%d", addrStr, port)

	log.Printf("Connecting to server:%s \n", remoteUDPAddr)
	conn, err := net.Dial("udp4", remoteUDPAddr)
	if err != nil {
		log.Printf("Error: UDP conn error: %v\n", err)
		return ""
	}
	defer conn.Close()
	// send to socket
	log.Printf("Sending to server:%s \n", remoteUDPAddr)
	fmt.Fprintf(conn, data)
	log.Printf("Sent to server:%s \n", remoteUDPAddr)
	// listen for reply
	// TODO: Get a correct end delimiter
	message, _ := bufio.NewReader(conn).ReadString(':')
	fmt.Print("Message from server: " + message)
	return message
}

func generateResponse(w http.ResponseWriter, r *http.Request,
	body string,
	isHealth bool) (int, error) {
	if isHealth {
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
	w.WriteHeader(http.StatusOK)

	if res.Body != nil {
		// Copy the body over
		io.CopyBuffer(w, res.Body, nil)
	}

	return res.StatusCode, nil
}
