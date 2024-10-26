package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
)

var (
	env     *Environment
	envOnce sync.Once
)

type EnvVar struct {
	Name     string
	Required bool
}

type Environment struct {
	Port string `env:"PORT"`
}

func GetEnv() *Environment {
	envOnce.Do(
		func() {
			var err error
			env, err = loadEnv()
			if err != nil {
				log.Fatal(err)
			}
		})
	return env
}

func loadEnv() (*Environment, error) {
	env := &Environment{}
	v := reflect.ValueOf(env).Elem()
	t := v.Type()

	var missingVars []string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("env")

		if tag == "" {
			continue
		}

		parts := splitTag(tag)
		name := parts[0]
		required := len(parts) > 1 && parts[1] == "required"

		value := os.Getenv(name)
		if required && value == "" {
			missingVars = append(missingVars, name)
		}

		v.Field(i).SetString(value)
	}

	if len(missingVars) > 0 {
		return nil, fmt.Errorf("missing required environment variables: \n%s", strings.Join(missingVars, "\n"))
	}
	return env, nil
}

func splitTag(tag string) []string {
	result := []string{}
	current := ""
	for _, char := range tag {
		if char == ',' {
			result = append(result, current)
			current = ""
		} else {
			current += string(char)
		}
	}
	result = append(result, current)
	return result
}
