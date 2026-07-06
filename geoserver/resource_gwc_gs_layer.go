package geoserver

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	gs "github.com/camptocamp/go-geoserver/client"
)

func resourceGwcGsLayer() *schema.Resource {
	return &schema.Resource{
		Create: resourceGwcGsLayerCreate,
		Read:   resourceGwcGsLayerRead,
		Update: resourceGwcGsLayerUpdate,
		Delete: resourceGwcGsLayerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGwcGsLayerImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the cache layer. Use as resource id.",
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Internal id of the GeoServer layer.",
			},
			"blobstore_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Blobstore to use for storing the tiles.",
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
		},
	}
}

func resourceGwcGsLayerCreate(d *schema.ResourceData, meta interface{}) error {
	layerName := d.Get("name").(string)

	log.Printf("[INFO] Creating GWC GS layer: %s", layerName)

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

	gsLayer := &gs.GwcGsLayer{
		Name:                 layerName,
		Enabled:              d.Get("enabled").(bool),
		Id:                   d.Get("id").(string),
		BlobStoreId:          d.Get("blobstore_id").(string),
		MimeFormats:          gs.MimeFormats{Formats: mimeFormats},
		GridSubsets:          gridSubsets,
		MetaTileDimensions:   metaTilesDimension,
		ExpireCacheDuration:  d.Get("expire_duration_cache").(int),
		ExpireClientDuration: d.Get("expire_duration_clients").(int),
		GutterSize:           d.Get("gutter_size").(int),
		CacheBypassAllowed:   d.Get("allow_cache_bypass").(bool),
	}

	err := client.CreateGwcGsLayer(layerName, gsLayer)
	if err != nil {
		return err
	}

	d.SetId(layerName)

	return resourceGwcGsLayerRead(d, meta)
}

func resourceGwcGsLayerRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing GWC GS Layer: %s", d.Id())

	gsLayerName := d.Id()

	client := meta.(*Config).GwcClient()

	gsLayer, err := client.GetGwcGsLayer(gsLayerName)
	if err != nil && !strings.Contains(strings.ToLower(err.Error()), "not found") {
		return err
	}

	if gsLayer == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", gsLayerName)
	d.Set("enabled", gsLayer.Enabled)
	d.Set("blobstore_id", gsLayer.BlobStoreId)
	d.Set("id", gsLayer.Id)
	d.Set("expire_duration_cache", gsLayer.ExpireCacheDuration)
	d.Set("expire_duration_clients", gsLayer.ExpireClientDuration)
	d.Set("gutter_size", gsLayer.GutterSize)
	d.Set("allow_cache_bypass", gsLayer.CacheBypassAllowed)
	d.Set("metatile_width", gsLayer.MetaTileDimensions[0])
	d.Set("metatile_height", gsLayer.MetaTileDimensions[1])

	var mimeFormats []string
	mimeFormats = append(mimeFormats, gsLayer.MimeFormats.Formats...)

	d.Set("mime_formats", mimeFormats)

	var gridsubsets []map[string]interface{}
	for _, value := range gsLayer.GridSubsets {
		gridsubsets = append(gridsubsets, map[string]interface{}{
			"name":             value.Name,
			"min_cached_level": value.MinCacheLevel,
			"max_cached_level": value.MaxCacheLevel,
		})
	}

	d.Set("grid_subset", gridsubsets)

	return nil
}

func resourceGwcGsLayerDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting GWC gs layer: %s", d.Id())

	gsLayerName := d.Id()

	client := meta.(*Config).GwcClient()

	err := client.DeleteGwcGsLayer(gsLayerName)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceGwcGsLayerUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating GWC gs layer: %s", d.Id())

	gsLayerName := d.Id()

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

	err := client.UpdateGwcGsLayer(gsLayerName, &gs.GwcGsLayer{
		Name:                 gsLayerName,
		Id:                   d.Get("id").(string),
		Enabled:              d.Get("enabled").(bool),
		BlobStoreId:          d.Get("blobstore_id").(string),
		MimeFormats:          gs.MimeFormats{Formats: mimeFormats},
		GridSubsets:          gridSubsets,
		MetaTileDimensions:   metaTilesDimension,
		ExpireCacheDuration:  d.Get("expire_duration_cache").(int),
		ExpireClientDuration: d.Get("expire_duration_clients").(int),
		GutterSize:           d.Get("gutter_size").(int),
		CacheBypassAllowed:   d.Get("allow_cache_bypass").(bool),
	})
	if err != nil {
		return err
	}

	return nil
}

func resourceGwcGsLayerImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	gsLayerName := d.Id()

	d.SetId(d.Id())
	d.Set("name", gsLayerName)

	log.Printf("[INFO] Importing GWC wms layer `%s`", gsLayerName)

	err := resourceGwcGsLayerRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
