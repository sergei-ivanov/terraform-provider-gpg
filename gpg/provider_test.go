package gpg

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	// Provider name for single configuration testing
	ProviderName = "gpg"
)

var testProviderFactories map[string]func() (*schema.Provider, error)

func init() {
	testProviderFactories = map[string]func() (*schema.Provider, error){
		ProviderName: func() (*schema.Provider, error) {
			return Provider(), nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("validating provider internally: %v", err)
	}
}

func TestProvider_HasResources(t *testing.T) {
	expectedResources := []string{
		"gpg_encrypted_message",
	}

	resources := Provider().ResourcesMap
	if len(expectedResources) != len(resources) {
		t.Errorf("There are an unexpected number of registered resources. Expected %v got %v", len(expectedResources), len(resources))
	}

	for _, resource := range expectedResources {
		if _, ok := resources[resource]; !ok {
			t.Errorf("An expected resource was not registered")
		}
		if resources[resource] == nil {
			t.Errorf("A resource cannot have a nil schema")
		}
	}
}

func TestProvider_HasDataSources(t *testing.T) {
	expectedDataSources := []string{}

	dataSources := Provider().DataSourcesMap
	if len(expectedDataSources) != len(dataSources) {
		t.Errorf("There are an unexpected number of registered data sources. Expected %v got %v", len(expectedDataSources), len(dataSources))
	}

	for _, resource := range expectedDataSources {
		if _, ok := dataSources[resource]; !ok {
			t.Errorf("An expected data source was not registered")
		}
		if dataSources[resource] == nil {
			t.Errorf("A data source cannot have a nil schema")
		}
	}
}
