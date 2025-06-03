package onload

import (
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/adapter/out/tidbrepository/table"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(applicationLoader,
		table.NewRunMigrations,
		tidbrepository.NewLoadStatuses)
}
func applicationLoader(
	runMigrations table.RunMigrations,
	LoadOrderStatuses tidbrepository.LoadStatuses) error {
	if err := runMigrations(); err != nil {
		return err
	}
	LoadOrderStatuses()
	return nil
}
