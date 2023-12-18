package geoserver

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	gs "github.com/camptocamp/go-geoserver/client"
)

func resourceGeoserverWmsLayer() *schema.Resource {
	return &schema.Resource{
		Create: resourceGeoserverWmsLayerCreate,
		Read:   resourceGeoserverWmsLayerRead,
		//		Update: resourceGeoserverWmsLayerUpdate,
		Delete: resourceGeoserverWmsLayerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGeoserverWmsLayerImport,
		},

		Schema: map[string]*schema.Schema{
			"workspace_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"wmsstore_name": {
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

func resourceGeoserverWmsLayerCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Creating Geoserver WMS layer: %s", d.Id())

	client := meta.(*Config).GeoserverClient()

	workspaceName := d.Get("workspace_name").(string)
	datastoreName := d.Get("wmsstore_name").(string)

	var metadata []*gs.WmsLayerMetadata
	for key, value := range d.Get("metadata").(map[string]interface{}) {
		metadata = append(metadata, &gs.WmsLayerMetadata{
			Key:   key,
			Value: value.(string),
		})
	}

	WmsLayer := &gs.WmsLayer{
		Name:             d.Get("name").(string),
		NativeName:       d.Get("native_name").(string),
		Enabled:          d.Get("enabled").(bool),
		ProjectionPolicy: d.Get("projection_policy").(string),
		Title:            d.Get("title").(string),
		Abstract:         d.Get("abstract").(string),
		NativeCRS: gs.WmsLayerCRS{
			Class: d.Get("native_crs_class").(string),
			Value: d.Get("native_crs_value").(string),
		},
		SRS: d.Get("srs").(string),
		NativeBoundingBox: gs.WmsLayerBoundingBox{
			MinX: d.Get("native_bounding_box_min_x").(float64),
			MaxX: d.Get("native_bounding_box_max_x").(float64),
			MinY: d.Get("native_bounding_box_min_y").(float64),
			MaxY: d.Get("native_bounding_box_max_y").(float64),
			CRS: gs.FeatureTypeCRS{
				Class: d.Get("native_bounding_box_crs_class").(string),
				Value: d.Get("native_bounding_box_crs_value").(string),
			},
		},
		LatLonBoundingBox: gs.WmsLayerBoundingBox{
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

	err := client.CreateWmsLayer(workspaceName, datastoreName, WmsLayer)
	if err != nil {
		client.DeleteWmsLayer(workspaceName, datastoreName, d.Get("name").(string), true)
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", workspaceName, datastoreName, d.Get("name").(string)))

	return resourceGeoserverWmsLayerRead(d, meta)
}

func resourceGeoserverWmsLayerRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Geoserver WmsLayer: %s", d.Id())

	splittedID := strings.Split(d.Id(), "/")
	workspaceName := splittedID[0]
	datastoreName := splittedID[1]
	WmsLayerName := splittedID[2]

	client := meta.(*Config).GeoserverClient()

	WmsLayer, err := client.GetWmsLayer(workspaceName, datastoreName, WmsLayerName)
	if err != nil && !strings.Contains(strings.ToLower(err.Error()), "not found") {
		return err
	}

	if WmsLayer == nil {
		d.SetId("")
		return nil
	}
	d.Set("workspace_name", workspaceName)
	d.Set("wmsstore_name", datastoreName)
	d.Set("name", WmsLayer.Name)
	d.Set("native_name", WmsLayer.NativeName)
	d.Set("enabled", WmsLayer.Enabled)
	d.Set("projection_policy", WmsLayer.ProjectionPolicy)
	d.Set("srs", WmsLayer.SRS)

	d.Set("native_crs_class", WmsLayer.NativeCRS.Class)
	d.Set("native_crs_value", WmsLayer.NativeCRS.Value)

	d.Set("native_bounding_box_min_x", WmsLayer.NativeBoundingBox.MinX)
	d.Set("native_bounding_box_max_x", WmsLayer.NativeBoundingBox.MaxX)
	d.Set("native_bounding_box_min_y", WmsLayer.NativeBoundingBox.MinY)
	d.Set("native_bounding_box_max_y", WmsLayer.NativeBoundingBox.MaxY)
	d.Set("native_bounding_box_crs_class", WmsLayer.NativeBoundingBox.CRS.Class)
	d.Set("native_bounding_box_crs_value", WmsLayer.NativeBoundingBox.CRS.Value)

	d.Set("lat_lon_bounding_box_min_x", WmsLayer.LatLonBoundingBox.MinX)
	d.Set("lat_lon_bounding_box_max_x", WmsLayer.LatLonBoundingBox.MaxX)
	d.Set("lat_lon_bounding_box_min_y", WmsLayer.LatLonBoundingBox.MinY)
	d.Set("lat_lon_bounding_box_max_y", WmsLayer.LatLonBoundingBox.MaxY)
	d.Set("lat_lon_bounding_box_crs_class", WmsLayer.LatLonBoundingBox.CRS.Class)
	d.Set("lat_lon_bounding_box_crs_value", WmsLayer.LatLonBoundingBox.CRS.Value)

	metadata := map[string]interface{}{}
	if WmsLayer.Metadata != nil {
		for _, entry := range WmsLayer.Metadata {
			metadata[entry.Key] = entry.Value
		}
		d.Set("metadata", metadata)
	}

	return nil
}

func resourceGeoserverWmsLayerDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Geoserver Feature Type: %s", d.Id())

	splittedID := strings.Split(d.Id(), "/")
	workspaceName := splittedID[0]
	datastoreName := splittedID[1]
	WmsLayerName := splittedID[2]

	client := meta.(*Config).GeoserverClient()

	err := client.DeleteWmsLayer(workspaceName, datastoreName, WmsLayerName, true)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceGeoserverWmsLayerUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Geoserver WmsLayer: %s", d.Id())

	splittedID := strings.Split(d.Id(), "/")
	workspaceName := splittedID[0]
	datastoreName := splittedID[1]
	WmsLayerName := splittedID[2]

	client := meta.(*Config).GeoserverClient()

	var metadata []*gs.WmsLayerMetadata
	for key, value := range d.Get("metadata").(map[string]interface{}) {
		metadata = append(metadata, &gs.WmsLayerMetadata{
			Key:   key,
			Value: value.(string),
		})
	}

	WmsLayer := &gs.WmsLayer{
		Name:             d.Get("name").(string),
		NativeName:       d.Get("native_name").(string),
		Enabled:          d.Get("enabled").(bool),
		ProjectionPolicy: d.Get("projection_policy").(string),
		Title:            d.Get("title").(string),
		Abstract:         d.Get("abstract").(string),
		NativeCRS: gs.WmsLayerCRS{
			Class: d.Get("native_crs_class").(string),
			Value: d.Get("native_crs_value").(string),
		},
		SRS: d.Get("srs").(string),
		NativeBoundingBox: gs.WmsLayerBoundingBox{
			MinX: d.Get("native_bounding_box_min_x").(float64),
			MaxX: d.Get("native_bounding_box_max_x").(float64),
			MinY: d.Get("native_bounding_box_min_y").(float64),
			MaxY: d.Get("native_bounding_box_max_y").(float64),
			CRS: gs.FeatureTypeCRS{
				Class: d.Get("native_bounding_box_crs_class").(string),
				Value: d.Get("native_bounding_box_crs_value").(string),
			},
		},
		LatLonBoundingBox: gs.WmsLayerBoundingBox{
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

	err := client.UpdateWmsLayer(workspaceName, datastoreName, WmsLayerName, WmsLayer)
	if err != nil {
		return err
	}

	return nil
}

func resourceGeoserverWmsLayerImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splittedID := strings.Split(d.Id(), "/")
	workspaceName := splittedID[0]
	datastoreName := splittedID[1]
	WmsLayerName := splittedID[2]

	d.SetId(d.Id())
	d.Set("workspace_name", workspaceName)
	d.Set("wmsstore_name", datastoreName)
	d.Set("name", WmsLayerName)

	log.Printf("[INFO] Importing Geoserver WmsLayer `%s` from workspace `%s`", datastoreName, workspaceName)

	err := resourceGeoserverWmsLayerRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
