package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"fmt"
	"net/http"
)

var port *int = flag.Int("port", 8800, "Port to listen")
var methods = []string{"play", "stop", "prev", "next", "favorite", "voteup", "votedown", "voteclear", "mute", "volup", "voldown"}

var ws_array = map[int]*websocket.Conn{}
var ws_id int = 0

// Main loop
func main() {
	flag.Parse()

	// Handling available methods
	for _, method := range methods {
		http.HandleFunc("/"+method, websocketProxy)
	}

	fmt.Printf("starting server :%d\n", *port)

	http.Handle("/", websocket.Handler(joinServer))
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}

// Handle proxing request http -> ws
func websocketProxy(rw http.ResponseWriter, req *http.Request) {
	method := req.URL.Path[1:]
	fmt.Printf("sending method '%s' to connected clients\n", method)

	for _, ws := range ws_array {
		websocket.Message.Send(ws, req.URL.Path[1:])
	}
}

// Handle WS socket joining party
func joinServer(ws *websocket.Conn) {
	ws_id++
	wid := ws_id
	ws_array[wid] = ws

	fmt.Printf("client #%d connected\n", wid)

	for {
		var buf string
		err := websocket.Message.Receive(ws, &buf)

		if err != nil {
			fmt.Printf("client #%d disconnected\n", wid)
			delete(ws_array, wid)
			ws.Close()
			break
		}
	}
}
