package handler

import (
	"log"
)

type Broker struct {
	Notifier       chan []byte          // Events are pushed to this channel by the main events-gathering routing
	NewClients     chan chan []byte     // new client connections
	ClosingClients chan chan []byte     // closed client connections
	Clients        map[chan []byte]bool // holds currently open connections
}

func NewBroker() (broker *Broker) {
	// Instantiate the broker
	broker = &Broker{
		Notifier:       make(chan []byte, 10000),
		NewClients:     make(chan chan []byte),
		ClosingClients: make(chan chan []byte),
		Clients:        make(map[chan []byte]bool),
	}

	// Set it running - listening and broadcasting events
	go broker.listen()
	return broker
}

func (broker *Broker) listen() {
	for {
		select {
		case s := <-broker.NewClients:
			// A client has connected
			// Register their message channel
			broker.Clients[s] = true
			log.Printf("Client added. %d registered clients", len(broker.Clients))

		case s := <-broker.ClosingClients:
			// A client has deattached and we want to
			// stop sending them messages
			delete(broker.Clients, s)
			log.Printf("Removed client. %d registered clients", len(broker.Clients))

		case event := <-broker.Notifier:
			// We got a new event from the outside
			// Send event to all connected clients
			for clientMessageChan := range broker.Clients {
				clientMessageChan <- event
			}
		}
	}
}
