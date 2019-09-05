package consistentcommit

type EnvProvider struct {
	DefaultEnv func(string) string
}

func (e EnvProvider) env(command Command) func(string) string {
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
