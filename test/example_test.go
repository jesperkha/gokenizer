package test

import (
	"os"
	"strings"
	"testing"

	"github.com/jesperkha/gokenizer"
	"github.com/joho/godotenv"
)

func parseEnv(file string) (kv map[string]string, err error) {
	kv = make(map[string]string)

	tokr := gokenizer.New()

	tokr.Class("key", "{var}")
	tokr.Class("value", "{string}", "{text}")

	tokr.Class("keyValue", "{ws}{key}{ws}={ws}{value}")
	tokr.Class("comment", "#{any}")

	tokr.Class("expression", "{comment}", "{keyValue}")

	tokr.Pattern("{expression}", func(t gokenizer.Token) error {
		keyval := t.Get("expression").Get("keyValue")

		if keyval.Length != 0 {
			key := keyval.Get("key").Lexeme
			value := keyval.Get("value").Lexeme

			// Convert value to env friendly text
			value = strings.ReplaceAll(value, "\"", "")    // Remove quotes
			value = strings.ReplaceAll(value, "\\n", "\n") // Put newlines

			kv[key] = value
		}

		return nil
	})

	// Run for each line
	for _, line := range strings.Split(file, "\n") {
		if err = tokr.Run(line); err != nil {
			return kv, err
		}
	}

	return kv, err
}

func TestExample(t *testing.T) {
	file, err := os.ReadFile("example.env")
	if err != nil {
		t.Fatal(err)
	}

	godotenv.Load("example.env")

	kv, err := parseEnv(string(file))
	if err != nil {
		t.Error(err)
	}

	for key, value := range kv {
		osv := os.Getenv(key)
		if osv != value {
			t.Errorf("expected '%s=%s', got '%s=%s'", key, value, key, osv)
		}
	}
}
