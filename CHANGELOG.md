# Changelog
All notable changes to this project will be documented in this file.

## [Unreleased]

### Added

### Changed

* [#15](https://github.com/bitgrip/cattlectl/issues/15) `--values` can be used multiple times.
  * all values files are merged
  * later flags have higher precedence

### Removed

### Fixed

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

[Unreleased]: https://github.com/bitgrip/cattlectl/compare/v1.1.0...HEAD
[1.1.0]: https://github.com/bitgrip/cattlectl/compare/v1.0.0...v1.1.0
