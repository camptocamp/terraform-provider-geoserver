package geoserver

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GEOSERVER_URL", ""),
				Description: descriptions["url"],
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GEOSERVER_USERNAME", ""),
				Description: descriptions["username"],
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GEOSERVER_PASSWORD", ""),
				Description: descriptions["password"],
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["insecure"],
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"geoserver_workspace":   resourceGeoserverWorkspace(),
			"geoserver_datastore":   resourceGeoserverDatastore(),
			"geoserver_featuretype": resourceGeoserverFeatureType(),
			"geoserver_style":       resourceGeoserverStyle(),
			"geoserver_layergroup":  resourceGeoserverLayerGroup(),
			"geoserver_resource":    resourceGeoserverResource(),
		},

		DataSourcesMap: map[string]*schema.Resource{},

		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"url":      "The Geoserver URL",
		"username": "Username to use for connection",
		"password": "Password to use for connection",
		"insecure": "Whether to verify the server's SSL certificate",
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	return &Config{
		URL:                d.Get("url").(string),
		Username:           d.Get("username").(string),
		Password:           d.Get("password").(string),
		InsecureSkipVerify: d.Get("insecure").(bool),
	}, nil
}
