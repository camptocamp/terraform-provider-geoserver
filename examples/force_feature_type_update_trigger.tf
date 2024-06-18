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
