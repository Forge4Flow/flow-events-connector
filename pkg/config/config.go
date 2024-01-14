package config

import (
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/onflow/flow-go-sdk/access/http"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/viper"
)

const (
	DefaultMySQLDatastoreMigrationSource     = "github://forge4flow/forge4flow-core/migrations/datastore/mysql"
	DefaultMySQLEventstoreMigrationSource    = "github://forge4flow/forge4flow-core/migrations/eventstore/mysql"
	DefaultPostgresDatastoreMigrationSource  = "github://forge4flow/forge4flow-core/migrations/datastore/postgres"
	DefaultPostgresEventstoreMigrationSource = "github://forge4flow/forge4flow-core/migrations/eventstore/postgres"
	DefaultSQLiteDatastoreMigrationSource    = "github://forge4flow/forge4flow-core/migrations/datastore/sqlite"
	DefaultSQLiteEventstoreMigrationSource   = "github://forge4flow/forge4flow-core/migrations/eventstore/sqlite"
	DefaultAuthenticationUserIdClaim         = "sub"
	DefaultSessionTokenLength                = 32
	DefaultSessionIdleTimeout                = 15 * time.Minute
	DefaultSessionExpTimeout                 = 24 * time.Hour
	DefaultAppIdentifier                     = "Forge4Flow IAM Service"
	DefaultFlowNetwork                       = "emulator"
	DefaultAutoRegister                      = false
	PrefixForge4Flow                         = "forge4flow"
	ConfigFileName                           = "forge4flow.yaml"
)

type Config interface {
	GetPort() int
	GetLogLevel() int8
	GetEnableAccessLog() bool
	GetAutoMigrate() bool
	GetDatastore() *DatastoreConfig
	GetEventstore() *EventstoreConfig
	GetAuthentication() *AuthConfig
	GetAppIdentifier() string
}

type Forge4FlowConfig struct {
	CoreInstall     bool              `mapstructure:"coreInstall"`
	Port            int               `mapstructure:"port"`
	LogLevel        int8              `mapstructure:"logLevel"`
	EnableAccessLog bool              `mapstructure:"enableAccessLog"`
	AutoMigrate     bool              `mapstructure:"autoMigrate"`
	Datastore       *DatastoreConfig  `mapstructure:"datastore"`
	Eventstore      *EventstoreConfig `mapstructure:"eventstore"`
	Authentication  *AuthConfig       `mapstructure:"authentication"`
	AppIdentifier   string            `mapstructure:"appIdentifier"`
	FlowNetwork     string            `mapstructure:"flowNetwork"`
	AdminAccount    string            `mapstructure:"adminAccount"`
}

func (forge4FlowConfig Forge4FlowConfig) GetPort() int {
	return forge4FlowConfig.Port
}

func (forge4FlowConfig Forge4FlowConfig) GetLogLevel() int8 {
	return forge4FlowConfig.LogLevel
}

func (forge4FlowConfig Forge4FlowConfig) GetEnableAccessLog() bool {
	return forge4FlowConfig.EnableAccessLog
}

func (forge4FlowConfig Forge4FlowConfig) GetAutoMigrate() bool {
	return forge4FlowConfig.AutoMigrate
}

func (forge4FlowConfig Forge4FlowConfig) GetDatastore() *DatastoreConfig {
	return forge4FlowConfig.Datastore
}

func (forge4FlowConfig Forge4FlowConfig) GetEventstore() *EventstoreConfig {
	return forge4FlowConfig.Eventstore
}

func (forge4FlowConfig Forge4FlowConfig) GetAuthentication() *AuthConfig {
	return forge4FlowConfig.Authentication
}

func (forge4FlowConfig Forge4FlowConfig) GetAppIdentifier() string {
	return forge4FlowConfig.AppIdentifier
}

func (forge4FlowConfig Forge4FlowConfig) GetFlowNetwork() string {
	return forge4FlowConfig.FlowNetwork
}

