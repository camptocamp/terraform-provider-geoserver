# Create a pregeneralized.xml resource at the root of the GeoServer directory
# wih the content of the provided file 
resource "geoserver_resource" "osm_pregen_cfg_file" {
  path = "pregeneralized"
  extension = "xml"
  resource = file("${path.module}/pregeneralized.xml")
}

