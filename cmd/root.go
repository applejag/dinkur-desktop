package cmd

import (
	"fmt"
	"os"

	"github.com/dinkur/dinkur-desktop/internal/console"
	"github.com/dinkur/dinkur-desktop/internal/license"
	"github.com/dinkur/dinkur-desktop/pkg/app"
	"github.com/dinkur/dinkur-desktop/pkg/config"
	"github.com/fatih/color"
	"github.com/iver-wharf/wharf-core/v2/pkg/logger"
	"github.com/iver-wharf/wharf-core/v2/pkg/logger/consolejson"
	"github.com/iver-wharf/wharf-core/v2/pkg/logger/consolepretty"
	"github.com/mattn/go-colorable"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfg     = config.Default
	cfgFile string

	log = logger.NewScoped("Dinkur desktop")
)

var rootFlags = struct {
	verbose bool

	showLicenseWarranty   bool
	showLicenseConditions bool
}{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "dinkur-desktop",
	Version: "0.1.0-preview",
	Short:   "Dinkur desktop",
	Long: license.Header + `
Track how you spend time on your entries with Dinkur.
`,
	SilenceErrors: true,
	SilenceUsage:  true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return readConfig(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		switch {
		case rootFlags.showLicenseWarranty:
			fmt.Println(license.Warranty)
		case rootFlags.showLicenseConditions:
			fmt.Println(license.Conditions)
		default:
			return app.Run(&cfg)
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Set up logger initially, before real config is read
	initLogger()

	err := rootCmd.Execute()
	if err != nil {
		log.Error().Messagef("Failed: %s", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initLogger)

	rootCmd.SetOut(colorable.NewColorableStdout())
	rootCmd.SetErr(colorable.NewColorableStderr())
	rootCmd.SetUsageTemplate(console.UsageTemplate())

	rootCmd.Flags().BoolVar(&rootFlags.showLicenseConditions, "license-c", false, "show program's license conditions")
	rootCmd.Flags().BoolVar(&rootFlags.showLicenseWarranty, "license-w", false, "show program's license warranty")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", cfgFile, "config file")

	rootCmd.PersistentFlags().Var(&cfg.Client, "client", `Dinkur client: "sqlite" or "grpc"`)

	rootCmd.PersistentFlags().String("sqlite.path", cfg.Sqlite.Path, "database file")
	rootCmd.PersistentFlags().Bool("sqlite.mkdir", cfg.Sqlite.Mkdir, "create directory for data if it doesn't exist")

	rootCmd.PersistentFlags().String("grpc.address", cfg.GRPC.Address, "address of connecting to Dinkur daemon gRPC API")
	rootCmd.PersistentFlags().String("daemon.bindAddress", cfg.Daemon.BindAddress, "bind address of hosting Dinkur daemon gRPC API")

	rootCmd.PersistentFlags().Var(&cfg.Log.Level, "log.level", `logging severity: "debug", "info", "warn", "error", or "panic"`)
	rootCmd.PersistentFlags().Var(&cfg.Log.Format, "log.format", `logging format: "pretty" or "json"`)
	rootCmd.PersistentFlags().Var(&cfg.Log.Color, "log.color", `logging colored output: "auto", "always", or "never"`)

	rootCmd.PersistentFlags().BoolVarP(&rootFlags.verbose, "verbose", "v", rootFlags.verbose, `enables debug logging (short for --log.level=debug)`)
}

func readConfig(cmd *cobra.Command) error {
	v := viper.New()
	if err := v.BindPFlags(cmd.Root().PersistentFlags()); err != nil {
		return err
	}

	var newCfg *config.Config
	var err error
	if cmd.Flag("config").Changed {
		newCfg, err = config.ReadFile(v, cfgFile)
	} else {
		newCfg, err = config.ReadAuto(v)
	}
	if err != nil {
		log.Warn().WithError(err).Message("Failed loading config. Continuing with default config.")
		cfg = config.Default
	} else {
		cfg = *newCfg
	}

	// Set up logger again, now that we've read in the new config
	initLogger()

	log.Debug().
		WithString("file", cfg.FileUsed()).
		Message("Loaded configuration.")

	return nil
}

func initLogger() {
	logger.ClearOutputs()
	level := logger.Level(cfg.Log.Level)
	if rootFlags.verbose {
		level = logger.LevelDebug
	}
	switch cfg.Log.Color {
	case config.LogColorAuto:
		// Do nothing, fatih/color is on auto by default
	case config.LogColorNever:
		color.NoColor = true
	case config.LogColorAlways:
		color.NoColor = false
	}
	if cfg.Log.Format == config.LogFormatPretty {
		prettyConf := consolepretty.DefaultConfig
		prettyConf.DisableDate = true
		prettyConf.DisableCaller = true
		prettyConf.Writer = colorable.NewColorableStderr()
		logger.AddOutput(level, consolepretty.New(prettyConf))
	} else {
		logger.AddOutput(level, consolejson.Default)
	}
}
