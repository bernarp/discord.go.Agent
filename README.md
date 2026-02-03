# Discord Bot Agent

## Installation
1. Clone the repository.
2. Install dependencies:
```bash
go mod download
```

## Running the Bot
Execute the following command from the root directory:
```bash
go run ./cmd
```

## Token Management
The bot looks for a Discord token in the following order:
1. **Environment Variables**: `BOT_TOKEN` in `.env.dev` or `.env`.
2. **Encrypted File**: `token.enc` in the root directory.
3. **Manual Input**: If no token is found, the bot will prompt for it in the terminal.

### token.enc
When a token is entered manually, it is encrypted using **AES-GCM** and tied to the local machine's hardware ID (`machineid`). 
* The file cannot be decrypted if moved to a different computer.
* To reset the token, delete `token.enc`.

## Configuration System

### config_df/
Contains default configuration files for modules. 
* If a required configuration file is missing or empty, the bot automatically generates a **placeholder** file with a default schema and instructions.
* Modules associated with placeholders remain disabled until valid data is provided.

### config_mrg/
Used for overriding values without modifying the defaults.
* Files must be prefixed with `MERGE.` (e.g., `MERGE.system.discord.template.yaml`).
* The system performs a **deep merge**: values in `config_mrg` take precedence over `config_df`.

### Hot-Reload
The bot monitors both directories for changes. Updating and saving a YAML file will trigger an automatic configuration reload and module state update without restarting the process.

## Project Structure
* `cmd/`: Application entry point and graceful shutdown logic.
* `config_df/`: Default YAML configurations.
* `config_mrg/`: Override YAML configurations.
* `internal/core/`:
    * `config_manager`: Logic for loading, merging, and watching YAML files.
    * `module_manager`: Lifecycle management (Enable/Disable) of bot modules.
    * `startup`: Token validation, hardware-bound encryption, and terminal prompts.
    * `eventbus`: Internal asynchronous event distribution.
* `internal/modules/`: Functional bot modules.
* `pkg/`: Shared utilities and context tracing.
* `logs/`: Automatically generated execution logs (JSON and Console).
