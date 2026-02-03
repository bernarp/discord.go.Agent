### Файл `./config_df/README.md`

```markdown
# Default Configurations (config_df)

This directory contains the base configuration files for all bot modules.

## Behavior
* **Automatic Generation**: If a required configuration file is missing or empty, the bot will automatically create a placeholder file with a default schema and instructions.
* **Module Lifecycle**: A module will remain in a `Disabled` state if its corresponding configuration file is a placeholder or contains invalid data.
* **Hot-Reload**: Changes made to files in this directory are detected automatically, triggering a module update without restarting the bot.

## Developer Note: Config Templates
The structure and default values of the generated placeholders are defined directly in the module's source code using the `ConfigTemplate()` method. 

Example from code:
```go
func (m *Module) ConfigTemplate() any {
    defaultEnabled := true
    return Config{
        Enabled: &defaultEnabled,
        LogDetails: LogDetails{
            Guild: true,
            Content: true,
        },
    }
}
```
When the bot generates a missing YAML file, it uses the values returned by this function as the initial content.

## Usage
The files currently present in this repository are examples. You can delete this folder or individual files; the bot will recreate them upon the next execution if they are required by registered modules.
