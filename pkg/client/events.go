package client

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"k8s.io/kubernetes/pkg/api"
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

func eventFormatter(e watch.Event) string {
	apiEvent := (e.Object).(*api.Event)
	msgFmt :=
		`
Namespace: %s
Kind     : %s
Component: %s
Host     : %s
Reason   : %s
`
	return fmt.Sprintf(
		msgFmt,
		apiEvent.ObjectMeta.Namespace,
		apiEvent.InvolvedObject.Kind,
		apiEvent.Source.Component,
		apiEvent.Source.Host,
		apiEvent.Reason)
}
