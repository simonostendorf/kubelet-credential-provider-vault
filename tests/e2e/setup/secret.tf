# vault in dev mode already has a secret backend mounted at `secret/`
# resource "vault_mount" "secret" {
#   path        = "secret"
#   type        = "kv-v2"
#   options = {
#     version = "2"
#     type    = "kv-v2"
#   }
# }

resource "vault_kv_secret_v2" "example" {
  mount               = "secret"
  name                = "example"
  delete_all_versions = true
  data_json = jsonencode({
    username = "username"
    password = "password"
  })
}
