// Package asuser exercises the postgrest module against the live local
// Supabase stack that scripts/integration-test.sh starts, as a signed-in user:
// every client here holds the publishable key, and the secret key is absent
// from the package's environment.
package asuser

import (
	"testing"

	"github.com/supabase/supabase-go/core"
	"github.com/supabase/supabase-go/core/configuration"
	"github.com/supabase/supabase-go/integration-testing/testkit"
	"github.com/supabase/supabase-go/postgrest"
)

// newIntegrationClient wires a client at the local Supabase stack started by
// scripts/integration-test.sh, failing the test when the required environment
// variables are absent. It forwards the optional configuration options so
// tests can exercise construction-time settings against the real stack.
func newIntegrationClient(t *testing.T, options ...configuration.Option) *postgrest.Client {
	t.Helper()
	projectURL, apiKey := testkit.Credentials(t)
	projectConfiguration, err := configuration.New(core.ModulePathPostgrest, projectURL, apiKey, options...)
	if err != nil {
		t.Fatalf("configuration.New: %v", err)
	}
	return postgrest.NewFromConfiguration(projectConfiguration)
}

// reservedValues holds each text value seeded into the "⚠ reserved ⚠" table,
// named as the row that holds it: single characters that URL, PostgREST or
// Postgres syntax reads specially or that test the text encoding, then
// longer forms that combine them, the empty string and the word NULL.
var reservedValues = []struct {
	name  string
	value string
}{
	{"just comma", ","},
	{"just dot", "."},
	{"just colon", ":"},
	{"just asterisk", "*"},
	{"just opening parenthesis", "("},
	{"just closing parenthesis", ")"},
	{"just dollar", "$"},
	{"just opening square bracket", "["},
	{"just closing square bracket", "]"},
	{"just semicolon", ";"},
	{"just opening brace", "{"},
	{"just closing brace", "}"},
	{"just caret", "^"},
	{"just percent", "%"},
	{"just left angle bracket", "<"},
	{"just right angle bracket", ">"},
	{"just equals", "="},
	{"just plus", "+"},
	{"just minus", "-"},
	{"just double quote", `"`},
	{"just single quote", "'"},
	{"just backslash", `\`},
	{"just backspace", "\b"},
	{"just form feed", "\f"},
	{"just newline", "\n"},
	{"just carriage return", "\r"},
	{"just tab", "\t"},
	{"just multitudinous", "众"},
	{"just poop", "💩"},
	{"a comma b", "a,b"},
	{"opening brace inside", "brace{inside"},
	{"closing brace inside", "brace}inside"},
	{"opening paren inside", "paren(open"},
	{"closing paren inside", "close)paren"},
	{"double quoted inside", `say "hi"`},
	{"double quoted", `"The IKEA Effect"`},
	{"backslash inside", `back\slash`},
	{"space padded", " padded "},
	{"empty", ""},
	{"wildcard version tag", "v1.2:rc*"},
	{"spelled NULL", "NULL"},
}
