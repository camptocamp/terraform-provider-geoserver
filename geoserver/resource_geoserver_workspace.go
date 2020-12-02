package geoserver

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceGeoserverWorkspace() *schema.Resource {
	return &schema.Resource{
		Create: resourceGeoserverWorkspaceCreate,
		Read:   resourceGeoserverWorkspaceRead,
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
			//"default": {
			//	Type:     schema.TypeBool,
			//	Optional: true,
			//	Default:  false,
			//},
		},
	}
}

func resourceGeoserverWorkspaceCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Creating Geoserver Workspace: %s", d.Id())

	client := meta.(*Config).Client()

	name := d.Get("name").(string)

	_, err := client.CreateWorkspace(name)
	if err != nil {
		return err
	}

	d.SetId(name)

	return resourceGeoserverWorkspaceRead(d, meta)
}

func resourceGeoserverWorkspaceRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Geoserver Workspace: %s", d.Id())

	client := meta.(*Config).Client()

	workspace, err := client.GetWorkspace(d.Id())
	if err != nil {
		return err
	}

	d.Set("name", workspace.Name)

	return nil
}

func resourceGeoserverWorkspaceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Geoserver Workspace: %s", d.Id())

	client := meta.(*Config).Client()

	_, err := client.DeleteWorkspace(d.Id(), true)
	if err != nil {
		return err
	}

	d.SetId("")

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
