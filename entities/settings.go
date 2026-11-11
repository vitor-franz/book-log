package entities

// RuntimeEnvironment identifies the environment in which the application runs.
type RuntimeEnvironment string

const (
	// EnvironmentDevelopment selects development behavior and colored debug logs.
	EnvironmentDevelopment RuntimeEnvironment = "dev"

	// EnvironmentHomologation selects staging behavior and JSON info logs.
	EnvironmentHomologation RuntimeEnvironment = "homo"

	// EnvironmentProduction selects production behavior and JSON info logs.
	EnvironmentProduction RuntimeEnvironment = "prod"
)

// ApplicationSettings mirrors the top-level tables in settings.toml.
type ApplicationSettings struct {
	// DatabaseSettings contains the MySQL connection and pool settings.
	DatabaseSettings DatabaseSettings `toml:"DatabaseSettings"`

	// JWTSettings contains token signing and authentication cookie configuration.
	JWTSettings JWTSettings `toml:"JWTSettings"`

	// RuntimeSettings contains the application execution environment.
	RuntimeSettings RuntimeSettings `toml:"RuntimeSettings"`

	// LoggerSettings contains the logging output configuration.
	LoggerSettings LoggerSettings `toml:"LoggerSettings"`
}

// DatabaseSettings mirrors the [DatabaseSettings] table in settings.toml.
// The named database must already exist before the application starts.
type DatabaseSettings struct {
	// Host is the MySQL server hostname or IP address.
	Host string `toml:"Host"`

	// Port is the MySQL server's TCP port.
	Port int64 `toml:"Port"`

	// Username is the MySQL account used for queries and migrations.
	Username string `toml:"Username"`

	// Password is the MySQL account password.
	Password string `toml:"Password"`

	// Name is the existing database that receives queries and migrations.
	Name string `toml:"Name"`
}

// JWTSettings mirrors the [JWTSettings] table in settings.toml.
type JWTSettings struct {
	// Secret is the HS256 signing secret and must contain at least 32 bytes.
	Secret string `toml:"Secret"`

	// InsecureCookie permits cookies over HTTP for local development only.
	// The default false value requires HTTPS.
	InsecureCookie bool `toml:"InsecureCookie"`
}

// RuntimeSettings mirrors the [RuntimeSettings] table in settings.toml.
type RuntimeSettings struct {
	// Environment selects the application environment: dev, homo, or prod.
	Environment RuntimeEnvironment `toml:"Environment"`
}

// LoggerSettings mirrors the [LoggerSettings] table in settings.toml.
type LoggerSettings struct {
	// OutputPath overrides the prod/homo output destination; empty uses stderr.
	OutputPath string `toml:"OutputPath"`
}
