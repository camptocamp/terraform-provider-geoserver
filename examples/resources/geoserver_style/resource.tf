resource "geoserver_style" "osm_all_regions" {
  name             = "all_regions"
  workspace_name   = geoserver_workspace.my_workspace.name
  filename         = "all_regions.css"
  style_definition = file("${path.module}/styles/all_regions.css")
  format           = "css"
  version          = "1.0.0"
}
