package geoserver

import (
	"crypto/tls"
	"log"
	"net/http"

	gs "github.com/hishamkaram/geoserver"
)

// Config is the configuration parameters for the Geoserver
type Config struct {
	URL                string
	Username           string
	Password           string
	InsecureSkipVerify bool
}

// Client creates a Geoserver client scoped to the global API
func (c *Config) Client() *gs.GeoServer {
	tspt := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: c.InsecureSkipVerify,
		},
	}

	client := gs.GetCatalog(c.URL, c.Username, c.Password)

	client.HttpClient.Transport = tspt

	log.Printf("[INFO] Geoserver Client configured")

	return client
}
