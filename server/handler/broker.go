package handler

import (
	"log"

	"github.com/prajwalad101/vincent/server/util"
)

type Client struct {
	id      string
	channel chan []byte
}

type Clients = map[string][](chan []byte)

type Event struct {
	data  []byte
	jobId string
}

type Broker struct {
	EventNotifier chan Event  // Events are pushed to this channel by the main events-gathering routing
	NewClient     chan Client // new client connections
	ClosingClient chan Client // closed client connections
	Clients       Clients     // holds currently open connections
}

func NewBroker() (broker *Broker) {
	// Instantiate the broker
	broker = &Broker{
		EventNotifier: make(chan Event, 1),
		NewClient:     make(chan Client),
		ClosingClient: make(chan Client),
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
				clients = [](chan []byte){client.channel}
			}

			broker.Clients[client.id] = clients

			log.Printf(
				"Adding a new client for job '%s'. %d registered clients.",
				client.id,
				len(broker.Clients[client.id]),
			)

		case closingClient := <-broker.ClosingClient:
			// check if any clients exist on that id
			clients := broker.Clients[closingClient.id]
			if clients != nil {
				// find the index of the closing client
				clientIndex := util.SliceIndex(len(clients), func(i int) bool {
					return clients[i] == closingClient.channel
				})
				// remove the client from list of registered clients
				clients[clientIndex] = clients[len(clients)-1]
				clients = clients[:len(clients)-1]

				broker.Clients[closingClient.id] = clients

				log.Printf(
					"Removing client for job '%s'. %d registered clients.",
					closingClient.id,
					len(broker.Clients[closingClient.id]),
				)
			} else {
				log.Printf("Client with job id '%s' does not exist", closingClient.id)
			}

		case event := <-broker.EventNotifier:
			// Send event to all connected clients that match the job id
			clients := broker.Clients[event.jobId]
			// send message to all clients
			for _, client := range clients {
				client <- event.data
			}
		}
	}
}
