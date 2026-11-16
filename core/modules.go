package core

// ModulePath represents a distinct module within the Supabase Go SDK, where that
// module is capable of emitting telemetry - that is, where that module can have its
// interactions with the service observed by the service.
type ModulePath string

const (
	// ModulePathRoot is the import path for the Supabase Go SDK's convenience
	// entry point module, composed of other domain modules. It is the
	// repository root, and the base of every other module's path.
	ModulePathRoot ModulePath = "github.com/supabase/supabase-go"

	// ModulePathPostgrest is the import path for the Supabase Go SDK's PostgREST
	// client module.
	ModulePathPostgrest ModulePath = ModulePathRoot + "/postgrest"

	// ModulePathAuth is the import path for the Supabase Go SDK's Auth client
	// module.
	ModulePathAuth ModulePath = ModulePathRoot + "/auth"
)
