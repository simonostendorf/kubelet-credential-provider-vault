resource "vault_auth_backend" "kubernetes" {
  path = "kubernetes"
  type = "kubernetes"
}

resource "vault_kubernetes_auth_backend_config" "this" {
  backend            = vault_auth_backend.kubernetes.path
  kubernetes_host    = var.kube_host
  kubernetes_ca_cert = base64decode(yamldecode(data.local_file.kubeconfig.content)["clusters"][0]["cluster"]["certificate-authority-data"])
  token_reviewer_jwt = data.kubernetes_secret_v1.vault_token_reviewer_token.data["token"]
}

resource "vault_kubernetes_auth_backend_role" "example" {
  backend                          = vault_auth_backend.kubernetes.path
  role_name                        = "example"
  bound_service_account_names      = ["*"]
  bound_service_account_namespaces = ["*"]
  token_policies                   = ["default", vault_policy.example.name]
}

data "local_file" "kubeconfig" {
  filename = var.kubeconfig_path
}

resource "kubernetes_service_account_v1" "vault_token_reviewer" {
  metadata {
    name      = "vault-token-reviewer"
    namespace = "default"
  }
}

resource "kubernetes_secret_v1" "vault_token_reviewer_token" {
  metadata {
    annotations = {
      "kubernetes.io/service-account.name" = kubernetes_service_account_v1.vault_token_reviewer.metadata.0.name
    }
    name      = kubernetes_service_account_v1.vault_token_reviewer.metadata.0.name
    namespace = kubernetes_service_account_v1.vault_token_reviewer.metadata.0.namespace
  }
  type                           = "kubernetes.io/service-account-token"
  wait_for_service_account_token = true
  depends_on = [
    kubernetes_service_account_v1.vault_token_reviewer
  ]
}

data "kubernetes_secret_v1" "vault_token_reviewer_token" {
  metadata {
    name      = kubernetes_secret_v1.vault_token_reviewer_token.metadata.0.name
    namespace = kubernetes_secret_v1.vault_token_reviewer_token.metadata.0.namespace
  }
  depends_on = [
    kubernetes_secret_v1.vault_token_reviewer_token
  ]
}

resource "kubernetes_cluster_role_binding_v1" "vault_token_reviewer" {
  metadata {
    name = kubernetes_service_account_v1.vault_token_reviewer.metadata.0.name
  }
  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = "system:auth-delegator"
  }
  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account_v1.vault_token_reviewer.metadata.0.name
    namespace = kubernetes_service_account_v1.vault_token_reviewer.metadata.0.namespace
  }
}
