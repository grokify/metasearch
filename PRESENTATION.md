---
marp: true
paginate: true
style: |
  .mermaid {
    display: flex;
    justify-content: center;
    align-items: center;
  }
  .mermaid svg {
    max-height: 400px;
    width: auto;
  }
---

<script type="module">
import mermaid from 'https://cdn.jsdelivr.net/npm/mermaid@10/dist/mermaid.esm.min.mjs';
mermaid.initialize({
  startOnLoad: true,
  theme: 'dark',
  themeVariables: {
    background: 'transparent',
    primaryColor: '#7c4dff',
    primaryTextColor: '#e8eaf6',
    primaryBorderColor: '#667eea',
    lineColor: '#b39ddb',
    secondaryColor: '#302b63',
    tertiaryColor: '#24243e'
  }
});
</script>

<!-- _paginate: false -->

# OmniSerp 🔍
## A Unified Search Engine Abstraction for AI Agents

---

<!-- _paginate: false -->
<!-- _class: lead -->

# Part 1
## The Problem 🎯

---

# The Problem ⚠️

## AI agents need real-time information

- 🧠 LLMs have knowledge cutoffs
- 📅 Static training data becomes stale
- 👤 Users expect current, accurate answers
- 🌐 Agents must access live web data

---

# The Challenge 🧩

## Search API fragmentation

- 🔀 Multiple providers: Serper.dev, SerpApi, Google, Bing
- 📋 Each has different APIs, parameters, response formats
- 🔧 Switching providers requires extensive code changes
- ⚙️ Different providers support different operations
- ❌ No standardized interface exists

---

# Why Web Search Matters for Agents 🤖

## Grounding AI in Reality

| Capability | Without Search | With Search |
|------------|----------------|-------------|
| Current events | ❌ Limited | ✅ Real-time |
| Fact verification | ⚠️ Hallucination risk | ✅ Verified |
| Research tasks | ❌ Incomplete | ✅ Comprehensive |
| Local information | ❌ None | ✅ Available |

---

# Web Search Enables 🚀

## Critical Agent Capabilities

- 📰 **Real-time information**: News, stock prices, weather
- ✅ **Fact verification**: Reduce hallucinations with sources
- 📚 **Research**: Academic papers, documentation, tutorials
- 📍 **Local context**: Places, maps, business reviews
- 👁️ **Visual search**: Image analysis via reverse search (Lens)
- 🛒 **Shopping**: Product comparisons and pricing

---

<!-- _paginate: false -->
<!-- _class: lead -->

# Part 2
## Search Providers 🤝

---

# Meet the Providers 🤝

## Serper.dev ⚡

- 🚀 Fast Google Search API
- 📦 Simple JSON responses
- 🔢 **12 search operations** supported
- 👁️ Includes Google Lens support
- 🔐 POST-based API with header authentication

---

# Meet the Providers 🤝

## SerpApi 🔬

- 📊 Comprehensive search API
- 🌍 Multiple search engines supported
- 🔢 **11 search operations** (no Lens)
- 🔑 GET-based API with query auth
- 📋 Structured, detailed responses

---

# Feature Comparison 📊

| Operation | Serper | SerpApi |
|-----------|:------:|:-------:|
| 🌐 Web Search | ✅ | ✅ |
| 📰 News | ✅ | ✅ |
| 🖼️ Images | ✅ | ✅ |
| 🎬 Videos | ✅ | ✅ |
| 📍 Places | ✅ | ✅ |
| 🗺️ Maps | ✅ | ✅ |
| 🛒 Shopping | ✅ | ✅ |
| 🎓 Scholar | ✅ | ✅ |
| 👁️ **Lens** | ✅ | ❌ |

---

<!-- _paginate: false -->
<!-- _class: lead -->

# Part 3
## The Solution 💡

---

# The Solution: OmniSerp 💡

## One interface, multiple providers

```go
// Same code works with any provider
client := omniserp.New()
result, _ := client.Search(ctx, params)

// Switch providers at runtime
client.SetEngine("serpapi")
```

---

# Architecture Overview 🏗️

<div class="mermaid">
flowchart TB
    A["🤖 Your Application / Agent"] --> B["📦 OmniSerp Client SDK<br/>Unified Interface<br/>Capability Checking<br/>Response Normalization"]
    B --> C["⚡ Serper<br/>Engine"]
    B --> D["🔬 SerpApi<br/>Engine"]
    C --> E["serper.dev"]
    D --> F["serpapi.com"]
    classDef app fill:#7c4dff,stroke:#667eea,color:#fff
    classDef sdk fill:#00bfa5,stroke:#00897b,color:#fff
    classDef engine fill:#ff7043,stroke:#e64a19,color:#fff
    classDef api fill:#78909c,stroke:#546e7a,color:#fff
    class A app
    class B sdk
    class C,D engine
    class E,F api
</div>

---

# Why a Generic Interface? 🔌

## Problem 1: Provider Lock-in 🔒

