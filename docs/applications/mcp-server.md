# MCP Server

The Model Context Protocol (MCP) server enables AI assistants to perform web searches through OmniSerp.

## Installation

```bash
go install github.com/plexusone/omniserp/cmd/mcp-omniserp@latest
```

Or build from source:

```bash
go build ./cmd/mcp-omniserp
```

## Configuration

Add to your Claude Desktop config (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "omniserp": {
      "command": "mcp-omniserp",
      "env": {
        "SERPER_API_KEY": "your_serper_api_key",
        "SEARCH_ENGINE": "serper"
      }
    }
  }
}
```

## Available Tools

The MCP server **dynamically registers only the tools supported by the current search engine backend**:

- When using **Serper**, all 12 tools are available including Lens search
- When using **SerpAPI**, 11 tools are available (Lens is excluded)

### Tool Categories

| Tool | Description | Serper | SerpAPI |
|------|-------------|:------:|:-------:|
| `google_search` | General web searches | ✓ | ✓ |
| `google_search_news` | Search news articles | ✓ | ✓ |
| `google_search_images` | Search for images | ✓ | ✓ |
| `google_search_videos` | Search for videos | ✓ | ✓ |
| `google_search_places` | Search for locations and businesses | ✓ | ✓ |
| `google_search_maps` | Search maps data | ✓ | ✓ |
| `google_search_reviews` | Search reviews | ✓ | ✓ |
| `google_search_shopping` | Search shopping/product listings | ✓ | ✓ |
| `google_search_scholar` | Search academic papers | ✓ | ✓ |
| `google_search_lens` | Visual search capabilities | ✓ | ✗ |
| `google_search_autocomplete` | Get search suggestions | ✓ | ✓ |
| `webpage_scrape` | Extract content from webpages | ✓ | ✓ |

All searches support parameters like location, language, country, and number of results.

## Server Logs

The MCP server logs which tools were registered and which were skipped:

```
2025/12/13 19:00:00 Using engine: serpapi v1.0.0
2025/12/13 19:00:00 Registered 11 tools: [google_search, google_search_news, ...]
2025/12/13 19:00:00 Skipped 1 unsupported tools: [google_search_lens]
```

## Secure Mode (Optional)

The MCP server supports optional secure credential management using VaultGuard. When a policy file exists, API keys are retrieved from the OS keychain instead of environment variables.

### Setup for Secure Mode

1. **Store your API key in the keychain:**

    ```bash
    security add-generic-password -s "omnivault" -a "SERPER_API_KEY" -w "your-key"
    ```

2. **Create a security policy** (`~/.vaultguard/policy.json`):

    ```json
    {
      "version": 1,
      "local": {
        "require_encryption": true,
        "min_security_score": 50
      }
    }
    ```

3. **Update your Claude Desktop config** (no `env` section needed):

    ```json
    {
      "mcpServers": {
        "omniserp": {
          "command": "mcp-omniserp"
        }
      }
    }
    ```

!!! note
    Without a policy file, the server works exactly as before using environment variables.

### Security Features

- **OS Keychain Integration**: API keys stored in macOS Keychain, Windows Credential Manager, or Linux Secret Service
- **Security Posture Checking**: Validates disk encryption, biometrics before credential access
- **Policy-Based Access**: Configure minimum security requirements
- **Graceful Fallback**: Uses environment variables when no policy exists

### Security Logging

When running in secure mode, the server logs security check results:

```
Security check passed: score=75, level=high
  Platform: darwin, Encrypted: true, Biometrics: true
SERPER_API_KEY retrieved from keychain successfully
```
