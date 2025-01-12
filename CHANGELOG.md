# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.3.1] - 2022-02-03
### Added
- Added CI using GitHub Actions.

### Changed
- Go 1.17 is now used for building.
- Exlucuded `darwin/arm`, `freebsd/arm64`, `windows/arm` and `windows/arm64` platforms from
release builds to speed up the release process, as they are not officially supported by
either Terraform itself or Terraform Registry.

### Fixed
- Fixed all issues found by linter.

## [0.3.0] - 2020-08-19
### Added
- This provider is now available via [Terraform Registry](https://registry.terraform.io/providers/invidian/gpg/latest).
- Added basic unit tests.

### Changed
- Migrated to use Terraform Plugin SDK.
- Changelog is now published in [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) format.
- Updated all dependencies to latest versions.

### Fixed
- Fixed `gpg_encrypted_message` resource destroying.
- Fixed found linter warnings.

## [0.2.1] - 2019-06-20
### Changed
- Sensitive fields now use `sensitive: true`, so they do not leak into Terraform plan.
- Resource `gpg_encrypted_message` now use SHA256 of message content as resource ID.

### Fixed
- Fixed `gpg_encrypted_message` resource update issues with Terraform 0.12.x.

### Removed
- Removed use of `SchemaStateFunc` for result, as it has no effect with Terraform 0.12.x.

## [0.2.0] - 2019-06-15
### Added
- Added Terraform 0.12 compatibility.

### Changed
- Changed code to standard Terraform provider layout.
- Added `Makefile` to document common tasks.

## [0.1.0] - 2019-04-12
### Added
- Initial release

[0.3.1]: https://github.com/invidian/terraform-provider-gpg/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/invidian/terraform-provider-gpg/compare/v0.2.1...v0.3.0
[0.2.1]: https://github.com/invidian/terraform-provider-gpg/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/invidian/terraform-provider-gpg/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/invidian/terraform-provider-gpg/releases/tag/v0.1.0
