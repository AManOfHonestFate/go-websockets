package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
)

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Connection") != "Upgrade" {
		return
	}
	if r.Header.Get("Upgrade") != "Websocket" {
		return
	}
	key := r.Header.Get("Sec-WebSocket-Key")
	if key == "" {
		return
	}

	sum := key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	hash := sha1.Sum([]byte(sum))
	str64 := base64.StdEncoding.EncodeToString(hash[:])

	hj, ok := w.(http.Hijacker)
	if !ok {
		return
	}
	conn, buf, err := hj.Hijack()
	if err != nil {
		return
	}
	defer conn.Close()

	_, err = buf.WriteString("HTTP/1.1 101 Switching Protocols\r\n")
	_, err = buf.WriteString("Upgrade: websocket\r\n")
	_, err = buf.WriteString("Connection: Upgrade\r\n")
	_, err = buf.WriteString("Sec-Websocket-Accept: " + str64 + "\r\n\r\n")
	err = buf.Flush()
	if err != nil {
		return
	}

	clientBuf := make([]byte, 1024)
	for {
		n, err := buf.Read(clientBuf)
		if err != nil {
			return
		}
		fmt.Println(clientBuf[:n])
	}
}

func main() {
	http.HandleFunc("/", websocketHandler)
	http.ListenAndServe("localhost:8080", nil)
}
