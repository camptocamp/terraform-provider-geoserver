---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "geoserver Provider"
subcategory: ""
description: |-
  
---

# geoserver Provider



## Example Usage

```terraform
provider "geoserver" {
    url = "http://localhost:9090/geoserver/cloud/rest"
    gwc_url = "http://localhost:9090/geoserver/cloud/gwc/rest"
    username = "admin"
    password = "geoserver"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `gwc_url` (String) The GeoWebCache URL
- `insecure` (Boolean) Whether to verify the server's SSL certificate
- `password` (String) Password to use for connection
- `url` (String) The Geoserver URL
- `username` (String) Username to use for connection
