# ICQ Bot API

## Installation

Go get: `go get gopkg.in/icq.v1`

Go mod / Go dep: `import "gopkg.in/icq.v1"`


## Working

Methods:

* SendMessage
* UploadFile

Webhooks to get updates

## Example

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"gopkg.in/icq.v1"
)

func main() {
	// New API object
	b := icq.NewAPI(os.Getenv("ICQ_TOKEN"))

	// Send message
	r, err := b.SendMessage("429950", "Hello, world!")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(r.State)

	// Send file
	f, err := os.Open("./example/icq.png")
	defer f.Close()
	if err != nil {
		log.Fatalln(err)
	}
	file, err := b.UploadFile("icq.png", f)
	if err != nil {
		log.Fatalln(err)
	}
	b.SendMessage("429950", file)

	// Webhook usage
	updates := make(chan icq.Update)
	errors := make(chan error)
	osSignal := make(chan os.Signal, 1)

	m := http.NewServeMux()
	m.HandleFunc("/webhook", b.GetWebhookHandler(updates, errors)) // Webhook sets here

	h := &http.Server{Addr: ":8080", Handler: m}
	go func() {
		log.Fatalln(h.ListenAndServe())
	}()
	signal.Notify(osSignal, os.Interrupt)
	signal.Notify(osSignal, os.Kill)
	for {
		select {
		case u := <-updates:
			log.Println("Incomming message", u)
			b.SendMessage(u.Update.From.ID, fmt.Sprintf("You sent me: %s", u.Update.Text))
			// ... process ICQ updates ...
		case err := <-errors:
			log.Fatalln(err)
		case sig := <-osSignal:
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			h.Shutdown(ctx)
			log.Fatalln("OS signal:", sig.String())
		}
	}
}
```
