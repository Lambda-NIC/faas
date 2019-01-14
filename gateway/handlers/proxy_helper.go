// Copyright (c) Sean Choi 2018. All rights reserved.

package handlers

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/Lambda-NIC/faas/gateway/requests"
)

/*
const localPort = 2222

var localIP = GetOutboundIP()

// GetOutboundIP Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
*/

func sendReceiveLambdaNic(addrStr string, port int,
	jobID int, data string) string {
	remoteUDPAddr := net.UDPAddr{IP: net.ParseIP(addrStr), Port: port}
	//localUDPAddr := net.UDPAddr{IP: localIP, Port: localPort}

	//log.Printf("Connecting from %s to %s \n",
	//	localUDPAddr.String(), remoteUDPAddr.String())
	//conn, err := net.DialUDP("udp4", &localUDPAddr, &remoteUDPAddr)
	conn, err := net.DialUDP("udp4", nil, &remoteUDPAddr)
	if err != nil {
		log.Printf("Error: UDP conn error: %v\n", err)
		return ""
	}
	defer conn.Close()

	// send to socket
	//log.Printf("Sending from %s to %s \n",
	//	localUDPAddr.String(), remoteUDPAddr.String())
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, uint32(jobID))
	dataBytes := append(bs, []byte(data)...)
	_, err = conn.Write(dataBytes)
	if err != nil {
		log.Printf("Error in sending to server\n")
		return ""
	}

	//log.Printf("Sent %d bytes to server:%s\n", len(dataBytes),
	// remoteUDPAddr.String())
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	msg := make([]byte, 32)
	n, err := conn.Read(msg)
	if err != nil {
		log.Printf("Error in receiving from server %s\n", err.Error())
		return ""
	}

	//fmt.Printf("Message from server: %d bytes: %s\n", n, string(msg[:n]))
	return string(msg[:n])
}

func generateResponse(w http.ResponseWriter, r *http.Request,
	serviceName string,
	imageName string,
	body string,
	isHealth bool) (int, error) {
	res := &http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		ContentLength: int64(len(body)),
		Request:       r,
		Header:        make(http.Header, 0),
	}
	if isHealth {
		function := requests.Function{
			Name:              serviceName,
			Replicas:          4,
			Image:             imageName,
			AvailableReplicas: 4,
			InvocationCount:   0,
		}
		functionBytes, _ := json.Marshal(function)
		res.Body = ioutil.NopCloser(bytes.NewBuffer(functionBytes))
		res.Header.Set("Content-Type", "application/json")
		res.ContentLength = int64(len(functionBytes))
	} else {
		res.Body = ioutil.NopCloser(bytes.NewBufferString(body))
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
