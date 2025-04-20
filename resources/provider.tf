terraform {
  required_providers {
    grafana = {
      source  = "grafana/grafana"
      version = "~> 1.40.0"
    }
  }
}

provider "kubernetes" {
  config_path = "~/.kube/config"
}

provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}
