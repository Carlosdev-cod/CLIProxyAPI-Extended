package helps

import "strings"

// antigravityTieredModels maps canonical base model ids exposed by the proxy
// to the upstream Antigravity tier used when the client does not select a
// reasoning level explicitly. Upstream only serves "<base>-<level>" ids for
// these families.
var antigravityTieredModels = map[string]string{
	"gemini-3.1-pro":   "low",
	"gemini-3.6-flash": "high",
	"gemini-3.7-flash": "high",
	"gemini-3.8-flash": "high",
}

// antigravityTierClamps defines the tiers each model family actually serves
// upstream and the nearest fallback for unsupported requests.
var antigravityTierClamps = map[string]map[string]string{
	"gemini-3.1-pro": {"medium": "low", "high": "low", "extra-low": "low"},
}

// AntigravityUpstreamModelID resolves a canonical base model id plus an
// optional reasoning level into the tiered model id upstream Antigravity
// expects. Non-tiered models and ids that already carry a tier suffix are
// returned unchanged.
func AntigravityUpstreamModelID(modelName, level string) string {
	base := strings.ToLower(strings.TrimSpace(modelName))
	defaultTier, tiered := antigravityTieredModels[base]
	if !tiered {
		return base
	}
	lvl := strings.ToLower(strings.TrimSpace(level))
	switch lvl {
	case "extra-low", "low", "medium", "high":
	default:
		return base + "-" + defaultTier
	}
	if clamp, ok := antigravityTierClamps[base]; ok {
		if fallback, bad := clamp[lvl]; bad {
			lvl = fallback
		}
	}
	return base + "-" + lvl
}
