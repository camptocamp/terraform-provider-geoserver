package geoserver

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	gs "github.com/camptocamp/go-geoserver/client"
)

func resourceGeoserverUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceGeoserverUserCreate,
		Read:   resourceGeoserverUserRead,
		Update: resourceGeoserverUserUpdate,
		Delete: resourceGeoserverUserDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGeoserverUserImport,
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the usergroup service. If empty the default geoserver one will be used. Used to compute the id of the resource.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the user. Used to compute the id of the resource.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Is the user account enabled?",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Password associated to the account.",
			},
		},
	}
}

func resourceGeoserverUserCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Creating Geoserver User: %s", d.Id())

	client := meta.(*Config).GeoserverClient()

	serviceName := d.Get("service_name").(string)

	user := &gs.User{
		Name:     d.Get("name").(string),
		Enabled:  d.Get("enabled").(bool),
		Password: d.Get("password").(string),
	}

	err := client.CreateUser(serviceName, user)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", serviceName, d.Get("name").(string)))

	return resourceGeoserverUserRead(d, meta)
}

func resourceGeoserverUserRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Geoserver User: %s", d.Id())

	splittedID := strings.Split(d.Id(), "/")
	serviceName := splittedID[0]
	userName := splittedID[1]

	client := meta.(*Config).GeoserverClient()

	user, err := client.GetUser(serviceName, userName)
	if err != nil && !strings.Contains(strings.ToLower(err.Error()), "not found") {
		return err
	}

	if user == nil {
		d.SetId("")
		return nil
	}
	d.Set("service_name", serviceName)
	d.Set("name", user.Name)
	d.Set("enabled", user.Enabled)
	d.Set("password", user.Password)

	return nil
}

func resourceGeoserverUserDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Geoserver User: %s", d.Id())

	splittedID := strings.Split(d.Id(), "/")
	serviceName := splittedID[0]
	userName := splittedID[1]

	client := meta.(*Config).GeoserverClient()

	err := client.DeleteUser(serviceName, userName)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceGeoserverUserUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Geoserver User: %s", d.Id())

	splittedID := strings.Split(d.Id(), "/")
	serviceName := splittedID[0]

	client := meta.(*Config).GeoserverClient()

	user := &gs.User{
		Name:     d.Get("name").(string),
		Enabled:  d.Get("enabled").(bool),
		Password: d.Get("password").(string),
	}

	errUpdateUser := client.UpdateUser(serviceName, user.Name, user)
	if errUpdateUser != nil {
		return errUpdateUser
	}

	return nil
}

func resourceGeoserverUserImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splittedID := strings.Split(d.Id(), "/")
	serviceName := splittedID[0]
	userName := splittedID[1]

	d.SetId(d.Id())
	d.Set("name", userName)
	d.Set("service_name", serviceName)

	log.Printf("[INFO] Importing Geoserver User `%s` in service `%s`", userName, serviceName)

	err := resourceGeoserverStyleRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
