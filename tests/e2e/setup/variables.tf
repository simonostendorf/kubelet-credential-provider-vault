variable "kubeconfig_path" {
  type = string
  default = "../kubeconfig.yaml"
}

variable "vault_addr" {
  type = string
  default = "http://127.0.0.1:8200"
}

variable "vault_root_token" {
  type = string
  default = "vault-root-token"
}
