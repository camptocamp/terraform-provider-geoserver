---
page_title: "Provider: GeoServer - How to force feature type"
description: |-
  A solution to trigger a feature type when working without explicit attributes
---
# How to force feature type

## The problem

When you declare featuretype resource into your terraform module and you don't explicitely declare the attributes available in your feature type, the GeoServer resource won't be updated if the underlying data structure is changed.

## The solution

The solution to manage is the following:

- Declare a input variable to be able to provide some version information into terraform
- Use a `terraform_data` resource to track into the terraform state how the version evolves over time
- Add a trigger on the related resources to update them when the version changes

Let's see an example on how to configure this.

## Input variable declaration

Just declare a standard terraform variable to inject into your terraform manifest a version information.

```terraform
variable "database_version" {
  type = string
}
```

You can provide the value of the variable like you want: default value change, using prompt, variable injection...

## Track the version changes into terraform

We need now to be able to track when the version changes in order to trigger things into terraform.

For that we will use a `terraform_data` resource (available since `v1.4`), which can store values which need to be followed.

```terraform
resource "terraform_data" "internal_version" {
  input = var.database_version
}
```

## Trigger other resource update on version change

The last part of the solution is to defined on each resource you want to update upon database changes a replace trigger.

It can be done like this:

```terraform
resource "geoserver_featuretype" "borehole" {
  workspace_name = geoserver_workspace.my_workspace.name
  datastore_name = geoserver_datastore.my_store.name
  name           = "borehole"
  native_name    = "borehole"
  enabled        = true

  lat_lon_bounding_box_max_x     = 180
  lat_lon_bounding_box_max_y     = 90
  lat_lon_bounding_box_min_x     = -180
  lat_lon_bounding_box_min_y     = -90
  lat_lon_bounding_box_crs_value = "EPSG:4326"

  native_bounding_box_max_x     = 180
  native_bounding_box_max_y     = 90
  native_bounding_box_min_x     = -180
  native_bounding_box_min_y     = -90
  native_bounding_box_crs_class = ""
  native_bounding_box_crs_value = "EPSG:4326"

  projection_policy = "FORCE_DECLARED"

  srs = "EPSG:4326"

  # The important part
  lifecycle {
    replace_triggered_by = [terraform_data.internal_version]
  }

}
```
