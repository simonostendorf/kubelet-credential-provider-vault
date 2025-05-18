terraform {
  required_version = "~> 1.9.1"
  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = "~> 4.8.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.36.0"
    }
    local = {
      source  = "hashicorp/local"
      version = "~> 2.5.3"
    }
  }
}

provider "vault" {
  address = var.vault_addr
  token   = var.vault_root_token
}

provider "kubernetes" {
  config_path = var.kubeconfig_path
}
