resource "geoserver_wms_store" "geoplateforme" {
  workspace_name = geoserver_workspace.ign.name
  name           = "geoplateforme"

  default = false
  enabled = true

  capabilities_url = "https://data.geopf.fr/wms-r/wms?SERVICE=WMS&amp;Version=1.3.0&amp;Request=GetCapabilities"

  max_connections = 10
  read_timeout = 20
  connection_timeout = 10
}
