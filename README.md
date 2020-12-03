Geoserver Terraform Provider
============================

[![Terraform Registry Version](https://img.shields.io/badge/dynamic/json?color=blue&label=registry&query=%24.version&url=https%3A%2F%2Fregistry.terraform.io%2Fv1%2Fproviders%2Fcamptocamp%geoserver)](https://registry.terraform.io/providers/camptocamp/geoserver)
[![Go Report Card](https://goreportcard.com/badge/github.com/camptocamp/terraform-provider-geoserver)](https://goreportcard.com/report/github.com/camptocamp/terraform-provider-geoserver)
[![By Camptocamp](https://img.shields.io/badge/by-camptocamp-fb7047.svg)](http://www.camptocamp.com)

This provider adds integration between Terraform and Geoserver.

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.10


Building The Provider
---------------------

Download the provider source code

```sh
$ go get github.com/camptocamp/terraform-provider-geoserver
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/camptocamp/terraform-provider-geoserver
$ make build
```

Installing the provider
-----------------------

After building the provider, install it using the Terraform instructions for [installing a third party provider](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins).

Example
----------------------

```hcl
provider geoserver {
  url = "https://geoserver.example.tld"   # or set $GEOSERVER_URL
  username = "admin"                      # or set $GEOSERVER_USERNAME
  password = "password"                   # or set $GEOSERVER_PASSWORD
  insecure = true
}

resource "geoserver_workspace" "foo" {
  name = "foo"
}

resource "geoserver_datastore" "default" {
  workspace_name = geoserver_workspace.foo.name
  name           = "default"
  type           = "postgis"
  host           = "pgmaster"
  port           = "5432"
  db_name        = "test"
  db_user        = "postgres"
  db_pass        = "postgres"
}
```
