package geoserver

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceGeoserverResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceGeoserverResourceCreate,
		Read:   resourceGeoserverResourceRead,
		Update: resourceGeoserverResourceUpdate,
		Delete: resourceGeoserverResourceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGeoserverResourceImport,
		},

		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"extension": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceGeoserverResourceCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Creating Geoserver Resource: %s", d.Id())

	client := meta.(*Config).GeoserverClient()

	err := client.CreateResource(d.Get("path").(string), d.Get("extension").(string), d.Get("resource").(string))
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s#%s", d.Get("path").(string), d.Get("extension").(string)))

	return resourceGeoserverResourceRead(d, meta)
}

func resourceGeoserverResourceRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Geoserver Resource: %s", d.Id())

	splittedID := strings.Split(d.Id(), "#")
	resourcePath := splittedID[0]
	resourceExt := splittedID[1]

	client := meta.(*Config).GeoserverClient()

	resourceContent, err := client.GetResource(resourcePath, resourceExt)
	if err != nil && !strings.Contains(strings.ToLower(err.Error()), "not found") {
		return err
	}

	d.Set("path", resourcePath)
	d.Set("extension", resourceExt)
	d.Set("resource", resourceContent)

	return nil
}

func resourceGeoserverResourceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Geoserver Resource: %s", d.Id())

	splittedID := strings.Split(d.Id(), "#")
	resourcePath := splittedID[0]
	resourceExt := splittedID[1]

	client := meta.(*Config).GeoserverClient()

	err := client.DeleteResource(fmt.Sprintf("%s.%s", resourcePath, resourceExt))
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceGeoserverResourceUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Geoserver Resource: %s", d.Id())

	client := meta.(*Config).GeoserverClient()
	err := client.UpdateResource(d.Get("path").(string), d.Get("extension").(string), d.Get("resource").(string))
	if err != nil {
		return err
	}

	return nil
}

func resourceGeoserverResourceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splittedID := strings.Split(d.Id(), "#")
	resourcePath := splittedID[0]
	resourceExt := splittedID[1]

	d.SetId(d.Id())
	d.Set("path", resourcePath)
	d.Set("extension", resourceExt)

	log.Printf("[INFO] Importing Geoserver Resource `%s`", resourcePath)

	err := resourceGeoserverResourceRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
