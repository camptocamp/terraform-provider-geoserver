# Example using a well known WMS service
resource "geoserver_gwc_wms_layer" "ign_orthophotos" {
    name = "ign_orthophotos"
    blobstore_id = geoserver_gwc_S3_blobstore.s3_blobstore.blobstore_id
    wms_url="https://wxs.ign.fr/essentiels/geoportail/r/wms?SERVICE=WMS"
    wms_layer="ORTHOIMAGERY.ORTHOPHOTOS"
    mime_formats=["image/png","image/jpeg"]
    grid_subsets=[geoserver_gwc_gridset.gridset_3857.name]
    metatile_height=4
    metatile_width=4
    expire_duration_clients=3600
}

# Example using a layer provided with a GeoServer instance
# The geoserver URL and the layer name are provided with variables
resource "geoserver_gwc_wms_layer" "nexsis_fdp_osm" {
    name = "nexsis_fdp_osm"
    blobstore_id = geoserver_gwc_S3_blobstore.s3_blobstore.blobstore_id
    wms_url=var.geoserver_wms
    wms_layer=var.fdp_osm_layer
    mime_formats=["image/png","image/jpeg"]
    grid_subsets=[geoserver_gwc_gridset.gridset_3857.name]
    metatile_height=4
    metatile_width=4
    expire_duration_clients=3600
}