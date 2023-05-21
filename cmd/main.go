package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Client struct {
	user_id string
	channel chan []byte
}

func main() {
	database := make(map[string]*Client)

	handler := http.NewServeMux()

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)
		user_id := r.URL.Query().Get("id")

		go func() {
			client := database[user_id]
			log.Println("client is...", client)
			client.channel <- []byte("data from -> " + time.Now().String() + "sending to " + user_id + "\n\n")
		}()

		w.Write([]byte("data from -> " + time.Now().String() + "sending to " + user_id))
	})

	handler.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)
		user_id := r.URL.Query().Get("id")
		client, ok := database[user_id]

		if ok {
			log.Println("Client already exists", client)
		} else {
			client = &Client{
				user_id: user_id,
				channel: make(chan []byte, 1),
			}
			database[user_id] = client
			log.Println("Adding new Client", database[user_id])
		}

		flasher, ok := w.(http.Flusher)

		if !ok {
			http.Error(w, "Streaming not supported!", http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)

		for data := range client.channel {
			log.Println("Incoming data: ", string(data))
			fmt.Fprintf(w, "data: %s\n\n", string(data))
			flasher.Flush()
		}
	})

	server := http.Server{
		Addr:    "0.0.0.0:8888",
		Handler: handler,
	}

	log.Println("Server ruuning on http://0.0.0.0:8888")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error while running server...", err)
	}
}