type DatastoreConfig struct {
	MySQL    *MySQLConfig    `mapstructure:"mysql"`
	Postgres *PostgresConfig `mapstructure:"postgres"`
	SQLite   *SQLiteConfig   `mapstructure:"sqlite"`
}

type MySQLConfig struct {
	Username           string `mapstructure:"username"`
	Password           string `mapstructure:"password"`
	Hostname           string `mapstructure:"hostname"`
	Database           string `mapstructure:"database"`
	MigrationSource    string `mapstructure:"migrationSource"`
	MaxIdleConnections int    `mapstructure:"maxIdleConnections"`
	MaxOpenConnections int    `mapstructure:"maxOpenConnections"`
}

type PostgresConfig struct {
	Username           string `mapstructure:"username"`
	Password           string `mapstructure:"password"`
	Hostname           string `mapstructure:"hostname"`
	Database           string `mapstructure:"database"`
	SSLMode            string `mapstructure:"sslmode"`
	MigrationSource    string `mapstructure:"migrationSource"`
	MaxIdleConnections int    `mapstructure:"maxIdleConnections"`
	MaxOpenConnections int    `mapstructure:"maxOpenConnections"`
}

type SQLiteConfig struct {
	Database           string `mapstructure:"database"`
	InMemory           bool   `mapstructure:"inMemory"`
	MigrationSource    string `mapstructure:"migrationSource"`
	MaxIdleConnections int    `mapstructure:"maxIdleConnections"`
	MaxOpenConnections int    `mapstructure:"maxOpenConnections"`
}

type EventstoreConfig struct {
	MySQL             *MySQLConfig    `mapstructure:"mysql"`
	Postgres          *PostgresConfig `mapstructure:"postgres"`
	SQLite            *SQLiteConfig   `mapstructure:"sqlite"`
	SynchronizeEvents bool            `mapstructure:"synchronizeEvents"`
}

type AuthConfig struct {
	ApiKey             string              `mapstructure:"apiKey"`
	AutoRegister       bool                `mapstructure:"autoRegister"`
	SessionTokenLength int64               `mapstructure:"sessionTokenLength"`
	SessionIdleTimeout time.Duration       `mapstructure:"sessionIdleTimeout"`
	SessionExpTimeout  int64               `mapstructure:"sessionExpTimeout"`
	Provider           *AuthProviderConfig `mapstructure:"providers"`
}

type AuthProviderConfig struct {
	Name          string `mapstructure:"name"`
	PublicKey     string `mapstructure:"publicKey"`
	UserIdClaim   string `mapstructure:"userIdClaim"`
	TenantIdClaim string `mapstructure:"tenantIdClaim"`
}

