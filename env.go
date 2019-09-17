package consistentcommit

import "github.com/theothertomelliott/consistentcommit/executor"

type EnvProvider struct {
	DefaultEnv func(string) string
}

func (e EnvProvider) env(command executor.Command) func(string) string {
	return func(key string) string {
		if value, exists := command.Env[key]; exists {
			return value
		}
		if e.DefaultEnv != nil {
			return e.DefaultEnv(key)
		}
		return ""
	}
}
