package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ihkokil/launchpad/cli/pkg/generator"
	"github.com/ihkokil/launchpad/cli/pkg/spec"
	"github.com/spf13/cobra"
)

var (
	renderServicesDir string
	renderOutputDir   string

	renderCmd = &cobra.Command{
		Use:   "render",
		Short: "Compile developer service specifications into production GitOps manifests",
		Long: `Scans all service contracts in services/*.yaml, validates them against platform schemas, 
and renders Flux HelmRelease and Sloth SLO manifests into gitops/clusters/local/apps/.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			services, err := spec.LoadAllServices(renderServicesDir)
			if err != nil {
				return fmt.Errorf("error loading services from %s: %w", renderServicesDir, err)
			}
			if len(services) == 0 {
				fmt.Printf("No service manifests found in %s\n", renderServicesDir)
				return nil
			}

			fmt.Printf("🔍 Discovered %d service specification(s):\n", len(services))
			var appPaths []string

			for _, s := range services {
				targetDir := filepath.Join(renderOutputDir, s.Name)
				if err := os.MkdirAll(targetDir, 0755); err != nil {
					return fmt.Errorf("failed to create directory %s: %w", targetDir, err)
				}

				// 1. Render HelmRelease
				hrData, err := generator.RenderHelmRelease(s)
				if err != nil {
					return err
				}
				hrPath := filepath.Join(targetDir, "helmrelease.yaml")
				if err := os.WriteFile(hrPath, hrData, 0644); err != nil {
					return fmt.Errorf("failed to write %s: %w", hrPath, err)
				}

				// 2. Render Sloth SLO specification
				sloData, err := generator.RenderSlothSLO(s)
				if err != nil {
					return err
				}
				sloPath := filepath.Join(targetDir, "slo.yaml")
				if err := os.WriteFile(sloPath, sloData, 0644); err != nil {
					return fmt.Errorf("failed to write %s: %w", sloPath, err)
				}

				// 3. Render app kustomization.yaml
				kustData := fmt.Sprintf("apiVersion: kustomize.config.k8s.io/v1beta1\nkind: Kustomization\nresources:\n  - helmrelease.yaml\n  - slo.yaml\n")
				kustPath := filepath.Join(targetDir, "kustomization.yaml")
				if err := os.WriteFile(kustPath, []byte(kustData), 0644); err != nil {
					return fmt.Errorf("failed to write %s: %w", kustPath, err)
				}

				fmt.Printf("  ✔ [%s] -> %s (HelmRelease + Sloth SLO)\n", s.Name, targetDir)
				appPaths = append(appPaths, s.Name)
			}

			// Generate top-level apps kustomization.yaml
			topKustHeader := "apiVersion: kustomize.config.k8s.io/v1beta1\nkind: Kustomization\nresources:\n"
			for _, p := range appPaths {
				topKustHeader += fmt.Sprintf("  - %s\n", p)
			}
			topKustPath := filepath.Join(renderOutputDir, "kustomization.yaml")
			if err := os.WriteFile(topKustPath, []byte(topKustHeader), 0644); err != nil {
				return fmt.Errorf("failed to write top-level kustomization %s: %w", topKustPath, err)
			}

			fmt.Printf("\n✨ Successfully rendered all GitOps artifacts in %s\n", renderOutputDir)
			return nil
		},
	}
)

func init() {
	renderCmd.Flags().StringVarP(&renderServicesDir, "services-dir", "s", "services", "Directory containing service specifications")
	renderCmd.Flags().StringVarP(&renderOutputDir, "out-dir", "o", "gitops/clusters/local/apps", "Output directory for GitOps manifests")
}
