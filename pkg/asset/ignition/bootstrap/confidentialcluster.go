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

// ConfidentialClusterConfig represents the clevis configuration for confidential clusters.
type ConfidentialClusterConfig struct {
	File   *asset.File
	Config *igntypes.Clevis
}

var _ asset.Asset = (*ConfidentialClusterConfig)(nil)

// Dependencies returns no dependencies.
func (c *ConfidentialClusterConfig) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate loads the trustee configuration from the JSON file and builds the Clevis config.
func (c *ConfidentialClusterConfig) Generate(_ context.Context, dependencies asset.Parents) error {
	configPath := os.Getenv(ConfidentialClusterConfigEnvVar)

	if configPath == "" {
		return nil
	}

	logrus.Infof("Loading confidential cluster trustee configuration from %s", configPath)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return errors.Wrapf(err, "failed to read confidential cluster config file %s", configPath)
	}

	// Parse the trustee-specific configuration
	var trusteeConfig TrusteeConfig
	if err := json.Unmarshal(data, &trusteeConfig); err != nil {
		return errors.Wrapf(err, "failed to parse confidential cluster config file %s", configPath)
	}

	// Convert the trustee config to JSON string for the Clevis custom config
	trusteeConfigJSON, err := json.Marshal(trusteeConfig)
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

	c.File = &asset.File{
		Filename: configPath,
		Data:     data,
	}

	logrus.Debug("Successfully loaded confidential cluster trustee configuration")
	return nil
}

func (c *ConfidentialClusterConfig) Name() string {
	return "Confidential Cluster Config"
}

func (c *ConfidentialClusterConfig) Load(asset.FileFetcher) (found bool, err error) {
	return false, nil
}
