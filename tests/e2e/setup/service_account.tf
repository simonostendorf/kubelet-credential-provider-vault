resource "kubernetes_service_account_v1" "example" {
  metadata {
    name      = "example"
    namespace = "default"
  }
}

resource "kubernetes_secret_v1" "example_token" {
  metadata {
    annotations = {
      "kubernetes.io/service-account.name" = kubernetes_service_account_v1.example.metadata.0.name
    }
    name     = kubernetes_service_account_v1.example.metadata.0.name
    namespace = kubernetes_service_account_v1.example.metadata.0.namespace
  }
  type                           = "kubernetes.io/service-account-token"
  wait_for_service_account_token = true
  depends_on = [
    kubernetes_service_account_v1.example
  ]
}

data "kubernetes_secret_v1" "example_token" {
  metadata {
    name      = kubernetes_secret_v1.example_token.metadata.0.name
    namespace = kubernetes_secret_v1.example_token.metadata.0.namespace
  }
  depends_on = [
    kubernetes_secret_v1.example_token
  ]
}
