module "aws-iam-identity-center" {
  source = "aws-ia/iam-identity-center/aws"

  sso_groups = {
    admin : {
      group_name        = "admin"
      group_description = "Admin IAM Identity Center Group"
    },
  }

  sso_users = {
    dano : {
      group_membership = ["admin"]
      user_name        = "dano"
      email            = "daniel@truestatus.cloud"
    },
  }

  permission_sets = {
    Administrator = {
      description      = "Provides AWS full access permissions.",
      session_duration = "PT12H",
      aws_managed_policies = [
        "arn:aws:iam::aws:policy/AdministratorAccess",
        "arn:aws:iam::aws:policy/job-function/Billing",
      ]
    },
    ReadOnly = {
      description      = "Provides AWS read-only permissions.",
      session_duration = "PT12H"
      aws_managed_policies = [
        "arn:aws:iam::aws:policy/ReadOnlyAccess",
      ]
    },
  }

  account_assignments = {
    Administrator : {
      principal_name  = "admin"
      principal_type  = "GROUP"
      principal_idp   = "INTERNAL"
      permission_sets = ["Administrator", "ReadOnly"]
      account_ids     = values(local.all_accounts)
    },
  }
}
