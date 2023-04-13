package geoserver

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	gs "github.com/camptocamp/go-geoserver/client"
)

func resourceGwcWmsLayer() *schema.Resource {
	return &schema.Resource{
		Create: resourceGwcWmsLayerCreate,
		Read:   resourceGwcWmsLayerRead,
		Update: resourceGwcWmsLayerUpdate,
		Delete: resourceGwcWmsLayerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGwcWmsLayerImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"blobstore_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"wms_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"wms_layer": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mime_formats": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"grid_subsets": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"metatile_height": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"metatile_width": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"allow_cache_bypass": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"expire_duration_cache": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"expire_duration_clients": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"gutter_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"backend_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  120,
			},
		},
	}
}

func resourceGwcWmsLayerCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Creating GWC wms layer: %s", d.Id())

	client := meta.(*Config).GwcClient()

	layerName := d.Get("name").(string)

	var mimeFormats []string
	for _, value := range d.Get("mime_formats").([]interface{}) {
		mimeFormats = append(mimeFormats,
			value.(string),
		)
	}

	var gridSubsets []*gs.GridSubset
	for _, value := range d.Get("grid_subsets").([]interface{}) {
		gridSubsets = append(gridSubsets,
			&gs.GridSubset{
				Name: value.(string),
			},
		)
	}

	var metaTilesDimension []int
	metaTilesDimension = append(metaTilesDimension, d.Get("metatile_width").(int))
	metaTilesDimension = append(metaTilesDimension, d.Get("metatile_height").(int))

	wmsLayer := &gs.GwcWmsLayer{
		Name:                 layerName,
		Enabled:              d.Get("enabled").(bool),
		BlobStoreId:          d.Get("blobstore_id").(string),
		MimeFormats:          gs.MimeFormats{Formats: mimeFormats},
		GridSubsets:          gridSubsets,
		MetaTileDimensions:   metaTilesDimension,
		ExpireCacheDuration:  d.Get("expire_duration_cache").(int),
		ExpireClientDuration: d.Get("expire_duration_clients").(int),
		GutterSize:           d.Get("gutter_size").(int),
		BackendTimeout:       d.Get("backend_timeout").(int),
		CacheBypassAllowed:   d.Get("allow_cache_bypass").(bool),
		WmsUrl:               d.Get("wms_url").(string),
		WmsLayer:             d.Get("wms_layer").(string),
	}

	err := client.CreateGwcWmsLayer(layerName, wmsLayer)
	if err != nil {
		return err
	}

	d.SetId(layerName)

	return resourceGwcWmsLayerRead(d, meta)
}

func resourceGwcWmsLayerRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing GWC WMS Layer: %s", d.Id())

	wmsLayerName := d.Id()

	client := meta.(*Config).GwcClient()

	wmsLayer, err := client.GetGwcWMSLayer(wmsLayerName)
	if err != nil && !strings.Contains(strings.ToLower(err.Error()), "not found") {
		return err
	}

	if wmsLayer == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", wmsLayerName)
	d.Set("enabled", wmsLayer.Enabled)
	d.Set("blobstore_id", wmsLayer.BlobStoreId)
	d.Set("expire_duration_cache", wmsLayer.ExpireCacheDuration)
	d.Set("expire_duration_clients", wmsLayer.ExpireClientDuration)
	d.Set("gutter_size", wmsLayer.GutterSize)
	d.Set("backend_timeout", wmsLayer.BackendTimeout)
	d.Set("allow_cache_bypass", wmsLayer.CacheBypassAllowed)
	d.Set("wms_url", wmsLayer.WmsUrl)
	d.Set("wms_layer", wmsLayer.WmsLayer)
	d.Set("metatile_width", wmsLayer.MetaTileDimensions[0])
	d.Set("metatile_height", wmsLayer.MetaTileDimensions[1])

	var mimeFormats []string
	mimeFormats = append(mimeFormats, wmsLayer.MimeFormats.Formats...)

	d.Set("mime_formats", mimeFormats)

	var gridsubsets []string
	for _, value := range wmsLayer.GridSubsets {
		gridsubsets = append(gridsubsets, value.Name)
	}
	d.Set("grid_subsets", gridsubsets)

	return nil
}

func resourceGwcWmsLayerDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting GWC wms layer: %s", d.Id())

	wmsLayerName := d.Id()

	client := meta.(*Config).GwcClient()

	err := client.DeleteGwcWmsLayer(wmsLayerName)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceGwcWmsLayerUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating GWC wms layer: %s", d.Id())

	wmsLayerName := d.Id()

	client := meta.(*Config).GwcClient()

	var mimeFormats []string
	for _, value := range d.Get("mime_formats").([]interface{}) {
		mimeFormats = append(mimeFormats,
			value.(string),
		)
	}

	var gridSubsets []*gs.GridSubset
	for _, value := range d.Get("grid_subsets").([]interface{}) {
		gridSubsets = append(gridSubsets,
			&gs.GridSubset{
				Name: value.(string),
			},
		)
	}

	var metaTilesDimension []int
	metaTilesDimension = append(metaTilesDimension, d.Get("metatile_width").(int))
	metaTilesDimension = append(metaTilesDimension, d.Get("metatile_height").(int))

	err := client.UpdateGwcWmsLayer(wmsLayerName, &gs.GwcWmsLayer{
		Name:                 wmsLayerName,
		Enabled:              d.Get("enabled").(bool),
		BlobStoreId:          d.Get("blobstore_id").(string),
		MimeFormats:          gs.MimeFormats{Formats: mimeFormats},
		GridSubsets:          gridSubsets,
		MetaTileDimensions:   metaTilesDimension,
		ExpireCacheDuration:  d.Get("expire_duration_cache").(int),
		ExpireClientDuration: d.Get("expire_duration_clients").(int),
		GutterSize:           d.Get("gutter_size").(int),
		BackendTimeout:       d.Get("backend_timeout").(int),
		CacheBypassAllowed:   d.Get("allow_cache_bypass").(bool),
		WmsUrl:               d.Get("wms_url").(string),
		WmsLayer:             d.Get("wms_layer").(string),
	})
	if err != nil {
		return err
	}

	return nil
}

func resourceGwcWmsLayerImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	wmsLayerName := d.Id()

	d.SetId(d.Id())
	d.Set("name", wmsLayerName)

	log.Printf("[INFO] Importing GWC wms layer `%s`", wmsLayerName)

	err := resourceGwcWmsLayerRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
