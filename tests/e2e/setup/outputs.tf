output "service_account_token" {
  value = kubernetes_secret_v1.example_token.data["token"]
  sensitive = true
}
