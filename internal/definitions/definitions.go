package definitions

type Logger interface {
	Error(args... interface{})
	Errorf(format string, args... interface{})
	Info(args... interface{})
	Infof(format string, args... interface{})
}

type Config interface {
	StringWithDefaults(key string, defaultValue string) string
	IntWithDefaults(key string, defaultValue int) int
}

type Service interface {
	Run(addr ...string)
}