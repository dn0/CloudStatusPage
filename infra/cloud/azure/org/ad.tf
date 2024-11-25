resource "azuread_user" "dano" {
  user_principal_name  = "admin@cloudstatus.onmicrosoft.com"
  mail_nickname        = "dano"
  mobile_phone         = ""
  show_in_address_list = false
}

resource "azuread_directory_role" "global_admin" {
  display_name = "Global Administrator"
}

resource "azuread_directory_role_assignment" "global_admin" {
  for_each = {
    (azuread_user.dano.user_principal_name) = azuread_user.dano.object_id,
  }

  role_id             = azuread_directory_role.global_admin.template_id # use template_id instead of object_id when referencing built-in roles
  principal_object_id = each.value
  directory_scope_id  = "/"
}

resource "azuread_group" "admin" {
  display_name     = "admin"
  security_enabled = true

  owners = [
    azuread_user.dano.object_id,
  ]

  members = [
    azuread_user.dano.object_id,
  ]
}
