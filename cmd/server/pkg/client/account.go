package client

import (
	"encoding/json"
	"log"

	"github.com/maplelm/dwarfwars/pkg/command"
)

func (f *Factory) RegisterAccount(logger *log.Logger, c Client, cmd *command.Command) error {
	var (
		cmdD command.DataWrapper[struct{}]
	)
	// Validate client
	if cmd.ClientID != c.Uid() {
		// cmd client id doesn't match client sent to funciton
	}
	// parse command based on format
	switch cmd.Format {
	case command.FormatJSON:
		err := json.Unmarshal(cmd.Data, cmdD)
		if err != nil {
			logger.Printf("Error Parsing command to register account from json: %s", err)
			return err
		}
	case command.FormatGLOB:
	case command.FormatCSV:
	default:
		// Unsupported format!
	}
	// Check that Command is valid
	// Check if account is unqiue
	// Create account
	// return command status to client
	// send client back to general dispatcher
}
