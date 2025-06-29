package actions

import (
	"context"
	"log/slog"

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/act3-ai/gitoci/pkg/apis"
	"github.com/act3-ai/gitoci/pkg/apis/gitoci.act3-ai.io/v1alpha1"

	"github.com/act3-ai/go-common/pkg/config"
)

var (
	// parts used to assemble config file name and default paths
	parts = []string{"gitoci", "config.yaml"}

	// List of patterns used to match a config file for schema validation
	FileMatch = config.DefaultConfigValidatePath(parts...)

	// Default config locations in descending priority order
	DefaultSearchPath = config.DefaultConfigSearchPath(parts...)

	// DefaultPath is the path we would save the configuration to if needed. In a sense it is the preferred configuration path.
	DefaultPath = config.EnvOr("GITOCI_CONFIG", config.DefaultConfigPath(parts...))
)

// Tool represents the base action
type Tool struct {
	version   string
	apiScheme *runtime.Scheme

	// ConfigFiles stores the search locations for the config file in ascending priority order
	ConfigFiles []string

	// Handles overrides for configuration
	ConfigOverrideFunctions []func(ctx context.Context, c *v1alpha1.Configuration) error
}

// NewTool creates a new Tool with default values
func NewTool(version string) *Tool {
	return &Tool{
		version:     version,
		apiScheme:   apis.NewScheme(),
		ConfigFiles: DefaultSearchPath,
	}
}

// Version returns the version (overwritten by main.version if needed)
func (action *Tool) Version() string {
	return action.version
}

// GetScheme returns the runtime scheme used for configuration file loading
func (action *Tool) GetScheme() *runtime.Scheme {
	return action.apiScheme
}

// AddConfigOverrideFunction adds an override function that will be called in GetConfig to edit config after loading
func (action *Tool) AddConfigOverrideFunction(overrideFunction ...func(ctx context.Context, c *v1alpha1.Configuration) error) {
	if action.ConfigOverrideFunctions == nil {
		action.ConfigOverrideFunctions = []func(ctx context.Context, c *v1alpha1.Configuration) error{}
	}
	action.ConfigOverrideFunctions = append(action.ConfigOverrideFunctions, overrideFunction...)
}

// GetConfig loads Configuration using the current Tool options
func (action *Tool) GetConfig(ctx context.Context) (c *v1alpha1.Configuration, err error) {
	c = &v1alpha1.Configuration{}

	err = config.Load(slog.Default(), action.GetScheme(), c, action.ConfigFiles)
	if err != nil {
		return c, err
	}

	defer slog.Debug("using config", "Configuration", c)

	// Loop through override functions, applying each to the configuration
	for _, overrideFunction := range action.ConfigOverrideFunctions {
		err = overrideFunction(ctx, c)
		if err != nil {
			return c, err
		}
	}

	return c, nil
}
