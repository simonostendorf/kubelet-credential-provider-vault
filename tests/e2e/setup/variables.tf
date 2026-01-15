variable "kubeconfig_path" {
  type    = string
  default = "../kubeconfig.yaml"
}

variable "kube_host" {
  description = "Kubernetes Host Address that will be used by Vault to connect to the TokenReviewer API of the Kubernetes Cluster. Address and cluster must be accessable by Vault."
  type        = string
  default     = "https://k3s.local:6443"
}

variable "vault_root_token" {
  type    = string
  default = "vault-root-token"
}
