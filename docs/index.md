# GPG Command Provider

The GPG provider allows to generate GPG encrypted message in ASCII-armored format using Terraform.

This provider uses built-in [openpgp](https://godoc.org/golang.org/x/crypto/openpgp) Golang library to perform GPG encryption. Currently the only supported option is encrypting message with public keys.

Managing GPG keyring or signing files is currently not implemented.

## Example Usage

```hcl
terraform {
  required_providers {
    gpg = {
      source  = "invidian/gpg"
      version = "0.3.0"
    }
  }
}

resource "gpg_encrypted_message" "example" {
  content     = "This is example of GPG encrypted message."
  public_keys = [
    var.gpg_public_key,
  ]
}

output "gpg_encrypted_message" {
  value = gpg_encrypted_message.example.result
}

variable "gpg_public_key" {
  default = <<EOF
-----BEGIN PGP PUBLIC KEY BLOCK-----

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
5hwt9fmwTMhEPu5m
=iJyi
-----END PGP PUBLIC KEY BLOCK-----
EOF
}
```

## Argument Reference

This provider currently takes no arguments.
