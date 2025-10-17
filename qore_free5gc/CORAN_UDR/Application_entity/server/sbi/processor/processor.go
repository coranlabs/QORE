package processor

import (
	"github.com/coranlabs/CORAN_UDR/Application_entity/database"
	"github.com/coranlabs/CORAN_UDR/Application_entity/pkg/app"
)

type Processor struct {
	app.App
	database.DbConnector
}

func NewProcessor(udr app.App) *Processor {
	return &Processor{
		App:         udr,
		DbConnector: database.NewDbConnector(udr.Config().Configuration.DbConnectorType),
	}
}
