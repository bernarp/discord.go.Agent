package config_manager

// Contract предоставляет типизированный доступ к именам конфигурационных файлов.
// Использование: config_manager.Contract.System.Discord.Template
var Contract = struct {
	System struct {
		Discord struct {
			Template  string
			Template2 string
			Tickets   string
			Moderator string
		}
	}
}{
	System: struct {
		Discord struct {
			Template  string
			Template2 string
			Tickets   string
			Moderator string
		}
	}{
		Discord: struct {
			Template  string
			Template2 string
			Tickets   string
			Moderator string
		}{
			Template:  "system.discord.template",
			Template2: "system.discord.template2",
			Tickets:   "system.discord.tickets",
			Moderator: "system.discord.moderator",
		},
	},
}
