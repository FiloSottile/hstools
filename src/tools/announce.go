package main

import (
	"fmt"
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

var html = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">

    <title>WebSocket Test</title>
    <script language="javascript" type="text/javascript">
        var wsUri = "wss://" + window.location.host + "/announce";

        function init() {
            websocket = new WebSocket(wsUri);
            websocket.onopen = function(evt) {
                websocket.send(window.location.hash.substr(1));
            };
            websocket.onclose = function(evt) {
                window.location.reload();
            };
            websocket.onmessage = function(evt) {
                onMessage(evt)
            };
            websocket.onerror = function(evt) {
                console.log(evt.data);
                window.location.reload();
            };
        }

        function onMessage(evt) {
            document.getElementsByTagName('body')[0].className = "on";
            console.log("on");
            setTimeout(function(){
                console.log("off");
                document.getElementsByTagName('body')[0].className = "off";
            }, 500);
        }

        window.addEventListener("load", init, false);
    </script>
    <style>
    .on {
        background-color: black;
    }
    </style>
</head>

<body></body>
</html>
`

func HTMLServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, html)
}

var channelsLock sync.RWMutex
var channels = make(map[*websocket.Conn]chan string)

func WSServer(ws *websocket.Conn) {
	onion := make([]byte, 16)
	if _, err := io.ReadFull(ws, onion); err != nil {
		log.Println(err)
		return
	}
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
		res, err := hstools.OnionToDescID(string(onion), time.Now())
		if err != nil {
			log.Println(err)
			return
		}
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
	if len(os.Args) < 4 {
		log.Fatal("usage: announce cert.pem key.pem logfile")
	}

	http.Handle("/", http.HandlerFunc(HTMLServer))
	http.Handle("/announce", websocket.Handler(WSServer))
	go func() {
		log.Fatal(http.ListenAndServeTLS(":14242", os.Args[1], os.Args[2], nil))
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
