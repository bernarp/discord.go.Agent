# Merge Configurations (config_mrg)

This directory is used for overriding default settings without modifying the files in `config_df`.

## Naming Convention
To override a configuration, create a file with the `MERGE.` prefix followed by the exact name of the base file.
* **Base file**: `system.discord.template.yaml`
* **Override file**: `MERGE.system.discord.template.yaml`

## Behavior
* **Deep Merge**: The system performs a deep merge of the two files. Values defined in `config_mrg` take precedence over those in `config_df`.
* **Partial Overrides**: You do not need to copy the entire configuration. You can specify only the specific keys you wish to change.
* **Hot-Reload**: Saving a merge file will trigger an immediate configuration update.

## Usage
This folder is optional. If no matching `MERGE.` file is found, the bot uses the defaults from `config_df`.