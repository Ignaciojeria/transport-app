package firebaseadminsdk

import (
	"context"
	"transport-app/app/shared/configuration"

	firebase "firebase.google.com/go/v4"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(NewFirebaseAdmin, configuration.NewConf)
}
func NewFirebaseAdmin(conf configuration.Conf) (*firebase.App, error) {
	return firebase.NewApp(context.Background(), &firebase.Config{
		ProjectID: conf.GOOGLE_PROJECT_ID,
	})
}
