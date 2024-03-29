package geoserver

import (
	"crypto/tls"
	"log"
	"net/http"

	gs "github.com/camptocamp/go-geoserver/client"
)

// Config is the configuration parameters for the Geoserver
type Config struct {
	URL                string
	GwcURL             string
	Username           string
	Password           string
	InsecureSkipVerify bool
}

// GeoserverClient creates a Geoserver client scoped to the global API
func (c *Config) GeoserverClient() *gs.Client {
	tspt := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: c.InsecureSkipVerify,
		},
	}

	client := &gs.Client{
		URL:      c.URL,
		Username: c.Username,
		Password: c.Password,
		HTTPClient: &http.Client{
			Transport: tspt,
		},
	}

	log.Printf("[INFO] Geoserver Client configured")

	return client
}

// Client creates a Geoserver client scoped to the global API
func (c *Config) GwcClient() *gs.Client {
	tspt := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: c.InsecureSkipVerify,
		},
	}

	client := &gs.Client{
		URL:      c.GwcURL,
		Username: c.Username,
		Password: c.Password,
		HTTPClient: &http.Client{
			Transport: tspt,
		},
	}

	log.Printf("[INFO] GeoWebCache Client configured")

	return client
}
