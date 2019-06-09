# Changelog
All notable changes to this project will be documented in this file.

## [Unreleased]

### Added

### Changed

### Removed

### Fixed

## [1.2.0]

### Added

* [#32](https://github.com/bitgrip/cattlectl/issues/32) Support `values_yaml` field as alternative to the `answers` field of apps.
* [#36](https://github.com/bitgrip/cattlectl/issues/36) Add template method `readTemplate` to include templates into descriptors

### Changed

* [#24](https://github.com/bitgrip/cattlectl/issues/24) Support cluster and project catalogs for apps
* [#38](https://github.com/bitgrip/cattlectl/pull/38) Verify API version to be compatible with current cattlectl version

### Removed

### Fixed

## [1.1.1]

### Added

* [#26](https://github.com/bitgrip/cattlectl/issues/26) bash completion
* [#28](https://github.com/bitgrip/cattlectl/issues/28) multi file includes
  * directory includes
  * pattern includes

### Changed

* [#15](https://github.com/bitgrip/cattlectl/issues/15) `--values` can be used multiple times.
  * all values files are merged
  * later flags have higher precedence

### Removed

### Fixed

* [#30](https://github.com/bitgrip/cattlectl/issues/30) Fix apply of workload.

## [1.1.0]

### Added

* [#13](https://github.com/bitgrip/cattlectl/issues/13) Add support for project workloads
* [#7](https://github.com/bitgrip/cattlectl/issues/7) Add support for descriptor includes
* [#14](https://github.com/bitgrip/cattlectl/issues/14) Feature/add to yaml function in template

### Changed

* [#18](https://github.com/bitgrip/cattlectl/issues/18) Add support to update project resources
* [#22](https://github.com/bitgrip/cattlectl/issues/22) Add support for self signed certs

### Removed

### Fixed

* Certificates do now respect namespace property
* DockerCredentials do now respect namespace property
* [#21](https://github.com/bitgrip/cattlectl/pull/21) fix readme typo

## 0.1.0 (Feb 15 2019)

* Initial release

[Unreleased]: https://github.com/bitgrip/cattlectl/compare/v1.2.0...HEAD
[1.2.0]: https://github.com/bitgrip/cattlectl/compare/v1.1.1...v1.2.0
[1.1.1]: https://github.com/bitgrip/cattlectl/compare/v1.1.0...v1.1.1
[1.1.0]: https://github.com/bitgrip/cattlectl/compare/v1.0.0...v1.1.0
