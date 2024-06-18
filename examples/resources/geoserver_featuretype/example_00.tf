resource "geoserver_featuretype" "bd_topo_canalisation" {
  workspace_name = geoserver_workspace.fdp.name
  datastore_name = geoserver_datastore.db-referentiels-carto-bd_topo-ign.name
  name           = "fdp_canalisation"
  native_name    = "canalisation"
  enabled        = true

  lat_lon_bounding_box_max_x     = 55.547359466552734
  lat_lon_bounding_box_max_y     = 51.05213928222656
  lat_lon_bounding_box_min_x     = -61.76982879638672
  lat_lon_bounding_box_min_y     = -21.289060592651367
  lat_lon_bounding_box_crs_value = "EPSG:4326"

  native_bounding_box_max_x     = 55.547359466552734
  native_bounding_box_max_y     = 51.05213928222656
  native_bounding_box_min_x     = -61.76982879638672
  native_bounding_box_min_y     = -21.289060592651367
  native_bounding_box_crs_value = "EPSG:4326"

  projection_policy = "FORCE_DECLARED"

  srs = "EPSG:4326"
}
