---
page_title: "Provider: GeoServer - How to generate feature types based on database tables"
description: |-
  This page explains how by using other TF providers we can have a GeoServer instance with
  discovery mechanism similar to pg_featserver or pg_tileserv.
---
# Feature type automatic declaration

## The problem

Several GIS open source projects like pg_featserv or pg_tileserv offers a nice automatic configuration mechanism from the tables found in a database. For some use case (like simple publishing of data with little configuration), it would be great to have something similar with GeoServer.

## The solution

To have a similar automatic configuration mechanism, we could use an other TF provider to retrieve the tables available in a database and use them to generate `geoserver_featuretype` resources.

For this example, we will use the `postgresql` provider available at `cyrilgdn/postgresql`.

We will need to do the following things:

- Register the provider
- Configure the provider
- Import the definifition of the database into TF
- Retrieve the data source for the tables
- Generate the feature types

### Register the provider

In the `versions.tf` file, add the provider definition with the geoserver provider definition:

```terraform
terraform {
  required_providers {
    geoserver = {
      source  = "camptocamp/geoserver"
      version = "0.0.21"
    }
    postgresql = {
      source = "cyrilgdn/postgresql"
      version = "1.25.0"
    }    
  }
}
```

### Configure the database access

The next step is to configure the provider in order to be able to connect to the database:

```terraform
provider "postgresql" {
  host            = "172.17.0.1"
  port            = 54334
  database        = "sig"
  username        = "sig"
  password        = "sig"
  sslmode         = "disable"
  connect_timeout = 15
}
```

Please refer to the provider documentation for the different configuration options.

### Import the database definition

The PG provider can be used to create a database from scratch. In our case we only want to read an existing database so we will import the resource definition into our TF context.

```terraform
# We declare the database we want to work with
resource "postgresql_database" "sigdb" {
  name = "sig"
}

# And we import the definition since we don't want to create it
import {
  to = postgresql_database.sigdb
  id = "sig"
}
```

### Get the table list

To get the table list, we only have to declare a TF data source of type `postgresql_tables` like this:

```terraform
data "postgresql_tables" "my_tables" {
  database = "sig"
  schemas = ["nexsis"] # we only want the tables from this schema
}
```

### Feature type template

The last step is to defined the template to be used for the configuration of our feature types, by using the `for_each` meta-argument.

```terraform
resource "geoserver_featuretype" "nexsis_feature_types" {
  for_each = toset([for tabledesc in data.postgresql_tables.my_tables.tables: tabledesc.object_name])

  workspace_name = "nexsis"
  datastore_name = "db-sig"
  name           = each.key
  native_name    = each.key
  enabled        = true

  lat_lon_bounding_box_max_x     = 180
  lat_lon_bounding_box_max_y     = 90
  lat_lon_bounding_box_min_x     = -180
  lat_lon_bounding_box_min_y     = -90
  lat_lon_bounding_box_crs_value = "EPSG:4326"

  native_bounding_box_max_x     = 180
  native_bounding_box_max_y     = 90
  native_bounding_box_min_x     = -180
  native_bounding_box_min_y     = -90
  native_bounding_box_crs_value = "EPSG:4326"

  projection_policy = "FORCE_DECLARED"

  srs = "EPSG:4326"
}
```

## Conclusion

With the described configuration, we are able to automatically maintain in GeoServer feature types for the tables defined in a database.

The mechanism has the following limits:

- we cannot automatically retrieve only the tables with at least a geometry column
- we cannot read the SRS from the geometry column
- we cannot detect schema changes on existing feature types