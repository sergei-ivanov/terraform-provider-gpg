package gpg

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown
}

func New() func() *schema.Provider {
	return func() *schema.Provider {
		return &schema.Provider{
			DataSourcesMap: map[string]*schema.Resource{
				// None at this point
			},
			ResourcesMap: map[string]*schema.Resource{
				"gpg_encrypted_message": resourceGPGEncryptedMessage(),
			},
		}
	}
}
