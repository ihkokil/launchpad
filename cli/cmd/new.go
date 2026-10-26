package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ihkokil/launchpad/cli/pkg/spec"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	newRuntime string
	newOwner   string
	newTier    string
	newPort    int
	newOutDir  string

	newCmd = &cobra.Command{
		Use:   "new <service-name>",
		Short: "Scaffold a new microservice specification and starter repository",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]
			svc := &spec.ServiceSpec{
				Name:    serviceName,
				Owner:   newOwner,
				Tier:    newTier,
				Runtime: newRuntime,
				Port:    newPort,
				Env: map[string]string{
					"NODE_ENV": "production",
				},
				Resources: spec.ResourcesSpec{
					CPU:    "100m",
					Memory: "128Mi",
				},
				SLO: spec.SLOSpec{
					Availability: 99.9,
					LatencyP95MS: 250,
				},
				Image: spec.ImageSpec{
					Repository: serviceName,
					Tag:        "v0.1.0",
				},
			}

			if err := svc.Validate(); err != nil {
				return fmt.Errorf("invalid service specification: %w", err)
			}

			servicesDir := "services"
			if newOutDir != "" {
				servicesDir = newOutDir
			}
			if err := os.MkdirAll(servicesDir, 0755); err != nil {
				return fmt.Errorf("failed to create services directory: %w", err)
			}

			targetFile := filepath.Join(servicesDir, fmt.Sprintf("%s.yaml", serviceName))
			data, err := yaml.Marshal(svc)
			if err != nil {
				return fmt.Errorf("failed to serialize service YAML: %w", err)
			}

			if err := os.WriteFile(targetFile, data, 0644); err != nil {
				return fmt.Errorf("failed to write %s: %w", targetFile, err)
			}

			fmt.Printf("✔ Successfully scaffolded service contract: %s\n", targetFile)
			fmt.Printf("  • Runtime: %s\n", newRuntime)
			fmt.Printf("  • Owner:   %s\n", newOwner)
			fmt.Printf("  • Tier:    %s (HA %s replicas)\n", newTier, map[string]string{"critical": "3+", "standard": "2"}[newTier])
			fmt.Printf("  • SLO:     %.1f%% availability, p95 < %dms\n", svc.SLO.Availability, svc.SLO.LatencyP95MS)
			fmt.Printf("\nNext step: Run 'launchpad render' to generate GitOps manifests.\n")
			return nil
		},
	}
)

func init() {
	newCmd.Flags().StringVarP(&newRuntime, "runtime", "r", "node", "Service runtime (node, go, php)")
	newCmd.Flags().StringVarP(&newOwner, "owner", "o", "platform-team", "Team owning the service")
	newCmd.Flags().StringVarP(&newTier, "tier", "t", "standard", "Service tier (standard, critical)")
	newCmd.Flags().IntVarP(&newPort, "port", "p", 8080, "Application port")
	newCmd.Flags().StringVar(&newOutDir, "out-dir", "services", "Directory to save service manifest")
}
