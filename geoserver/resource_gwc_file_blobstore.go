package geoserver

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	gs "github.com/camptocamp/go-geoserver/client"
)

func resourceGwcFileBlobstore() *schema.Resource {
	return &schema.Resource{
		Create: resourceGwcFileBlobstoreCreate,
		Read:   resourceGwcFileBlobstoreRead,
		Update: resourceGwcFileBlobstoreUpdate,
		Delete: resourceGwcFileBlobstoreDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGwcFileBlobstoreImport,
		},

		Schema: map[string]*schema.Schema{
			"blobstore_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identifier of the blobstore. Use as resource id.",
			},
			"base_directory": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Base directory to store the tiles.",
			},
			"file_system_block_size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Block size of files.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Is the blobstore enabled? Default to true.",
			},
		},
	}
}

func resourceGwcFileBlobstoreCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Creating GWC File Blobstore: %s", d.Id())

	client := meta.(*Config).GwcClient()

	blobstoreName := d.Get("blobstore_id").(string)
	blobstore := &gs.BlobstoreFile{
		Id:                  blobstoreName,
		Enabled:             d.Get("enabled").(bool),
		BaseDirectory:       d.Get("base_directory").(string),
		FileSystemBlockSize: d.Get("file_system_block_size").(int),
	}

	err := client.CreateBlobstoreFile(blobstoreName, blobstore)
	if err != nil {
		return err
	}

	d.SetId(blobstoreName)

	return resourceGwcFileBlobstoreRead(d, meta)
}

func resourceGwcFileBlobstoreRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing GWC File Blobstore: %s", d.Id())

	blobstoreID := d.Id()

	client := meta.(*Config).GwcClient()

	blobstore, err := client.GetBlobstoreFile(blobstoreID)
	if err != nil && !strings.Contains(strings.ToLower(err.Error()), "not found") {
		return err
	}

	if blobstore == nil {
		d.SetId("")
		return nil
	}

	d.Set("blobstore_id", blobstore.Id)
	d.Set("base_directory", blobstore.BaseDirectory)
	d.Set("file_system_block_size", blobstore.FileSystemBlockSize)
	d.Set("enabled", blobstore.Enabled)

	return nil
}

func resourceGwcFileBlobstoreDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting GWC File Blobstore: %s", d.Id())

	blobstoreID := d.Id()

	client := meta.(*Config).GwcClient()

	err := client.DeleteBlobstoreFile(blobstoreID)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceGwcFileBlobstoreUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating GWC File Blobstore: %s", d.Id())

	blobstoreID := d.Id()

	client := meta.(*Config).GwcClient()

	err := client.UpdateBlobstoreFile(blobstoreID, &gs.BlobstoreFile{
		Id:                  blobstoreID,
		Enabled:             d.Get("enabled").(bool),
		BaseDirectory:       d.Get("base_directory").(string),
		FileSystemBlockSize: d.Get("file_system_block_size").(int),
	},
	)

	if err != nil {
		return err
	}

	return nil
}

func resourceGwcFileBlobstoreImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	blobstoreID := d.Id()

	d.SetId(d.Id())
	d.Set("blobstore_id", blobstoreID)

	log.Printf("[INFO] Importing GWC File Blobstore `%s`", blobstoreID)

	err := resourceGwcFileBlobstoreRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
