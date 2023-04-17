package geoserver

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	gs "github.com/camptocamp/go-geoserver/client"
)

func resourceGeoserverWorkspace() *schema.Resource {
	return &schema.Resource{
		Create: resourceGeoserverWorkspaceCreate,
		Read:   resourceGeoserverWorkspaceRead,
		Update: resourceGeoserverWorkspaceUpdate,
		Delete: resourceGeoserverWorkspaceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGeoserverWorkspaceImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"isolated": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceGeoserverWorkspaceCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Creating Geoserver Workspace: %s", d.Id())

	client := meta.(*Config).GeoserverClient()

	name := d.Get("name").(string)

	err := client.CreateWorkspace(&gs.Workspace{
		Name:     name,
		Isolated: d.Get("isolated").(bool),
	}, d.Get("default").(bool))
	if err != nil {
		return err
	}

	d.SetId(name)

	return resourceGeoserverWorkspaceRead(d, meta)
}

func resourceGeoserverWorkspaceRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Geoserver Workspace: %s", d.Id())

	client := meta.(*Config).GeoserverClient()

	workspace, err := client.GetWorkspace(d.Id())
	if err != nil && !strings.Contains(strings.ToLower(err.Error()), "not found") {
		return err
	}

	if workspace == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", workspace.Name)
	d.Set("isolated", workspace.Isolated)

	return nil
}

func resourceGeoserverWorkspaceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Geoserver Workspace: %s", d.Id())

	client := meta.(*Config).GeoserverClient()

	err := client.DeleteWorkspace(d.Id(), true)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceGeoserverWorkspaceUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Geoserver Workspace: %s", d.Id())

	client := meta.(*Config).GeoserverClient()

	err := client.UpdateWorkspace(d.Id(), &gs.Workspace{
		Name:     d.Get("name").(string),
		Isolated: d.Get("isolated").(bool),
	})
	if err != nil {
		return err
	}

	return nil
}

func resourceGeoserverWorkspaceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	d.SetId(d.Id())
	d.Set("name", d.Id())

	log.Printf("[INFO] Importing Geoserver Workspace `%s`", d.Id())

	err := resourceGeoserverWorkspaceRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
