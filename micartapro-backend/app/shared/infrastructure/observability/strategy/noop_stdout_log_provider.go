package strategy

import (
	"log/slog"
	"micartapro/app/shared/configuration"
	"os"

	otelslogjson "github.com/go-slog/otelslog"
)

func NoOpStdoutLogProvider(conf configuration.Conf) *slog.Logger {
	return slog.New(otelslogjson.NewHandler(slog.NewJSONHandler(os.Stdout, nil))).With(
		slog.String("env", conf.ENVIRONMENT),
		slog.String("version", conf.VERSION),
		slog.String("service", conf.PROJECT_NAME),
	)
}
