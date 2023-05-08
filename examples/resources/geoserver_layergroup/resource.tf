# Example 1. Styles are in the same workspace and referenced by name
resource "geoserver_layergroup" "fdp_normal" {
  name                   = "fdp_normal"
  title                  = "fdp_normal"
  workspace_name = geoserver_workspace.my_workspace.name
  bounding_box_crs_class = "projected"
  bounding_box_crs_value = "EPSG:4326"
  bounding_box_min_x      = -63.16
  bounding_box_min_y      = -21.40
  bounding_box_max_x      = 55.84
  bounding_box_max_y      = 51.20

  layers {
    name  = geoserver_featuretype.bd_carto_zone_occupation_sol_vue_foret.name
    style = geoserver_style.zone_occupation_sol_vue_foret.name
  }

  layers {
    name  = geoserver_featuretype.bd_carto_zone_occupation_sol_vue_bati.name
    style = geoserver_style.zone_occupation_sol_vue_bati.name
  }

  layers {
    name  = geoserver_featuretype.bd_carto_zone_reglementee_touristique.name
    style = geoserver_style.zone_reglementee_touristique.name
  }

  layers {
    name  = geoserver_featuretype.bd_carto_zone_occupation_sol.name
    style = geoserver_style.zone_occupation_sol.name
  }

  layers {
    name  = geoserver_featuretype.bd_topo_zone_de_vegetation.name
    style = geoserver_style.zone_de_vegetation.name
  }

  layers {
    name  = geoserver_featuretype.bd_carto_zone_hydrographique_texture.name
    style = geoserver_style.zone_hydrographique_texture.name
  }

  layers {
    name  = geoserver_featuretype.bd_topo_zone_d_habitation.name
    style = geoserver_style.zone_d_habitation_2.name
  }

  layers {
    name  = geoserver_featuretype.bd_topo_zone_d_habitation.name
    style = geoserver_style.zone_d_habitation_1.name
  }


}

# Example 2. Styles are stored in the same workspace but we generate qualified names for the reference
resource "geoserver_layergroup" "osm_fdp_normal" {
  name                   = "fdp_normal"
  title                  = "fdp_normal"
  workspace_name = geoserver_workspace.my_workspace.name
  bounding_box_crs_class = "projected"
  bounding_box_crs_value = "EPSG:3857"
  bounding_box_max_x     = 20237886
  bounding_box_max_y     = 20237886
  bounding_box_min_x     = -20237886
  bounding_box_min_y     = -20237886

  layers {
    name = format("%s:%s",geoserver_workspace.osm.name,geoserver_featuretype.osm_simplified_water_polygons.name)
    style = format("%s:%s",geoserver_workspace.osm.name,geoserver_style.osm_simplified_water.name)
  }
   layers {
    name = format("%s:%s",geoserver_workspace.osm.name,geoserver_featuretype.osm_water_polygons.name)
    style = format("%s:%s",geoserver_workspace.osm.name,geoserver_style.osm_water.name)
  }
  layers {
    name = format("%s:%s",geoserver_workspace.osm.name,geoserver_featuretype.osm_land_polygons.name)
    style = format("%s:%s",geoserver_workspace.osm.name,geoserver_style.osm_coast_poly.name)
  }
}
