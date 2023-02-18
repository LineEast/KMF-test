package configuration

import (
	"os"

	"github.com/pelletier/go-toml"

	"kmfRedirect/internal/database"
	"kmfRedirect/internal/server"
)

type (
	Configuration struct {
		Database *database.Configuration `toml:"database"`
		Server   *server.Configuration   `toml:"server"`
	}
)

func New() (configuration *Configuration, err error) {
	file, err := os.OpenFile(".configuration.toml", os.O_RDONLY, 0)
	if err != nil {
		return
	}

	configuration = new(Configuration)

	err = toml.NewDecoder(file).Decode(configuration)
	if err != nil {
		return
	}

	return
}
