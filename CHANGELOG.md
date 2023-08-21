Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased] - yyyy-mm-dd
 
Here we write upgrading notes for brands. It's a team effort to make them as
straightforward as possible.
 
### Added

- [CA-1125](https://cmcom.atlassian.net/browse/CA-1125)
  Added changes for Magstripe transactions inside the domain

### Changed

### Fixed

- [CA-1154](https://cmcom.atlassian.net/browse/CA-1154)
  Fixed check if CAVV for Visa is a numeric value
- [CA-1156](https://cmcom.atlassian.net/browse/CA-1156)
    ECI 910 invalid for Mastercard


## [v0.0.61] - 2023-07-01
 
### Added

- Added test cards for acceptance tests


## [v0.0.60] - 2023-06-16
 
### Added

- [CA-975](https://cmcom.atlassian.net/browse/CA-975)
  Adding the tests from postman collection to authorization
- [CA-1144](https://cmcom.atlassian.net/browse/CA-1144)
  Download and parse Visa ARDEF table

### Fixed

- [CA-1105](https://cmcom.atlassian.net/browse/CA-1105)
  Fix for returning correct error when decoding CAVV

## [v0.0.59] - 2023-05-15

Here we write upgrading notes for brands. It's a team effort to make them as
straightforward as possible.

### Fixed

- [CA-1120](https://cmcom.atlassian.net/browse/CA-1120)
  Comment out Visa BIN parsing. This blocks starting the authorization repo.

## [v0.0.58] - 2023-05-15

Here we write upgrading notes for brands. It's a team effort to make them as
straightforward as possible.

### Added

- [CA-1080](https://cmcom.atlassian.net/browse/CA-1080)
  Added visa bin recognition

- [CA-833](https://cmcom.atlassian.net/browse/CA-833)
  Add echo endpoints to Swagger docs

- [CA-1081](https://cmcom.atlassian.net/browse/CA-1081)
  Added bin-blocking

- [CA-1106](https://cmcom.atlassian.net/browse/CA-1106)
  Added validation currency check

- [CA-1075](https://cmcom.atlassian.net/browse/CA-1075)
  Added check for acquiring EEA only

### Changed
- [CA-1100](https://cmcom.atlassian.net/browse/CA-1100)
  Removed validation for authorizationType, made it optional for refunds and authorizations.

- [CA-1092](https://cmcom.atlassian.net/browse/CA-1092)
  Removed topic authorization.presentment.cleared.v1 and dangling dependencies

- [CA-244](https://cmcom.atlassian.net/browse/CA-244)
  Moved currency and country to platform and imported new version of platform.

### Fixed

- [CA-1093](https://cmcom.atlassian.net/browse/CA-1093)
  Reversal always Full, even if Partial Capture exists

## [v0.0.57] - 2023-05-02

### Added

- [CA-1067](https://cmcom.atlassian.net/browse/CA-1067)
  Mask 8-digit BINs.

### Changed
- [CA-862](https://cmcom.atlassian.net/browse/CA-862)
  Added new version of platform which contains the single place of mastercard category codes.
  Changed all the dependencies to point to the MCCS on platform.

### Fixed

 - [CA-1089](https://cmcom.atlassian.net/browse/CA-1089)
   Add panic recovery to PubSub subscribers


## [v0.0.56] - 2023-03-30

Fix critical bug in starting the authorization service.

### Fixed
- [CA-1074](https://cmcom.atlassian.net/browse/CA-1074)
  Fix path to BIN ranges and limit card file look back loop

## v0.0.55 - 2023-03-29

Expose echo endpoints for the NOC to monitor and fix showstopper bug in refunds.

### Added
 - [CA-833](https://cmcom.atlassian.net/browse/CA-833)
    Send echo messages to VISA EAS and Mastercard MIP from HA-Proxy and monitoring endpoint (`GET /v1/echo/<scheme>`)
 - [CA-832](https://cmcom.atlassian.net/browse/CA-832)
    Gather throughput metrics and expose them via `GET /v1/metrics`.

### Changed
 - [CA-226](https://cmcom.atlassian.net/browse/CA-226)
    Properly log Tokenization errors; do not return error details to client
 - [CA-190](https://cmcom.atlassian.net/browse/CA-190)
    Only accept API keys that are at least 32 bits
 - [CA-212](https://cmcom.atlassian.net/browse/CA-212)
    Move generic code from Authorization source repository into a shared library

### Fixed
 - [CA-267](https://cmcom.atlassian.net/browse/CA-267)
    Downgrade database migrations
 - [CA-1071](https://cmcom.atlassian.net/browse/CA-1071)
    Publish POS data for clearing


## [v0.0.54] - 2023-03-15

Last pre-1.0 version; baseline for the changelog.

