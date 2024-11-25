data "azurerm_client_config" "current" {}

locals {
  env = substr(var.environment, 0, 1)

  regions = {
    westus3 = {
      subnet = {
        ip_range = "10.10.102.0/24"
      }
    }
    westus2 = {
      subnet = {
        ip_range = "10.10.104.0/24"
      }
    }
    westus = {
      subnet = {
        ip_range = "10.10.106.0/24"
      }
    }
    westcentralus = {
      subnet = {
        ip_range = "10.10.108.0/24"
      }
    }
    southcentralus = {
      subnet = {
        ip_range = "10.10.110.0/24"
      }
    }
    northcentralus = {
      subnet = {
        ip_range = "10.10.112.0/24"
      }
    }
    eastus2 = {
      subnet = {
        ip_range = "10.10.114.0/24"
      }
    }
    eastus = {
      subnet = {
        ip_range = "10.10.116.0/24"
      }
    }
    centralus = {
      subnet = {
        ip_range = "10.10.118.0/24"
      }
    }
    brazilsoutheast = {
      subnet = {
        ip_range = "10.10.120.0/24"
      }
    }
    brazilsouth = {
      subnet = {
        ip_range = "10.10.122.0/24"
      }
    }
    uaenorth = {
      subnet = {
        ip_range = "10.10.124.0/24"
      }
    }
    uaecentral = {
      subnet = {
        ip_range = "10.10.126.0/24"
      }
    }
    qatarcentral = {
      subnet = {
        ip_range = "10.10.128.0/24"
      }
    }
    westeurope = {
      subnet = {
        ip_range = "10.10.130.0/24"
      }
    }
    ukwest = {
      subnet = {
        ip_range = "10.10.132.0/24"
      }
    }
    uksouth = {
      subnet = {
        ip_range = "10.10.134.0/24"
      }
    }
    switzerlandwest = {
      subnet = {
        ip_range = "10.10.136.0/24"
      }
    }
    switzerlandnorth = {
      subnet = {
        ip_range = "10.10.138.0/24"
      }
    }
    swedencentral = {
      subnet = {
        ip_range = "10.10.140.0/24"
      }
    }
    polandcentral = {
      subnet = {
        ip_range = "10.10.142.0/24"
      }
    }
    norwaywest = {
      subnet = {
        ip_range = "10.10.144.0/24"
      }
    }
    norwayeast = {
      subnet = {
        ip_range = "10.10.146.0/24"
      }
    }
    northeurope = {
      subnet = {
        ip_range = "10.10.148.0/24"
      }
    }
    italynorth = {
      subnet = {
        ip_range = "10.10.150.0/24"
      }
    }
    germanywestcentral = {
      subnet = {
        ip_range = "10.10.152.0/24"
      }
    }
    germanynorth = {
      subnet = {
        ip_range = "10.10.154.0/24"
      }
    }
    francesouth = {
      subnet = {
        ip_range = "10.10.156.0/24"
      }
    }
    francecentral = {
      subnet = {
        ip_range = "10.10.158.0/24"
      }
    }
    canadaeast = {
      subnet = {
        ip_range = "10.10.160.0/24"
      }
    }
    canadacentral = {
      subnet = {
        ip_range = "10.10.162.0/24"
      }
    }
    westindia = {
      subnet = {
        ip_range = "10.10.164.0/24"
      }
    }
    southindia = {
      subnet = {
        ip_range = "10.10.166.0/24"
      }
    }
    southeastasia = {
      subnet = {
        ip_range = "10.10.168.0/24"
      }
    }
    koreasouth = {
      subnet = {
        ip_range = "10.10.170.0/24"
      }
    }
    koreacentral = {
      subnet = {
        ip_range = "10.10.172.0/24"
      }
    }
    japanwest = {
      subnet = {
        ip_range = "10.10.174.0/24"
      }
    }
    japaneast = {
      subnet = {
        ip_range = "10.10.176.0/24"
      }
    }
    eastasia = {
      subnet = {
        ip_range = "10.10.178.0/24"
      }
    }
    centralindia = {
      subnet = {
        ip_range = "10.10.180.0/24"
      }
    }
    australiasoutheast = {
      subnet = {
        ip_range = "10.10.182.0/24"
      }
    }
    australiaeast = {
      subnet = {
        ip_range = "10.10.184.0/24"
      }
    }
    australiacentral2 = {
      subnet = {
        ip_range = "10.10.186.0/24"
      }
    }
    australiacentral = {
      subnet = {
        ip_range = "10.10.188.0/24"
      }
    }
    southafricawest = {
      subnet = {
        ip_range = "10.10.190.0/24"
      }
    }
    southafricanorth = {
      subnet = {
        ip_range = "10.10.192.0/24"
      }
    }
    spaincentral = {
      subnet = {
        ip_range = "10.10.194.0/24"
      }
    }
    mexicocentral = {
      subnet = {
        ip_range = "10.10.196.0/24"
      }
    }
    israelcentral = {
      subnet = {
        ip_range = "10.10.198.0/24"
      }
    }
  }

  subnets = { for region, cfg in local.regions : region => cfg.subnet if contains(var.regions, region) }

  regions_cartesian = { for group in setproduct(var.regions, var.regions) : "${group[0]}_${group[1]}" => { src = group[0], dst = group[1] } if group[0] != group[1] }
}
