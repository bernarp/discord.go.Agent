
# Discord Bot Agent

## Installation

```bash
go mod download

```

## Running the Bot

```bash
go run ./cmd

```

## Environment and Secrets Management

The bot requires four core parameters: **Bot Token**, **Application ID**, **Guild ID**, and **API Port**. It resolves them in this order:

### 1. Environment Variables (.env) — Recommended

The system checks for a `.env` file in the root directory:

* `BOT_TOKEN`
* `APP_ID`
* `GUILD_ID`
* `PORT`
* `PREFIX`

### 2. Encrypted Files (.enc) — Native Hardware Binding

If environment variables are missing, the bot uses the following hardware-bound encrypted files:

* **`token.enc`**: Discord Bot Token.
* **`appid.enc`**: Discord Application ID (required for command sync).
* **`guildid.enc`**: Target Guild ID (for instant command updates).
* **`port.enc`**: API Server Port.
* **`prefix.enc`**: Prefix commands. 

**Technical Specifications:**

* **Encryption**: AES-GCM.
* **Key Derivation**: Tied to the local hardware ID (`machineid`).
* **Portability**: These files are **non-portable**. They cannot be decrypted on a different machine.
* **Initial Setup**: If neither `.env` nor `.enc` files exist, the bot will prompt you to enter these values in the terminal and generate the `.enc` files automatically.

---

## Configuration System

### config_df/

* **Default Settings**: Contains default YAML for each module.
* **Placeholders**: If a config is missing, the bot generates a template. Modules stay disabled until the placeholder is filled.

### config_mrg/

* **Overrides**: Uses the `MERGE.` prefix (e.g., `MERGE.system.discord.yaml`).
* **Deep Merge**: Values here override defaults in `config_df`.

### Hot-Reload

The bot watches both directories using `fsnotify`. Saving a YAML file triggers an immediate config update and module state refresh without a restart.

---

## Interaction Orchestration

The bot uses a decoupled pipeline to handle Discord interactions:

* **Interaction Manager**: The orchestrator. It listens for `InteractionCreate` via the `EventBus`, logs metadata, and routes requests.
* **Commands Handler**:
* Automatically triggers a **Deferred Response** ("Thinking...").
* Validates module status before execution.
* Uses **Webhook API** (`InteractionResponseEdit`) for the final execution result.



---

## Project Structure

* `cmd/`: App entry point and graceful shutdown.
* `internal/core/`:
* `config_manager`: Loading, merging, and watching configs.
* `module_manager`: Lifecycle and dependency management.
* `eventbus`: Asynchronous internal event distribution.
* `startup`: Secret validation, terminal prompts, and hardware-bound encryption.


* `internal/client/`: Discord session, orchestrator, and command handlers.
* `internal/modules/`: Functional bot modules (e.g., `template`).

---
