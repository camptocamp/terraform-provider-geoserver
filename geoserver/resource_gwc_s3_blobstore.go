package geoserver

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	gs "github.com/camptocamp/go-geoserver/client"
)

func resourceGwcS3Blobstore() *schema.Resource {
	return &schema.Resource{
		Create: resourceGwcS3BlobstoreCreate,
		Read:   resourceGwcS3BlobstoreRead,
		Update: resourceGwcS3BlobstoreUpdate,
		Delete: resourceGwcS3BlobstoreDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGwcS3BlobstoreImport,
		},

		Schema: map[string]*schema.Schema{
			"blobstore_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identifier of the blobstore. Use as resource id.",
			},
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Bucket used by the blobstore.",
			},
			"bucket_access_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Access key to the bucket.",
			},
			"bucket_secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Secret key associated to the access key.",
			},
			"prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Prefix to add to the path for storing the tiles. Default is empty",
			},
			"access_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "PUBLIC",
				Description: "Access type to the bucket. Could be PUBLIC or PRIVATE. Default is PUBLIC",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "End-point to use to access the bucket, when not using AWS S3",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Is the blobstore enabled? Default to true.",
			},
			"default": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Is the blobstore the default one? Default to false.",
			},
			"use_https": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Use https to contact the S3 bucket. Default to false.",
			},
			"use_gzip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Compress tiles with gzip. Default to false.",
			},
			"max_connections": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     50,
				Description: "Maximum number of parallel connections. Default to 50.",
			},
		},
	}
}

func resourceGwcS3BlobstoreCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Creating GWC S3 Blobstore: %s", d.Id())

	client := meta.(*Config).GwcClient()

	blobstoreName := d.Get("blobstore_id").(string)
	blobstore := &gs.BlobstoreS3{
		Id:             blobstoreName,
		Bucket:         d.Get("bucket").(string),
		Prefix:         d.Get("prefix").(string),
		AwsAccessKey:   d.Get("bucket_access_key").(string),
		AwsSecretKey:   d.Get("bucket_secret_key").(string),
		Access:         d.Get("access_type").(string),
		Endpoint:       d.Get("endpoint").(string),
		MaxConnections: d.Get("max_connections").(int),
		UseHTTPS:       d.Get("use_https").(bool),
		UseGzip:        d.Get("use_gzip").(bool),
		Enabled:        d.Get("enabled").(bool),
		Default:        d.Get("default").(bool),
	}

	err := client.CreateBlobstoreS3(blobstoreName, blobstore)
	if err != nil {
		return err
	}

	d.SetId(blobstoreName)

	return resourceGwcS3BlobstoreRead(d, meta)
}

func resourceGwcS3BlobstoreRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing GWC S3 Blobstore: %s", d.Id())

	blobstoreID := d.Id()

	client := meta.(*Config).GwcClient()

	blobstore, err := client.GetBlobstoreS3(blobstoreID)
	if err != nil && !strings.Contains(strings.ToLower(err.Error()), "not found") {
		return err
	}

	if blobstore == nil {
		d.SetId("")
		return nil
	}

	d.Set("blobstore_id", blobstore.Id)
	d.Set("bucket", blobstore.Bucket)
	d.Set("prefix", blobstore.Prefix)
	d.Set("bucket_access_key", blobstore.AwsAccessKey)
	d.Set("bucket_secret_key", blobstore.AwsSecretKey)
	d.Set("access_type", blobstore.Access)
	d.Set("endpoint", blobstore.Endpoint)
	d.Set("max_connections", blobstore.MaxConnections)
	d.Set("use_https", blobstore.UseHTTPS)
	d.Set("use_gzip", blobstore.UseGzip)
	d.Set("enabled", blobstore.Enabled)
	d.Set("default", blobstore.Default)

	return nil
}

func resourceGwcS3BlobstoreDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting GWC S3 Blobstore: %s", d.Id())

	blobstoreID := d.Id()

	client := meta.(*Config).GwcClient()

	err := client.DeleteBlobstoreS3(blobstoreID)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceGwcS3BlobstoreUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating GWC S3 Blobstore: %s", d.Id())

	blobstoreID := d.Id()

	client := meta.(*Config).GwcClient()

	err := client.UpdateBlobstoreS3(blobstoreID, &gs.BlobstoreS3{
		Id:             blobstoreID,
		Bucket:         d.Get("bucket").(string),
		Prefix:         d.Get("prefix").(string),
		AwsAccessKey:   d.Get("bucket_access_key").(string),
		AwsSecretKey:   d.Get("bucket_secret_key").(string),
		Access:         d.Get("access_type").(string),
		Endpoint:       d.Get("endpoint").(string),
		MaxConnections: d.Get("max_connections").(int),
		UseHTTPS:       d.Get("use_https").(bool),
		UseGzip:        d.Get("use_gzip").(bool),
		Enabled:        d.Get("enabled").(bool),
		Default:        d.Get("default").(bool),
	},
	)

	if err != nil {
		return err
	}

	return nil
}

func resourceGwcS3BlobstoreImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	blobstoreID := d.Id()

	d.SetId(d.Id())
	d.Set("blobstore_id", blobstoreID)

	log.Printf("[INFO] Importing GWC S3 Blobstore `%s`", blobstoreID)

	err := resourceGwcS3BlobstoreRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
