package temporal_client

import (
	"log"

	"github.com/knoxgao67/temporal-dispatcher-demo/config"
	"go.temporal.io/sdk/client"
)

var c client.Client

func GetClient() client.Client {
	return c
}

// Init Call after config.Init.
func Init() {
	var err error
	c, err = client.Dial(client.Options{
		HostPort: config.GetConfig().TemporalHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client.", err)
	}
}

func Close() {
	c.Close()
}
