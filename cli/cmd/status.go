package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var (
	statusNamespace string

	statusCmd = &cobra.Command{
		Use:   "status <service-name>",
		Short: "Display runtime health, pod distribution, and SLO metrics for a service",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]
			fmt.Printf("🔍 Checking runtime status for service '%s' in namespace '%s'...\n\n", serviceName, statusNamespace)

			// 1. Check Deployments
			deployOut, err := exec.Command("kubectl", "get", "deployment", serviceName, "-n", statusNamespace, "-o", "wide").CombinedOutput()
			if err != nil {
				fmt.Printf("⚠️  Deployment '%s' not found or cluster offline:\n%s\n", serviceName, string(deployOut))
			} else {
				fmt.Println("📦 Deployment Status:")
				fmt.Println(strings.TrimSpace(string(deployOut)))
				fmt.Println()
			}

			// 2. Check Pods with topology
			podsOut, err := exec.Command("kubectl", "get", "pods", "-l", fmt.Sprintf("app.kubernetes.io/name=%s", serviceName), "-n", statusNamespace, "-o", "wide").CombinedOutput()
			if err == nil && len(podsOut) > 0 {
				fmt.Println("🚀 Running Pod Instances (Multi-Node Distribution):")
				fmt.Println(strings.TrimSpace(string(podsOut)))
				fmt.Println()
			}

			// 3. Check HPA
			hpaOut, err := exec.Command("kubectl", "get", "hpa", serviceName, "-n", statusNamespace).CombinedOutput()
			if err == nil && len(hpaOut) > 0 {
				fmt.Println("📈 Horizontal Pod Autoscaler:")
				fmt.Println(strings.TrimSpace(string(hpaOut)))
				fmt.Println()
			}

			// 4. Check Sloth PrometheusServiceLevel
			sloOut, err := exec.Command("kubectl", "get", "prometheusservicelevel", fmt.Sprintf("%s-slo", serviceName), "-n", statusNamespace).CombinedOutput()
			if err == nil && len(sloOut) > 0 {
				fmt.Println("🎯 SLO Health & Error Budget Status:")
				fmt.Println(strings.TrimSpace(string(sloOut)))
			}

			return nil
		},
	}
)

func init() {
	statusCmd.Flags().StringVarP(&statusNamespace, "namespace", "n", "default", "Kubernetes namespace")
}
