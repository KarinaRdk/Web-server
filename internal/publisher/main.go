package main

import (
	"log"
	"os"
	"sync"
	"time"

	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

func main() {

	// Connect Options.
	opts := []nats.Option{nats.Name("NATS Streaming Example Publisher")}

	// Connect to NATS
	nc, err := nats.Connect(stan.DefaultNatsURL, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	sc, err := stan.Connect("test-cluster", "stan-pub", stan.NatsConn(nc))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, stan.DefaultNatsURL)
	}
	defer sc.Close()

	msg, err := os.ReadFile("../../broker/messages/model.json")
	if err != nil {
		log.Println("Cannot read file:", err)
	}

	subj := "foo"

	ch := make(chan bool)
	var glock sync.Mutex
	var guid string
	acb := func(lguid string, err error) {
		glock.Lock()
		log.Printf("Received ACK for guid %s\n", lguid)
		defer glock.Unlock()
		if err != nil {
			log.Fatalf("Error in server ack for guid %s: %v\n", lguid, err)
		}
		if lguid != guid {
			log.Fatalf("Expected a matching guid in ack callback, got %s vs %s\n", lguid, guid)
		}
		ch <- true
	}

	glock.Lock()
	guid, err = sc.PublishAsync(subj, msg, acb)

	if err != nil {
		log.Fatalf("Error during async publish: %v\n", err)
	}
	glock.Unlock()
	if guid == "" {
		log.Fatal("Expected non-empty guid to be returned.")
	}
	log.Printf("Published [%s] : '%s' [guid: %s]\n", subj, msg, guid)

	select {
	case <-ch:
		break
	case <-time.After(5 * time.Second):
		log.Fatal("timeout")
	}

}
