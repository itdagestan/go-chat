package main

import (
	"context"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"time"
)

type Server struct {
	conns map[string]*websocket.Conn
}

func NewServer() *Server {
	return &Server{
		conns: make(map[string]*websocket.Conn),
	}
}

func (s *Server) echoHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Star send")
	chatIds, ok := r.URL.Query()["chatId"]

	if !ok && len(chatIds) < 1 {
		log.Println("Not found chat ID")
		return
	}
	chatId := chatIds[0]
	log.Println("chat id " + chatId)
	conn, ok := s.conns[chatId]
	if !ok {
		log.Println("chat not found")
		c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			InsecureSkipVerify: true,
		})
		if err != nil {
			log.Println(err)
			return
		}
		defer func(c *websocket.Conn, code websocket.StatusCode, reason string) {
			err := c.Close(code, reason)
			if err != nil {
				log.Fatalln(err)
			}
		}(c, websocket.StatusInternalError, "the sky is falling")
		s.conns[chatId] = c
		conn = c
	} else {

	}
	log.Println("chat found")

	for {
		log.Println("start loop")
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
		defer cancel()

		_, message, err := conn.Read(ctx)
		log.Printf("Received error %v", err)
		if err != nil {
			break
		}
		log.Printf("Received %v", message)

		err = conn.Write(ctx, websocket.MessageText, message)
		log.Printf("Sended error %v", err)
		if err != nil {
			break
		}
		log.Printf("Sended %v", message)
	}
}

func main() {
	server := NewServer()
	address := "localhost:1234"
	http.HandleFunc("/echo", server.echoHandler)
	log.Printf("Starting server, go to http://%s/ to try it out!", address)
	http.Handle("/", http.FileServer(http.Dir("assets")))
	err := http.ListenAndServe(address, nil)
	log.Fatal(err)
}
