package handler

import (
	"log"
)

type Client struct {
	id      string
	channel ClientChan
}

type ClientChan struct {
	messageChan chan []byte
	closeChan   chan bool
}

type Clients = map[string][]ClientChan

type Event struct {
	data  []byte
	jobId string
}

type Broker struct {
	EventNotifier chan Event  // Events are pushed to this channel by the main events-gathering routing
	NewClient     chan Client // new client connections
	ClosingClient chan string // closed client connections
	Clients       Clients     // holds currently open connections
}

func NewBroker() (broker *Broker) {
	// Instantiate the broker
	broker = &Broker{
		EventNotifier: make(chan Event, 1),
		NewClient:     make(chan Client),
		ClosingClient: make(chan string),
		Clients:       make(Clients),
	}

	// Set it running - listening and broadcasting events
	go broker.listen()
	return broker
}

func (broker *Broker) listen() {
	for {
		select {
		case client := <-broker.NewClient:
			// check if there are clients on that id
			clients := broker.Clients[client.id]
			if clients != nil {
				// append the channel
				clients = append(clients, client.channel)
			} else {
				// create a new client entry on that id
				clients = [](ClientChan){client.channel}
			}

			broker.Clients[client.id] = clients

			log.Printf(
				"Adding a new client for job '%s'. %d registered clients.",
				client.id,
				len(broker.Clients[client.id]),
			)

		case id := <-broker.ClosingClient:
			clients := broker.Clients[id]
			if clients == nil {
				log.Printf("Client with job id '%s' does not exist", id)
			} else {
				// send a message on close channel
				for _, client := range clients {
					client.closeChan <- true
				}
				// delete the entry from the clients
				delete(broker.Clients, id)
				log.Printf("Removing all clients for job '%s'.", id)
			}

		case event := <-broker.EventNotifier:
			// Send event to all connected clients that match the job id
			clients := broker.Clients[event.jobId]
			for _, client := range clients {
				client.messageChan <- event.data
			}
		}
	}
}
