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
	reliable     = flag.Bool("reliable", true, "Wait for the publisher confirmation before exiting")
)

func init() {
	flag.Parse()
}

func main() {
	if err := publish(*uri, *exchangeName, *exchangeType, *routingKey, *body, *reliable); err != nil {
		log.Fatalf("%s", err)
	}
	log.Printf("published %dB OK", len(*body))
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

	log.Printf("Declare exchange, publishig %dB body (%q)", len(body), body)
	if err = ch.Publish(
		exchangeType,
		routingKey,
		false,
		false,
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(body),
			DeliveryMode:    amqp.Transient,
			Priority:        0,
		}); err != nil {
		return fmt.Errorf("Exchange publish: %s", err)
	}
	return nil
}

func confirmOne(confirms <-chan amqp.Confirmation) {
	log.Printf("Waiting for confirmation of the publishing")
	if confirmed := <-confirms; confirmed.Ack {
		log.Printf("Confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
	} else {
		log.Printf("Failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}
