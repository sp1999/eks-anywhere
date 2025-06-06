package v1alpha1_test

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/aws/eks-anywhere/pkg/api/v1alpha1"
)

func TestValidateCreateOIDCConfigSuccess(t *testing.T) {
	ctx := context.Background()
	c := oidcConfig()
	c.Spec.ClientId = "test"
	c.Spec.IssuerUrl = "https://test.com"
	o := NewWithT(t)

	o.Expect(c.ValidateCreate(ctx, &c)).Error().To(Succeed())
}

func TestClusterValidateCreateInvalidOIDCConfig(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name   string
		config v1alpha1.OIDCConfig
		err    string
	}{
		{
			name: "No clientID",
			config: v1alpha1.OIDCConfig{
				Spec: v1alpha1.OIDCConfigSpec{
					ClientId: "",
				},
			},
			err: "clientId is required",
		},
		{
			name: "Null issuerID",
			config: v1alpha1.OIDCConfig{
				Spec: v1alpha1.OIDCConfigSpec{
					ClientId:  "test",
					IssuerUrl: "",
				},
			},
			err: "issuerUrl is required",
		},
		{
			name: "Invalid issuer url",
			config: v1alpha1.OIDCConfig{
				Spec: v1alpha1.OIDCConfigSpec{
					ClientId:  "test",
					IssuerUrl: "invalid-url",
				},
			},
			err: "invalid URI for request",
		},
		{
			name: "Issuer url, non https",
			config: v1alpha1.OIDCConfig{
				Spec: v1alpha1.OIDCConfigSpec{
					ClientId:  "test",
					IssuerUrl: "http://test.com",
				},
			},
			err: "issuerUrl should have HTTPS scheme",
		},
		{
			name: "Extra required claims",
			config: v1alpha1.OIDCConfig{
				Spec: v1alpha1.OIDCConfigSpec{
					ClientId:  "test",
					IssuerUrl: "https://test.com",
					RequiredClaims: []v1alpha1.OIDCConfigRequiredClaim{
						{
							Claim: "claim1",
							Value: "val1",
						},
						{
							Claim: "claim2",
							Value: "val2",
						},
					},
				},
			},
			err: "only one OIDConfig requiredClaim is supported at this time",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(tt.config.ValidateCreate(ctx, &tt.config)).Error().To(MatchError(ContainSubstring(tt.err)))
		})
	}
}

func TestValidateUpdateOIDCClientIdMgmtCluster(t *testing.T) {
	ctx := context.Background()
	ocOld := oidcConfig()
	ocOld.Spec.ClientId = "test"
	c := ocOld.DeepCopy()

	c.Spec.ClientId = "test2"
	o := NewWithT(t)
	o.Expect(c.ValidateUpdate(ctx, &ocOld, c)).Error().To(MatchError(ContainSubstring("OIDCConfig: Forbidden: config is immutable")))
}

func TestValidateUpdateOIDCGroupsClaimMgmtCluster(t *testing.T) {
	ctx := context.Background()
	ocOld := oidcConfig()
	ocOld.Spec.GroupsClaim = "test"
	c := ocOld.DeepCopy()

	c.Spec.GroupsClaim = "test2"
	o := NewWithT(t)
	o.Expect(c.ValidateUpdate(ctx, &ocOld, c)).Error().To(MatchError(ContainSubstring("OIDCConfig: Forbidden: config is immutable")))
}

func TestValidateUpdateOIDCGroupsPrefixMgmtCluster(t *testing.T) {
	ctx := context.Background()
	ocOld := oidcConfig()
	ocOld.Spec.GroupsPrefix = "test"
	c := ocOld.DeepCopy()

	c.Spec.GroupsPrefix = "test2"
	o := NewWithT(t)
	o.Expect(c.ValidateUpdate(ctx, &ocOld, c)).Error().To(MatchError(ContainSubstring("OIDCConfig: Forbidden: config is immutable")))
}

func TestValidateUpdateOIDCIssuerUrlMgmtCluster(t *testing.T) {
	ctx := context.Background()
	ocOld := oidcConfig()
	ocOld.Spec.IssuerUrl = "test"
	c := ocOld.DeepCopy()

	c.Spec.IssuerUrl = "test2"
	o := NewWithT(t)
	o.Expect(c.ValidateUpdate(ctx, &ocOld, c)).Error().To(MatchError(ContainSubstring("OIDCConfig: Forbidden: config is immutable")))
}

