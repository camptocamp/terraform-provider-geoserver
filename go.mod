module github.com/camptocamp/terraform-provider-geoserver

replace github.com/camptocamp/go-geoserver => /home/agacon/camptocamp/go-geoserver

go 1.15

require (
	github.com/camptocamp/go-geoserver v0.0.0-20240205195341-221dfa93f919
	github.com/hashicorp/terraform-plugin-sdk v1.17.2
)
