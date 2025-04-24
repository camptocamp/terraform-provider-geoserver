package geoserver

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	gs "github.com/camptocamp/go-geoserver/client"
)

func resourceGeoserverWmtsLayer() *schema.Resource {
	return &schema.Resource{
		Create: resourceGeoserverWmtsLayerCreate,
		Read:   resourceGeoserverWmtsLayerRead,
		//		Update: resourceGeoserverWmsLayerUpdate,
		Delete: resourceGeoserverWmtsLayerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGeoserverWmtsLayerImport,
		},

		Schema: map[string]*schema.Schema{
			"workspace_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"wmtsstore_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"native_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},
			"projection_policy": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"title": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"abstract": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"native_crs_class": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"native_crs_value": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"srs": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"native_bounding_box_min_x": {
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: true,
			},
			"native_bounding_box_max_x": {
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: true,
			},
			"native_bounding_box_min_y": {
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: true,
			},
			"native_bounding_box_max_y": {
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: true,
			},
			"native_bounding_box_crs_class": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"native_bounding_box_crs_value": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"lat_lon_bounding_box_min_x": {
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: true,
			},
			"lat_lon_bounding_box_max_x": {
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: true,
			},
			"lat_lon_bounding_box_min_y": {
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: true,
			},
			"lat_lon_bounding_box_max_y": {
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: true,
			},
			"lat_lon_bounding_box_crs_class": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"lat_lon_bounding_box_crs_value": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
			},
		},
	}
}

func resourceGeoserverWmtsLayerCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Creating Geoserver WMTS layer: %s", d.Id())

	client := meta.(*Config).GeoserverClient()

	workspaceName := d.Get("workspace_name").(string)
	datastoreName := d.Get("wmsstore_name").(string)

	var metadata []*gs.WmtsLayerMetadata
	for key, value := range d.Get("metadata").(map[string]interface{}) {
		metadata = append(metadata, &gs.WmtsLayerMetadata{
			Key:   key,
			Value: value.(string),
		})
	}

	WmtsLayer := &gs.WmtsLayer{
		Name:             d.Get("name").(string),
		NativeName:       d.Get("native_name").(string),
		Enabled:          d.Get("enabled").(bool),
		ProjectionPolicy: d.Get("projection_policy").(string),
		Title:            d.Get("title").(string),
		Abstract:         d.Get("abstract").(string),
		NativeCRS: gs.WmtsLayerCRS{
			Class: d.Get("native_crs_class").(string),
			Value: d.Get("native_crs_value").(string),
		},
		SRS: d.Get("srs").(string),
		NativeBoundingBox: gs.WmtsLayerBoundingBox{
			MinX: d.Get("native_bounding_box_min_x").(float64),
			MaxX: d.Get("native_bounding_box_max_x").(float64),
			MinY: d.Get("native_bounding_box_min_y").(float64),
			MaxY: d.Get("native_bounding_box_max_y").(float64),
			CRS: gs.FeatureTypeCRS{
				Class: d.Get("native_bounding_box_crs_class").(string),
				Value: d.Get("native_bounding_box_crs_value").(string),
			},
		},
		LatLonBoundingBox: gs.WmtsLayerBoundingBox{
			MinX: d.Get("lat_lon_bounding_box_min_x").(float64),
			MaxX: d.Get("lat_lon_bounding_box_max_x").(float64),
			MinY: d.Get("lat_lon_bounding_box_min_y").(float64),
			MaxY: d.Get("lat_lon_bounding_box_max_y").(float64),
			CRS: gs.FeatureTypeCRS{
				Class: d.Get("lat_lon_bounding_box_crs_class").(string),
				Value: d.Get("lat_lon_bounding_box_crs_value").(string),
			},
		},
		Metadata: metadata,
	}

	err := client.CreateWmtsLayer(workspaceName, datastoreName, WmtsLayer)
	if err != nil {
		client.DeleteWmtsLayer(workspaceName, datastoreName, d.Get("name").(string), true)
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", workspaceName, datastoreName, d.Get("name").(string)))

	return resourceGeoserverWmtsLayerRead(d, meta)
}

func resourceGeoserverWmtsLayerRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Geoserver WmtsLayer: %s", d.Id())

	splittedID := strings.Split(d.Id(), "/")
	workspaceName := splittedID[0]
	datastoreName := splittedID[1]
	WmtsLayerName := splittedID[2]

	client := meta.(*Config).GeoserverClient()

	WmtsLayer, err := client.GetWmtsLayer(workspaceName, datastoreName, WmtsLayerName)
	if err != nil && !strings.Contains(strings.ToLower(err.Error()), "not found") {
		return err
	}

	if WmtsLayer == nil {
		d.SetId("")
		return nil
	}
	d.Set("workspace_name", workspaceName)
	d.Set("wmsstore_name", datastoreName)
	d.Set("name", WmtsLayer.Name)
	d.Set("native_name", WmtsLayer.NativeName)
	d.Set("enabled", WmtsLayer.Enabled)
	d.Set("projection_policy", WmtsLayer.ProjectionPolicy)
	d.Set("srs", WmtsLayer.SRS)

	d.Set("native_crs_class", WmtsLayer.NativeCRS.Class)
	d.Set("native_crs_value", WmtsLayer.NativeCRS.Value)

	d.Set("native_bounding_box_min_x", WmtsLayer.NativeBoundingBox.MinX)
	d.Set("native_bounding_box_max_x", WmtsLayer.NativeBoundingBox.MaxX)
	d.Set("native_bounding_box_min_y", WmtsLayer.NativeBoundingBox.MinY)
	d.Set("native_bounding_box_max_y", WmtsLayer.NativeBoundingBox.MaxY)
	d.Set("native_bounding_box_crs_class", WmtsLayer.NativeBoundingBox.CRS.Class)
	d.Set("native_bounding_box_crs_value", WmtsLayer.NativeBoundingBox.CRS.Value)

	d.Set("lat_lon_bounding_box_min_x", WmtsLayer.LatLonBoundingBox.MinX)
	d.Set("lat_lon_bounding_box_max_x", WmtsLayer.LatLonBoundingBox.MaxX)
	d.Set("lat_lon_bounding_box_min_y", WmtsLayer.LatLonBoundingBox.MinY)
	d.Set("lat_lon_bounding_box_max_y", WmtsLayer.LatLonBoundingBox.MaxY)
	d.Set("lat_lon_bounding_box_crs_class", WmtsLayer.LatLonBoundingBox.CRS.Class)
	d.Set("lat_lon_bounding_box_crs_value", WmtsLayer.LatLonBoundingBox.CRS.Value)

	metadata := map[string]interface{}{}
	if WmtsLayer.Metadata != nil {
		for _, entry := range WmtsLayer.Metadata {
			metadata[entry.Key] = entry.Value
		}
		d.Set("metadata", metadata)
	}

	return nil
}

func resourceGeoserverWmtsLayerDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Geoserver WMTS Layer: %s", d.Id())

	splittedID := strings.Split(d.Id(), "/")
	workspaceName := splittedID[0]
	datastoreName := splittedID[1]
	WmtsLayerName := splittedID[2]

	client := meta.(*Config).GeoserverClient()

	err := client.DeleteWmtsLayer(workspaceName, datastoreName, WmtsLayerName, true)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceGeoserverWmtsLayerUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Geoserver WmtsLayer: %s", d.Id())

	splittedID := strings.Split(d.Id(), "/")
	workspaceName := splittedID[0]
	datastoreName := splittedID[1]
	WmtsLayerName := splittedID[2]

	client := meta.(*Config).GeoserverClient()

	var metadata []*gs.WmtsLayerMetadata
	for key, value := range d.Get("metadata").(map[string]interface{}) {
		metadata = append(metadata, &gs.WmtsLayerMetadata{
			Key:   key,
			Value: value.(string),
		})
	}

	WmtsLayer := &gs.WmtsLayer{
		Name:             d.Get("name").(string),
		NativeName:       d.Get("native_name").(string),
		Enabled:          d.Get("enabled").(bool),
		ProjectionPolicy: d.Get("projection_policy").(string),
		Title:            d.Get("title").(string),
		Abstract:         d.Get("abstract").(string),
		NativeCRS: gs.WmtsLayerCRS{
			Class: d.Get("native_crs_class").(string),
			Value: d.Get("native_crs_value").(string),
		},
		SRS: d.Get("srs").(string),
		NativeBoundingBox: gs.WmtsLayerBoundingBox{
			MinX: d.Get("native_bounding_box_min_x").(float64),
			MaxX: d.Get("native_bounding_box_max_x").(float64),
			MinY: d.Get("native_bounding_box_min_y").(float64),
			MaxY: d.Get("native_bounding_box_max_y").(float64),
			CRS: gs.FeatureTypeCRS{
				Class: d.Get("native_bounding_box_crs_class").(string),
				Value: d.Get("native_bounding_box_crs_value").(string),
			},
		},
		LatLonBoundingBox: gs.WmtsLayerBoundingBox{
			MinX: d.Get("lat_lon_bounding_box_min_x").(float64),
			MaxX: d.Get("lat_lon_bounding_box_max_x").(float64),
			MinY: d.Get("lat_lon_bounding_box_min_y").(float64),
			MaxY: d.Get("lat_lon_bounding_box_max_y").(float64),
			CRS: gs.FeatureTypeCRS{
				Class: d.Get("lat_lon_bounding_box_crs_class").(string),
				Value: d.Get("lat_lon_bounding_box_crs_value").(string),
			},
		},
		Metadata: metadata,
	}

	err := client.UpdateWmtsLayer(workspaceName, datastoreName, WmtsLayerName, WmtsLayer)
	if err != nil {
		return err
	}

	return nil
}

func resourceGeoserverWmtsLayerImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splittedID := strings.Split(d.Id(), "/")
	workspaceName := splittedID[0]
	datastoreName := splittedID[1]
	WmtsLayerName := splittedID[2]

	d.SetId(d.Id())
	d.Set("workspace_name", workspaceName)
	d.Set("wmtsstore_name", datastoreName)
	d.Set("name", WmtsLayerName)

	log.Printf("[INFO] Importing Geoserver WmtsLayer `%s` from workspace `%s`", datastoreName, workspaceName)

	err := resourceGeoserverWmtsLayerRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
