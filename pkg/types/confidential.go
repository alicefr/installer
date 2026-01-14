package types

// ConfidentialCluster represents the clevis and attestation configuration for confidential clusters.
type ConfidentialCluster struct {
	// Clevis is the configuration for LUKS encryption using Clevis.
	// +optional
	Clevis *Clevis `json:"clevis,omitempty"`

	// Attestation is the configuration for attestation key generation and registration.
	// +optional
	Attestation *Attestation `json:"attestation,omitempty"`
}

// Clevis represents the clevis configuration for LUKS encryption.
type Clevis struct {
	// Trustee is the trustee server configuration for Clevis.
	Trustee Trustee `json:"trustee"`
}

// Trustee represents the trustee server configuration.
type Trustee struct {
	// Servers is the list of trustee servers.
	Servers []TrusteeServer `json:"servers"`

	// Path is the trustee path.
	Path string `json:"path"`
}

// TrusteeServer represents a single trustee server configuration.
type TrusteeServer struct {
	// URL is the trustee server URL.
	URL string `json:"url"`

	// Cert is the trustee server certificate.
	Cert string `json:"cert"`
}

// Attestation represents the attestation configuration.
type Attestation struct {
	// AttestationKey is the attestation key configuration.
	AttestationKey AttestationKey `json:"attestation_key,omitempty"`
}

// AttestationKey represents the attestation key configuration.
type AttestationKey struct {
	// Registration is the registration configuration for the attestation key.
	Registration Registration `json:"registration,omitempty"`
}

// Registration represents the registration configuration.
type Registration struct {
	// URL is the registration server URL.
	URL string `json:"url,omitempty"`

	// Certificate is the registration server certificate.
	Certificate string `json:"certificate,omitempty"`
}
