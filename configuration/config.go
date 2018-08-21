package configuration

import (
	"github.com/Appliscale/tyr/tyrsession/clientfactory"
	"github.com/Appliscale/tyr/tyrsession/sessionfactory"
)

type Config struct {
	Regions        *[]string
	Service        string
	Profile        string
	SessionFactory *sessionfactory.SessionFactory
	ClientFactory  *clientfactory.ClientFactory
}
