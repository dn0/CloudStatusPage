resource "aws_s3_bucket" "artifacts" {
  bucket = "cloudstatus-artifacts-${local.region}"
}

resource "aws_s3_bucket_policy" "artifacts" {
  bucket = aws_s3_bucket.artifacts.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Principal = {
        AWS = [for _, account_id in local.probe_accounts : "arn:aws:iam::${account_id}:root"]
      }
      Effect = "Allow"
      Action = [
        "s3:List*",
        "s3:Get*",
      ]
      Resource = [
        aws_s3_bucket.artifacts.arn,
        "${aws_s3_bucket.artifacts.arn}/*",
      ]
    }]
  })
}

resource "aws_s3_bucket_public_access_block" "artifacts" {
  bucket                  = aws_s3_bucket.artifacts.id
  block_public_acls       = "true"
  block_public_policy     = "true"
  ignore_public_acls      = "true"
  restrict_public_buckets = "true"
}
