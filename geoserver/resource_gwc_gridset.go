package geoserver

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	gs "github.com/camptocamp/go-geoserver/client"
)

func resourceGwcGridset() *schema.Resource {
	return &schema.Resource{
		Create: resourceGwcGridsetCreate,
		Read:   resourceGwcGridsetRead,
		Update: resourceGwcGridsetUpdate,
		Delete: resourceGwcGridsetDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGwcGridsetImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the gridset. Use as resource id.",
			},
			"meters_per_unit": {
				Type:        schema.TypeFloat,
				Required:    true,
				Description: "Conversion to/from meter",
			},
			"pixel_size": {
				Type:        schema.TypeFloat,
				Required:    true,
				Description: "Resolution",
			},
			"tile_height": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Hauteur d'une tuile",
			},
			"tile_width": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Largeur d'une tuile",
			},
			"srs": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Identifier of the SRS",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"align_top_left": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Is the top left corner aligned? Default to true.",
			},
			"y_coordinate_first": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Does the Y coordinate come first? Default to true.",
			},
			"scales": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Scales of the set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the level (0,1,2,...)",
						},
						"denominator": {
							Type:        schema.TypeFloat,
							Required:    true,
							Description: "Scale denominator",
						},
					},
				},
			},
			"extent_min_x": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"extent_max_x": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"extent_min_y": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"extent_max_y": {
				Type:     schema.TypeFloat,
				Required: true,
			},
		},
	}
}

func resourceGwcGridsetCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Creating GWC gridset: %s", d.Id())

	client := meta.(*Config).GwcClient()

	gridsetName := d.Get("name").(string)

	var scaleNames []string
	var scaleDenominators []float64
	for _, value := range d.Get("scales").([]interface{}) {
		v := value.(map[string]interface{})
		scaleNames = append(scaleNames,
			v["name"].(string),
		)
		scaleDenominators = append(scaleDenominators,
			v["denominator"].(float64),
		)
	}

	gridset := &gs.Gridset{
		Name:              gridsetName,
		Description:       d.Get("description").(string),
		AlignTopLeft:      d.Get("align_top_left").(bool),
		MetersPerUnit:     d.Get("meters_per_unit").(float64),
		PixelSize:         d.Get("pixel_size").(float64),
		TileHeight:        d.Get("tile_height").(int),
		TileWidth:         d.Get("tile_width").(int),
		YCoordinateFirst:  d.Get("y_coordinate_first").(bool),
		Extent:            []float64{d.Get("extent_min_x").(float64), d.Get("extent_min_y").(float64), d.Get("extent_max_x").(float64), d.Get("extent_max_y").(float64)},
		ScaleNames:        gs.ScaleNames{ScaleName: scaleNames},
		ScaleDenominators: gs.ScaleDenominators{ScaleDenominator: scaleDenominators},
		Srs:               gs.SRS{SrsNumber: d.Get("srs").(int)},
	}

	err := client.CreateGridset(gridsetName, gridset)
	if err != nil {
		return err
	}

	d.SetId(gridsetName)

	return resourceGwcGridsetRead(d, meta)
}

func resourceGwcGridsetRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing GWC Gridset: %s", d.Id())

	gridsetName := d.Id()

	client := meta.(*Config).GwcClient()

	gridSet, err := client.GetGridset(gridsetName)
	if err != nil && !strings.Contains(strings.ToLower(err.Error()), "not found") {
		return err
	}

	if gridSet == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", gridsetName)
	d.Set("description", gridSet.Description)
	d.Set("align_top_left", gridSet.AlignTopLeft)
	d.Set("meters_per_unit", gridSet.MetersPerUnit)
	d.Set("pixel_size", gridSet.PixelSize)
	d.Set("tile_height", gridSet.TileHeight)
	d.Set("tile_width", gridSet.TileWidth)
	d.Set("y_coordinate_first", gridSet.YCoordinateFirst)
	d.Set("extent_min_x", gridSet.Extent[0])
	d.Set("extent_min_y", gridSet.Extent[1])
	d.Set("extent_max_x", gridSet.Extent[2])
	d.Set("extent_max_y", gridSet.Extent[3])
	d.Set("srs", gridSet)

	var scales []map[string]interface{}
	for index, value := range gridSet.ScaleNames.ScaleName {
		scales = append(scales, map[string]interface{}{
			"name":        value,
			"denominator": gridSet.ScaleDenominators.ScaleDenominator[index],
		})
	}
	d.Set("scales", scales)

	return nil
}

func resourceGwcGridsetDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting GWC gridset: %s", d.Id())

	gridsetName := d.Id()

	client := meta.(*Config).GwcClient()

	err := client.DeleteGridset(gridsetName)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceGwcGridsetUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating GWC gridset: %s", d.Id())

	gridsetName := d.Id()

	client := meta.(*Config).GwcClient()

	var scaleNames []string
	var scaleDenominators []float64
	for _, value := range d.Get("scales").([]interface{}) {
		v := value.(map[string]interface{})
		scaleNames = append(scaleNames,
			v["name"].(string),
		)
		scaleDenominators = append(scaleDenominators,
			v["denominator"].(float64),
		)
	}

	err := client.UpdateGridset(gridsetName, &gs.Gridset{
		Name:              gridsetName,
		Description:       d.Get("description").(string),
		AlignTopLeft:      d.Get("align_top_left").(bool),
		MetersPerUnit:     d.Get("meters_per_unit").(float64),
		PixelSize:         d.Get("pixel_size").(float64),
		TileHeight:        d.Get("tile_height").(int),
		TileWidth:         d.Get("tile_width").(int),
		YCoordinateFirst:  d.Get("y_coordinate_first").(bool),
		Extent:            []float64{d.Get("extent_min_x").(float64), d.Get("extent_min_y").(float64), d.Get("extent_max_x").(float64), d.Get("extent_max_y").(float64)},
		ScaleNames:        gs.ScaleNames{ScaleName: scaleNames},
		ScaleDenominators: gs.ScaleDenominators{ScaleDenominator: scaleDenominators},
		Srs:               gs.SRS{SrsNumber: d.Get("srs").(int)},
	})
	if err != nil {
		return err
	}

	return nil
}

func resourceGwcGridsetImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	gridsetName := d.Id()

	d.SetId(d.Id())
	d.Set("name", gridsetName)

	log.Printf("[INFO] Importing GWC gridset `%s`", gridsetName)

	err := resourceGwcGridsetRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
