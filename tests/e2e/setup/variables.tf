variable "kubeconfig_path" {
  description = "Path to the kubeconfig file used by the Kubernetes provider to connect to the Kubernetes cluster."
  type        = string
  default     = "../kubeconfig.yaml"
}

variable "kube_host" {
  description = "Kubernetes Host Address that will be used by Vault to connect to the TokenReviewer API of the Kubernetes Cluster. Address and cluster must be accessible by Vault. Defaults to https://k3s.local:6443 for k3s setup defined in docker-compose.yaml."
  type        = string
  default     = "https://k3s.local:6443"
}

variable "vault_root_token" {
  description = "Vault root token used to initialize the Vault provider."
  type        = string
  default     = "vault-root-token"
}
