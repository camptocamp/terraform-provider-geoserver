package geoserver

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	gs "github.com/hishamkaram/geoserver"
)

func resourceGeoserverDatastore() *schema.Resource {
	return &schema.Resource{
		Create: resourceGeoserverDatastoreCreate,
		Read:   resourceGeoserverDatastoreRead,
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
				ForceNew: true,
			},
			"host": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"db_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"db_user": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"db_pass": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceGeoserverDatastoreCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Creating Geoserver Datastore: %s", d.Id())

	client := meta.(*Config).Client()

	workspaceName := d.Get("workspace_name").(string)
	datastore := gs.DatastoreConnection{
		Name:   d.Get("name").(string),
		Host:   d.Get("host").(string),
		Port:   d.Get("port").(int),
		DBName: d.Get("db_name").(string),
		DBUser: d.Get("db_user").(string),
		DBPass: d.Get("db_pass").(string),
		Type:   d.Get("type").(string),
	}

	_, err := client.CreateDatastore(datastore, workspaceName)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s.%s", workspaceName, d.Get("name").(string)))

	return resourceGeoserverDatastoreRead(d, meta)
}

func resourceGeoserverDatastoreRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Geoserver Datastore: %s", d.Id())

	splittedID := strings.Split(d.Id(), ".")
	workspaceName := splittedID[0]
	datastoreName := splittedID[1]

	client := meta.(*Config).Client()

	datastore, err := client.GetDatastoreDetails(workspaceName, datastoreName)
	if err != nil {
		return err
	}

	d.Set("workspace_name", datastore.Workspace.Name)
	d.Set("name", datastore.Name)

	for _, entry := range datastore.ConnectionParameters.Entry {
		if entry.Key == "host" {
			d.Set("host", entry.Value)
		}
		if entry.Key == "port" {
			d.Set("port", entry.Value)
		}
		if entry.Key == "database" {
			d.Set("db_name", entry.Value)
		}
		if entry.Key == "user" {
			d.Set("db_user", entry.Value)
		}
		if entry.Key == "passwd" {
			d.Set("db_pass", entry.Value)
		}
		if entry.Key == "dbtype" {
			d.Set("type", entry.Value)
		}
	}

	return nil
}

func resourceGeoserverDatastoreDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Geoserver Datastore: %s", d.Id())

	splittedID := strings.Split(d.Id(), ".")
	workspaceName := splittedID[0]
	datastoreName := splittedID[1]

	client := meta.(*Config).Client()

	_, err := client.DeleteDatastore(workspaceName, datastoreName, true)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceGeoserverDatastoreImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splittedID := strings.Split(d.Id(), ".")
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
