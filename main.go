package main

import (
	"bytes"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	url    = `ws://localhost:8080/guacamole/websocket-tunnel?token=996D0343862534374488C7AFEAEF08AB09908AA5BA3BFA7333B5B25ED3DE4812&GUAC_DATA_SOURCE=postgresql&GUAC_ID=0d029c1f-d63d-45ed-a2b6-9a5ea9eef127&GUAC_TYPE=a&GUAC_WIDTH=2248&GUAC_HEIGHT=990&GUAC_DPI=192&GUAC_TIMEZONE=America%2FChicago&GUAC_AUDIO=audio%2FL8&GUAC_AUDIO=audio%2FL16&GUAC_IMAGE=image%2Fjpeg&GUAC_IMAGE=image%2Fpng&GUAC_IMAGE=image%2Fwebp`
	cookie = "PUT COOKIES HERE"
	load   = 100
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	wg := sync.WaitGroup{}
	wg.Add(load)

	log.Println("Starting", load, "clients")

	for i := 0; i < load; i++ {
		go func(j int) {
			connect(j)
			wg.Done()
		}(i)
		time.Sleep(time.Second)
	}

	wg.Wait()
}

func connect(i int) {
	headers := http.Header{"Cookie": []string{
		cookie,
	}}

	log.Println("Dialing", i)
	conn, res, err := websocket.DefaultDialer.Dial(url, headers)
	if err != nil {
		if res != nil {
			data, _ := ioutil.ReadAll(res.Body)
			log.Println(i, string(data))
		}
		panic(err)
	}
	defer conn.Close()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Println(i, err)
			return
		}

		if bytes.Contains(data, []byte("error")) || bytes.Contains(data, []byte("disconnect")) {
			log.Println(string(data))
			return
		} else if bytes.HasPrefix(data, []byte("4.sync")) {
			if err = conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Panicln(i, err)
			}
		}
	}
}
