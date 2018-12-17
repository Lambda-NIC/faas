// Copyright (c) Sean Choi 2018. All rights reserved.

package handlers

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"
)

const udpPacketSize = 10

func sendReceiveLambdaNic(addrStr string, port int, data string) string {
	var wg sync.WaitGroup
	var inbound string
	remoteUDPAddr := net.UDPAddr{IP: net.ParseIP(addrStr), Port: port}

	log.Printf("Listing to port:%d \n", 2222)
	conn, err := net.ListenPacket("udp4", ":2222")
	if err != nil {
		log.Printf("Error: UDP conn error: %v", err)
		return ""
	}
	defer conn.Close()

	wg.Add(2)
	go func() {
		defer wg.Done()
		_, err := conn.WriteTo([]byte(data), &remoteUDPAddr)
		if err != nil {
			log.Printf("Error: UDP write error: %v", err)
		} else {
			log.Printf("Wrote: %s to %s:%d", data, addrStr, port)
		}
	}()

	go func() {
		defer wg.Done()
		b := make([]byte, udpPacketSize)
		for {
			n, _, err := conn.ReadFrom(b)
			if err != nil {
				log.Printf("Error: UDP read error: %v", err)
				continue
			}
			b2 := make([]byte, udpPacketSize)
			copy(b2, b)
			inbound = string(b2[:n])
			return
		}
	}()
	wg.Wait()
	return inbound
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
