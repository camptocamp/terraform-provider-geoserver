package geoserver

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	gs "github.com/camptocamp/go-geoserver/client"
)

func resourceGeoserverFeatureType() *schema.Resource {
	return &schema.Resource{
		Create: resourceGeoserverFeatureTypeCreate,
		Read:   resourceGeoserverFeatureTypeRead,
		Update: resourceGeoserverFeatureTypeUpdate,
		Delete: resourceGeoserverFeatureTypeDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGeoserverFeatureTypeImport,
		},

		Schema: map[string]*schema.Schema{
			"workspace_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"datastore_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"native_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"projection_policy": {
				Type:     schema.TypeString,
				Required: true,
			},
			"title": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"abstract": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"native_crs_class": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"native_crs_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"srs": {
				Type:     schema.TypeString,
				Required: true,
			},
			"native_bounding_box_min_x": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"native_bounding_box_max_x": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"native_bounding_box_min_y": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"native_bounding_box_max_y": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"native_bounding_box_crs_class": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"native_bounding_box_crs_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lat_lon_bounding_box_min_x": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"lat_lon_bounding_box_max_x": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"lat_lon_bounding_box_min_y": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"lat_lon_bounding_box_max_y": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"lat_lon_bounding_box_crs_class": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lat_lon_bounding_box_crs_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attribute": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"min_occurs": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"max_occurs": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"nillable": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"binding": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceGeoserverFeatureTypeCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Creating Geoserver FeatureType: %s", d.Id())

	client := meta.(*Config).Client()

	workspaceName := d.Get("workspace_name").(string)
	datastoreName := d.Get("datastore_name").(string)

	var attributes []*gs.FeatureTypeAttribute
	for _, value := range d.Get("attribute").(*schema.Set).List() {
		v := value.(map[string]interface{})
		attributes = append(attributes, &gs.FeatureTypeAttribute{
			Name:      v["name"].(string),
			Nillable:  v["nillable"].(bool),
			Binding:   v["binding"].(string),
			MinOccurs: v["min_occurs"].(int),
			MaxOccurs: v["max_occurs"].(int),
		})
	}

	featureType := &gs.FeatureType{
		Name:             d.Get("name").(string),
		NativeName:       d.Get("native_name").(string),
		Enabled:          d.Get("enabled").(bool),
		ProjectionPolicy: d.Get("projection_policy").(string),
		Title:            d.Get("title").(string),
		Abstract:         d.Get("abstract").(string),
		NativeCRS: gs.CRSWrapper{
			Class: d.Get("native_crs_class").(string),
			Value: d.Get("native_crs_value").(string),
		},
		SRS: d.Get("srs").(string),
		NativeBoundingBox: gs.BoundingBox{
			MinX: d.Get("native_bounding_box_min_x").(float64),
			MaxX: d.Get("native_bounding_box_max_x").(float64),
			MinY: d.Get("native_bounding_box_min_y").(float64),
			MaxY: d.Get("native_bounding_box_max_y").(float64),
			CRS: gs.CRSWrapper{
				Class: d.Get("native_bounding_box_crs_class").(string),
				Value: d.Get("native_bounding_box_crs_value").(string),
			},
		},
		LatLonBoundingBox: gs.BoundingBox{
			MinX: d.Get("lat_lon_bounding_box_min_x").(float64),
			MaxX: d.Get("lat_lon_bounding_box_max_x").(float64),
			MinY: d.Get("lat_lon_bounding_box_min_y").(float64),
			MaxY: d.Get("lat_lon_bounding_box_max_y").(float64),
			CRS: gs.CRSWrapper{
				Class: d.Get("lat_lon_bounding_box_crs_class").(string),
				Value: d.Get("lat_lon_bounding_box_crs_value").(string),
			},
		},
		Attributes: gs.FeatureTypeAttributesList{
			Attribute: attributes,
		},
	}

	err := client.CreateFeatureType(workspaceName, datastoreName, featureType)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", workspaceName, datastoreName, d.Get("name").(string)))

	return resourceGeoserverFeatureTypeRead(d, meta)
}

func resourceGeoserverFeatureTypeRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Geoserver FeatureType: %s", d.Id())

	splittedID := strings.Split(d.Id(), "/")
	workspaceName := splittedID[0]
	datastoreName := splittedID[1]
	featureTypeName := splittedID[2]

	client := meta.(*Config).Client()

	featureType, err := client.GetFeatureType(workspaceName, datastoreName, featureTypeName)
	if err != nil {
		return err
	}

	d.Set("workspace_name", workspaceName)
	d.Set("datastore_name", datastoreName)
	d.Set("name", featureType.Name)
	d.Set("native_name", featureType.NativeName)
	d.Set("enabled", featureType.Enabled)
	d.Set("projection_policy", featureType.ProjectionPolicy)
	d.Set("srs", featureType.SRS)

	d.Set("native_bounding_box_min_x", featureType.NativeBoundingBox.MinX)
	d.Set("native_bounding_box_max_x", featureType.NativeBoundingBox.MaxX)
	d.Set("native_bounding_box_min_y", featureType.NativeBoundingBox.MinY)
	d.Set("native_bounding_box_max_y", featureType.NativeBoundingBox.MaxY)
	d.Set("native_bounding_box_crs_class", featureType.NativeBoundingBox.CRS.Class)
	d.Set("native_bounding_box_crs_value", featureType.NativeBoundingBox.CRS.Value)

	d.Set("lat_lon_bounding_box_min_x", featureType.LatLonBoundingBox.MinX)
	d.Set("lat_lon_bounding_box_max_x", featureType.LatLonBoundingBox.MaxX)
	d.Set("lat_lon_bounding_box_min_y", featureType.LatLonBoundingBox.MinY)
	d.Set("lat_lon_bounding_box_max_y", featureType.LatLonBoundingBox.MaxY)
	d.Set("lat_lon_bounding_box_crs_class", featureType.LatLonBoundingBox.CRS.Class)
	d.Set("lat_lon_bounding_box_crs_value", featureType.LatLonBoundingBox.CRS.Value)

	var attributes []map[string]interface{}
	for _, value := range featureType.Attributes.Attribute {
		attributes = append(attributes, map[string]interface{}{
			"name":       value.Name,
			"nillable":   value.Nillable,
			"binding":    value.Binding,
			"min_occurs": value.MinOccurs,
			"max_occurs": value.MaxOccurs,
		})
	}
	d.Set("attribute", attributes)

	return nil
}

func resourceGeoserverFeatureTypeDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Geoserver Feature Type: %s", d.Id())

	splittedID := strings.Split(d.Id(), "/")
	workspaceName := splittedID[0]
	datastoreName := splittedID[1]
	featureTypeName := splittedID[2]

	client := meta.(*Config).Client()

	err := client.DeleteFeatureType(workspaceName, datastoreName, featureTypeName, true)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceGeoserverFeatureTypeUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Geoserver FeatureType: %s", d.Id())

	splittedID := strings.Split(d.Id(), "/")
	workspaceName := splittedID[0]
	datastoreName := splittedID[1]
	featureTypeName := splittedID[2]

	client := meta.(*Config).Client()

	var attributes []*gs.FeatureTypeAttribute
	for _, value := range d.Get("attribute").(*schema.Set).List() {
		v := value.(map[string]interface{})
		attributes = append(attributes, &gs.FeatureTypeAttribute{
			Name:      v["name"].(string),
			Nillable:  v["nillable"].(bool),
			Binding:   v["binding"].(string),
			MinOccurs: v["min_occurs"].(int),
			MaxOccurs: v["max_occurs"].(int),
		})
	}

	featureType := &gs.FeatureType{
		Name:             d.Get("name").(string),
		NativeName:       d.Get("native_name").(string),
		Enabled:          d.Get("enabled").(bool),
		ProjectionPolicy: d.Get("projection_policy").(string),
		Title:            d.Get("title").(string),
		Abstract:         d.Get("abstract").(string),
		NativeCRS: gs.CRSWrapper{
			Class: d.Get("native_crs_class").(string),
			Value: d.Get("native_crs_value").(string),
		},
		SRS: d.Get("srs").(string),
		NativeBoundingBox: gs.BoundingBox{
			MinX: d.Get("native_bounding_box_min_x").(float64),
			MaxX: d.Get("native_bounding_box_max_x").(float64),
			MinY: d.Get("native_bounding_box_min_y").(float64),
			MaxY: d.Get("native_bounding_box_max_y").(float64),
			CRS: gs.CRSWrapper{
				Class: d.Get("native_bounding_box_crs_class").(string),
				Value: d.Get("native_bounding_box_crs_value").(string),
			},
		},
		LatLonBoundingBox: gs.BoundingBox{
			MinX: d.Get("lat_lon_bounding_box_min_x").(float64),
			MaxX: d.Get("lat_lon_bounding_box_max_x").(float64),
			MinY: d.Get("lat_lon_bounding_box_min_y").(float64),
			MaxY: d.Get("lat_lon_bounding_box_max_y").(float64),
			CRS: gs.CRSWrapper{
				Class: d.Get("lat_lon_bounding_box_crs_class").(string),
				Value: d.Get("lat_lon_bounding_box_crs_value").(string),
			},
		},
		Attributes: gs.FeatureTypeAttributesList{
			Attribute: attributes,
		},
	}

	err := client.UpdateFeatureType(workspaceName, datastoreName, featureTypeName, featureType)
	if err != nil {
		return err
	}

	return nil
}

func resourceGeoserverFeatureTypeImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splittedID := strings.Split(d.Id(), "/")
	workspaceName := splittedID[0]
	datastoreName := splittedID[1]
	featureTypeName := splittedID[2]

	d.SetId(d.Id())
	d.Set("workspace_name", workspaceName)
	d.Set("datastore_name", datastoreName)
	d.Set("name", featureTypeName)

	log.Printf("[INFO] Importing Geoserver FeatureType `%s` in workspace `%s`", datastoreName, workspaceName)

	err := resourceGeoserverFeatureTypeRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
