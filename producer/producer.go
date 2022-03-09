package main

import (
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

var (
	uri          = flag.String("uri", "amqp://guest:guest@localhost:5672/", "AMQP uri")
	exchangeName = flag.String("exchage", "test-exchange", "Durable exchange name")
	exchangeType = flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
	routingKey   = flag.String("key", "test-key", "AMQP routing key")
	body         = flag.String("body", "foobar", "Body of message")
	reliable     = flag.String("reliable", "true", "Wait for the publisher confirmation before exiting")
)

func init() {
	flag.Parse()
}

func main() {

}

func publish(amqpURI, exchange, exchangeType, routingKey, body string, reliable bool) error {
	log.Printf("Dialing %q", amqpURI)
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		return fmt.Errorf("Dial: %s", err)
	}
	defer conn.Close()

	log.Printf("Got connection, getting channel")
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	log.Printf("Got chanel, declaring %q Exchange (%q)", exchangeType, exchange)
	if err := ch.ExchangeDeclare(
		exchange,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("Exchange declare: %s", err)
	}

	if reliable {
		log.Printf("Enabling publishing confirms.")
		if err := ch.Confirm(false); err != nil {
			return fmt.Errorf("Channel could not be put into confirm mode: %s", err)
		}
		confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))
		defer confirmOne(confirms)
	}

}
