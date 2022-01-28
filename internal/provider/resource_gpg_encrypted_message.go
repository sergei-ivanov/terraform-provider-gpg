package provider

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

func resourceGPGEncryptedMessage() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "GPG-encrypted message",

		CreateContext: resourceGPGEncryptedMessageCreate,
		// These 2 functions below do nothing, but must be implemented.
		ReadContext:   resourceGPGEncryptedMessageRead,
		DeleteContext: resourceGPGEncryptedMessageDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"content": {
				Description: "Text to be encrypted.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Sensitive:   true,
				StateFunc:   sha256sum,
			},
			"public_keys": {
				Description: "A list of GPG public keys, in ASCII-armored format, which will be used to encrypt the `content`.",
				Type:        schema.TypeList,
				MinItems:    1,
				ForceNew:    true,
				Required:    true,
				Elem: &schema.Schema{
					Description: "GPG public key, in ASCII-armored format.",
					Type:        schema.TypeString,
					ForceNew:    true,
					StateFunc: func(val interface{}) string {
						recipient, err := entityFromString(val.(string))
						if err != nil {
							// We only keep KeyId in state, as we want to keep it small and also
							// we always read public keys anyway. If public key is malformed,
							// creation of resource will fail anyway, so it's fine to set it here.
							return "MALFORMED KEY"
						}

						// Instead of full ASCII-armored key, write only KeyId to state.
						return recipient.PrimaryKey.KeyIdString()
					},
				},
			},
			"result": {
				Description: "GPG-encrypted `content`, in ASCII-armored format.",
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Sensitive:   true,
			},
		},
	}
}

func getRecipients(d *schema.ResourceData) ([]*openpgp.Entity, error) {
	// Store recipients for encryption.
	recipients := []*openpgp.Entity{}

	// Iterate over public keys, decode, parse, collect their IDs and add to recipients list.
	for i, pk := range d.Get("public_keys").([]interface{}) {
		recipient, err := entityFromString(pk.(string))
		if err != nil {
			return nil, fmt.Errorf("decoding public key #%d: %w", i, err)
		}

		recipients = append(recipients, recipient)
	}

	return recipients, nil
}

func savePublicKeys(d *schema.ResourceData, recipients []*openpgp.Entity) error {
	// Store ID of each public key, to store them in state (StateFunc does not work for TypeList for some reason).
	pksIDs := []string{}

	for _, recipient := range recipients {
		pksIDs = append(pksIDs, recipient.PrimaryKey.KeyIdString())
	}

	if err := d.Set("public_keys", pksIDs); err != nil {
		return fmt.Errorf("setting %q property: %w", "public_keys", err)
	}

	return nil
}

func encryptMessage(recipients []*openpgp.Entity, message string, destination io.Writer) error {
	wcEncrypt, err := openpgp.Encrypt(destination, recipients, nil, &openpgp.FileHints{IsBinary: true}, nil)
	if err != nil {
		return fmt.Errorf("encrypting message: %w", err)
	}

	if _, err := io.Copy(wcEncrypt, strings.NewReader(message)); err != nil {
		return fmt.Errorf("writing content to buffer: %w", err)
	}

	if err := wcEncrypt.Close(); err != nil {
		return fmt.Errorf("closing encrypted message: %w", err)
	}

	return nil
}

func encryptAndEncodeMessage(recipients []*openpgp.Entity, message string) (string, error) {
	var buf bytes.Buffer

	// We produce output in ASCII-armor format.
	wcEncode, err := armor.Encode(&buf, "PGP MESSAGE", nil)
	if err != nil {
		return "", fmt.Errorf("encoding message: %w", err)
	}

	if err := encryptMessage(recipients, message, wcEncode); err != nil {
		return "", fmt.Errorf("encrypting message: %w", err)
	}

	if err := wcEncode.Close(); err != nil {
		return "", fmt.Errorf("closing encoded message: %w", err)
	}

	return buf.String(), nil
}

func resourceGPGEncryptedMessageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	recipients, err := getRecipients(d)
	if err != nil {
		return diag.Errorf("getting recipients: %s", err)
	}

	if err := savePublicKeys(d, recipients); err != nil {
		return diag.Errorf("saving public keys: %s", err)
	}

	encryptedMessage, err := encryptAndEncodeMessage(recipients, d.Get("content").(string))
	if err != nil {
		return diag.Errorf("encrypting message: %s", err)
	}

	if err := d.Set("result", encryptedMessage); err != nil {
		return diag.Errorf("setting %q property: %s", "result", err)
	}

	// Calculate SHA-256 checksum of message for ID.
	d.SetId(sha256sum(encryptedMessage))

	return nil
}

func resourceGPGEncryptedMessageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceGPGEncryptedMessageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")

	return nil
}

func entityFromString(key string) (*openpgp.Entity, error) {
	block, err := armor.Decode(strings.NewReader(key))
	if err != nil {
		return nil, fmt.Errorf("decoding public key: %w", err)
	}

	recipient, err := openpgp.ReadEntity(packet.NewReader(block.Body))
	if err != nil {
		return nil, fmt.Errorf("parsing public key: %w", err)
	}

	return recipient, nil
}

func sha256sum(data interface{}) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(data.(string))))
}