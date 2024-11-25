locals {
  regions = {
    africa-south1 = {
      multi_region = "asia"
      subnet = {
        ip_range = "10.10.102.0/24"
      }
    }
    asia-east1 = {
      multi_region = "asia"
      subnet = {
        ip_range = "10.10.104.0/24"
      }
    }
    asia-east2 = {
      multi_region = "asia"
      subnet = {
        ip_range = "10.10.106.0/24"
      }
    }
    asia-northeast1 = {
      multi_region = "asia"
      subnet = {
        ip_range = "10.10.108.0/24"
      }
    }
    asia-northeast2 = {
      multi_region = "asia"
      subnet = {
        ip_range = "10.10.110.0/24"
      }
    }
    asia-northeast3 = {
      multi_region = "asia"
      subnet = {
        ip_range = "10.10.112.0/24"
      }
    }
    asia-south1 = {
      multi_region = "asia"
      subnet = {
        ip_range = "10.10.114.0/24"
      }
    }
    asia-south2 = {
      multi_region = "asia"
      subnet = {
        ip_range = "10.10.116.0/24"
      }
    }
    asia-southeast1 = {
      multi_region = "asia"
      subnet = {
        ip_range = "10.10.118.0/24"
      }
    }
    asia-southeast2 = {
      multi_region = "asia"
      subnet = {
        ip_range = "10.10.120.0/24"
      }
    }
    australia-southeast1 = {
      multi_region = "asia"
      subnet = {
        ip_range = "10.10.122.0/24"
      }
    }
    australia-southeast2 = {
      multi_region = "asia"
      subnet = {
        ip_range = "10.10.124.0/24"
      }
    }
    europe-central2 = {
      multi_region = "europe"
      subnet = {
        ip_range = "10.10.126.0/24"
      }
    }
    europe-north1 = {
      multi_region = "europe"
      subnet = {
        ip_range = "10.10.128.0/24"
      }
    }
    europe-southwest1 = {
      multi_region = "europe"
      subnet = {
        ip_range = "10.10.130.0/24"
      }
    }
    europe-west1 = {
      multi_region = "europe"
      subnet = {
        ip_range = "10.10.132.0/24"
      }
    }
    europe-west10 = {
      multi_region = "europe"
      subnet = {
        ip_range = "10.10.134.0/24"
      }
    }
    europe-west12 = {
      multi_region = "europe"
      subnet = {
        ip_range = "10.10.136.0/24"
      }
    }
    europe-west2 = {
      multi_region = "europe"
      subnet = {
        ip_range = "10.10.138.0/24"
      }
    }
    europe-west3 = {
      multi_region = "europe"
      subnet = {
        ip_range = "10.10.140.0/24"
      }
    }
    europe-west4 = {
      multi_region = "europe"
      subnet = {
        ip_range = "10.10.142.0/24"
      }
    }
    europe-west6 = {
      multi_region = "europe"
      subnet = {
        ip_range = "10.10.144.0/24"
      }
    }
    europe-west8 = {
      multi_region = "europe"
      subnet = {
        ip_range = "10.10.146.0/24"
      }
    }
    europe-west9 = {
      multi_region = "europe"
      subnet = {
        ip_range = "10.10.148.0/24"
      }
    }
    me-central1 = {
      multi_region = "asia"
      subnet = {
        ip_range = "10.10.150.0/24"
      }
    }
    me-central2 = {
      multi_region = "asia"
      subnet = {
        ip_range = "10.10.152.0/24"
      }
    }
    me-west1 = {
      multi_region = "asia"
      subnet = {
        ip_range = "10.10.154.0/24"
      }
    }
    northamerica-northeast1 = {
      multi_region = "us"
      subnet = {
        ip_range = "10.10.156.0/24"
      }
    }
    northamerica-northeast2 = {
      multi_region = "us"
      subnet = {
        ip_range = "10.10.158.0/24"
      }
    }
    southamerica-east1 = {
      multi_region = "us"
      subnet = {
        ip_range = "10.10.160.0/24"
      }
    }
    southamerica-west1 = {
      multi_region = "us"
      subnet = {
        ip_range = "10.10.162.0/24"
      }
    }
    us-central1 = {
      multi_region = "us"
      subnet = {
        ip_range = "10.10.164.0/24"
      }
    }
    us-east1 = {
      multi_region = "us"
      subnet = {
        ip_range = "10.10.166.0/24"
      }
    }
    us-east4 = {
      multi_region = "us"
      subnet = {
        ip_range = "10.10.168.0/24"
      }
    }
    us-east5 = {
      multi_region = "us"
      subnet = {
        ip_range = "10.10.170.0/24"
      }
    }
    us-south1 = {
      multi_region = "us"
      subnet = {
        ip_range = "10.10.172.0/24"
      }
    }
    us-west1 = {
      multi_region = "us"
      subnet = {
        ip_range = "10.10.174.0/24"
      }
    }
    us-west2 = {
      multi_region = "us"
      subnet = {
        ip_range = "10.10.176.0/24"
      }
    }
    us-west3 = {
      multi_region = "us"
      subnet = {
        ip_range = "10.10.178.0/24"
      }
    }
    us-west4 = {
      multi_region = "us"
      subnet = {
        ip_range = "10.10.180.0/24"
      }
    }
  }

  subnets = { for region, cfg in local.regions : region => cfg.subnet if contains(var.regions, region) }
}