func NewConfig() Forge4FlowConfig {
	viper.SetConfigFile(ConfigFileName)
	viper.SetDefault("coreInstall", false)
	viper.SetDefault("port", 8000)
	viper.SetDefault("logLevel", zerolog.DebugLevel)
	viper.SetDefault("enableAccessLog", true)
	viper.SetDefault("autoMigrate", false)
	viper.SetDefault("datastore.mysql.migrationSource", DefaultMySQLDatastoreMigrationSource)
	viper.SetDefault("datastore.postgres.migrationSource", DefaultPostgresDatastoreMigrationSource)
	viper.SetDefault("datastore.sqlite.migrationSource", DefaultSQLiteDatastoreMigrationSource)
	viper.SetDefault("eventstore.mysql.migrationSource", DefaultMySQLEventstoreMigrationSource)
	viper.SetDefault("eventstore.postgres.migrationSource", DefaultPostgresEventstoreMigrationSource)
	viper.SetDefault("eventstore.sqlite.migrationSource", DefaultSQLiteEventstoreMigrationSource)
	viper.SetDefault("eventstore.synchronizeEvents", false)
	viper.SetDefault("authentication.provider.userIdClaim", DefaultAuthenticationUserIdClaim)
	viper.SetDefault("authentication.sessionTokenLength", DefaultSessionTokenLength)
	viper.SetDefault("authentication.sessionIdleTimeout", DefaultSessionIdleTimeout)
	viper.SetDefault("authentication.sessionExpTimeout", DefaultSessionExpTimeout)
	viper.SetDefault("appIdentifier", DefaultAppIdentifier)
	viper.SetDefault("flowNetwork", DefaultFlowNetwork)
	viper.SetDefault("authentication.autoRegister", DefaultAutoRegister)

	// If config file exists, use it
	_, err := os.ReadFile(ConfigFileName)
	if err == nil {
		if err := viper.ReadInConfig(); err != nil {
			log.Fatal().Err(err).Msg("Error while reading forge4flow.yaml. Shutting down.")
		}
	} else {
		if os.IsNotExist(err) {
			log.Info().Msg("Could not find forge4flow.yaml. Attempting to use environment variables.")
		} else {
			log.Fatal().Err(err).Msg("Error while reading forge4flow.yaml. Shutting down.")
		}
	}

	var config Forge4FlowConfig
	// If available, use env vars for config
	for _, fieldName := range getFlattenedStructFields(reflect.TypeOf(config)) {
		envKey := strings.ToUpper(fmt.Sprintf("%s_%s", PrefixForge4Flow, strings.ReplaceAll(fieldName, ".", "_")))
		envVar := os.Getenv(envKey)
		if envVar != "" {
			viper.Set(fieldName, envVar)
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal().Err(err).Msg("Error while creating config. Shutting down.")
	}

	// Configure logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.DurationFieldUnit = time.Millisecond
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.SetGlobalLevel(zerolog.Level(config.GetLogLevel()))
	if zerolog.GlobalLevel() == zerolog.DebugLevel {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	if config.GetAuthentication() == nil || config.GetAuthentication().ApiKey == "" {
		log.Fatal().Msg("Must provide an API key to authenticate incoming requests to Warrant.")
	}

	// Check the selected network host value and assign it to the appropriate constant
	flowNetwork := config.GetFlowNetwork()
	switch flowNetwork {
	case "emulator":
		config.FlowNetwork = http.EmulatorHost
	case "testnet":
		config.FlowNetwork = http.TestnetHost
	case "mainnet":
		config.FlowNetwork = http.MainnetHost
	default:
		if !isValidFlowURL(flowNetwork) {
			log.Fatal().Msgf("Invalid flowNetwork parameter: %s - valid options are: emulator, testnet, mainnet, or valid Access Note HTTP URL", flowNetwork)
		}
	}

	return config
}

func isValidFlowURL(flowURL string) bool {
	_, err := url.ParseRequestURI(flowURL)

	// TODO: Check is valid access node

	return err == nil
}

func getFlattenedStructFields(t reflect.Type) []string {
	return getFlattenedStructFieldsHelper(t, []string{})
}

func getFlattenedStructFieldsHelper(t reflect.Type, prefixes []string) []string {
	unwrappedT := t
	if t.Kind() == reflect.Pointer {
		unwrappedT = t.Elem()
	}

	flattenedFields := make([]string, 0)
	for i := 0; i < unwrappedT.NumField(); i++ {
		field := unwrappedT.Field(i)
		fieldName := field.Tag.Get("mapstructure")
		switch field.Type.Kind() {
		case reflect.Struct, reflect.Pointer:
			flattenedFields = append(flattenedFields, getFlattenedStructFieldsHelper(field.Type, append(prefixes, fieldName))...)
		default:
			flattenedField := fieldName
			if len(prefixes) > 0 {
				flattenedField = fmt.Sprintf("%s.%s", strings.Join(prefixes, "."), fieldName)
			}
			flattenedFields = append(flattenedFields, flattenedField)
		}
	}

	return flattenedFields
}
