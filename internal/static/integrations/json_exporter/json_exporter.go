package json_exporter

import (
	"fmt"
	"github.com/go-kit/log"
	"github.com/grafana/alloy/internal/static/integrations"
	integrations_v2 "github.com/grafana/alloy/internal/static/integrations/v2"
	"github.com/grafana/alloy/internal/static/integrations/v2/metricsutils"
	"github.com/prometheus-community/json_exporter/exporter"
	"os"
	"path/filepath"
)

var DefaultConfig = Config{
	ConfigPath: "",
}

// Config controls the json_exporter integration.
type Config struct {
	ConfigPath string `yaml:"config_path,omitempty"`
}

// UnmarshalYAML implements yaml.Unmarshaler for Config
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = DefaultConfig

	type plain Config
	return unmarshal((*plain)(c))
}

// Name returns the name of the integration that this config represents.
func (c *Config) Name() string {
	return "json_exporter"
}

// InstanceKey returns the config file path the json_exporter being queried.
func (c *Config) InstanceKey(_ string) (string, error) {
	rootDir, err := os.Getwd()

	if err != nil {
		return "", fmt.Errorf("could not found root dir")
	}
	if c.ConfigPath == "" {
		return "", fmt.Errorf("no configuration file specified")
	}

	configFilePath := filepath.Join(rootDir, c.ConfigPath)

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return "", fmt.Errorf("configuration file %s does not exist", configFilePath)
	}

	return configFilePath, err
}

// NewIntegration creates a new json_exporter
func (c *Config) NewIntegration(logger log.Logger) (integrations.Integration, error) {
	return New(logger, c)
}

func init() {
	integrations.RegisterIntegration(&Config{})
	integrations_v2.RegisterLegacy(&Config{}, integrations_v2.TypeMultiplex, metricsutils.NewNamedShim("json"))
}

// New creates a new json_exporter integration.
func New(logger log.Logger, c *Config) (integrations.Integration, error) {
	//conf, err := json_config.LoadConfig(c.ConfigPath)
	//
	//if err != nil {
	//	return nil, fmt.Errorf("can not load json_exporter config. error: %w", errors.Unwrap(err))
	//}

	jsonExporter := exporter.JSONMetricCollector{
		Logger: logger,
	}

	//jsonExporter := json_exporter. {
	//	APIMetrics: json_exporter.AddMetrics(),
	//	Config:     conf,
	//}

	return integrations.NewCollectorIntegration(
		c.Name(),
		integrations.WithCollectors(&jsonExporter),
	), nil

}
