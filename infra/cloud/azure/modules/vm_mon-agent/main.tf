locals {
  env = substr(var.environment, 0, 1)
  tags = merge({
    cost-center = "mon-agent"
    app_name    = "mon-agent"
    app_version = var.app_version
  }, var.tags)
  startup_script_env = {
    ARTIFACTS_CONTAINER_URL = var.artifacts_container_url
    VAULT_URL               = var.vault_url
    ENV = merge({
      LOG_LEVEL  = "info"
      LOG_FORMAT = "text"
    }, var.env)
  }
}

resource "azurerm_orchestrated_virtual_machine_scale_set" "this" {
  name                        = "mon-agent-${local.env}-${var.region}-r${var.revision}"
  location                    = var.region
  resource_group_name         = var.rg_name
  sku_name                    = var.vm_size
  priority                    = var.vm_spot ? "Spot" : "Regular"
  eviction_policy             = var.vm_spot ? "Delete" : null
  instances                   = 1
  platform_fault_domain_count = 1
  zones                       = var.zones
  user_data_base64            = base64encode(templatefile("${path.module}/startup_script.sh", local.startup_script_env))
  tags                        = local.tags

  identity {
    type         = "UserAssigned"
    identity_ids = [var.user_identity_id]
  }

  os_profile {
    linux_configuration {
      computer_name_prefix            = "mon-agent-"
      provision_vm_agent              = true
      disable_password_authentication = true
      admin_username                  = "azureuser"
      admin_ssh_key {
        username   = "azureuser"
        public_key = var.ssh_key
      }
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "None"
  }

  source_image_reference {
    publisher = "resf"
    offer     = "rockylinux-x86_64"
    sku       = "9-base"
    version   = "latest"
  }

  plan {
    name      = "9-base"
    publisher = "resf"
    product   = "rockylinux-x86_64"

  }

  network_interface {
    name    = "mon-agent-nic"
    primary = true
    # network_security_group_id = the subnet is already associated with a security group

    ip_configuration {
      name      = "mon-agent-ip"
      primary   = true
      subnet_id = var.subnet_id
      public_ip_address {
        name     = "mon-agent-public-ip"
        sku_name = "Standard_Regional"
      }
    }
  }

  automatic_instance_repair {
    enabled      = true
    grace_period = "PT10M"
  }

  extension {
    name                               = "AzureHealthExtension"
    publisher                          = "Microsoft.ManagedServices"
    type                               = "ApplicationHealthLinux"
    type_handler_version               = "2.0"
    auto_upgrade_minor_version_enabled = true
    settings = jsonencode({
      protocol          = "http"
      port              = 8000
      requestPath       = "/healthz"
      intervalInSeconds = 5
      numberOfProbes    = 2
      gracePeriod       = 180
    })
  }

  # TODO: costs and permissions
  #extension {
  #  name                                = "AzureMonitorLinuxAgent"
  #  publisher                           = "Microsoft.Azure.Monitor"
  #  type                                = "AzureMonitorLinuxAgent"
  #  type_handler_version                = "1.0"
  #  auto_upgrade_minor_version_enabled  = true
  #}

  lifecycle {
    create_before_destroy = true
    ignore_changes        = [os_disk]
  }
}

resource "null_resource" "vmss_update" {
  # Flexible scale set don't have any rolling update settings yet
  # https://learn.microsoft.com/en-us/azure/virtual-machine-scale-sets/virtual-machine-scale-sets-configure-rolling-upgrades
  provisioner "local-exec" {
    when    = create
    command = "az vmss update-instances --subscription ${element(split("/", azurerm_orchestrated_virtual_machine_scale_set.this.id), 2)} --instance-ids '*' --name ${azurerm_orchestrated_virtual_machine_scale_set.this.name} --resource-group ${azurerm_orchestrated_virtual_machine_scale_set.this.resource_group_name} && az vmss reimage --ids ${azurerm_orchestrated_virtual_machine_scale_set.this.id} --no-wait"
  }

  triggers = {
    app_version = var.app_version
  }

  depends_on = [
    azurerm_orchestrated_virtual_machine_scale_set.this
  ]
}
