variable "vm_size" {
  type    = string
  default = "Standard_B1ls"
}

variable "vm_spot" {
  type    = bool
  default = false
}

variable "zones" {
  type    = list(string)
  default = null
}

variable "environment" {
  type = string
}

variable "region" {
  type = string
}

variable "rg_name" {
  type = string
}

variable "user_identity_id" {
  type = string
}

variable "subnet_id" {
  type = string
}

variable "vault_url" {
  type = string
}

variable "artifacts_container_url" {
  type = string
}

variable "app_version" {
  type = string
}

variable "env" {
  type = map(string)
}

variable "ssh_key" {
  type    = string
  default = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIDl0M37evHnmKlLkYg8gqVWtH94qJCaSgBuVhgC3Res7"
}

variable "revision" {
  type    = number
  default = 0
}

variable "tags" {
  type    = map(string)
  default = {}
}
