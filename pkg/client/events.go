package client

import (
	"log"
	"os"
	"os/signal"

	"k8s.io/kubernetes/pkg/watch"
)

// EventLoop process events in infinitive loop, apply handler function to each event
// Stop when receive interrupt signal
func (c *Client) EventLoop(w watch.Interface, handler func(watch.Event) error) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	defer signal.Stop(signals)

	for {
		select {
		case event, ok := <-w.ResultChan():
			if !ok {
				return
			}
			if err := handler(event); err != nil {
				log.Println(err)
				w.Stop()
			}
		case <-signals:
			log.Println("Catched signal, quit normally.")
			w.Stop()
		}
	}

}
