resource "geoserver_gwc_gs_layer" "basemap" {
    name = "my_workspace::my_basemap"
    id = "LayerGroupInfo-01KW94WAF9DH2DEHSVN4QCE7ZT"
    blobstore_id = geoserver_gwc_S3_blobstore.s3_blobstore.blobstore_id
    mime_formats=["image/png","image/jpeg"]
    grid_subsets=[geoserver_gwc_gridset.gridset_3857.name]
    metatile_height=4
    metatile_width=4
    expire_duration_clients=3600
}
