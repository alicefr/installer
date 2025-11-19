package bootstrap

import (
	"context"
	"encoding/json"
	"os"

	igntypes "github.com/coreos/ignition/v2/config/v3_6_experimental/types"
	"github.com/openshift/installer/pkg/asset"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/utils/ptr"
)

const (
	ConfidentialClusterConfigEnvVar = "OPENSHIFT_INSTALL_CONFIDENTIAL_CLUSTER_CONFIG"
)

// TrusteeConfig represents the trustee-specific configuration.
type TrusteeConfig struct {
	Servers []TrusteeServer `json:"servers"`
	Path    string          `json:"path"`
}

// TrusteeServer represents a trustee server configuration.
type TrusteeServer struct {
	URL  string `json:"url"`
	Cert string `json:"cert"`
}

// ConfidentialClusterConfigJSON represents the JSON structure for confidential cluster configuration.
type ConfidentialClusterConfigJSON struct {
	Clevis      *ClevisConfig      `json:"clevis,omitempty"`
	Attestation *AttestationConfig `json:"attestation,omitempty"`
}

// ClevisConfig represents the clevis configuration in JSON.
type ClevisConfig struct {
	TrusteeConfig TrusteeConfig `json:"trustee"`
}

// AttestationConfig represents the attestation configuration in JSON.
type AttestationConfig struct {
	AttestationKey AttestationKeyConfig `json:"attestation_key,omitempty"`
}

// AttestationKeyConfig represents the attestation key configuration in JSON.
type AttestationKeyConfig struct {
	Registration RegistrationConfig `json:"registration,omitempty"`
}

// RegistrationConfig represents the registration configuration in JSON.
type RegistrationConfig struct {
	URL         string `json:"url,omitempty"`
	Certificate string `json:"certificate,omitempty"`
}

// ConfidentialClusterConfig represents the clevis and attestation configuration for confidential clusters.
type ConfidentialClusterConfig struct {
	File        *asset.File
	Config      *igntypes.Clevis
	Attestation *igntypes.Attestation
}

var _ asset.Asset = (*ConfidentialClusterConfig)(nil)

// Dependencies returns no dependencies.
func (c *ConfidentialClusterConfig) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate loads the confidential cluster configuration from the JSON file and builds the Clevis and Attestation configs.
func (c *ConfidentialClusterConfig) Generate(_ context.Context, dependencies asset.Parents) error {
	configPath := os.Getenv(ConfidentialClusterConfigEnvVar)

	if configPath == "" {
		return nil
	}

	logrus.Infof("Loading confidential cluster configuration from %s", configPath)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return errors.Wrapf(err, "failed to read confidential cluster config file %s", configPath)
	}

	// Parse the confidential cluster configuration
	var config ConfidentialClusterConfigJSON
	if err := json.Unmarshal(data, &config); err != nil {
		return errors.Wrapf(err, "failed to parse confidential cluster config file %s", configPath)
	}

	// Build the Clevis configuration if present
	if config.Clevis != nil {
		// Convert the trustee config to JSON string for the Clevis custom config
		trusteeConfigJSON, err := json.Marshal(config.Clevis.TrusteeConfig)
		if err != nil {
			return errors.Wrap(err, "failed to marshal trustee config to JSON")
		}

		// Build the Clevis configuration with trustee pin
		needsNetwork := true
		trusteeConfigStr := string(trusteeConfigJSON)
		c.Config = &igntypes.Clevis{
			Custom: igntypes.ClevisCustom{
				Pin:          ptr.To("trustee"),
				Config:       ptr.To(trusteeConfigStr),
				NeedsNetwork: ptr.To(needsNetwork),
			},
		}
		logrus.Debug("Successfully loaded Clevis configuration")
	}

	// Build the Attestation configuration if present
	if config.Attestation != nil {
		c.Attestation = &igntypes.Attestation{
			AttestationKey: igntypes.AttestationKey{
				Registration: igntypes.Registration{
					Url:         ptr.To(config.Attestation.AttestationKey.Registration.URL),
					Certificate: ptr.To(config.Attestation.AttestationKey.Registration.Certificate),
				},
			},
		}
		logrus.Debug("Successfully loaded Attestation configuration")
	}

	c.File = &asset.File{
		Filename: configPath,
		Data:     data,
	}

	logrus.Debug("Successfully loaded confidential cluster configuration")
	return nil
}

func (c *ConfidentialClusterConfig) Name() string {
	return "Confidential Cluster Config"
}

func (c *ConfidentialClusterConfig) Load(asset.FileFetcher) (found bool, err error) {
	return false, nil
}
