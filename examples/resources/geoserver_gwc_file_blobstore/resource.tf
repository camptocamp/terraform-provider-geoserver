resource "geoserver_gwc_file_blobstore" "file_blobstore" {
  blobstore_id           = "file_blobstore"
  base_directory         = "/mnt/cache/geowebcache"
  file_system_block_size = 4096
}
