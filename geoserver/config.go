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

func CreateConfig(URL string,
	GwcURL string,
	Username string,
	Password string,
	InsecureSkipVerify bool) (interface{}, error) {
	// test GeoServer URL
	if len(URL) > 0 {
		_, err := http.Get(URL)
		if err != nil {
			log.Print(err.Error())
			return nil, err
		}
	} else {
		log.Printf("[INFO] No GeoServer URL configured")
	}

	// test GWC URL
	if len(GwcURL) > 0 {
		_, err := http.Get(GwcURL)
		if err != nil {
			log.Print(err.Error())
			return nil, err
		}
	} else {
		log.Printf("[INFO] No GeoWebCache URL configured")
	}
	return &Config{
		URL:                URL,
		GwcURL:             GwcURL,
		Username:           Username,
		Password:           Password,
		InsecureSkipVerify: InsecureSkipVerify,
	}, nil
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
