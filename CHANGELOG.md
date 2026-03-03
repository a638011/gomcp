# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Repository template enhancements
- GitHub Actions CI/CD workflows
- Comprehensive documentation
- Issue and PR templates
- Contributing guidelines
- Security policy

## [0.1.0] - 2024-12-01

### Added
- Initial release of gomcp
- Complete MCP 2025-06-18 specification support
- **Tools**: 4 example tools with structured outputs
  - `multiply_numbers` - Number multiplication with structured output
  - `code_review` - Generate code review analysis
  - `logo` - Display server logo
  - `long_operation` - Demonstrates progress & cancellation
- **Prompts**: 3 reusable prompt templates
  - `code_review` - Code review template
  - `git_commit` - Git commit message generator
  - `debug_help` - Debugging assistance
- **Resources**: 6 static/dynamic resources
  - `project://info` - Server information
  - `project://status` - System status
  - `docs://quickstart` - Quick start guide
  - `docs://api-reference` - API documentation
  - `config://template` - Configuration template
  - `config://env-vars` - Environment variables
- **Transport Support**
  - HTTP/SSE transport for web clients
  - stdio transport for Claude Desktop
  - Cursor IDE compatibility mode
- **MCP Features**
  - Roots - Filesystem root definitions
  - Completion - Structured tool outputs with JSON schemas
  - Logging - Server→client log notifications (8 levels)
  - Pagination - Cursor-based pagination (max 100/page)
  - Sampling - Server→client LLM requests
  - Elicitation - Server→user data requests
  - Progress - Real-time progress notifications
  - Cancellation - Request cancellation support
  - Ping - Health checks
- **Infrastructure**
  - Docker and Podman support
  - Docker Compose configuration
  - Makefile with comprehensive targets
  - Hot reload development mode (air)
  - Production build with version embedding
- **Security**
  - OAuth authentication support
  - Session management
  - CORS configuration
  - SSL/TLS support

### Security
- Non-root Docker container execution
- Environment-based secret management

---

## Version History

| Version | Date | Highlights |
|---------|------|------------|
| 0.1.0 | 2024-12-01 | Initial release with full MCP 2025-06-18 support |

[Unreleased]: https://github.com/NP-compete/gomcp/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/NP-compete/gomcp/releases/tag/v0.1.0
