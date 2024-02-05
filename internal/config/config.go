package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/viper"
)

const (
	DefaultDatastoreMigrationSource = "github://forge4flow/flow-events-connector/internal/database/migrations"
	DefaultUserAgent                = "forge4flow/flow-events-connector@v0.0.1"
	DefaultTopic                    = "flow-events"
	DefaultRebuildInterval          = 5 * time.Minute
	DefaultRebuildTimeout           = 60 * time.Second
	DefaultTopicAnnotationDelimiter = ","
	DefaultGatewayBasicAuth         = true
	DefaultGatewayURL               = "http://gateway:8080"
	DefaultPrintResponseBody        = true
	DefaultMainnetAccessNode        = "access.mainnet.nodes.onflow.org:9000"
	DefaultTestnetAccessNode        = "access.testnet.nodes.onflow.org:9000"
	DefaultCrescendoAccessNode      = "access.crescendo.nodes.onflow.org:9000"
)

type FlowEventsConnectorConfig struct {
	LogLevel    int8             `mapstructure:"logLevel"`
	AutoMigrate bool             `mapstructure:"autoMigrate"`
	Datastore   *DatastoreConfig `mapstructure:"datastore"`
	Flow        *FlowConfig      `mapstructure:"flow"`
	ControllerConfig
}

type DatastoreConfig struct {
	Username           string `mapstructure:"username"`
	Password           string `mapstructure:"password"`
	Hostname           string `mapstructure:"hostname"`
	Database           string `mapstructure:"database"`
	SSLMode            string `mapstructure:"sslmode"`
	MigrationSource    string `mapstructure:"migrationSource"`
	MaxIdleConnections int    `mapstructure:"maxIdleConnections"`
	MaxOpenConnections int    `mapstructure:"maxOpenConnections"`
}

type ControllerConfig struct {
	UpstreamTimeout          time.Duration `mapstructure:"upstreamTimeout"`
	GatewayURL               string        `mapstructure:"gatewayURL"`
	PrintResponse            bool          `mapstructure:"printResponse"`
	PrintResponseBody        bool          `mapstructure:"printResponseBody"`
	PrintRequestBody         bool          `mapstructure:"printRequestBody"`
	RebuildInterval          time.Duration `mapstructure:"rebuildInterval"`
	RebuildTimeout           time.Duration `mapstructure:"rebuildTimeout"`
	TopicAnnotationDelimiter string        `mapstructure:"topicAnnontationDelimiter"`
	AsyncFunctionInvocation  bool          `mapstructure:"asyncFunctionInvocation"`
	PrintSync                bool          `mapstructure:"printSync"`
	ContentType              string        `mapstructure:"contentType"`
	BasicAuth                bool          `mapstructure:"basicAuth"`
	UserAgent                string        `mapstructure:"userAgent"`
	Topic                    string        `mapstructure:"topic"`
}

type FlowConfig struct {
	MainnetAccessNode   string `mapstructure:"mainnetAccessNode"`
	TestnetAccessNode   string `mapstructure:"testnetAccessNode"`
	CrescendoAccessNode string `mapstructure:"crescendoAccessNode"`
	EmulatorAccessNode  string `mapstructure:"emulatorAccessNode"`
	UseCrescendo        bool   `mapstructure:"useCrescendo"`
	UseEmulator         bool   `mapstructure:"useEmulator"`
}

func NewConfig() FlowEventsConnectorConfig {
	viper.SetDefault("logLevel", zerolog.DebugLevel)
	viper.SetDefault("autoMigrate", false)
	viper.SetDefault("datastore.migrationSource", DefaultDatastoreMigrationSource)
	viper.SetDefault("controller.userAgent", DefaultUserAgent)
	viper.SetDefault("controller.topic", DefaultTopic)
	viper.SetDefault("rebuildInterval", DefaultRebuildInterval)
	viper.SetDefault("rebuildTimeout", DefaultRebuildTimeout)
	viper.SetDefault("topicAnnontationDelimiter", DefaultTopicAnnotationDelimiter)
	viper.SetDefault("basicAuth", DefaultGatewayBasicAuth)
	viper.SetDefault("gatewayURL", DefaultGatewayURL)
	viper.SetDefault("printResponseBody", DefaultPrintResponseBody)
	viper.SetDefault("flow.mainnetAccessNode", DefaultMainnetAccessNode)
	viper.SetDefault("flow.testnetAccessNode", DefaultTestnetAccessNode)
	viper.SetDefault("flow.crescendoAccessNode", DefaultCrescendoAccessNode)
	viper.SetDefault("flow.useCrescendo", false)
	viper.SetDefault("flow.useEmulator", false)

	var config FlowEventsConnectorConfig
	// If available, use env vars for config
	for _, fieldName := range getFlattenedStructFields(reflect.TypeOf(config)) {
		envKey := strings.ToUpper(strings.ReplaceAll(fieldName, ".", "_"))
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
	zerolog.SetGlobalLevel(zerolog.Level(config.LogLevel))
	if zerolog.GlobalLevel() == zerolog.DebugLevel {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	if config.Flow.UseEmulator && config.Flow.EmulatorAccessNode == "" {
		log.Fatal().Msg("You must specify EmulatorAccessNode when UseEmulator is set to true")
	}

	return config
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
