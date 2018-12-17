// Copyright (c) Sean Choi 2018. All rights reserved.

package handlers

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

func sendReceiveLambdaNic(addrStr string, port int, data string) string {
	remoteUDPAddr := net.UDPAddr{IP: net.ParseIP(addrStr), Port: port}

	log.Printf("Connecting to server:%s \n", remoteUDPAddr.String())
	conn, err := net.DialUDP("udp4", nil, &remoteUDPAddr)
	if err != nil {
		log.Printf("Error: UDP conn error: %v\n", err)
		return ""
	}
	defer conn.Close()

	// send to socket
	log.Printf("Sending to server:%s \n", remoteUDPAddr.String())
	n, err := conn.Write([]byte(data))
	if err != nil {
		log.Printf("Error in sending to server\n")
		return ""
	}
	log.Printf("Sent %d bytes to server:%s\n", n, remoteUDPAddr.String())
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	msg := make([]byte, 512)
	n, err = conn.Read(msg)
	if err != nil {
		log.Printf("Error in receiving from server\n")
		return ""
	}
	fmt.Printf("Message from server: %d bytes: %s\n", n, string(msg[:n]))
	return string(msg[:n])
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
