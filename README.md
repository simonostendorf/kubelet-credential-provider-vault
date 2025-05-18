# Kubelet Credential Provider Vault

A [Kubernetes Kubelet Image Credential Provider](https://kubernetes.io/docs/tasks/administer-cluster/kubelet-credential-provider/) for [HashiCorp Vault](https://www.hashicorp.com/en/products/vault).

> [!CAUTION]
> This credential provider **relies on the new `KubeletServiceAccountTokenForCredentialProviders` feature** which was introduced in Kubernetes 1.33 in **alpha state**.

# Kubelet Credential Provider Vault

This repository contains the Kubelet Credential Provider Vault - a go binary functioning as kubelet plugin to provide image registry credentials stored in HashiCorp Vault.

## Usage

The plugin reads the `CredentialProviderRequest` via stdin, fetches the credentials from Vault and writes the `CredentialProviderResponse` to stdout.

The plugin will authenticate against Vault using the provided service account token (**Attention: this is only available in Kubernetes 1.33+**) using the [kubernetes auth method](https://developer.hashicorp.com/vault/docs/auth/kubernetes).

The plugin will fetch the credentials from the provided secret mount and secret name.
The secret must be structured as follows:

```json
{
	"password": "my-password",
	"username": "my-username"
}
```

### Supported CredentialProvider APIs

The plugin supports the following versions of the `CredentialProviderRequest`:

- [`credentialprovider.kubelet.k8s.io/v1`](https://kubernetes.io/docs/reference/config-api/kubelet-credentialprovider.v1/) (>= Kubernetes v1.26)

### Needed Feature Flags

And again: this plugin **relies on the new `KubeletServiceAccountTokenForCredentialProviders` feature** which was introduced in Kubernetes 1.33 in **alpha state**.

To use this credential provider, you need to enable the feature gate in your cluster configuration.

### Configuration

The application supports different ways for configuration.

- command line flags
- environment variables
- configuration file (yaml)
- .env file

The following configuration options are available:

| Flag                           | Description                                                                                                     | Environment Variable         | Config File Path           | Required | Default                                   |
| ------------------------------ | --------------------------------------------------------------------------------------------------------------- | ---------------------------- | -------------------------- | -------- | ----------------------------------------- |
| `--config`                     | configuration file to use. If not set, the application will look for `./kubelet-credential-provider-vault.yaml` | -                            | -                          | no       | -                                         |
| `--log-file`                   | file the logger will write to                                                                                   | `LOG_FILE`                   | `log.file`                 | no       | `./kubelet-credential-provider-vault.log` |
| `--log-level`                  | log level to use. Possible values: debug, info, warn, error                                                     | `LOG_LEVEL`                  | `log.level`                | no       | `info`                                    |
| `--log-enabled`                | enable or disable logging                                                                                       | `LOG_ENABLED`                | `log.enabled`              | no       | `true`                                    |
| `--vault-addr`                 | address of the Vault server                                                                                     | `VAULT_ADDR`                 | `vault.addr`               | yes      | -                                         |
| `--vault-insecure-skip-verify` | skip TLS verification of the Vault server                                                                       | `VAULT_INSECURE_SKIP_VERIFY` | `vault.insecureSkipVerify` | no       | `false`                                   |
| `--vault-auth-method`          | name of the auth method to use. Possible values: kubernetes                                                     | `VAULT_AUTH_METHOD`          | `vault.auth.method`        | no       | `kubernetes`                              |
| `--vault-auth-mount`           | name of the auth mount to use                                                                                   | `VAULT_AUTH_MOUNT`           | `vault.auth.mount`         | yes      | -                                         |
| `--vault-auth-role`            | name of the auth role to use                                                                                    | `VAULT_AUTH_ROLE`            | `vault.auth.role`          | yes      | -                                         |
| `--vault-secret-mount`         | name of the secret mount to use                                                                                 | `VAULT_SECRET_MOUNT`         | `vault.secret.mount`       | yes      | -                                         |
| `--vault-secret-name`          | name of the secret to use                                                                                       | `VAULT_SECRET_NAME`          | `vault.secret.name`        | yes      | -                                         |

### Usage with kubelet

The binary must be added to all nodes in the cluster (where the kubelet is running).

The folder where the binary is located must be references in the `--image-credential-provider-bin-dir` kubelet flag.
The name of the binary must be the same as in the following configuration file.

Example configuration file:
You can read more about the configuration file in the [official documentation](https://kubernetes.io/docs/tasks/administer-cluster/kubelet-credential-provider/).

```yaml
apiVersion: kubelet.config.k8s.io/v1
kind: CredentialProviderConfig
providers:
  - apiVersion: credentialprovider.kubelet.k8s.io/v1
    name: vault
    matchImages:
      - "registry.example.com"
    defaultCacheDuration: 1m
    args:
      - --vault-addr="https://vault.example.com"
      - --vault-auth-method="kubernetes"
      - --vault-auth-mount="kubernetes"
      - --vault-auth-role="example"
      - --vault-secret-mount="secret"
      - --vault-secret-path="example"
    env: []
    tokenAttributes:
      serviceAccountTokenAudience: "example"
      requireServiceAccount: true
```

The configuration file must also be placed on all nodes and must be referenced in the `--image-credential-provider-config` kubelet flag.

## Releases

The plugin will be released using [Semantic Versioning](https://semver.org/).

The releases are available on the [Releases](https://github.com/simonostendorf/kubelet-credential-provider-vault/releases) page.
You can find the release notes for each release there.

### Binaries

The following binaries are available:

- linux/amd64
- linux/arm64

## Contributing and Development

Contributions are welcome!
Maybe I will add a `CONTRIBUTING.md` file in the future to provide more information about the contribution process.
Please feel free to open an issue or a pull request if you have any questions or suggestions.
If you want to contribute to the project, please fork the repository and create a new branch for your changes.

## License

This project is licensed under the [MIT License](LICENSE.md).

## Security Policy

Please see our [Security Policy](SECURITY.md) for information about supported versions and how to report security vulnerabilities.

## References

More information about the CredentialProvider API can be found here:

- [https://pkg.go.dev/k8s.io/kubelet/pkg/apis/credentialprovider](https://pkg.go.dev/k8s.io/kubelet/pkg/apis/credentialprovider)
- [https://kubernetes.io/docs/reference/config-api/kubelet-credentialprovider.v1/](https://kubernetes.io/docs/reference/config-api/kubelet-credentialprovider.v1/)
- [https://kubernetes.io/docs/tasks/administer-cluster/kubelet-credential-provider/](https://kubernetes.io/docs/tasks/administer-cluster/kubelet-credential-provider/)
- [https://hyperconnect.github.io/2022/02/21/no-more-image-pull-secrets.html](https://hyperconnect.github.io/2022/02/21/no-more-image-pull-secrets.html)
- [https://github.com/kubernetes/enhancements/tree/master/keps/sig-node/2133-kubelet-credential-providers#credential-provider-configuration](https://github.com/kubernetes/enhancements/tree/master/keps/sig-node/2133-kubelet-credential-providers#credential-provider-configuration)
- [https://github.com/adisky/sample-credential-provider](https://github.com/adisky/sample-credential-provider)
- [https://github.com/kubernetes/cloud-provider-aws/tree/master/cmd/ecr-credential-provider](https://github.com/kubernetes/cloud-provider-aws/tree/master/cmd/ecr-credential-provider)
