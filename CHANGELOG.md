# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v1.0.0] - 2026-02-07

### Added
- Initial release of SnowRedis - a distributed ID generator based on the Snowflake algorithm with Redis coordination
- Support for Redis-based coordination to ensure unique ID generation across multiple nodes
- Builder pattern for easy configuration of RedisSnowflake instances
- Support for strict mode with Redis assistance to prevent duplicate IDs
- Automatic ID allocation through Redis for datacenter and worker IDs
- Comprehensive test suite covering concurrent ID generation, performance, and Redis coordination

### Changed
- Renamed project from "Snowflake-Redis" to "SnowRedis"
- Standardized constants naming convention to uppercase
- Unified comment style using multi-line block comments with JavaDoc-style @param and @return tags
- Updated to use go-redis v9 (github.com/redis/go-redis/v9)
- Simplified project structure by removing internal directory layer
- Updated module path to github.com/sunquakes/snowredis
- Lowered Go version requirement from 1.25.3 to 1.19
- Refactored Redis client interface to use simpler Client name instead of RedisClient
- Migrated from SetNX/Incr/Del methods to standard Redis client interface

### Fixed
- Various linting issues for better code quality and consistency
- Clock rollback protection in ID generation algorithm
- Memory usage optimizations in ID generation
- Documentation inconsistencies between English and Chinese versions

### Removed
- Deprecated methods and configurations
- Old mutex-based generation approach (replaced with Redis coordination)
- Internal directory layer for simplified structure
- Unused parameters and functions

### Security
- Added Apache License 2.0