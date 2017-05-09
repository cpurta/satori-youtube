package satori

import (
	"encoding/json"
	"log"

	"github.com/cpurta/satori/satori-youtube/cmd/internal/config"
	"github.com/satori-com/satori-rtm-sdk-go/rtm"
	"github.com/satori-com/satori-rtm-sdk-go/rtm/auth"
)

type SatoriPublisher struct {
	AppKey   string
	Endpoint string
	Channel  string
	Role     string
	Secret   string
	Client   *rtm.RTM
	publish  chan json.RawMessage
	shutdown bool
}

func NewPublisher(config *config.Config, pub chan json.RawMessage) *SatoriPublisher {
	authorization := auth.New(config.SatoriRole, config.SatoriSecret)

	client, _ := rtm.New(config.SatoriEndpoint, config.SatoriAppKey, rtm.Options{
		AuthProvider: authorization,
	})

	return &SatoriPublisher{
		AppKey:   config.SatoriAppKey,
		Endpoint: config.SatoriEndpoint,
		Channel:  config.SatoriChannel,
		Role:     config.SatoriRole,
		Secret:   config.SatoriSecret,
		Client:   client,
		publish:  pub,
		shutdown: false,
	}
}

func (publisher *SatoriPublisher) Start() {
	publisher.Client.Start()
}

func (publisher *SatoriPublisher) Shutdown() {
	publisher.shutdown = true
}

func (publisher *SatoriPublisher) Publish() {
	for !publisher.shutdown {
		select {
		case message := <-publisher.publish:
			if publisher.Client.IsConnected() {
				publisher.Client.Publish(publisher.Channel, message)
			}
		default:
			// do nothing
		}
	}

	publisher.Client.Stop()

	log.Println("Publisher stopped publishing")
}
