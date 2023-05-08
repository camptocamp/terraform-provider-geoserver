# All S3 elements are provided through K8s secrets
resource "geoserver_gwc_S3_blobstore" "s3_blobstore" {
  blobstore_id="s3_blobstore"
  bucket=data.kubernetes_secret.s3_credentials.data.AWS_BUCKET
  bucket_access_key=data.kubernetes_secret.s3_credentials.data.AWS_ACCESS_KEY_ID
  bucket_secret_key=data.kubernetes_secret.s3_credentials.data.AWS_SECRET_ACCESS_KEY
  prefix=var.env_name
  endpoint=data.kubernetes_secret.s3_credentials.data.AWS_ENDPOINT
  use_https=true
}