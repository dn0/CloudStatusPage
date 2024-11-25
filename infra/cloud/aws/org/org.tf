resource "aws_organizations_organization" "org" {
  feature_set = "ALL"

  aws_service_access_principals = [
    "sso.amazonaws.com",
  ]

  enabled_policy_types = [
    "SERVICE_CONTROL_POLICY",
  ]
}

resource "aws_organizations_policy" "default" {
  name = "Default"
  type = "SERVICE_CONTROL_POLICY"

  content = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid      = "FullAWSAccess"
        Effect   = "Allow"
        Action   = "*"
        Resource = "*"
      },
      {
        Sid    = "DenyMemberAccountInstances"
        Effect = "Deny"
        Action = [
          "sso:CreateInstance"
        ],
        Resource = "*"
      },
    ]
  })
}

resource "aws_organizations_policy_attachment" "org" {
  policy_id = aws_organizations_policy.default.id
  target_id = aws_organizations_organization.org.roots[0].id
}

resource "aws_organizations_policy_attachment" "account" {
  for_each = local.all_accounts

  policy_id = aws_organizations_policy.default.id
  target_id = each.value
}
