resource "geoserver_datastore" "db-referentiels-carto-osm" {
  workspace_name = geoserver_workspace.my_workspace.name
  name           = "db-referentiels-carto-osm-jndi"

  default = false
  enabled = true

  connection_params = {
    "Estimated extends"                          = "true"
    "fetch size"                                 = "1000"
    "encode functions"                           = "true"
    "Expose primary keys"                        = "true"
    "Support on the fly geometry simplification" = "true"
    "Batch insert size"                          = "1"
    "preparedStatements"                         = "false"
    "Method used to simplify geometries"         = "FAST"
    "dbtype"                                     = "postgis"
    "Loose bbox"                                 = "true"
    "namespace"                                  = "nexsis"
    "jndiReferenceName"                          = "java:comp/env/jdbc/referentiels"
    "schema"                                     = "osm"
  }
}
