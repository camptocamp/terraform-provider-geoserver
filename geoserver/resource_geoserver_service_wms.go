package geoserver

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	gs "github.com/camptocamp/go-geoserver/client"
)

func resourceGeoServerServiceWms() *schema.Resource {
	return &schema.Resource{
		Description: "Manage global WMS configuration. The configuration is a singleton so Create is similar to Update and Delete has no effect.",
		Create:      resourceGeoServerServiceWmsCreate,
		Read:        resourceGeoServerServiceWmsRead,
		Update:      resourceGeoServerServiceWmsUpdate,
		Delete:      resourceGeoServerServiceWmsDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGeoServerServiceWmsImport,
		},

		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Is the WMS service enabled?",
			},
			"title": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Title of the service",
			},
			"maintainer": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Maintainer of the service",
			},
			"abstract": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the service",
			},
			"access_constraints": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specific access constraints of the service",
			},
			"fees": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Fee details for the service",
			},
			"online_resource": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Additional online resources about the service",
			},
			"schema_base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "http://schemas.opengis.net",
				Description: "The base url for the schemas describing the service.",
			},
			"is_verbose": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Flag indicating if the service should be verbose or not",
			},
			"use_bbox_foreach_crs": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Flag indicating if watermarking is enabled",
			},
			"watermark_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Flag indicating if watermarking is enabled",
			},
			"watermark_position": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "TOP_LEFT",
				Description: "Position of the watermark. Authorized values are : TOP_LEFT, TOP_CENTER, TOP_RIGHT, MID_LEFT, MID_CENTER, MID_RIGHT, BOT_LEFT, BOT_CENTER, BOT_RIGHT",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					allowed_values := []string{"TOP_LEFT", "TOP_CENTER", "TOP_RIGHT", "MID_LEFT", "MID_CENTER", "MID_RIGHT", "BOT_LEFT", "BOT_CENTER", "BOT_RIGHT"}
					if !slices.Contains(allowed_values, v) {
						errs = append(errs, fmt.Errorf("%q must be one of this values %q, got: %q", key, strings.Join(allowed_values, ","), v))
					}
					return
				},
			},
			"watermark_transparency": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
				Description: "The transparency of the watermark logo, ranging from 0 to 255",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 || v > 255 {
						errs = append(errs, fmt.Errorf("%q must be one between 0 and 255, got: %q", key, v))
					}
					return
				},
			},
			"interpolation": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Nearest",
				Description: "Interpolation strategy. Authorized values are : Nearest, Bilinear, Bicubic",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					allowed_values := []string{"Nearest", "Bilinear", "Bicubic"}
					if !slices.Contains(allowed_values, v) {
						errs = append(errs, fmt.Errorf("%q must be one of this values %q, got: %q", key, strings.Join(allowed_values, ","), v))
					}
					return
				},
			},
			"is_cite_compliant": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Status of service CITE compliance",
			},
			"maximum_buffer": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Maximum search radius for GetFeatureInfo",
			},
			"is_dynamic_styling_disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "status of dynamic styling (SLD and SLD_BODY params) allowance",
			},
			"is_getfeatureinfo_mimetype_checking_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Flag indicating if getFeatureInfo MIME type checking is enabled",
			},
			"maximum_request_memory": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Max amount of memory, in kilobytes, that each WMS request can allocate (each output format will make a best effort attempt to respect it, but there are no guarantees). 0 indicates no limit.",
			},
			"maximum_rendering_errors": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Max number of rendering errors that will be tolerated before stating the rendering operation failed by throwing a service exception back to the client",
			},
			"maximum_rendering_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Max time, in seconds, a WMS request is allowed to spend rendering the map. Various output formats will do a best effort to respect it (raster formats, for example, will account just rendering time, but not image encoding time).",
			},
			"workspace_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the workspace owning the style. Used to compute the id of the resource.",
			},
			"supported_versions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "	The versions of the service that are available.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"is_getmap_mimetype_checking_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Flag indicating if getMap MIME type checking is enabled.",
			},
			"is_features_reprojection_disabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"maximum_requested_dimension_values": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  100,
			},
			"is_cache_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"cache_maximum_entries": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1000,
			},
			"cache_maximum_entry_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  51200,
			},
			"remote_style_max_request_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  60000,
			},
			"remote_style_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30000,
			},
			"is_default_group_style_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"is_transform_feature_info_disabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"is_autoescape_templatevalues_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"root_layer_title": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

func resourceGeoServerServiceWmsCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Configuring WMS Service")

	client := meta.(*Config).GeoserverClient()

	var keywords []string
	for _, value := range d.Get("keywords").([]interface{}) {
		keywords = append(keywords, value.(string))
	}

	var versions []*gs.ServiceVersion
	for _, value := range d.Get("supported_versions").([]interface{}) {
		versions = append(versions, &gs.ServiceVersion{Version: value.(string)})
	}

	var serviceMetadata []*gs.ServiceWmsMetadata
	for key, value := range d.Get("metadata").(map[string]interface{}) {
		serviceMetadata = append(serviceMetadata, &gs.ServiceWmsMetadata{
			Key:   key,
			Value: value.(string),
		})
	}

	workspaceName := d.Get("workspace_name").(string)

	serviceConfiguration := &gs.ServiceWms{
		Name:              "WMS",
		IsEnabled:         d.Get("enabled").(bool),
		Title:             d.Get("title").(string),
		Maintainer:        d.Get("maintainer").(string),
		Abstract:          d.Get("abstract").(string),
		AccessConstraints: d.Get("access_constraints").(string),
		OnlineResource:    d.Get("online_resource").(string),
		IsVerbose:         d.Get("is_verbose").(bool),
		Watermark: gs.Watermark{
			IsEnabled:    d.Get("watermark_enabled").(bool),
			Position:     d.Get("watermark_position").(string),
			Transparency: d.Get("watermark_transparency").(int),
		},
		Interpolation:                           d.Get("interpolation").(string),
		IsCiteCompliant:                         d.Get("is_cite_compliant").(bool),
		MaximumBuffer:                           d.Get("maximum_buffer").(int),
		IsDynamicStylingDisabled:                d.Get("is_dynamic_styling_disabled").(bool),
		Metadata:                                serviceMetadata,
		Keywords:                                gs.ServiceWmsKeywords{Keywords: keywords},
		IsGetFeatureInfoMimeTypeCheckingEnabled: d.Get("is_getfeatureinfo_mimetype_checking_enabled").(bool),
		MaximumRequestMemory:                    d.Get("maximum_request_memory").(int),
		Fees:                                    d.Get("fees").(string),
		MaximumRenderingErrors:                  d.Get("maximum_rendering_errors").(int),
		MaximumRenderingTime:                    d.Get("maximum_rendering_time").(int),
		SupportedVersions:                       versions,
		SchemaBaseURL:                           d.Get("schema_base_url").(string),
		UseBBOXForEachCRS:                       d.Get("use_bbox_foreach_crs").(bool),
		IsGetMapMimeTypeCheckingEnabled:         d.Get("is_getmap_mimetype_checking_enabled").(bool),
		IsFeaturesReprojectionDisabled:          d.Get("is_features_reprojection_disabled").(bool),
		MaximumRequestedDimensionValues:         d.Get("maximum_requested_dimension_values").(int),
		CacheConfiguration: gs.CacheConfiguration{
			IsEnabled:    d.Get("is_cache_enabled").(bool),
			MaxEntrySize: d.Get("cache_maximum_entry_size").(int),
			MaxEntries:   d.Get("cache_maximum_entries").(int),
		},
		RemoteStyleMaxRequestTime:         d.Get("remote_style_max_request_time").(int),
		RemoteStyleTimeout:                d.Get("remote_style_timeout").(int),
		IsDefaultGroupStyleEnabled:        d.Get("is_default_group_style_enabled").(bool),
		IsTransformFeatureInfoDisabled:    d.Get("is_transform_feature_info_disabled").(bool),
		IsAutoEscapeTemplateValuesEnabled: d.Get("is_autoescape_templatevalues_enabled").(bool),
		RootLayerTitle:                    d.Get("root_layer_title").(string),
	}

	err := client.UpdateServiceWMS(workspaceName, serviceConfiguration)
	if err != nil {
		return err
	}

	if d.Get("workspace_name").(string) == "" {
		d.SetId("wms_service_configuration")
	} else {
		d.SetId(fmt.Sprintf("wms_service_configuration/%s", workspaceName))
	}

	return resourceGeoServerServiceWmsRead(d, meta)
}

func resourceGeoServerServiceWmsRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing WMS Service configuration: %s", d.Id())

	var workspaceName string
	if strings.Contains(d.Id(), "/") {
		splittedID := strings.Split(d.Id(), "/")
		workspaceName = splittedID[1]
	}

	client := meta.(*Config).GeoserverClient()

	wmsConfiguration, err := client.GetServiceWMS(workspaceName)

	if err != nil && !strings.Contains(strings.ToLower(err.Error()), "not found") {
		return err
	}

	if wmsConfiguration == nil {
		d.SetId("")
		return nil
	}

	d.Set("enabled", wmsConfiguration.IsEnabled)
	d.Set("title", wmsConfiguration.Title)
	d.Set("maintainer", wmsConfiguration.Maintainer)
	d.Set("abstract", wmsConfiguration.Abstract)
	d.Set("access_constraints", wmsConfiguration.AccessConstraints)
	d.Set("fees", wmsConfiguration.Fees)
	d.Set("online_resource", wmsConfiguration.OnlineResource)
	d.Set("schema_base_url", wmsConfiguration.SchemaBaseURL)
	d.Set("is_verbose", wmsConfiguration.IsVerbose)
	d.Set("use_bbox_foreach_crs", wmsConfiguration.UseBBOXForEachCRS)
	d.Set("watermark_enabled", wmsConfiguration.Watermark.IsEnabled)
	d.Set("watermark_position", wmsConfiguration.Watermark.Position)
	d.Set("watermark_transparency", wmsConfiguration.Watermark.Transparency)
	d.Set("interpolation", wmsConfiguration.Interpolation)
	d.Set("is_cite_compliant", wmsConfiguration.IsCiteCompliant)
	d.Set("maximum_buffer", wmsConfiguration.MaximumBuffer)
	d.Set("is_dynamic_styling_disabled", wmsConfiguration.IsDynamicStylingDisabled)
	d.Set("is_getfeatureinfo_mimetype_checking_enabled", wmsConfiguration.IsGetFeatureInfoMimeTypeCheckingEnabled)
	d.Set("maximum_request_memory", wmsConfiguration.MaximumRequestMemory)
	d.Set("maximum_rendering_errors", wmsConfiguration.MaximumRenderingErrors)
	d.Set("maximum_rendering_time", wmsConfiguration.MaximumRenderingTime)
	if wmsConfiguration.Workspace != nil {
		d.Set("workspace_name", wmsConfiguration.Workspace.Name)
	}
	d.Set("is_getmap_mimetype_checking_enabled", wmsConfiguration.IsGetMapMimeTypeCheckingEnabled)
	d.Set("is_features_reprojection_disabled", wmsConfiguration.IsFeaturesReprojectionDisabled)
	d.Set("maximum_requested_dimension_values", wmsConfiguration.MaximumRequestedDimensionValues)
	d.Set("is_cache_enabled", wmsConfiguration.CacheConfiguration.IsEnabled)
	d.Set("cache_maximum_entries", wmsConfiguration.CacheConfiguration.MaxEntries)
	d.Set("cache_maximum_entry_size", wmsConfiguration.CacheConfiguration.MaxEntrySize)
	d.Set("remote_style_max_request_time", wmsConfiguration.RemoteStyleMaxRequestTime)
	d.Set("remote_style_timeout", wmsConfiguration.RemoteStyleTimeout)
	d.Set("is_default_group_style_enabled", wmsConfiguration.IsDefaultGroupStyleEnabled)
	d.Set("is_transform_feature_info_disabled", wmsConfiguration.IsTransformFeatureInfoDisabled)
	d.Set("is_autoescape_templatevalues_enabled", wmsConfiguration.IsAutoEscapeTemplateValuesEnabled)
	d.Set("root_layer_title", wmsConfiguration.RootLayerTitle)
	d.Set("is_transform_feature_info_disabled", wmsConfiguration.IsTransformFeatureInfoDisabled)

	var supportedVersions []string
	for _, value := range wmsConfiguration.SupportedVersions {
		supportedVersions = append(supportedVersions, value.Version)
	}

	d.Set("supported_versions", supportedVersions)

	var keywords []string
	keywords = append(keywords, wmsConfiguration.Keywords.Keywords...)
	d.Set("keywords", keywords)

	metadata := map[string]string{}
	for _, entry := range wmsConfiguration.Metadata {
		metadata[entry.Key] = entry.Value
	}
	d.Set("metadata", metadata)

	return nil
}

func resourceGeoServerServiceWmsDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Disk Quota cannot be deleted. Skipping action")

	d.SetId("")

	return nil
}

func resourceGeoServerServiceWmsUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating WMS Service configuration: %s", d.Id())

	client := meta.(*Config).GeoserverClient()

	var keywords []string
	for _, value := range d.Get("keywords").([]interface{}) {
		keywords = append(keywords, value.(string))
	}

	var versions []*gs.ServiceVersion
	for _, value := range d.Get("supported_versions").([]interface{}) {
		versions = append(versions, &gs.ServiceVersion{Version: value.(string)})
	}

	var serviceMetadata []*gs.ServiceWmsMetadata
	for key, value := range d.Get("metadata").(map[string]interface{}) {
		serviceMetadata = append(serviceMetadata, &gs.ServiceWmsMetadata{
			Key:   key,
			Value: value.(string),
		})
	}

	workspaceName := d.Get("workspace_name").(string)

	serviceConfiguration := &gs.ServiceWms{
		Name:              "WMS",
		IsEnabled:         d.Get("enabled").(bool),
		Title:             d.Get("title").(string),
		Maintainer:        d.Get("maintainer").(string),
		Abstract:          d.Get("abstract").(string),
		AccessConstraints: d.Get("access_constraints").(string),
		OnlineResource:    d.Get("online_resource").(string),
		IsVerbose:         d.Get("is_verbose").(bool),
		Watermark: gs.Watermark{
			IsEnabled:    d.Get("watermark_enabled").(bool),
			Position:     d.Get("watermark_position").(string),
			Transparency: d.Get("watermark_transparency").(int),
		},
		Interpolation:                           d.Get("interpolation").(string),
		IsCiteCompliant:                         d.Get("is_cite_compliant").(bool),
		MaximumBuffer:                           d.Get("maximum_buffer").(int),
		IsDynamicStylingDisabled:                d.Get("is_dynamic_styling_disabled").(bool),
		Metadata:                                serviceMetadata,
		Keywords:                                gs.ServiceWmsKeywords{Keywords: keywords},
		IsGetFeatureInfoMimeTypeCheckingEnabled: d.Get("is_getfeatureinfo_mimetype_checking_enabled").(bool),
		MaximumRequestMemory:                    d.Get("maximum_request_memory").(int),
		Fees:                                    d.Get("fees").(string),
		MaximumRenderingErrors:                  d.Get("maximum_rendering_errors").(int),
		MaximumRenderingTime:                    d.Get("maximum_rendering_time").(int),
		SupportedVersions:                       versions,
		SchemaBaseURL:                           d.Get("schema_base_url").(string),
		UseBBOXForEachCRS:                       d.Get("use_bbox_foreach_crs").(bool),
		IsGetMapMimeTypeCheckingEnabled:         d.Get("is_getmap_mimetype_checking_enabled").(bool),
		IsFeaturesReprojectionDisabled:          d.Get("is_features_reprojection_disabled").(bool),
		MaximumRequestedDimensionValues:         d.Get("maximum_requested_dimension_values").(int),
		CacheConfiguration: gs.CacheConfiguration{
			IsEnabled:    d.Get("is_cache_enabled").(bool),
			MaxEntrySize: d.Get("cache_maximum_entry_size").(int),
			MaxEntries:   d.Get("cache_maximum_entries").(int),
		},
		RemoteStyleMaxRequestTime:         d.Get("remote_style_max_request_time").(int),
		RemoteStyleTimeout:                d.Get("remote_style_timeout").(int),
		IsDefaultGroupStyleEnabled:        d.Get("is_default_group_style_enabled").(bool),
		IsTransformFeatureInfoDisabled:    d.Get("is_transform_feature_info_disabled").(bool),
		IsAutoEscapeTemplateValuesEnabled: d.Get("is_autoescape_templatevalues_enabled").(bool),
		RootLayerTitle:                    d.Get("root_layer_title").(string),
	}

	err := client.UpdateServiceWMS(workspaceName, serviceConfiguration)
	if err != nil {
		return err
	}

	return nil
}

func resourceGeoServerServiceWmsImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	var workspaceName string
	if strings.Contains(d.Id(), "/") {
		splittedID := strings.Split(d.Id(), "/")
		workspaceName = splittedID[1]
	}
	d.SetId(d.Id())
	d.Set("workspace_name", workspaceName)

	log.Printf("[INFO] Importing WMS Service configuration: %s", d.Id())

	err := resourceGeoServerServiceWmsRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
