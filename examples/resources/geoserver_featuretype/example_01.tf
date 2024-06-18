resource "geoserver_featuretype" "river" {
  workspace_name = geoserver_workspace.my_workspace.name
  datastore_name = geoserver_datastore.my_datastore.name
  name           = "river"
  native_name    = "river"
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


  metadata = {
    "JDBC_VIRTUAL_TABLE" = "\n      \u003cvirtualTable\u003e\n        \u003cname\u003eriver\u003c/name\u003e\n        \u003csql\u003eselect\u0026#xd;\nfo_water_withdrawal.bu_code as bu_code,\u0026#xd;\nasset_gid as asset_gid,\u0026#xd;\nasset_gid as asset_gid_pk,\u0026#xd;\noperational_contract_gid as contract_gid,\u0026#xd;\noperational_contract_id as contract_id,\u0026#xd;\nasset_name as installation_name,\u0026#xd;\nwater_withdrawal_type_level_2_name as asset_subtype_name,\u0026#xd;\nfunctional_location_geometry as asset_geometry\u0026#xd;\nfrom public.fo_water_withdrawal\u0026#xd;\nLEFT JOIN fo_water_resource\u0026#xd;\nON fo_water_withdrawal.resource_gid = fo_water_resource.resource_gid\u0026#xd;\nwhere  fo_water_withdrawal.water_withdrawal_type_level_1_code = \u0026apos;SurfaceWaterPoint\u0026apos;\u0026#xd;\nAND fo_water_resource.resource_type_code = \u0026apos;WaterCourseSection\u0026apos;\n\u003c/sql\u003e\n        \u003cescapeSql\u003efalse\u003c/escapeSql\u003e\n        \u003ckeyColumn\u003easset_gid_pk\u003c/keyColumn\u003e\n        \u003cgeometry\u003e\n          \u003cname\u003easset_geometry\u003c/name\u003e\n          \u003ctype\u003ePoint\u003c/type\u003e\n          \u003csrid\u003e4326\u003c/srid\u003e\n        \u003c/geometry\u003e\n      \u003c/virtualTable\u003e\n    "
  }
}