func TestValidateUpdateOIDCUsernameClaimMgmtCluster(t *testing.T) {
	ctx := context.Background()
	ocOld := oidcConfig()
	ocOld.Spec.UsernameClaim = "test"
	c := ocOld.DeepCopy()

	c.Spec.UsernameClaim = "test2"
	o := NewWithT(t)
	o.Expect(c.ValidateUpdate(ctx, &ocOld, c)).Error().To(MatchError(ContainSubstring("OIDCConfig: Forbidden: config is immutable")))
}

func TestValidateUpdateOIDCUsernamePrefixMgmtCluster(t *testing.T) {
	ctx := context.Background()
	ocOld := oidcConfig()
	ocOld.Spec.UsernamePrefix = "test"
	c := ocOld.DeepCopy()

	c.Spec.UsernamePrefix = "test2"
	o := NewWithT(t)
	o.Expect(c.ValidateUpdate(ctx, &ocOld, c)).Error().To(MatchError(ContainSubstring("OIDCConfig: Forbidden: config is immutable")))
}

func TestValidateUpdateOIDCRequiredClaimsMgmtCluster(t *testing.T) {
	ctx := context.Background()
	ocOld := oidcConfig()
	ocOld.Spec.RequiredClaims = []v1alpha1.OIDCConfigRequiredClaim{{Claim: "test", Value: "value"}}
	c := ocOld.DeepCopy()

	c.Spec.RequiredClaims = []v1alpha1.OIDCConfigRequiredClaim{{Claim: "test", Value: "value2"}}
	o := NewWithT(t)
	o.Expect(c.ValidateUpdate(ctx, &ocOld, c)).Error().To(MatchError(ContainSubstring("OIDCConfig: Forbidden: config is immutable")))
}

func TestValidateUpdateOIDCRequiredClaimsMultipleMgmtCluster(t *testing.T) {
	ctx := context.Background()
	ocOld := oidcConfig()
	ocOld.Spec.RequiredClaims = []v1alpha1.OIDCConfigRequiredClaim{{Claim: "test", Value: "value"}}
	c := ocOld.DeepCopy()

	c.Spec.RequiredClaims = append(c.Spec.RequiredClaims, v1alpha1.OIDCConfigRequiredClaim{
		Claim: "test2",
		Value: "value2",
	})
	o := NewWithT(t)
	o.Expect(c.ValidateUpdate(ctx, &ocOld, c)).Error().To(MatchError(ContainSubstring("OIDCConfig: Forbidden: config is immutable")))
}

func TestClusterValidateUpdateOIDCclientIdMutableUpdateNameWorkloadCluster(t *testing.T) {
	ctx := context.Background()
	ocOld := oidcConfig()
	ocOld.Spec.ClientId = "test"
	ocOld.SetManagedBy("test")
	c := ocOld.DeepCopy()

	c.Spec.ClientId = "test2"
	o := NewWithT(t)
	o.Expect(c.ValidateUpdate(ctx, &ocOld, c)).Error().To(Succeed())
}

func TestValidateUpdateOIDCClientIdWorkloadCluster(t *testing.T) {
	ctx := context.Background()
	ocOld := oidcConfig()
	ocOld.Spec.ClientId = "test"
	ocOld.SetManagedBy("test")

	c := ocOld.DeepCopy()

	c.Spec.ClientId = "test2"
	o := NewWithT(t)
	o.Expect(c.ValidateUpdate(ctx, &ocOld, c)).Error().To(Succeed())
}

func TestValidateUpdateOIDCGroupsClaimWorkloadCluster(t *testing.T) {
	ctx := context.Background()
	ocOld := oidcConfig()
	ocOld.Spec.GroupsClaim = "test"
	ocOld.SetManagedBy("test")

	c := ocOld.DeepCopy()

	c.Spec.GroupsClaim = "test2"
	o := NewWithT(t)
	o.Expect(c.ValidateUpdate(ctx, &ocOld, c)).Error().To(Succeed())
}

func TestValidateUpdateOIDCGroupsPrefixWorkloadCluster(t *testing.T) {
	ctx := context.Background()
	ocOld := oidcConfig()
	ocOld.Spec.GroupsPrefix = "test"
	ocOld.SetManagedBy("test")

	c := ocOld.DeepCopy()

	c.Spec.GroupsPrefix = "test2"
	o := NewWithT(t)
	o.Expect(c.ValidateUpdate(ctx, &ocOld, c)).Error().To(Succeed())
}

func TestValidateUpdateOIDCIssuerUrlWorkloadCluster(t *testing.T) {
	ctx := context.Background()
	ocOld := oidcConfig()
	ocOld.Spec.IssuerUrl = "test"
	ocOld.SetManagedBy("test")

	c := ocOld.DeepCopy()

	c.Spec.IssuerUrl = "test2"
	o := NewWithT(t)
	o.Expect(c.ValidateUpdate(ctx, &ocOld, c)).Error().To(Succeed())
}

