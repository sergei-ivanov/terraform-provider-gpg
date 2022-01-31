package gpg_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/invidian/terraform-provider-gpg/gpg"
)

func TestProvider(t *testing.T) {
	t.Parallel()

	provider, ok := gpg.Provider().(*schema.Provider)
	if !ok {
		t.Fatalf("Got unexpected provider type, expected %T, got %T", &schema.Provider{}, gpg.Provider())
	}

	if err := provider.InternalValidate(); err != nil {
		t.Fatalf("validating provider internally: %v", err)
	}
}
