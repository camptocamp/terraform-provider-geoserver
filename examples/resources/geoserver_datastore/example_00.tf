resource "geoserver_datastore" "db-referentiels-carto-osm-pregen" {
  workspace_name = geoserver_workspace.my_workspace.name
  name           = "db-referentiels-carto-osm-pregen"

  default = false
  enabled = true

  connection_params = {
    "RepositoryClassName"                  = "org.geoserver.data.gen.DSFinderRepository"
    "GeneralizationInfosProviderParam"     = "pregeneralized.xml"
    "GeneralizationInfosProviderClassName" = "org.geoserver.data.gen.info.GeneralizationInfosProviderImpl"
    "namespace"                            = "http://geoserver.org/osm"
  }
}