func TestValidateUpdateOIDCUsernameClaimWorkloadCluster(t *testing.T) {
	ctx := context.Background()
	ocOld := oidcConfig()
	ocOld.Spec.UsernameClaim = "test"
	ocOld.SetManagedBy("test")

	c := ocOld.DeepCopy()

	c.Spec.UsernameClaim = "test2"
	o := NewWithT(t)
	o.Expect(c.ValidateUpdate(ctx, &ocOld, c)).Error().To(Succeed())
}

func TestValidateUpdateOIDCUsernamePrefixWorkloadCluster(t *testing.T) {
	ctx := context.Background()
	ocOld := oidcConfig()
	ocOld.Spec.UsernamePrefix = "test"
	ocOld.SetManagedBy("test")

	c := ocOld.DeepCopy()

	c.Spec.UsernamePrefix = "test2"
	o := NewWithT(t)
	o.Expect(c.ValidateUpdate(ctx, &ocOld, c)).Error().To(Succeed())
}

func TestValidateUpdateOIDCRequiredClaimsWorkloadCluster(t *testing.T) {
	ctx := context.Background()
	ocOld := oidcConfig()
	ocOld.Spec.RequiredClaims = []v1alpha1.OIDCConfigRequiredClaim{{Claim: "test", Value: "value"}}
	ocOld.SetManagedBy("test")

	c := ocOld.DeepCopy()

	c.Spec.RequiredClaims = []v1alpha1.OIDCConfigRequiredClaim{{Claim: "test", Value: "value2"}}
	o := NewWithT(t)
	o.Expect(c.ValidateUpdate(ctx, &ocOld, c)).Error().To(Succeed())
}

func TestValidateUpdateOIDCRequiredClaimsMultipleWorkloadCluster(t *testing.T) {
	ctx := context.Background()
	ocOld := oidcConfig()
	ocOld.Spec.RequiredClaims = []v1alpha1.OIDCConfigRequiredClaim{{Claim: "test", Value: "value"}}
	ocOld.SetManagedBy("test")

	c := ocOld.DeepCopy()

	c.Spec.RequiredClaims = append(c.Spec.RequiredClaims, v1alpha1.OIDCConfigRequiredClaim{
		Claim: "test2",
		Value: "value2",
	})
	o := NewWithT(t)
	o.Expect(c.ValidateUpdate(ctx, &ocOld, c)).Error().To(Succeed())
}

func oidcConfig() v1alpha1.OIDCConfig {
	return v1alpha1.OIDCConfig{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{Annotations: make(map[string]string, 1)},
		Spec:       v1alpha1.OIDCConfigSpec{},
		Status:     v1alpha1.OIDCConfigStatus{},
	}
}

func TestOIDCConfigValidateCreateCastFail(t *testing.T) {
	g := NewWithT(t)

	// Create a different type that will cause the cast to fail
	wrongType := &v1alpha1.Cluster{}

	// Create the config object that implements CustomValidator
	config := &v1alpha1.OIDCConfig{}

	// Call ValidateCreate with the wrong type
	warnings, err := config.ValidateCreate(context.TODO(), wrongType)

	// Verify that an error is returned
	g.Expect(warnings).To(BeNil())
	g.Expect(err).To(HaveOccurred())
	g.Expect(err.Error()).To(ContainSubstring("expected an OIDCConfig"))
}

func TestOIDCConfigValidateUpdateCastFail(t *testing.T) {
	g := NewWithT(t)

	// Create a different type that will cause the cast to fail
	wrongType := &v1alpha1.Cluster{}

	// Create the config object that implements CustomValidator
	config := &v1alpha1.OIDCConfig{}

	// Call ValidateUpdate with the wrong type
	warnings, err := config.ValidateUpdate(context.TODO(), &v1alpha1.OIDCConfig{}, wrongType)

	// Verify that an error is returned
	g.Expect(warnings).To(BeNil())
	g.Expect(err).To(HaveOccurred())
	g.Expect(err.Error()).To(ContainSubstring("expected an OIDCConfig"))
}

func TestOIDCConfigValidateDeleteCastFail(t *testing.T) {
	g := NewWithT(t)

	// Create a different type that will cause the cast to fail
	wrongType := &v1alpha1.Cluster{}

	// Create the config object that implements CustomValidator
	config := &v1alpha1.OIDCConfig{}

	// Call ValidateDelete with the wrong type
	warnings, err := config.ValidateDelete(context.TODO(), wrongType)

	// Verify that an error is returned
	g.Expect(warnings).To(BeNil())
	g.Expect(err).To(HaveOccurred())
	g.Expect(err.Error()).To(ContainSubstring("expected an OIDCConfig"))
}
