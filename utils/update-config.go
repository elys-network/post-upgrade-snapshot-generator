package utils

import (
	"log"

	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func UpdateConfig(homePath, dbEngine string) {
	// Path to config files
	configPath := homePath + "/config/config.toml"
	appPath := homePath + "/config/app.toml"
	clientPath := homePath + "/config/client.toml"

	// Update config.toml for cors_allowed_origins
	Sed("s/^cors_allowed_origins =.*/cors_allowed_origins = [\\\"*\\\"]/", configPath)

	// Update config.toml for timeout_broadcast_tx_commit
	Sed("s/^timeout_broadcast_tx_commit =.*/timeout_broadcast_tx_commit = \\\"120s\\\"/", configPath)

	// Update config.toml for db_backend
	Sed("s/^db_backend =.*/db_backend = \\\""+dbEngine+"\\\"/", configPath)

	// Update app.toml for enabling the API server
	Sed("/^# Enable defines if the API server should be enabled./{n;s/enable = false/enable = true/;}", appPath)

	// Update app.toml for app-db-backend
	Sed("s/^app\\-db\\-backend =.*/app\\-db\\-backend = \\\""+dbEngine+"\\\"/", appPath)

	// Update app.toml for gas-to-suggest
	Sed("s/^gas\\-to\\-suggest =.*/gas\\-to\\-suggest = 300000/", appPath)

	// Update client.toml for keyring-backend
	Sed("s/^keyring\\-backend =.*/keyring\\-backend = \\\"test\\\"/", clientPath)

	log.Printf(types.ColorYellow + "config files have been updated successfully.")
}
