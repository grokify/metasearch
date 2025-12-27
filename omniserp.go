package omniserp

import (
	"fmt"
	"os"
)

// GetDefaultEngine returns the default engine based on environment variable
// Falls back to "serper" if SEARCH_ENGINE is not set or the specified engine is not available
func GetDefaultEngine(registry *Registry) (Engine, error) {
	selectedEngine := os.Getenv("SEARCH_ENGINE")
	if selectedEngine == "" {
		selectedEngine = "serper" // Default to Serper
	}

	engine, exists := registry.Get(selectedEngine)
	if !exists {
		availableEngines := registry.List()
		if len(availableEngines) == 0 {
			return nil, fmt.Errorf("no search engines available")
		}

		// Try to fallback to serper, otherwise use the first available
		if fallbackEngine, exists := registry.Get("serper"); exists {
			return fallbackEngine, fmt.Errorf("engine '%s' not found, falling back to 'serper'. Available engines: %v", selectedEngine, availableEngines)
		}

		// Use first available engine
		firstEngine := availableEngines[0]
		engine, _ = registry.Get(firstEngine)
		return engine, fmt.Errorf("engine '%s' not found, falling back to '%s'. Available engines: %v", selectedEngine, firstEngine, availableEngines)
	}

	return engine, nil
}

// GetEngineInfo returns information about an engine
type EngineInfo struct {
	Name           string   `json:"name"`
	Version        string   `json:"version"`
	SupportedTools []string `json:"supported_tools"`
}

// GetEngineInfo returns information about a specific engine
func GetEngineInfo(engine Engine) EngineInfo {
	return EngineInfo{
		Name:           engine.GetName(),
		Version:        engine.GetVersion(),
		SupportedTools: engine.GetSupportedTools(),
	}
}

// GetAllEngineInfo returns information about all registered engines
func GetAllEngineInfo(registry *Registry) map[string]EngineInfo {
	engines := registry.GetAll()
	info := make(map[string]EngineInfo)

	for name, engine := range engines {
		info[name] = GetEngineInfo(engine)
	}

	return info
}
