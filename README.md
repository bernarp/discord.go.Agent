# DiscordBotAgent: Core Infrastructure

The **Configuration Manager** and **Functional Modules** form the core infrastructure of the DiscordBotAgent. The following sections detail their implementation and interaction.

---

## Configuration Manager

The `config_manager` package handles the loading, validation, and monitoring of configuration files. It operates independently of the business logic, providing data to registered consumers.

### File Loading and Merging
The `Manager` struct initializes with two directory paths: a default path (`config_df`) and a merge path (`config_mrg`). The loading process uses a deep merge strategy:

* **Base Load:** The system reads the base YAML file from the default directory.
* **Merge Override:** It checks the merge directory for a file with the `MERGE.` prefix. If found, the system reads this file and recursively overwrites values in the base map using `deepMerge`.
* **Unmarshal:** The merged map is unmarshaled into the target struct using a strict decoder (`KnownFields(true)`) to detect undefined fields.

### Registration and Placeholders
Modules register their configurations via `Manager.Register`. This method accepts a config name, a template struct, and an update callback.

* If the configuration file does not exist or is empty, `createPlaceholder` generates a file using the provided template and writes a header explaining how to enable the module.
* The function returns `ErrPlaceholderCreated`, signaling the caller that the module cannot start immediately.

### Validation
Before a configuration is applied, it passes through `ConfigValidator`. This component uses the `go-playground/validator` library to enforce rules defined in struct tags (e.g., `validate:"required"`, `validate:"gte=0"`). Invalid configurations prevent the update from proceeding.

### Hot-Reloading
The manager utilizes `fsnotify` to watch both configuration directories. Upon detecting a `Write`, `Create`, or `Remove` event:

1.  **Debounce:** A timer delays execution (200ms) to prevent multiple triggers during a single file save.
2.  **Reload:** The system calls `reloadConfig`, which re-executes the load, merge, and validation steps.
3.  **Callback:** If valid, the new configuration replaces the old one in the registry, and the registered callback function is executed. If invalid, the error is logged, and the current state remains unchanged.

---

## Functional Modules

The `module_manager` package controls the lifecycle of bot features. It enforces a specific interface and manages dependencies between modules.

### Module Interface
Components must implement the `Module` interface defined in `internal/core/module_manager/base.go`:

* `Name()`: Returns the module identifier.
* `ConfigKey()`: Returns the filename identifier for the configuration.
* `ConfigTemplate()`: Returns the default configuration struct.
* `OnEnable(ctx, cfg)`: Executed when the module starts.
* `OnDisable(ctx)`: Executed when the module stops.
* `OnConfigUpdate(ctx, cfg)`: Executed when configuration changes occur at runtime.

### State Management
The manager tracks the state of each module: `disabled`, `enabled`, `error`, or `dependency_disabled`.

* **Registration:** `Manager.Register` scans the module struct using reflection (`scanDependencies`) to identify fields that implement the `Module` interface. These are recorded as dependencies.
* **Enabling:** `tryEnable` verifies that all recorded dependencies are registered and in the enabled state. If a dependency is missing or disabled, the target module transitions to `dependency_disabled` or `error`.
* **Disabling:** When a module is disabled (manually or via config error), `disableDependents` recursively finds and disables any modules that depend on it.

---

## Integration and Workflow

The `app.go` file initializes both managers and links them.

1.  **Initialization:** The `module_manager` is initialized with a reference to the `config_manager`.
2.  **Binding:** When a module registers with `module_manager`, it registers its config with `config_manager`, passing `m.onConfigUpdate` as the callback.
3.  **Runtime Flow:**
    * When `config_manager` detects a file change, it invokes the callback.
    * `module_manager` receives the update.
    * If the config is valid, it calls `module.OnConfigUpdate`.
    * If the config is invalid, it calls `module.OnDisable` and propagates the disable signal to downstream dependencies.

### Example Implementation
The `template2` module demonstrates dependency injection. It holds a reference to `template.Module` in its struct. During registration, the `module_manager` detects this field via reflection, ensuring `template2` only enables if `template` is active.