package main

import (
	"chat-app/chat"
	"chat-app/config"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {return true},
}

var (
	cfg config.Config
)

func init() {
	cfg = config.Cfg
}

func main() {
	http.Handle("/chat/", http.HandlerFunc(websocketHandler))
	server := http.Server{Addr: ":" + cfg.Port, Handler: nil}
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed  {
			log.Fatal("failed to start the server", err)
		}
	}()

	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT)
	<-exit

	log.Println("exit signalled")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chat.Cleanup()
	server.Shutdown(ctx)

	log.Println("chat app exited")
}

func websocketHandler(rw http.ResponseWriter, req *http.Request) {
	user := strings.TrimPrefix(req.URL.Path, "/chat/")

	peer, err := upgrader.Upgrade(rw, req, nil)
	if err != nil {
		log.Fatal("websocket conn failed", err)
	}

	chatSession := chat.NewChatSession(user, peer)
	chatSession.Start()
}