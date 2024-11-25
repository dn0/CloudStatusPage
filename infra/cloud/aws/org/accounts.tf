resource "aws_organizations_account" "this" {
  for_each = local.all_accounts

  name  = "cloudstatus-${each.key}"
  email = "daniel+aws-${each.key}@truestatus.cloud"
}
