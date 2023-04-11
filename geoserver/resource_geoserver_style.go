package geoserver

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	gs "github.com/camptocamp/go-geoserver/client"
)

func resourceGeoserverStyle() *schema.Resource {
	return &schema.Resource{
		Create: resourceGeoserverStyleCreate,
		Read:   resourceGeoserverStyleRead,
		Update: resourceGeoserverStyleUpdate,
		Delete: resourceGeoserverStyleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGeoserverStyleImport,
		},

		Schema: map[string]*schema.Schema{
			"workspace_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"filename": {
				Type:     schema.TypeString,
				Required: true,
			},
			"format": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"style_definition": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceGeoserverStyleCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Creating Geoserver Style: %s", d.Id())

	client := meta.(*Config).GeoserverClient()

	workspaceName := d.Get("workspace_name").(string)

	style := &gs.Style{
		Name:     d.Get("name").(string),
		FileName: d.Get("filename").(string),
		Format:   d.Get("format").(string),
		Version:  &gs.LanguageVersion{Version: d.Get("version").(string)},
	}

	err := client.CreateStyle(workspaceName, style)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", workspaceName, d.Get("name").(string)))

	errUpdateStyle := client.UpdateStyleContent(workspaceName, style, d.Get("style_definition").(string))
	if errUpdateStyle != nil {
		return errUpdateStyle
	}

	return resourceGeoserverStyleRead(d, meta)
}

func resourceGeoserverStyleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Geoserver Style: %s", d.Id())

	splittedID := strings.Split(d.Id(), "/")
	workspaceName := splittedID[0]
	styleName := splittedID[1]

	client := meta.(*Config).GeoserverClient()

	style, err := client.GetStyle(workspaceName, styleName)
	if err != nil && !strings.Contains(strings.ToLower(err.Error()), "not found") {
		return err
	}

	if style == nil {
		d.SetId("")
		return nil
	}
	d.Set("workspace_name", workspaceName)
	d.Set("name", style.Name)
	d.Set("filename", style.FileName)
	d.Set("format", style.Format)
	d.Set("version", style.Version.Version)

	styleContent, errContent := client.GetStyleFile(workspaceName, styleName, style.Format, style.Version.Version)
	if errContent != nil {
		return errContent
	}

	d.Set("style_definition", styleContent)

	return nil
}

func resourceGeoserverStyleDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Geoserver Style: %s", d.Id())

	splittedID := strings.Split(d.Id(), "/")
	workspaceName := splittedID[0]
	styleName := splittedID[1]

	client := meta.(*Config).GeoserverClient()

	err := client.DeleteStyle(workspaceName, styleName, true, true)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceGeoserverStyleUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Geoserver Style: %s", d.Id())

	splittedID := strings.Split(d.Id(), "/")
	workspaceName := splittedID[0]

	client := meta.(*Config).GeoserverClient()

	style := &gs.Style{
		Name:     d.Get("name").(string),
		FileName: d.Get("filename").(string),
		Format:   d.Get("format").(string),
		Version:  &gs.LanguageVersion{Version: d.Get("version").(string)},
	}

	errUpdateStyle := client.UpdateStyleContent(workspaceName, style, d.Get("style_definition").(string))
	if errUpdateStyle != nil {
		return errUpdateStyle
	}

	return nil
}

func resourceGeoserverStyleImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splittedID := strings.Split(d.Id(), "/")
	workspaceName := splittedID[0]
	styleName := splittedID[1]

	d.SetId(d.Id())
	d.Set("workspace_name", workspaceName)
	d.Set("name", styleName)

	log.Printf("[INFO] Importing Geoserver Style `%s` in workspace `%s`", styleName, workspaceName)

	err := resourceGeoserverStyleRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
