# CLI Tool

The OmniSerp CLI provides a command-line interface for quick searches.

## Installation

```bash
go build ./cmd/omniserp
```

Or install globally:

```bash
go install github.com/plexusone/omniserp/cmd/omniserp@latest
```

## Basic Usage

```bash
# Set API key
export SERPER_API_KEY="your_api_key"

# Basic search (specify engine and query)
./omniserp -e serper -q "golang programming"

# Or use long flags
./omniserp --engine serpapi --query "golang programming"

# With SerpAPI
export SERPAPI_API_KEY="your_api_key"
./omniserp -e serpapi -q "golang programming"
```

## Options

| Flag | Long Flag | Description | Required |
|------|-----------|-------------|----------|
| `-e` | `--engine` | Search engine (serper, serpapi) | Yes |
| `-q` | `--query` | Search query | Yes |

## Output

The CLI outputs JSON-formatted search results:

```json
{
  "data": {
    "organic": [
      {
        "title": "The Go Programming Language",
        "link": "https://go.dev/",
        "snippet": "..."
      }
    ]
  }
}
```
