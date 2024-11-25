resource "aws_s3_bucket" "mon-probe" {
  bucket = "${var.account}-${local.region}"

  tags = {
    cost-center = "mon-probe"
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "mon-probe" {
  bucket = aws_s3_bucket.mon-probe.id

  rule {
    id     = "test/"
    status = "Enabled"

    expiration {
      days = 1
    }

    filter {
      prefix = "test/"
    }
  }
}