Without abstraction:
- ⛓️ Code tightly coupled to one provider
- 🔄 Switching requires rewriting everything
- 🧪 Testing requires live API calls

With OmniSerp:
- ✨ Swap providers with one line change
- 🎭 Mock engines for testing
- 🛡️ Graceful fallbacks

---

# Why a Generic Interface? 🔌

## Problem 2: Feature Parity ⚖️

```go
// Check capabilities at runtime
if client.SupportsOperation(OpSearchLens) {
    result, _ := client.SearchLens(ctx, params)
} else {
    // Fallback strategy
}
```

🤖 Agents can adapt behavior based on available features

---

# Why a Generic Interface? 🔌

## Problem 3: Response Inconsistency 🔀

Raw responses differ between providers:

```json
// Serper
{"answerBox": {...}}

// SerpApi
{"answer_box": {...}}
```

✅ OmniSerp normalizes to a unified structure

---

# Normalized Responses 📐

## Engine-agnostic data structures

```go
normalized, _ := client.SearchNormalized(ctx, params)

// Same fields regardless of provider
for _, result := range normalized.OrganicResults {
    fmt.Println(result.Title)
    fmt.Println(result.Link)
    fmt.Println(result.Snippet)
}
```

---

# The Engine Interface ⚙️

## Clean abstraction

```go
type Engine interface {
    GetName() string
    GetSupportedTools() []string

    Search(ctx, params) (*SearchResult, error)
    SearchNews(ctx, params) (*SearchResult, error)
    SearchImages(ctx, params) (*SearchResult, error)
    // ... 12 operations total
}
```

➕ Adding a new provider = implement this interface

---

<!-- _paginate: false -->
<!-- _class: lead -->

# Part 4
## Real-World Applications 🌍

---

# Real-World Applications 🌍

## 1. MCP Server for AI Assistants 🤖

```bash
# Claude Desktop integration
./mcp-omniserp
# Registers 12 tools for Serper
# Registers 11 tools for SerpApi
```

💬 Enables Claude to search the web natively

---

# Real-World Applications 🌍

## 2. Multi-Provider Redundancy 🛡️

```go
result, err := client.Search(ctx, params)
if err != nil {
    client.SetEngine("serpapi") // Fallback
    result, _ = client.Search(ctx, params)
}
```

✅ Never fail because one API is down

---

# Real-World Applications 🌍

## 3. Cost Optimization 💰

- 📊 Use different providers for different query types
- 📈 Route high-volume queries to cheaper provider
- ⭐ Reserve premium features for when needed

---

# Real-World Applications 🌍

## 4. Secure Credentials with VaultGuard 🔐

```json
// ~/.vaultguard/policy.json
{
  "version": 1,
  "local": {
    "require_encryption": true,
    "min_security_score": 50
  }
}
```

🔑 API keys stored in OS keychain, not environment variables

---

# VaultGuard Integration 🛡️

## Enterprise-grade security for AI agents

- 🔐 **OS Keychain**: macOS Keychain, Windows Credential Manager, Linux Secret Service
- 📊 **Security Posture**: Validates disk encryption, biometrics before credential access
- 📜 **Policy-Based**: Configure minimum security requirements
- 🔄 **Graceful Fallback**: Uses env vars when no policy exists

🏢 Built on **PlexusOne VaultGuard** security framework

---

<!-- _paginate: false -->
<!-- _class: lead -->

# Part 5
## Getting Started 🚀

---

# Benefits Summary ✨

## For Agent Developers

- 🎯 **Simplicity**: One API to learn
- 🔄 **Flexibility**: Switch providers easily
- 🛡️ **Reliability**: Built-in fallbacks
- 🔮 **Future-proof**: Add providers without code changes
- 🔒 **Type safety**: Strongly-typed Go interfaces

---

# Getting Started 🚀

## Quick setup

```bash
export SERPER_API_KEY="your-key"
# or
export SERPAPI_API_KEY="your-key"
```

```go
import "github.com/plexusone/omniserp/client"

c, _ := client.New()
result, _ := c.Search(ctx, omniserp.SearchParams{
    Query: "latest AI news",
})
```

---

# The Future 🔮

## Extensible by design

- ➕ Easy to add new providers (Bing, DuckDuckGo, etc.)
- 🔌 Plugin architecture ready
- 👥 Community contributions welcome
- 🌱 Growing ecosystem of AI tools

---

# Summary 📋

## OmniSerp solves

1. 🔀 **Search API fragmentation** → unified interface
2. 🔒 **Provider lock-in** → runtime switching
3. 📐 **Response inconsistency** → normalization layer
4. ⚖️ **Feature gaps** → capability checking
5. 🤖 **Agent integration** → MCP server included
6. 🔐 **Credential security** → VaultGuard integration

---

<!-- _paginate: false -->

# Thank You 🙏

## OmniSerp 🔍

**github.com/plexusone/omniserp**

🌐 Unified search for the AI agent era
