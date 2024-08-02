package geoserver

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	gs "github.com/camptocamp/go-geoserver/client"
)

func resourceGeoserverUrlCheck() *schema.Resource {
	return &schema.Resource{
		Create: resourceGeoserverUrlCheckCreate,
		Read:   resourceGeoserverUrlCheckRead,
		Update: resourceGeoserverUrlCheckUpdate,
		Delete: resourceGeoserverUrlCheckDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGeoserverUrlCheckImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the check. Use as resource id.",
			},
			"regex": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Regular expression to evaluate the URL check.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				ForceNew:    true,
				Description: "Declare the check as enabled. Default value: true.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the check.",
			},
		},
	}
}

func resourceGeoserverUrlCheckCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Creating Geoserver UrlCheck: %s", d.Id())

	client := meta.(*Config).GeoserverClient()

	name := d.Get("name").(string)

	err := client.CreateRegExUrlCheck(name, &gs.RegexUrlCheck{
		Name:        name,
		IsEnabled:   d.Get("enabled").(bool),
		Description: d.Get("description").(string),
		Regex:       d.Get("regex").(string),
	})
	if err != nil {
		return err
	}

	d.SetId(name)

	return resourceGeoserverUrlCheckRead(d, meta)
}

func resourceGeoserverUrlCheckRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Geoserver UrlCheck: %s", d.Id())

	client := meta.(*Config).GeoserverClient()

	urlcheck, err := client.GetRegExUrlCheck(d.Id())
	if err != nil && !strings.Contains(strings.ToLower(err.Error()), "not found") {
		return err
	}

	if urlcheck == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", urlcheck.Name)
	d.Set("enabled", urlcheck.IsEnabled)
	d.Set("description", urlcheck.Description)
	d.Set("regex", urlcheck.Regex)

	return nil
}

func resourceGeoserverUrlCheckDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Geoserver UrlCheck: %s", d.Id())

	client := meta.(*Config).GeoserverClient()

	err := client.DeleteUrlCheck(d.Id())
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceGeoserverUrlCheckUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Geoserver UrlCheck: %s", d.Id())

	client := meta.(*Config).GeoserverClient()

	err := client.UpdateRegExUrlCheck(d.Id(), &gs.RegexUrlCheck{
		Name:        d.Get("name").(string),
		IsEnabled:   d.Get("enabled").(bool),
		Description: d.Get("description").(string),
		Regex:       d.Get("regex").(string),
	})
	if err != nil {
		return err
	}

	return nil
}

func resourceGeoserverUrlCheckImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	d.SetId(d.Id())
	d.Set("name", d.Id())

	log.Printf("[INFO] Importing Geoserver UrlCheck `%s`", d.Id())

	err := resourceGeoserverUrlCheckRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
