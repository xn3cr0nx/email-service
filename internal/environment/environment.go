package environment

var env *Env

type Env struct {
	Host         string
	ServiceName  string
	TemplateDir  string
	Port         int
	Debug        bool
	Env          string
	Rest         bool
	Concurrency  int
	Sender       string
	FrontendHost string
	Queue        string

	Backend string

	// redis related variables
	RedisHost     string
	RedisPort     int
	RedisPassword string
	RedisDB       int

	// kafka related variables
	KafkaAddresses []string
	KafkaTopic     string
	KafkaGroup     string

	// nats related variables
	NatsHost    string
	NatsPort    int
	NatsSubject string

	OtelExporterJaegerEnable     bool
	OtelExporterJaegerAgentHost  string
	OtelExporterJaegerAgentPort  int
	OtelExporterPrometheusEnable bool
	OtelExporterPrometheusPort   int
}

func New() *Env {
	return new(Env)
}

// Set assign the shared global environment object
func Set(e *Env) {
	env = e
}

// Get returns the shared global environment object
func Get() *Env {
	return env
}

// String return the actual environment (development, staging, production) of the endpoint
func (e *Env) String() string {
	return e.Env
}

// IsProduction returns true if the environment is production
func (e *Env) IsProduction() bool {
	return e.Env == "production"
}

// IsStaging returns true if the environment is staging
func (e *Env) IsStaging() bool {
	return e.Env == "staging"
}

// IsTest returns true if the environment is test
func (e *Env) IsTest() bool {
	return e.Env == "test"
}

// IsDevelopment returns true if the environment is development. If no environment is
// specified it is considered development by default.
func (e *Env) IsDevelopment() bool {
	return e.Env == "development" || e.Env == ""
}
