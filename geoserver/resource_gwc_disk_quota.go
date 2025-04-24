package geoserver

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	gs "github.com/camptocamp/go-geoserver/client"
)

func resourceGwcDiskQuota() *schema.Resource {
	return &schema.Resource{
		Description: "Manage GeoWebCache Disk Quota Policy. The policy is a singleton so Create is similar to Update and Delete has no effect. There is a bug on the GWC REST API which prevents reading the existing configuration",
		Create:      resourceGwcDiskQuotaCreate,
		Read:        resourceGwcDiskQuotaRead,
		Update:      resourceGwcDiskQuotaUpdate,
		Delete:      resourceGwcDiskQuotaDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGwcDiskQuotaImport,
		},

		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Is the disk quota mechanism enabled? Default is false.",
			},
			"cache_cleanup_frequency": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Duration to wait between each cleanup execution. Unit depends on cache_cleanup_units",
			},
			"cache_cleanup_units": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Time unit used to express cleanup frequency. Authorized values are : SECONDS, MINUTES, HOURS, DAYS",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					allowed_values := []string{"SECONDS", "MINUTES", "HOURS", "DAYS"}
					if !slices.Contains(allowed_values, v) {
						errs = append(errs, fmt.Errorf("%q must be one of this values %q, got: %q", key, strings.Join(allowed_values, ","), v))
					}
					return
				},
			},
			"maximum_concurrent_cleanup": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The amount of threads to use when processing the disk quota",
			},
			"global_expiration_policy_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Policy to use when tiles have to be removed. Authorized values are : LRU, LFU",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					allowed_values := []string{"LRU", "LFU"}
					if !slices.Contains(allowed_values, v) {
						errs = append(errs, fmt.Errorf("%q must be one of this values %q, got: %q", key, strings.Join(allowed_values, ","), v))
					}
					return
				},
			},
			"global_quota_value": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Numeric value for global quota.",
			},
			"global_quota_units": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of unit quantified by value. Authorized values are : B, KiB, MiB, GiB, TiB",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					allowed_values := []string{"B", "KiB", "MiB", "GiB", "TiB"}
					if !slices.Contains(allowed_values, v) {
						errs = append(errs, fmt.Errorf("%q must be one of this values %q, got: %q", key, strings.Join(allowed_values, ","), v))
					}
					return
				},
			},
			"layer_quota": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The disk quota to apply for a given layer.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"layer": {
							Type:     schema.TypeString,
							Required: true,
						},
						"expiration_policy_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Policy to use when tiles have to be removed. Authorized values are : LRU, LFU",
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(string)
								allowed_values := []string{"LRU", "LFU"}
								if !slices.Contains(allowed_values, v) {
									errs = append(errs, fmt.Errorf("%q must be one of this values %q, got: %q", key, strings.Join(allowed_values, ","), v))
								}
								return
							},
						},
						"quota_value": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Numeric value for layer quota.",
						},
						"quota_units": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The type of unit quantified by value. Authorized values are : B, KiB, MiB, GiB, TiB",
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(string)
								allowed_values := []string{"B", "KiB", "MiB", "GiB", "TiB"}
								if !slices.Contains(allowed_values, v) {
									errs = append(errs, fmt.Errorf("%q must be one of this values %q, got: %q", key, strings.Join(allowed_values, ","), v))
								}
								return
							},
						},
					},
				},
			},
		},
	}
}

func resourceGwcDiskQuotaCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Configuring GWC Disk Quota")

	client := meta.(*Config).GwcClient()

	var layerDiskQuotas []*gs.GwcLayerQuota
	for _, value := range d.Get("layer_quota").(*schema.Set).List() {
		v := value.(map[string]interface{})
		layerDiskQuotas = append(layerDiskQuotas, &gs.GwcLayerQuota{
			Layer:                v["layer"].(string),
			ExpirationPolicyName: v["expiration_policy_name"].(string),
			Quota: gs.GwcQuota{
				Value: v["quota_value"].(int),
				Units: v["quota_units"].(string),
			},
		})
	}

	diskQuota := &gs.GwcQuotaConfiguration{
		Enabled:                    d.Get("enabled").(bool),
		CacheCleanUpFrequency:      d.Get("cache_cleanup_frequency").(int),
		CacheCleanUpUnits:          d.Get("cache_cleanup_units").(string),
		MaxConcurrentCleanUps:      d.Get("maximum_concurrent_cleanup").(int),
		GlobalExpirationPolicyName: d.Get("global_expiration_policy_name").(string),
		GlobalQuota: gs.GwcQuota{
			Value: d.Get("global_quota_value").(int),
			Units: d.Get("global_quota_units").(string),
		},
		LayersQuotas: layerDiskQuotas,
	}

	err := client.UpdateGwcQuotaConfiguration(diskQuota)
	if err != nil {
		return err
	}

	d.SetId("gwc_disk_quota_singleton")

	return resourceGwcDiskQuotaRead(d, meta)
}

func resourceGwcDiskQuotaRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing GWC Disk Quota")

	client := meta.(*Config).GwcClient()

	diskQuota, err := client.GetGwcQuotaConfiguration()
	if err != nil && !strings.Contains(strings.ToLower(err.Error()), "not found") {
		return err
	}

	if diskQuota == nil {
		d.SetId("")
		return nil
	}

	d.Set("enabled", diskQuota.Enabled)
	d.Set("cache_cleanup_frequency", diskQuota.CacheCleanUpFrequency)
	d.Set("cache_cleanup_units", diskQuota.CacheCleanUpUnits)
	d.Set("maximum_concurrent_cleanup", diskQuota.MaxConcurrentCleanUps)
	d.Set("global_quota_value", diskQuota.GlobalQuota.Value)
	d.Set("global_quota_units", diskQuota.GlobalQuota.Units)
	d.Set("global_expiration_policy_name", diskQuota.GlobalExpirationPolicyName)

	var layerQuotas []map[string]interface{}
	for _, value := range diskQuota.LayersQuotas {
		layerQuotas = append(layerQuotas, map[string]interface{}{
			"layer":                  value.Layer,
			"expiration_policy_name": value.ExpirationPolicyName,
			"quota_value":            value.Quota.Value,
			"quota_units":            value.Quota.Units,
		})
	}
	d.Set("layer_quota", layerQuotas)

	return nil
}

func resourceGwcDiskQuotaDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Disk Quota cannot be deleted. Skipping action")

	d.SetId("")

	return nil
}

func resourceGwcDiskQuotaUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] GWC Disk Quota Configuration")

	client := meta.(*Config).GwcClient()

	var layerDiskQuotas []*gs.GwcLayerQuota
	for _, value := range d.Get("layer_quota").(*schema.Set).List() {
		v := value.(map[string]interface{})
		layerDiskQuotas = append(layerDiskQuotas, &gs.GwcLayerQuota{
			Layer:                v["layer"].(string),
			ExpirationPolicyName: v["expiration_policy_name"].(string),
			Quota: gs.GwcQuota{
				Value: v["quota_value"].(int),
				Units: v["quota_units"].(string),
			},
		})
	}

	diskQuota := &gs.GwcQuotaConfiguration{
		Enabled:                    d.Get("enabled").(bool),
		CacheCleanUpFrequency:      d.Get("cache_cleanup_frequency").(int),
		CacheCleanUpUnits:          d.Get("cache_cleanup_units").(string),
		MaxConcurrentCleanUps:      d.Get("maximum_concurrent_cleanup").(int),
		GlobalExpirationPolicyName: d.Get("global_expiration_policy_name").(string),
		GlobalQuota: gs.GwcQuota{
			Value: d.Get("global_quota_value").(int),
			Units: d.Get("global_quota_units").(string),
		},
		LayersQuotas: layerDiskQuotas,
	}

	err := client.UpdateGwcQuotaConfiguration(diskQuota)
	if err != nil {
		return err
	}

	return nil
}

func resourceGwcDiskQuotaImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	d.SetId("gwc_disk_quota_singleton")

	log.Printf("[INFO] Importing GWC Disk Quota Configuration")

	err := resourceGwcDiskQuotaRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
