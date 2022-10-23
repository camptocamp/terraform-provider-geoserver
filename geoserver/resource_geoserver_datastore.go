package geoserver

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	gs "github.com/camptocamp/go-geoserver/client"
)

func resourceGeoserverDatastore() *schema.Resource {
	return &schema.Resource{
		Create: resourceGeoserverDatastoreCreate,
		Read:   resourceGeoserverDatastoreRead,
		Update: resourceGeoserverDatastoreUpdate,
		Delete: resourceGeoserverDatastoreDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGeoserverDatastoreImport,
		},

		Schema: map[string]*schema.Schema{
			"workspace_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"default": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"connection_params": {
				Type:     schema.TypeMap,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceGeoserverDatastoreCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Creating Geoserver Datastore: %s", d.Id())

	connectionParameters := []*gs.DatastoreConnectionParameter{}
	for key, value := range d.Get("connection_params").(map[string]interface{}) {
		connectionParameters = append(connectionParameters, &gs.DatastoreConnectionParameter{
			Key:   key,
			Value: value.(string),
		})
	}

	client := meta.(*Config).Client()

	workspaceName := d.Get("workspace_name").(string)
	datastore := &gs.Datastore{
		Name:                 d.Get("name").(string),
		Description:          d.Get("description").(string),
		Enabled:              d.Get("enabled").(bool),
		Default:              d.Get("default").(bool),
		ConnectionParameters: connectionParameters,
	}

	err := client.CreateDatastore(workspaceName, datastore)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", workspaceName, d.Get("name").(string)))

	return resourceGeoserverDatastoreRead(d, meta)
}

func resourceGeoserverDatastoreRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Geoserver Datastore: %s", d.Id())

	splittedID := strings.Split(d.Id(), "/")
	workspaceName := splittedID[0]
	datastoreName := splittedID[1]

	client := meta.(*Config).Client()

	datastore, err := client.GetDatastore(workspaceName, datastoreName)
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

	connectionParameters := map[string]string{}
	for _, entry := range datastore.ConnectionParameters {
		connectionParameters[entry.Key] = entry.Value
	}

	d.Set("connection_params", connectionParameters)

	return nil
}

func resourceGeoserverDatastoreDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Geoserver Datastore: %s", d.Id())

	splittedID := strings.Split(d.Id(), "/")
	workspaceName := splittedID[0]
	datastoreName := splittedID[1]

	client := meta.(*Config).Client()

	err := client.DeleteDatastore(workspaceName, datastoreName, true)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceGeoserverDatastoreUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Geoserver Datastore: %s", d.Id())

	splittedID := strings.Split(d.Id(), "/")
	workspaceName := splittedID[0]
	datastoreName := splittedID[1]

	connectionParameters := []*gs.DatastoreConnectionParameter{}
	for key, value := range d.Get("connection_params").(map[string]interface{}) {
		connectionParameters = append(connectionParameters, &gs.DatastoreConnectionParameter{
			Key:   key,
			Value: value.(string),
		})
	}

	client := meta.(*Config).Client()

	err := client.UpdateDatastore(workspaceName, datastoreName, &gs.Datastore{
		Name:                 d.Get("name").(string),
		Description:          d.Get("description").(string),
		Enabled:              d.Get("enabled").(bool),
		Default:              d.Get("default").(bool),
		ConnectionParameters: connectionParameters,
	})
	if err != nil {
		return err
	}

	return nil
}

func resourceGeoserverDatastoreImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splittedID := strings.Split(d.Id(), "/")
	workspaceName := splittedID[0]
	datastoreName := splittedID[1]

	d.SetId(d.Id())
	d.Set("workspace_name", workspaceName)
	d.Set("name", datastoreName)

	log.Printf("[INFO] Importing Geoserver Datastore `%s` in workspace `%s`", datastoreName, workspaceName)

	err := resourceGeoserverDatastoreRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
