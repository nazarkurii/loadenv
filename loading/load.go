package loading

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator"
	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func Do[T any](envPath ...string) (*T, error) {
	var k = koanf.New(".")

	if len(envPath) != 0 {
		if err := k.Load(file.Provider(envPath[0]), dotenv.Parser()); err != nil {
			return nil, fmt.Errorf("failed to load env file: %w", err)
		}
	} else {
		k.Load(env.Provider(".", env.Opt{}), nil)
	}

	for key, value := range k.All() {
		key = strings.ReplaceAll(strings.ToLower(key), "_", ".")

		k.Set(key, value)
	}

	var cfg T
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := validator.New().Struct(&cfg); err != nil {
		return nil, fmt.Errorf("config validation error: %w", err)
	}

	return &cfg, nil
}
