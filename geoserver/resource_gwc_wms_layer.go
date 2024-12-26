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
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the cache layer. Use as resource id.",
			},
			"blobstore_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Blobstore to use for storing the tiles.",
			},
			"wms_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "URL of the server hosting the WMS.",
			},
			"wms_layer": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identifier of the layer to cache.",
			},
			"wms_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "1.3.0",
				Description: "WMS version to use when requesting service. Default to 1.3.0",
			},
			"vendor_parameters": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Additional vendor parameters to the service.",
			},
			"background_color": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Background color when requesting maps.",
			},
			"mime_formats": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of the mime formats upported for this cache.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"grid_subset": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "List of the grids supported for this cache.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"max_cached_level": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"min_cached_level": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"metatile_height": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Height of the meta tiles.",
			},
			"metatile_width": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Width of the meta tiles.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Is the layer cache enabled? Default to true.",
			},
			"allow_cache_bypass": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Allow bypass of the cache. Default to false.",
			},
			"expire_duration_cache": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Server cache expire duration. Default to 0.",
			},
			"expire_duration_clients": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Client cache expire duration. Default to 0.",
			},
			"gutter_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Size of the gutter to use for the meta-tiles. Default to 0.",
			},
			"backend_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     120,
				Description: "Timeout of the backend. Default to 120.",
			},
			"transparent": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Request tiles with transparent background? Default to true.",
			},
		},
	}
}

func resourceGwcWmsLayerCreate(d *schema.ResourceData, meta interface{}) error {
	layerName := d.Get("name").(string)

	log.Printf("[INFO] Creating GWC wms layer: %s", layerName)

	client := meta.(*Config).GwcClient()

	var mimeFormats []string
	for _, value := range d.Get("mime_formats").([]interface{}) {
		mimeFormats = append(mimeFormats,
			value.(string),
		)
	}

	var gridSubsets []*gs.GridSubset
	for _, value := range d.Get("grid_subset").(*schema.Set).List() {
		v := value.(map[string]interface{})
		gridSubset := &gs.GridSubset{
			Name: v["name"].(string),
		}

		minCacheLevel, ok := v["min_cached_level"]
		if ok {
			gridSubset.MinCacheLevel = minCacheLevel.(int)
		}
		maxCacheLevel, ok2 := v["max_cached_level"]
		if ok2 {
			gridSubset.MaxCacheLevel = maxCacheLevel.(int)
		}

		gridSubsets = append(gridSubsets, gridSubset)
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
		Transparent:          d.Get("transparent").(bool),
		WmsVersion:           d.Get("wms_version").(string),
		BgColor:              d.Get("background_color").(string),
		VendorParameters:     d.Get("vendor_parameters").(string),
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
	d.Set("wms_version", wmsLayer.WmsVersion)
	d.Set("vendor_parameters", wmsLayer.VendorParameters)
	d.Set("background_color", wmsLayer.BgColor)
	d.Set("transparent", wmsLayer.Transparent)
	d.Set("metatile_width", wmsLayer.MetaTileDimensions[0])
	d.Set("metatile_height", wmsLayer.MetaTileDimensions[1])

	var mimeFormats []string
	mimeFormats = append(mimeFormats, wmsLayer.MimeFormats.Formats...)

	d.Set("mime_formats", mimeFormats)

	var gridsubsets []map[string]interface{}
	for _, value := range wmsLayer.GridSubsets {
		gridsubsets = append(gridsubsets, map[string]interface{}{
			"name":             value.Name,
			"min_cached_level": value.MinCacheLevel,
			"max_cached_level": value.MaxCacheLevel,
		})
	}

	d.Set("grid_subset", gridsubsets)

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
	for _, value := range d.Get("grid_subset").(*schema.Set).List() {
		v := value.(map[string]interface{})
		gridSubset := &gs.GridSubset{
			Name: v["name"].(string),
		}

		minCacheLevel, ok := v["min_cached_level"]
		if ok {
			gridSubset.MinCacheLevel = minCacheLevel.(int)
		}
		maxCacheLevel, ok2 := v["max_cached_level"]
		if ok2 {
			gridSubset.MaxCacheLevel = maxCacheLevel.(int)
		}

		gridSubsets = append(gridSubsets, gridSubset)
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
		Transparent:          d.Get("transparent").(bool),
		WmsVersion:           d.Get("wms_version").(string),
		BgColor:              d.Get("background_color").(string),
		VendorParameters:     d.Get("vendor_parameters").(string),
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
