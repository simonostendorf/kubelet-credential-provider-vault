resource "vault_policy" "example" {
  name   = "example"
  policy = data.vault_policy_document.example.hcl
}

data "vault_policy_document" "example" {
  rule {
    path         = "secret/data/example"
    capabilities = ["read"]
  }
}
