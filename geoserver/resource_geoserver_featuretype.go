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
			"use_custom_attributes": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"attribute": {
				Type:     schema.TypeSet,
				Optional: true,
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
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceGeoserverFeatureTypeCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Creating Geoserver FeatureType: %s", d.Get("name").(string))

	client := meta.(*Config).GeoserverClient()

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

	var metadata []*gs.FeatureTypeMetadata
	for key, value := range d.Get("metadata").(map[string]interface{}) {
		metadata = append(metadata, &gs.FeatureTypeMetadata{
			Key:   key,
			Value: value.(string),
		})
	}

	featureType := &gs.FeatureType{
		Name:             d.Get("name").(string),
		NativeName:       d.Get("native_name").(string),
		Enabled:          d.Get("enabled").(bool),
		ProjectionPolicy: d.Get("projection_policy").(string),
		Title:            d.Get("title").(string),
		Abstract:         d.Get("abstract").(string),
		NativeCRS: gs.FeatureTypeCRS{
			Class: d.Get("native_crs_class").(string),
			Value: d.Get("native_crs_value").(string),
		},
		SRS: d.Get("srs").(string),
		NativeBoundingBox: gs.BoundingBox{
			MinX: d.Get("native_bounding_box_min_x").(float64),
			MaxX: d.Get("native_bounding_box_max_x").(float64),
			MinY: d.Get("native_bounding_box_min_y").(float64),
			MaxY: d.Get("native_bounding_box_max_y").(float64),
			CRS: gs.FeatureTypeCRS{
				Class: d.Get("native_bounding_box_crs_class").(string),
				Value: d.Get("native_bounding_box_crs_value").(string),
			},
		},
		LatLonBoundingBox: gs.BoundingBox{
			MinX: d.Get("lat_lon_bounding_box_min_x").(float64),
			MaxX: d.Get("lat_lon_bounding_box_max_x").(float64),
			MinY: d.Get("lat_lon_bounding_box_min_y").(float64),
			MaxY: d.Get("lat_lon_bounding_box_max_y").(float64),
			CRS: gs.FeatureTypeCRS{
				Class: d.Get("lat_lon_bounding_box_crs_class").(string),
				Value: d.Get("lat_lon_bounding_box_crs_value").(string),
			},
		},
		Attributes: attributes,
		Metadata:   metadata,
	}

	err := client.CreateFeatureType(workspaceName, datastoreName, featureType)
	if err != nil {
		client.DeleteFeatureType(workspaceName, datastoreName, d.Get("name").(string), true)
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", workspaceName, datastoreName, d.Get("name").(string)))
	d.Set("use_custom_attributes", len(attributes) > 0)

	return resourceGeoserverFeatureTypeRead(d, meta)
}

func resourceGeoserverFeatureTypeRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Geoserver FeatureType: %s", d.Id())

	splittedID := strings.Split(d.Id(), "/")
	workspaceName := splittedID[0]
	datastoreName := splittedID[1]
	featureTypeName := splittedID[2]

	client := meta.(*Config).GeoserverClient()

	featureType, err := client.GetFeatureType(workspaceName, datastoreName, featureTypeName)
	if err != nil && !strings.Contains(strings.ToLower(err.Error()), "not found") {
		return err
	}

	if featureType == nil {
		d.SetId("")
		return nil
	}
	d.Set("workspace_name", workspaceName)
	d.Set("datastore_name", datastoreName)
	d.Set("name", featureType.Name)
	d.Set("native_name", featureType.NativeName)
	d.Set("enabled", featureType.Enabled)
	d.Set("projection_policy", featureType.ProjectionPolicy)
	d.Set("srs", featureType.SRS)

	d.Set("native_crs_class", featureType.NativeCRS.Class)
	d.Set("native_crs_value", featureType.NativeCRS.Value)

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
	if d.Get("use_custom_attributes").(bool) {
		for _, value := range featureType.Attributes {
			attributes = append(attributes, map[string]interface{}{
				"name":       value.Name,
				"nillable":   value.Nillable,
				"binding":    value.Binding,
				"min_occurs": value.MinOccurs,
				"max_occurs": value.MaxOccurs,
			})
		}
		d.Set("attribute", attributes)
	}

	metadata := map[string]interface{}{}
	if featureType.Metadata != nil {
		for _, entry := range featureType.Metadata {
			metadata[entry.Key] = entry.Value
		}
		d.Set("metadata", metadata)
	}

	return nil
}

func resourceGeoserverFeatureTypeDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Geoserver Feature Type: %s", d.Id())

	splittedID := strings.Split(d.Id(), "/")
	workspaceName := splittedID[0]
	datastoreName := splittedID[1]
	featureTypeName := splittedID[2]

	client := meta.(*Config).GeoserverClient()

	layer, errget := client.GetLayer(workspaceName, featureTypeName)
	if layer != nil {
		// If there is a matching layer into geoserver, delete it first
		err1 := client.DeleteLayer(workspaceName, featureTypeName, true)
		if err1 != nil {
			return err1
		}

		err2 := client.DeleteFeatureType(workspaceName, datastoreName, featureTypeName, false)
		if err2 != nil {
			return err2
		}
	} else if errget.Error() == "not found" {
		// If not, delete only the feature type
		err2 := client.DeleteFeatureType(workspaceName, datastoreName, featureTypeName, true)
		if err2 != nil {
			return err2
		}
	} else {
		return errget
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

	client := meta.(*Config).GeoserverClient()

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

	var metadata []*gs.FeatureTypeMetadata
	for key, value := range d.Get("metadata").(map[string]interface{}) {
		metadata = append(metadata, &gs.FeatureTypeMetadata{
			Key:   key,
			Value: value.(string),
		})
	}

	featureType := &gs.FeatureType{
		Name:             d.Get("name").(string),
		NativeName:       d.Get("native_name").(string),
		Enabled:          d.Get("enabled").(bool),
		ProjectionPolicy: d.Get("projection_policy").(string),
		Title:            d.Get("title").(string),
		Abstract:         d.Get("abstract").(string),
		NativeCRS: gs.FeatureTypeCRS{
			Class: d.Get("native_crs_class").(string),
			Value: d.Get("native_crs_value").(string),
		},
		SRS: d.Get("srs").(string),
		NativeBoundingBox: gs.BoundingBox{
			MinX: d.Get("native_bounding_box_min_x").(float64),
			MaxX: d.Get("native_bounding_box_max_x").(float64),
			MinY: d.Get("native_bounding_box_min_y").(float64),
			MaxY: d.Get("native_bounding_box_max_y").(float64),
			CRS: gs.FeatureTypeCRS{
				Class: d.Get("native_bounding_box_crs_class").(string),
				Value: d.Get("native_bounding_box_crs_value").(string),
			},
		},
		LatLonBoundingBox: gs.BoundingBox{
			MinX: d.Get("lat_lon_bounding_box_min_x").(float64),
			MaxX: d.Get("lat_lon_bounding_box_max_x").(float64),
			MinY: d.Get("lat_lon_bounding_box_min_y").(float64),
			MaxY: d.Get("lat_lon_bounding_box_max_y").(float64),
			CRS: gs.FeatureTypeCRS{
				Class: d.Get("lat_lon_bounding_box_crs_class").(string),
				Value: d.Get("lat_lon_bounding_box_crs_value").(string),
			},
		},
		Attributes: attributes,
		Metadata:   metadata,
	}

	sync_attributes := !d.Get("use_custom_attributes").(bool)
	err := client.UpdateFeatureType(workspaceName, datastoreName, featureTypeName, featureType, sync_attributes)
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
	d.Set("use_custom_attributes", true)

	log.Printf("[INFO] Importing Geoserver FeatureType `%s` in workspace `%s`", datastoreName, workspaceName)

	err := resourceGeoserverFeatureTypeRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
