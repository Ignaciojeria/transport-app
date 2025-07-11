package natsconn

import (
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/database"
	"transport-app/app/shared/infrastructure/observability"
	"transport-app/app/shared/infrastructure/resendcli"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func init() {
	ioc.Registry(NewJetStream,
		NewConn,
		resendcli.NewClient,
		observability.NewObservability,
		database.NewConnectionFactory,
		configuration.NewConf,
		configuration.NewTiDBConfiguration)
}
func NewJetStream(
	conn *nats.Conn,
	resend resendcli.ResendClient,
	obs observability.Observability,
	connFactory database.ConnectionFactory,
	conf configuration.Conf,
	dbconf configuration.DBConfiguration) (jetstream.JetStream, error) {
	obs.Logger.Info("checking resend health")
	if conf.RESEND_API_KEY != "" {
		obs.Logger.Info("resend health check")
		err := resend.HealthCheck()
		if err != nil {
			return nil, err
		}
	}
	if dbconf.DB_STRATEGY != "disabled" {
		obs.Logger.Info("checking database health")
		sqlDB, err := connFactory.DB.DB()
		if err != nil {
			return nil, err
		}
		err = sqlDB.Ping()
		if err != nil {
			return nil, err
		}
	}
	obs.Logger.Info("stablishing connection to jetstream")
	return jetstream.New(conn)
}
