package tls

import (
	"crypto/x509"
	"crypto/x509/pkix"

	"github.com/openshift/installer/pkg/asset"
)

// KubeletCertKey is the asset that generates the kubelet key/cert pair.
// [DEPRECATED]
type KubeletCertKey struct {
	SignedCertKey
}

var _ asset.Asset = (*KubeletCertKey)(nil)

// Dependencies returns the dependency of the the cert/key pair, which includes
// the parent CA, and install config if it depends on the install config for
// DNS names, etc.
func (a *KubeletCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{
		&KubeCA{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *KubeletCertKey) Generate(dependencies asset.Parents) error {
	kubeCA := &KubeCA{}
	dependencies.Get(kubeCA)

	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: "system:serviceaccount:openshift-machine-config-operator:node-bootstrapper", Organization: []string{"system:serviceaccounts:openshift-machine-config-operator"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		Validity:     ValidityOneDay,
	}

	return a.SignedCertKey.Generate(cfg, kubeCA, "kubelet", DoNotAppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *KubeletCertKey) Name() string {
	return "Certificate (system:serviceaccount:openshift-machine-config-operator:node-bootstrapper)"
}

// KubeletCSRSignerCertKey is a key/cert pair that signs the kubelet client certs.
type KubeletCSRSignerCertKey struct {
	SelfSignedCertKey
}

var _ asset.WritableAsset = (*KubeletCSRSignerCertKey)(nil)

// Dependencies returns the dependency of the root-ca, which is empty.
func (c *KubeletCSRSignerCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates the root-ca key and cert pair.
func (c *KubeletCSRSignerCertKey) Generate(parents asset.Parents) error {
	cfg := &CertCfg{
		Subject:   pkix.Name{CommonName: "kubelet-signer", OrganizationalUnit: []string{"openshift"}},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:  ValidityOneDay,
		IsCA:      true,
	}

	return c.SelfSignedCertKey.Generate(cfg, "kubelet-signer")
}

// Name returns the human-friendly name of the asset.
func (c *KubeletCSRSignerCertKey) Name() string {
	return "Certificate (kubelet-signer)"
}

// KubeletClientCABundle is the asset the generates the kubelet-client-ca-bundle,
// which contains all the individual client CAs.
type KubeletClientCABundle struct {
	CertBundle
}

var _ asset.Asset = (*KubeletClientCABundle)(nil)

// Dependencies returns the dependency of the cert bundle.
func (a *KubeletClientCABundle) Dependencies() []asset.Asset {
	return []asset.Asset{
		&KubeletCSRSignerCertKey{},
	}
}

// Generate generates the cert bundle based on its dependencies.
func (a *KubeletClientCABundle) Generate(deps asset.Parents) error {
	var certs []CertInterface
	for _, asset := range a.Dependencies() {
		deps.Get(asset)
		certs = append(certs, asset.(CertInterface))
	}
	return a.CertBundle.Generate("kubelet-client-ca-bundle", certs...)
}

// Name returns the human-friendly name of the asset.
func (a *KubeletClientCABundle) Name() string {
	return "Certificate (kubelet-client-ca-bundle)"
}

// KubeletServingCABundle is the asset the generates the kubelet-serving-ca-bundle,
// which contains all the individual client CAs.
type KubeletServingCABundle struct {
	CertBundle
}

var _ asset.Asset = (*KubeletServingCABundle)(nil)

// Dependencies returns the dependency of the cert bundle.
func (a *KubeletServingCABundle) Dependencies() []asset.Asset {
	return []asset.Asset{
		&KubeletCSRSignerCertKey{},
	}
}

// Generate generates the cert bundle based on its dependencies.
func (a *KubeletServingCABundle) Generate(deps asset.Parents) error {
	var certs []CertInterface
	for _, asset := range a.Dependencies() {
		deps.Get(asset)
		certs = append(certs, asset.(CertInterface))
	}
	return a.CertBundle.Generate("kubelet-serving-ca-bundle", certs...)
}

// Name returns the human-friendly name of the asset.
func (a *KubeletServingCABundle) Name() string {
	return "Certificate (kubelet-serving-ca-bundle)"
}

// KubeletClientCertKey is the asset that generates the key/cert pair for kubelet client to apiserver.
type KubeletClientCertKey struct {
	SignedCertKey
}

var _ asset.Asset = (*KubeletClientCertKey)(nil)

// Dependencies returns the dependency of the the cert/key pair, which includes
// the parent CA, and install config if it depends on the install config for
// DNS names, etc.
func (a *KubeletClientCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{
		&KubeletCSRSignerCertKey{},
		&KubeAPIServerLBSignerCertKey{},
		&KubeAPIServerLocalhostSignerCertKey{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *KubeletClientCertKey) Generate(dependencies asset.Parents) error {
	ca := &KubeletCSRSignerCertKey{}
	dependencies.Get(ca)

	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: "system:serviceaccount:openshift-machine-config-operator:node-bootstrapper", Organization: []string{"system:serviceaccounts:openshift-machine-config-operator"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		Validity:     ValidityOneDay,
	}

	return a.SignedCertKey.Generate(cfg, ca, "kubelet-client", DoNotAppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *KubeletClientCertKey) Name() string {
	return "Certificate (kubelet-client)"
}
