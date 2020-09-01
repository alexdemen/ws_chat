package main

import (
	"github.com/alexdemen/ws_chat/domain"
	"github.com/alexdemen/ws_chat/handler/ws"
	"net/http"
	"time"
)

func main() {
	sender := domain.NewSender()

	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/admin/stats", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(sender.Stats()))
	})

	mux := http.NewServeMux()
	mux.Handle("/admin/", adminMux)
	mux.Handle("/ws", ws.NewHandler(sender))
	mux.HandleFunc("/message", func(writer http.ResponseWriter, request *http.Request) {
		text := request.URL.Query().Get("text")
		sender.SendMessage(domain.Message{Text: text})
		writer.WriteHeader(http.StatusOK)
	})

	s := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := s.ListenAndServe(); err != nil {

	}
}
