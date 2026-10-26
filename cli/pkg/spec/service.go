package spec

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v3"
)

var (
	validNameRegex = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)
)

type ServiceSpec struct {
	Name      string            `yaml:"name"`
	Owner     string            `yaml:"owner"`
	Tier      string            `yaml:"tier"`    // standard | critical
	Runtime   string            `yaml:"runtime"` // node | go | php
	Port      int               `yaml:"port"`
	Env       map[string]string `yaml:"env,omitempty"`
	Resources ResourcesSpec     `yaml:"resources"`
	SLO       SLOSpec           `yaml:"slo"`
	Image     ImageSpec         `yaml:"image,omitempty"`
}

type ResourcesSpec struct {
	CPU    string `yaml:"cpu"`
	Memory string `yaml:"memory"`
}

type SLOSpec struct {
	Availability float64 `yaml:"availability"`   // e.g. 99.9
	LatencyP95MS int     `yaml:"latency_p95_ms"` // e.g. 250
}

type ImageSpec struct {
	Repository string `yaml:"repository,omitempty"`
	Tag        string `yaml:"tag,omitempty"`
}

func (s *ServiceSpec) Validate() error {
	if s.Name == "" {
		return errors.New("service 'name' is required")
	}
	if !validNameRegex.MatchString(s.Name) {
		return fmt.Errorf("service 'name' '%s' is invalid: must match DNS-1123 subdomain (lowercase alphanumeric and hyphens)", s.Name)
	}
	if s.Owner == "" {
		return errors.New("service 'owner' is required")
	}
	if s.Tier != "standard" && s.Tier != "critical" {
		return fmt.Errorf("invalid tier '%s': must be either 'standard' or 'critical'", s.Tier)
	}
	if s.Runtime != "node" && s.Runtime != "go" && s.Runtime != "php" {
		return fmt.Errorf("invalid runtime '%s': must be one of 'node', 'go', 'php'", s.Runtime)
	}
	if s.Port <= 0 || s.Port > 65535 {
		return fmt.Errorf("invalid port %d: must be between 1 and 65535", s.Port)
	}
	if s.Resources.CPU == "" || s.Resources.Memory == "" {
		return errors.New("both 'resources.cpu' and 'resources.memory' are required")
	}
	if s.SLO.Availability <= 0 || s.SLO.Availability >= 100 {
		return fmt.Errorf("invalid availability SLO %f: must be between 0 and 100", s.SLO.Availability)
	}
	if s.SLO.LatencyP95MS <= 0 {
		return fmt.Errorf("invalid latency_p95_ms %d: must be greater than 0", s.SLO.LatencyP95MS)
	}
	return nil
}

func LoadServiceSpec(path string) (*ServiceSpec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read service spec %s: %w", path, err)
	}
	var spec ServiceSpec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("failed to parse YAML in %s: %w", path, err)
	}
	if err := spec.Validate(); err != nil {
		return nil, fmt.Errorf("validation error in %s: %w", path, err)
	}
	return &spec, nil
}

func LoadAllServices(dir string) ([]*ServiceSpec, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var services []*ServiceSpec
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := filepath.Ext(entry.Name())
		if ext == ".yaml" || ext == ".yml" {
			spec, err := LoadServiceSpec(filepath.Join(dir, entry.Name()))
			if err != nil {
				return nil, err
			}
			services = append(services, spec)
		}
	}
	return services, nil
}
