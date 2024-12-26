package geoserver

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	gs "github.com/camptocamp/go-geoserver/client"
)

func resourceGeoserverLayerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGeoserverLayerGroupCreate,
		Read:   resourceGeoserverLayerGroupRead,
		Update: resourceGeoserverLayerGroupUpdate,
		Delete: resourceGeoserverLayerGroupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGeoserverLayerGroupImport,
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
			"mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "SINGLE",
			},
			"title": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"abstract": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bounding_box_min_x": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"bounding_box_max_x": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"bounding_box_min_y": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"bounding_box_max_y": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"bounding_box_crs_class": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bounding_box_crs_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"metadatalink": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"metadatatype": {
							Type:     schema.TypeString,
							Required: true,
						},
						"content": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"layers": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"style": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "layer",
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(string)
								if v != "layer" && v != "layerGroup" {
									errs = append(errs, fmt.Errorf("%q must be 'layer' or 'layerGroup', got: %q", key, v))
								}
								return
							},
						},
					},
				},
			},
			"keywords": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceGeoserverLayerGroupCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Creating Geoserver LayerGroup: %s", d.Id())

	client := meta.(*Config).GeoserverClient()

	workspaceName := d.Get("workspace_name").(string)

	var metadatas []*gs.MetadataLink
	for _, value := range d.Get("metadatalink").(*schema.Set).List() {
		v := value.(map[string]interface{})
		metadatas = append(metadatas, &gs.MetadataLink{
			Type:         v["type"].(string),
			MetadataType: v["metadatatype"].(string),
			Content:      v["content"].(string),
		})
	}

	var keywords []string
	for _, value := range d.Get("keywords").([]interface{}) {
		keywords = append(keywords, value.(string))
	}

	var layers []*gs.LayerRef
	var styles []*gs.StyleRef
	for _, value := range d.Get("layers").([]interface{}) {
		v := value.(map[string]interface{})
		layers = append(layers, &gs.LayerRef{
			Type: v["type"].(string),
			Name: v["name"].(string),
		})
		styles = append(styles, &gs.StyleRef{
			Name: v["style"].(string),
		})
	}

	layerGroup := &gs.LayerGroup{
		Name:         d.Get("name").(string),
		Mode:         d.Get("mode").(string),
		Title:        d.Get("title").(string),
		Abstract:     d.Get("abstract").(string),
		Publishables: layers,
		Styles:       styles,
		Bounds: &gs.BoundingBox{
			MinX: d.Get("bounding_box_min_x").(float64),
			MaxX: d.Get("bounding_box_max_x").(float64),
			MinY: d.Get("bounding_box_min_y").(float64),
			MaxY: d.Get("bounding_box_max_y").(float64),
			CRS: gs.FeatureTypeCRS{
				Class: d.Get("bounding_box_crs_class").(string),
				Value: d.Get("bounding_box_crs_value").(string),
			},
		},
		MetadataLinks: metadatas,
		Keywords:      gs.LayerGroupKeywords{Keywords: keywords},
	}

	err := client.CreateGroup(workspaceName, layerGroup)
	if err != nil {
		client.DeleteGroup(workspaceName, d.Get("name").(string))
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", workspaceName, d.Get("name").(string)))
	return resourceGeoserverLayerGroupRead(d, meta)
}

func resourceGeoserverLayerGroupRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Geoserver Layer Group: %s", d.Id())

	splittedID := strings.Split(d.Id(), "/")
	workspaceName := splittedID[0]
	groupName := splittedID[1]

	client := meta.(*Config).GeoserverClient()

	layerGroup, err := client.GetGroup(workspaceName, groupName)
	if err != nil && !strings.Contains(strings.ToLower(err.Error()), "not found") {
		return err
	}

	if layerGroup == nil {
		d.SetId("")
		return nil
	}

	d.Set("workspace_name", workspaceName)
	d.Set("name", layerGroup.Name)
	d.Set("mode", layerGroup.Mode)
	d.Set("title", layerGroup.Title)
	d.Set("abstract", layerGroup.Abstract)

	d.Set("bounding_box_min_x", layerGroup.Bounds.MinX)
	d.Set("bounding_box_max_x", layerGroup.Bounds.MaxX)
	d.Set("bounding_box_min_y", layerGroup.Bounds.MinY)
	d.Set("bounding_box_max_y", layerGroup.Bounds.MaxY)
	d.Set("bounding_box_crs_class", layerGroup.Bounds.CRS.Class)
	d.Set("bounding_box_crs_value", layerGroup.Bounds.CRS.Value)

	var metadataLinks []map[string]interface{}
	for _, value := range layerGroup.MetadataLinks {
		metadataLinks = append(metadataLinks, map[string]interface{}{
			"type":         value.Type,
			"metadatatype": value.MetadataType,
			"content":      value.Content,
		})
	}
	d.Set("metadatalink", metadataLinks)

	var keywords []map[string]interface{}
	for _, value := range layerGroup.Keywords.Keywords {
		keywords = append(keywords, map[string]interface{}{
			"string": value,
		})
	}
	d.Set("keywords", keywords)

	var layers []map[string]interface{}
	for index, value := range layerGroup.Publishables {
		layers = append(layers, map[string]interface{}{
			"name":  value.Name,
			"style": layerGroup.Styles[index].Name,
		})
	}
	d.Set("layers", layers)

	return nil
}

func resourceGeoserverLayerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Geoserver LayerGroup: %s", d.Id())

	splittedID := strings.Split(d.Id(), "/")
	workspaceName := splittedID[0]
	layerGroupName := splittedID[1]

	client := meta.(*Config).GeoserverClient()

	err := client.DeleteGroup(workspaceName, layerGroupName)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceGeoserverLayerGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Geoserver LayerGroup: %s", d.Id())

	splittedID := strings.Split(d.Id(), "/")
	workspaceName := splittedID[0]

	client := meta.(*Config).GeoserverClient()

	var metadatas []*gs.MetadataLink
	for _, value := range d.Get("metadatalink").(*schema.Set).List() {
		v := value.(map[string]interface{})
		metadatas = append(metadatas, &gs.MetadataLink{
			Type:         v["type"].(string),
			MetadataType: v["metadatatype"].(string),
			Content:      v["content"].(string),
		})
	}

	var keywords []string
	for _, value := range d.Get("keywords").([]interface{}) {
		keywords = append(keywords, value.(string))
	}

	var layers []*gs.LayerRef
	var styles []*gs.StyleRef
	for _, value := range d.Get("layers").([]interface{}) {
		v := value.(map[string]interface{})
		layers = append(layers, &gs.LayerRef{
			Type: "layer",
			Name: v["name"].(string),
		})
		styles = append(styles, &gs.StyleRef{
			Name: v["style"].(string),
		})
	}

	layerGroup := &gs.LayerGroup{
		Name:         d.Get("name").(string),
		Mode:         d.Get("mode").(string),
		Title:        d.Get("title").(string),
		Abstract:     d.Get("abstract").(string),
		Publishables: layers,
		Styles:       styles,
		Bounds: &gs.BoundingBox{
			MinX: d.Get("bounding_box_min_x").(float64),
			MaxX: d.Get("bounding_box_max_x").(float64),
			MinY: d.Get("bounding_box_min_y").(float64),
			MaxY: d.Get("bounding_box_max_y").(float64),
			CRS: gs.FeatureTypeCRS{
				Class: d.Get("bounding_box_crs_class").(string),
				Value: d.Get("bounding_box_crs_value").(string),
			},
		},
		MetadataLinks: metadatas,
		Keywords:      gs.LayerGroupKeywords{Keywords: keywords},
	}

	errUpdateGroup := client.UpdateGroup(workspaceName, layerGroup)
	if errUpdateGroup != nil {
		return errUpdateGroup
	}

	return nil
}

func resourceGeoserverLayerGroupImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splittedID := strings.Split(d.Id(), "/")
	workspaceName := splittedID[0]
	layerGroupName := splittedID[1]

	d.SetId(d.Id())
	d.Set("workspace_name", workspaceName)
	d.Set("name", layerGroupName)

	log.Printf("[INFO] Importing Geoserver LayerGroup `%s` in workspace `%s`", layerGroupName, workspaceName)

	err := resourceGeoserverLayerGroupRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
