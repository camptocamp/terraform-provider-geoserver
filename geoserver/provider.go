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
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GEOSERVER_URL", ""),
				Description: descriptions["url"],
			},
			"gwc_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GEOWEBCACHE_URL", ""),
				Description: descriptions["gwc_url"],
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GEOSERVER_USERNAME", ""),
				Description: descriptions["username"],
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
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
			"geoserver_workspace":        resourceGeoserverWorkspace(),
			"geoserver_datastore":        resourceGeoserverDatastore(),
			"geoserver_featuretype":      resourceGeoserverFeatureType(),
			"geoserver_style":            resourceGeoserverStyle(),
			"geoserver_layergroup":       resourceGeoserverLayerGroup(),
			"geoserver_resource":         resourceGeoserverResource(),
			"geoserver_gwc_S3_blobstore": resourceGwcS3Blobstore(),
			"geoserver_gwc_gridset":      resourceGwcGridset(),
			"geoserver_gwc_wms_layer":    resourceGwcWmsLayer(),
			"geoserver_wms_store":        resourceGeoserverWmsStore(),
			"geoserver_wms_layer":        resourceGeoserverWmsLayer(),
		},

		DataSourcesMap: map[string]*schema.Resource{},

		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"url":      "The Geoserver URL",
		"gwc_url":  "The GeoWebCache URL",
		"username": "Username to use for connection",
		"password": "Password to use for connection",
		"insecure": "Whether to verify the server's SSL certificate",
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	return &Config{
		URL:                d.Get("url").(string),
		GwcURL:             d.Get("gwc_url").(string),
		Username:           d.Get("username").(string),
		Password:           d.Get("password").(string),
		InsecureSkipVerify: d.Get("insecure").(bool),
	}, nil
}
