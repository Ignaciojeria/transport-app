package supabasecli

import (
	"micartapro/app/shared/configuration"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	supabase "github.com/supabase-community/supabase-go"
)

func init() {
	ioc.Registry(NewSupabaseClient, configuration.NewConf)
}

func NewSupabaseClient(conf configuration.Conf) (*supabase.Client, error) {
	return supabase.NewClient(
		conf.SUPABASE_PROJECT_URL,
		conf.SUPABASE_BACKEND_API_KEY,
		&supabase.ClientOptions{})
}
