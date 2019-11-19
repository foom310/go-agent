package nrnats_test

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/newrelic/go-agent/v3/integrations/nrnats"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func currentTransaction() *newrelic.Transaction { return nil }

func ExampleStartPublishSegment() {
	nc, _ := nats.Connect(nats.DefaultURL)
	txn := currentTransaction()
	subject := "testing.subject"

	// Start the Publish segment
	seg := nrnats.StartPublishSegment(txn, nc, subject)
	err := nc.Publish(subject, []byte("Hello World"))
	if nil != err {
		panic(err)
	}
	// Manually end the segment
	seg.End()
}

func ExampleStartPublishSegment_defer() {
	nc, _ := nats.Connect(nats.DefaultURL)
	txn := currentTransaction()
	subject := "testing.subject"

	// Start the Publish segment and defer End till the func returns
	defer nrnats.StartPublishSegment(txn, nc, subject).End()
	m, err := nc.Request(subject, []byte("request"), time.Second)
	if nil != err {
		panic(err)
	}
	fmt.Println("Received reply message:", string(m.Data))
}

var clusterID, clientID string

// StartPublishSegment can be used with a NATS Streamming Connection as well
// (https://github.com/nats-io/stan.go).  Use the `NatsConn()` method on the
// `stan.Conn` interface (https://godoc.org/github.com/nats-io/stan#Conn) to
// access the `nats.Conn` object.
func ExampleStartPublishSegment_stan() {
	sc, _ := stan.Connect(clusterID, clientID)
	txn := currentTransaction()
	subject := "testing.subject"

	defer nrnats.StartPublishSegment(txn, sc.NatsConn(), subject).End()
	sc.Publish(subject, []byte("Hello World"))
}
