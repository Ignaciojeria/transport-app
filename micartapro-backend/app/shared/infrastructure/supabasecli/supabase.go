package supabasecli

import (
	"micartapro/app/shared/configuration"

	ioc "github.com/Ignaciojeria/ioc"
	supabase "github.com/supabase-community/supabase-go"
)

func init() {
	ioc.Register(NewSupabaseClient)
}

func NewSupabaseClient(conf configuration.Conf) (*supabase.Client, error) {
	return supabase.NewClient(
		conf.SUPABASE_PROJECT_URL,
		conf.SUPABASE_BACKEND_API_KEY,
		&supabase.ClientOptions{})
}
