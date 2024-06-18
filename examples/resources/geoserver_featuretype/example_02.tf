resource "geoserver_featuretype" "bd_topo_batiment" {
  workspace_name = geoserver_workspace.my_workspace.name
  datastore_name = geoserver_datastore.db-referentiels-carto-bd_topo-ign.name
  name           = "fdp_batiment"
  native_name    = "batiment"
  enabled        = true

  lat_lon_bounding_box_max_x     = 55.83018112182617
  lat_lon_bounding_box_max_y     = 51.08795166015625
  lat_lon_bounding_box_min_x     = -63.152530670166016
  lat_lon_bounding_box_min_y     = -21.387481689453125
  lat_lon_bounding_box_crs_value = "EPSG:4326"

  native_bounding_box_max_x     = 55.83018112182617
  native_bounding_box_max_y     = 51.08795166015625
  native_bounding_box_min_x     = -63.152530670166016
  native_bounding_box_min_y     = -21.387481689453125
  native_bounding_box_crs_value = "EPSG:4326"

  projection_policy = "FORCE_DECLARED"

  srs = "EPSG:4326"

  attribute {
    name       = "geometrie"
    min_occurs = 0
    max_occurs = 1
    nillable   = true
    binding    = "org.locationtech.jts.geom.Geometry"
  }
  attribute {
    name       = "cleabs"
    min_occurs = 1
    max_occurs = 1
    nillable   = false
    binding    = "java.lang.String"
  }
  attribute {
    name       = "nature"
    min_occurs = 0
    max_occurs = 1
    nillable   = true
    binding    = "java.lang.String"
  }
}
