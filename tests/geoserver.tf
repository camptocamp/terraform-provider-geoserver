provider geoserver {
  url = "http://localhost:8080/geoserver"
  username = "admin"
  password = "geoserver"
  insecure = true
}

resource "geoserver_workspace" "foo" {
  name = "foo"
  namespace = "http://foo.info"

}
