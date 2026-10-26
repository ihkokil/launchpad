package spec

import (
	"testing"
)

func TestServiceSpecValidate(t *testing.T) {
	tests := []struct {
		name    string
		spec    ServiceSpec
		wantErr bool
	}{
		{
			name: "valid standard service",
			spec: ServiceSpec{
				Name:    "orders-api",
				Owner:   "team-checkout",
				Tier:    "standard",
				Runtime: "node",
				Port:    8080,
				Resources: ResourcesSpec{
					CPU:    "100m",
					Memory: "128Mi",
				},
				SLO: SLOSpec{
					Availability: 99.9,
					LatencyP95MS: 250,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid name with capital letters",
			spec: ServiceSpec{
				Name:    "OrdersAPI",
				Owner:   "team-checkout",
				Tier:    "standard",
				Runtime: "node",
				Port:    8080,
				Resources: ResourcesSpec{
					CPU:    "100m",
					Memory: "128Mi",
				},
				SLO: SLOSpec{
					Availability: 99.9,
					LatencyP95MS: 250,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid tier",
			spec: ServiceSpec{
				Name:    "orders-api",
				Owner:   "team-checkout",
				Tier:    "unsupported-tier",
				Runtime: "node",
				Port:    8080,
				Resources: ResourcesSpec{
					CPU:    "100m",
					Memory: "128Mi",
				},
				SLO: SLOSpec{
					Availability: 99.9,
					LatencyP95MS: 250,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid runtime",
			spec: ServiceSpec{
				Name:    "orders-api",
				Owner:   "team-checkout",
				Tier:    "standard",
				Runtime: "ruby",
				Port:    8080,
				Resources: ResourcesSpec{
					CPU:    "100m",
					Memory: "128Mi",
				},
				SLO: SLOSpec{
					Availability: 99.9,
					LatencyP95MS: 250,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid availability percentage",
			spec: ServiceSpec{
				Name:    "orders-api",
				Owner:   "team-checkout",
				Tier:    "standard",
				Runtime: "go",
				Port:    8080,
				Resources: ResourcesSpec{
					CPU:    "100m",
					Memory: "128Mi",
				},
				SLO: SLOSpec{
					Availability: 105.0,
					LatencyP95MS: 250,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.spec.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

