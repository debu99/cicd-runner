package vault

import (
	"fmt"

	"github.com/debu99/cicd-runner/common"
	"github.com/debu99/cicd-runner/helpers/secrets"
	"github.com/debu99/cicd-runner/helpers/vault/service"
)

const (
	resolverName = "vault"
)

var newVaultService = service.NewVault

type resolver struct {
	secret common.Secret
}

func newResolver(secret common.Secret) common.SecretResolver {
	return &resolver{
		secret: secret,
	}
}

func (v *resolver) Name() string {
	return resolverName
}

func (v *resolver) IsSupported() bool {
	return v.secret.Vault != nil
}

func (v *resolver) Resolve() (string, error) {
	if !v.IsSupported() {
		return "", secrets.NewResolvingUnsupportedSecretError(resolverName)
	}

	secret := v.secret.Vault

	url := secret.Server.URL

	s, err := newVaultService(url, secret)
	if err != nil {
		return "", err
	}

	data, err := s.GetField(secret, secret)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", data), nil
}

func init() {
	common.GetSecretResolverRegistry().Register(newResolver)
}
