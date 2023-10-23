package handler

import (
	"log"
)

type Client struct {
	channel chan []byte
	id      string
}

type Event struct {
	data  []byte
	jobId string
}

type Broker struct {
	EventNotifier  chan Event        // Events are pushed to this channel by the main events-gathering routing
	NewClients     chan Client       // new client connections
	ClosingClients chan Client       // closed client connections
	Clients        map[Client]string // holds currently open connections
}

func NewBroker() (broker *Broker) {
	// Instantiate the broker
	broker = &Broker{
		EventNotifier:  make(chan Event, 1),
		NewClients:     make(chan Client),
		ClosingClients: make(chan Client),
		Clients:        make(map[Client]string),
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
			broker.Clients[s] = "testjobid"
			log.Printf("Client added. %d registered clients", len(broker.Clients))

		case s := <-broker.ClosingClients:
			// A client has deattached and we want to
			// stop sending them messages
			delete(broker.Clients, s)
			log.Printf("Removed client. %d registered clients", len(broker.Clients))

		case event := <-broker.EventNotifier:
			// Send event to all connected clients that match the job id
			for client := range broker.Clients {
				if client.id == event.jobId {
					client.channel <- event.data
				}
			}
		}
	}
}
