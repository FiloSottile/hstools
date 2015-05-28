package main

import (
	"hstools"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ActiveState/tail"
	"golang.org/x/net/websocket"
)

var channelsLock sync.RWMutex
var channels = make(map[*websocket.Conn]chan string)

func WSServer(ws *websocket.Conn) {
	onion := make([]byte, 16)
	if _, err := io.ReadFull(ws, onion); err != nil {
		log.Println(err)
		return
	}
	log.Printf("%s", string(onion))
	res, err := hstools.OnionToDescID(string(onion), time.Now())
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("%s %s", hstools.ToBase32(res[0]), hstools.ToBase32(res[1]))
	announce := make(chan string)
	channelsLock.Lock()
	channels[ws] = announce
	channelsLock.Unlock()
	defer func() {
		channelsLock.Lock()
		delete(channels, ws)
		channelsLock.Unlock()
	}()
	for desc := range announce {
		if hstools.ToBase32(res[0]) != desc && hstools.ToBase32(res[1]) != desc {
			continue
		}
		if _, err := ws.Write([]byte("1")); err != nil {
			log.Println(err)
			return
		}
	}
}

func main() {
	defer tail.Cleanup()
	if len(os.Args) < 4 {
		log.Fatal("usage: announce cert.pem key.pem logfile")
	}

	http.Handle("/announce", websocket.Handler(WSServer))
	go func() {
		log.Fatal(http.ListenAndServeTLS("0.0.0.0:14242", os.Args[1], os.Args[2], nil))
	}()

	t, err := tail.TailFile(os.Args[3], tail.Config{
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: os.SEEK_END,
		},
		MustExist: true,
		Follow:    true,
		ReOpen:    true,
	})
	if err != nil {
		log.Fatal(err)
	}
	for line := range t.Lines {
		if strings.Contains(line.Text, "Got a v2 rendezvous descriptor request for ID") {
			i := strings.Index(line.Text, `'"`) + 2
			log.Println(line.Text[i : i+32])
			channelsLock.RLock()
			for _, ch := range channels {
				ch <- line.Text[i : i+32]
			}
			channelsLock.RUnlock()
		}
	}
	log.Fatal(t.Err())
}
