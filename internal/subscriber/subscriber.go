package subscriber

import (
	"TestWebServer/internal/config"
	"TestWebServer/internal/model"
	"TestWebServer/internal/storage"
	"TestWebServer/internal/validation"
	"encoding/json"
	"fmt"

	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/nats-io/stan.go/pb"
	"log"
	"os"
	"os/signal"
)

type Subscriber struct {
	pg *storage.Database
	// cache *redis, *cache
}

func NewSubscriber(s *storage.Database) *Subscriber {
	return &Subscriber{pg: s}
}

func (s *Subscriber) Receive() {
	cfg := config.MustLoad()
	fmt.Println(cfg)
	db, err := storage.InitDatabase(*cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()
	// Connect Options.
	opts := []nats.Option{nats.Name("NATS Streaming Example Subscriber")}

	// Connect to NATS
	nc, err := nats.Connect(stan.DefaultNatsURL, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	sc, err := stan.Connect("test-cluster", "stan-sub", stan.NatsConn(nc),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, stan.DefaultNatsURL)
	}
	log.Printf("Connected to %s clusterID: [%s] clientID: [%s]\n", stan.DefaultNatsURL, "test-cluster", "test-client")

	// Process Subscriber Options.
	startOpt := stan.StartAt(pb.StartPosition_NewOnly)
	startOpt = stan.DeliverAllAvailable()

	_, err = sc.QueueSubscribe("foo", "", s.work, startOpt, stan.DurableName(""))
	if err != nil {
		sc.Close()
		log.Fatal(err)
	}
	log.Printf("Listening on [%s], clientID=[%s], qgroup=[%s] durable=[%s]\n", "test-cluster", "", "")

	// Wait for a SIGINT (perhaps triggered by user with CTRL-C)
	// Run cleanup when signal is received
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			fmt.Printf("\nReceived an interrupt, unsubscribing and closing connection...\n\n")
			sc.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}

func (s *Subscriber) work(msg *stan.Msg) {
	order := model.Order{}
	log.Printf("Received: %s\n", string(msg.Data))
	err := json.Unmarshal(msg.Data, &order)
	if err != nil {
		log.Println("Error of Unmarshalling a message ", err.Error())
		return
	}
	fmt.Printf("%#v\n", order)
	if !validation.Check(order) {
		log.Println("Validation failed")
		return
	}

	s.pg.InsertOrder(order)
	s.pg.InsertDelivery(*order.OrderUid, order.Delivery)
	s.pg.InsertPayment(*order.OrderUid, order.Payment)
	s.pg.InsertItems(*order.OrderUid, order.Items)
}
