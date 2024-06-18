resource "geoserver_datastore" "db-referentiels-carto-osm" {
  workspace_name = geoserver_workspace.my_workspace.name
  name           = "db-referentiels-carto-osm"

  default = false
  enabled = true

  connection_params = {
    "Batch insert size"                          = "1"
    "Connection timeout"                         = "5"
    "Estimated extends"                          = "true"
    "Evictor run periodicity"                    = "300"
    "Evictor tests per run"                      = "3"
    "Expose primary keys"                        = "true"
    "Loose bbox"                                 = "true"
    "Max connection idle time"                   = "300"
    "Max open prepared statements"               = "50"
    "Method used to simplify geometries"         = "FAST"
    "SSL mode"                                   = "DISABLE"
    "Support on the fly geometry simplification" = "false"
    "Test while idle"                            = "true"
    "create database"                            = "false"
    "database"                                   = var.referentiels_carto_db_config.DATABASE_NAME
    "dbtype"                                     = "postgis"
    "encode functions"                           = "true"
    "fetch size"                                 = "2000"
    "host"                                       = var.referentiels_carto_db_config.HOST
    "max connections"                            = "15"
    "min connections"                            = "1"
    "namespace"                                  = "nexsis"
    "passwd"                                     = var.referentiels_carto_db_config.PASSWORD
    "port"                                       = "5432"
    "preparedStatements"                         = "false"
    "schema"                                     = "osm"
    "user"                                       = var.referentiels_carto_db_config.ROLE
    "validate connections"                       = "true"
  }
}