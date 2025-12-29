# Release Notes v0.7.0

**Release Date:** December 29, 2025

## New Features

### Secure Credential Management in MCP Server

The `mcp-omniserp` MCP server now includes optional secure credential management using VaultGuard. This provides enterprise-grade security for API key handling while maintaining full backward compatibility.

**Key Features:**

- 🔐 **OS Keychain Integration**: API keys can be stored in the OS keychain (macOS Keychain, Windows Credential Manager, Linux Secret Service) instead of environment variables
- 🛡️ **Security Posture Checking**: Validates system security before granting credential access
- 📜 **Policy-Based Access Control**: Configure security requirements via `~/.agentplexus/policy.json`
- 🔄 **Graceful Fallback**: Automatically uses environment variables when no policy is configured (backward compatible)

**How It Works:**

- **Without a policy file**: Works exactly as before - reads API keys from environment variables
- **With a policy file**: Retrieves API keys from the OS keychain with security posture validation

**Standard Mode (environment variables):**
```bash
export SERPER_API_KEY="your-key"
export SEARCH_ENGINE="serper"  # optional, defaults to serper
./mcp-omniserp
```

**Secure Mode (OS keychain + policy):**

1. Store your API key in the keychain:
   ```bash
   security add-generic-password -s "omnivault" -a "SERPER_API_KEY" -w "your-key"
   ```

2. Create a security policy (`~/.agentplexus/policy.json`):
   ```json
   {
     "version": 1,
     "local": {
       "require_encryption": true,
       "min_security_score": 50
     }
   }
   ```

3. Run the server:
   ```bash
   ./mcp-omniserp
   ```

**Security Logging:**
When running in secure mode, the server logs security check results including:
- Security score and level
- Platform information
- Disk encryption status
- Biometrics configuration status

**New Dependencies:**
- `github.com/agentplexus/omnivault-keyring` v0.1.0 - OS keychain integration
- `github.com/agentplexus/vaultguard` v0.2.0 - Security posture checking and policy management

### New Client SDK Functions

- `client.NewWithRegistry()` - Create a client with a pre-configured registry
- `serpapi.NewWithAPIKey()` - Create a SerpAPI engine with a specific API key

## Changes Since v0.6.0

### Refactoring
- Renamed MCP server from `mcpserver` to `mcp-omniserp` to follow the `mcp-` prefix naming convention
- Transferred repository to `agentplexus` organization

### Documentation
- Added icons to README overview feature list for improved readability
- Updated README formatting and project structure documentation
- Added Marp presentation (`PRESENTATION.md`)
- Added documentation site (`docs/index.html`)

## Migration Guide

### MCP Server Rename

If you're using the MCP server, update your configuration:

**Before:**
```json
{
  "mcpServers": {
    "omniserp": {
      "command": "mcpserver"
    }
  }
}
```

**After:**
```json
{
  "mcpServers": {
    "omniserp": {
      "command": "mcp-omniserp",
      "env": {
        "SERPER_API_KEY": "your_serper_api_key"
      }
    }
  }
}
```

### Upgrading to Secure Credentials (Optional)

To upgrade from environment variables to secure keychain storage:

1. Add your API key to the OS keychain
2. Create a policy file at `~/.agentplexus/policy.json`
3. Remove the `env` section from your MCP config (credentials now come from keychain)

```json
{
  "mcpServers": {
    "omniserp": {
      "command": "mcp-omniserp"
    }
  }
}
```
