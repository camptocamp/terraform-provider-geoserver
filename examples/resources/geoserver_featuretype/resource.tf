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
  attribute {
    name       = "usage_1"
    min_occurs = 0
    max_occurs = 1
    nillable   = true
    binding    = "java.lang.String"
  }
  attribute {
    name       = "usage_2"
    min_occurs = 0
    max_occurs = 1
    nillable   = true
    binding    = "java.lang.String"
  }
  attribute {
    name       = "construction_legere"
    min_occurs = 0
    max_occurs = 1
    nillable   = true
    binding    = "java.lang.Boolean"
  }
  attribute {
    name       = "etat_de_l_objet"
    min_occurs = 0
    max_occurs = 1
    nillable   = true
    binding    = "java.lang.String"
  }
  attribute {
    name       = "date_creation"
    min_occurs = 0
    max_occurs = 1
    nillable   = true
    binding    = "java.sql.Timestamp"
  }
  attribute {
    name       = "date_modification"
    min_occurs = 0
    max_occurs = 1
    nillable   = true
    binding    = "java.sql.Timestamp"
  }
  attribute {
    name       = "date_d_apparition"
    min_occurs = 0
    max_occurs = 1
    nillable   = true
    binding    = "java.sql.Date"
  }
  attribute {
    name       = "date_de_confirmation"
    min_occurs = 0
    max_occurs = 1
    nillable   = true
    binding    = "java.sql.Date"
  }
  attribute {
    name       = "sources"
    min_occurs = 0
    max_occurs = 1
    nillable   = true
    binding    = "java.lang.String"
  }
  attribute {
    name       = "identifiants_sources"
    min_occurs = 0
    max_occurs = 1
    nillable   = true
    binding    = "java.lang.String"
  }
  attribute {
    name       = "precision_planimetrique"
    min_occurs = 0
    max_occurs = 1
    nillable   = true
    binding    = "java.math.BigDecimal"
  }
  attribute {
    name       = "precision_altimetrique"
    min_occurs = 0
    max_occurs = 1
    nillable   = true
    binding    = "java.math.BigDecimal"
  }
  attribute {
    name       = "nombre_de_logements"
    min_occurs = 0
    max_occurs = 1
    nillable   = true
    binding    = "java.lang.Integer"
  }
  attribute {
    name       = "nombre_d_etages"
    min_occurs = 0
    max_occurs = 1
    nillable   = true
    binding    = "java.lang.Integer"
  }
  attribute {
    name       = "materiaux_des_murs"
    min_occurs = 0
    max_occurs = 1
    nillable   = true
    binding    = "java.lang.String"
  }
  attribute {
    name       = "materiaux_de_la_toiture"
    min_occurs = 0
    max_occurs = 1
    nillable   = true
    binding    = "java.lang.String"
  }
  attribute {
    name       = "hauteur"
    min_occurs = 0
    max_occurs = 1
    nillable   = true
    binding    = "java.math.BigDecimal"
  }
  attribute {
    name       = "altitude_minimale_sol"
    min_occurs = 0
    max_occurs = 1
    nillable   = true
    binding    = "java.math.BigDecimal"
  }
  attribute {
    name       = "altitude_minimale_toit"
    min_occurs = 0
    max_occurs = 1
    nillable   = true
    binding    = "java.math.BigDecimal"
  }
  attribute {
    name       = "altitude_maximale_toit"
    min_occurs = 0
    max_occurs = 1
    nillable   = true
    binding    = "java.math.BigDecimal"
  }
  attribute {
    name       = "altitude_maximale_sol"
    min_occurs = 0
    max_occurs = 1
    nillable   = true
    binding    = "java.math.BigDecimal"
  }
  attribute {
    name       = "origine_du_batiment"
    min_occurs = 0
    max_occurs = 1
    nillable   = true
    binding    = "java.lang.String"
  }
  attribute {
    name       = "appariement_fichiers_fonciers"
    min_occurs = 0
    max_occurs = 1
    nillable   = true
    binding    = "java.lang.String"
  }
}
