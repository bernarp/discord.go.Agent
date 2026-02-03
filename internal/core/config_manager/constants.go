package config_manager

import "time"

const (
	ExtensionYaml    = ".yaml"
	PrefixMerge      = "MERGE."
	DebounceDuration = 200 * time.Millisecond

	PlaceholderHeader = `# ==============================================================================
# CONFIGURATION PLACEHOLDER
# ==============================================================================
# This file was automatically generated because it was missing.
# The module associated with this configuration is currently DISABLED.
#
# TO ENABLE THE MODULE:
# 1. Fill in the required values below.
# 2. Save the file.
# 3. Restart the application or wait for the hot-reload system to detect changes.
# ==============================================================================
`
)
