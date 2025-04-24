package geoserver

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	gs "github.com/camptocamp/go-geoserver/client"
)

func resourceGeoserverWmsStore() *schema.Resource {
	return &schema.Resource{
		Create: resourceGeoserverWmsStoreCreate,
		Read:   resourceGeoserverWmsStoreRead,
		Update: resourceGeoserverWmsStoreUpdate,
		Delete: resourceGeoserverWmsStoreDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGeoserverWmsStoreImport,
		},

		Schema: map[string]*schema.Schema{
			"workspace_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the workspace owning the WMS store. Used to compute the id of the resource.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the WMS store. Used to compute the id of the resource.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Description of the WMS store. Default value is empty.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Mark the WMS store as enabled. Default value is true.",
			},
			"default": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Mark the WMS store as default. Default value is false.",
			},
			"disable_connection_on_failure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Don't try to connect to remote server if failure occurs. Default value is false.",
			},
			"capabilities_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "URL of the remote WMS server capability URL.",
			},
			"max_connections": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     6,
				Description: "Number of maximum parallel connections allowed to the remote server. Default value is 6",
			},
			"read_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60,
				Description: "Number of seconds before considering a read request in timeout. Default value is 60.",
			},
			"connection_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     30,
				Description: "Number of seconds before considering a connection request in timeout. Default value is 30.",
			},
		},
	}
}

func resourceGeoserverWmsStoreCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Creating Geoserver WMS store: %s", d.Id())

	client := meta.(*Config).GeoserverClient()

	workspaceName := d.Get("workspace_name").(string)
	datastore := gs.NewWmsStore()
	datastore.Name = d.Get("name").(string)
	datastore.Description = d.Get("description").(string)
	datastore.Enabled = d.Get("enabled").(bool)
	datastore.Default = d.Get("default").(bool)
	datastore.DisableConnectionOnFailure = d.Get("disable_connection_on_failure").(bool)
	datastore.CapabilitiesUrl = d.Get("capabilities_url").(string)
	datastore.MaxConnections = d.Get("max_connections").(int)
	datastore.ReadTimeOut = d.Get("read_timeout").(int)
	datastore.ConnectTimeOut = d.Get("connection_timeout").(int)

	err := client.CreateWmStore(workspaceName, datastore)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", workspaceName, d.Get("name").(string)))

	return resourceGeoserverWmsStoreRead(d, meta)
}

func resourceGeoserverWmsStoreRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Geoserver WMS Store: %s", d.Id())

	splittedID := strings.Split(d.Id(), "/")
	workspaceName := splittedID[0]
	datastoreName := splittedID[1]

	client := meta.(*Config).GeoserverClient()

	datastore, err := client.GetWmsStore(workspaceName, datastoreName)
	if err != nil && !strings.Contains(strings.ToLower(err.Error()), "not found") {
		return err
	}

	if datastore == nil {
		d.SetId("")
		return nil
	}

	d.Set("workspace_name", datastore.Workspace.Name)
	d.Set("name", datastore.Name)
	d.Set("description", datastore.Description)
	d.Set("enabled", datastore.Enabled)
	d.Set("default", datastore.Default)
	d.Set("capabilities_url", datastore.CapabilitiesUrl)
	d.Set("max_connections", datastore.MaxConnections)
	d.Set("read_timeout", datastore.ReadTimeOut)
	d.Set("connection_timeout", datastore.ConnectTimeOut)

	return nil
}

func resourceGeoserverWmsStoreDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Geoserver WMS Store: %s", d.Id())

	splittedID := strings.Split(d.Id(), "/")
	workspaceName := splittedID[0]
	datastoreName := splittedID[1]

	client := meta.(*Config).GeoserverClient()

	err := client.DeleteWmsStore(workspaceName, datastoreName, true)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceGeoserverWmsStoreUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Geoserver WMS Store: %s", d.Id())

	splittedID := strings.Split(d.Id(), "/")
	workspaceName := splittedID[0]
	datastoreName := splittedID[1]

	client := meta.(*Config).GeoserverClient()

	datastore := gs.NewWmsStore()
	datastore.Name = d.Get("name").(string)
	datastore.Description = d.Get("description").(string)
	datastore.Enabled = d.Get("enabled").(bool)
	datastore.Default = d.Get("default").(bool)
	datastore.DisableConnectionOnFailure = d.Get("disable_connection_on_failure").(bool)
	datastore.CapabilitiesUrl = d.Get("capabilities_url").(string)
	datastore.MaxConnections = d.Get("max_connections").(int)
	datastore.ReadTimeOut = d.Get("read_timeout").(int)
	datastore.ConnectTimeOut = d.Get("connection_timeout").(int)

	err := client.UpdateWmsStore(workspaceName, datastoreName, datastore)
	if err != nil {
		return err
	}

	return nil
}

func resourceGeoserverWmsStoreImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splittedID := strings.Split(d.Id(), "/")
	workspaceName := splittedID[0]
	datastoreName := splittedID[1]

	d.SetId(d.Id())
	d.Set("workspace_name", workspaceName)
	d.Set("name", datastoreName)

	log.Printf("[INFO] Importing Geoserver WMS Store `%s` from workspace `%s`", datastoreName, workspaceName)

	err := resourceGeoserverWmsStoreRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}

