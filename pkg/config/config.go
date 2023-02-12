// SPDX-FileCopyrightText: 2022 Risk.Ident GmbH <contact@riskident.com>
// SPDX-FileCopyrightText: 2023 Kalle Fagerberg
//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by the
// Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License for
// more details.
//
// You should have received a copy of the GNU General Public License along
// with this program.  If not, see <http://www.gnu.org/licenses/>.

package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/dinkur/dinkur-desktop/internal/casing"
	"github.com/dinkur/dinkur/pkg/config"
	"github.com/invopop/jsonschema"
	"github.com/iver-wharf/wharf-core/v2/pkg/logger"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var log = logger.NewScoped("config")

var Path string

var Header = `# This file is automatically managed by Dinkur desktop.
# Recommended to only change settings via the application's settings panel.
# Any settings in this file may get overridden.

# yaml-language-server: $schema=https://github.com/dinkur/dinkur/raw/main/dinkur.schema.json`

var Default = Config{
	fileUsed: "(embedded defaults)",

	ExitOnWindowClose: false,
	Sqlite: Sqlite{
		Path:  config.Default.Sqlite.Path,
		Mkdir: true,
	},
	Client: ClientTypeSqlite,
	GRPC: GRPC{
		Address: "localhost:59122",
	},
	Daemon: Daemon{
		Enabled:     false,
		BindAddress: "localhost:59122",
	},
	Log: Log{
		Format: LogFormatPretty,
		Level:  LogLevel(logger.LevelInfo),
		Color:  LogColorAuto,
	},
}

func init() {
	cfgPath, err := os.UserConfigDir()
	if err != nil {
		panic(fmt.Errorf("resolve user config directory: %w", err))
	}
	Path = filepath.Join(cfgPath, "dinkur-desktop.yaml")
}

type Config struct {
	fileUsed string

	ExitOnWindowClose bool

	Client ClientType
	Sqlite Sqlite
	GRPC   GRPC
	Daemon Daemon

	Log Log
}

func (c *Config) FileUsed() string {
	return c.fileUsed
}

type Sqlite struct {
	// Path is the file path of where to store the sqlite database file, i.e
	// the file containing all the time-tracked entries.
	Path string
	// Mkdir will enable creating any missing directories for the data
	// directory, if set to true. Will fail if directories don't exist and
	// this is set to false.
	Mkdir bool
}

type GRPC struct {
	// Address defines which IP/hostname and port to reach the API on.
	Address string
}

type Daemon struct {
	Enabled bool
	// BindAddress defines which IP/hostname and port to serve the gRPC API on.
	// Can be set to 0.0.0.0 as IP to allow access from any IP.
	BindAddress string
}

type Log struct {
	// Format defines how the logs are printed to the console, either "pretty"
	// for human readable, or "json" for machine readable.
	Format LogFormat
	// Level defines the logging severity level. All log messages below this
	// config will not be logged.
	Level LogLevel
	// Color defines if the log output should be colored. Defaults to "auto",
	// where it will only use colors if it detects interactive TTY, but
	// options "always" and "never" can override this.
	Color LogColor
}

type jsonSchemaInterface interface {
	JSONSchema() *jsonschema.Schema
}

func (c *Config) Save() error {
	file, err := os.Create(Path)
	if err != nil {
		return err
	}
	defer file.Close()
	fmt.Fprintln(file, Header)
	enc := yaml.NewEncoder(file)
	defer enc.Close()
	return enc.Encode(c)
}

func ReadAuto(v *viper.Viper) (*Config, error) {
	if err := AddDefaults(v); err != nil {
		return nil, err
	}
	if err := AddAuto(v); err != nil {
		return nil, err
	}
	var cfg Config
	if err := Unmarshal(v, &cfg); err != nil {
		return nil, fmt.Errorf("decoding config file: %w", err)
	}
	return &cfg, nil
}

func ReadFile(v *viper.Viper, file string) (*Config, error) {
	if err := AddDefaults(v); err != nil {
		return nil, err
	}
	if err := AddFile(v, file); err != nil {
		return nil, err
	}
	var cfg Config
	if err := Unmarshal(v, &cfg); err != nil {
		return nil, fmt.Errorf("decoding config file: %w", err)
	}
	return &cfg, nil
}

func AddAuto(v *viper.Viper) error {
	if err := AddFile(v, Path); err != nil {
		if os.IsNotExist(err) || errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return nil
		}
		return err
	}
	return nil
}

// AddFile uses the [Default] configs and then merges the config from a
// specific file. The function will error if the file does not exist.
func AddFile(v *viper.Viper, file string) error {
	v.SetConfigName("dinkur-desktop")
	v.SetConfigType("yaml")
	v.SetConfigFile(file)
	return v.MergeInConfig()
}

func AddDefaults(v *viper.Viper) error {
	v.SetConfigType("yaml")
	b, err := yaml.Marshal(Default)
	if err != nil {
		return err
	}
	return v.MergeConfig(bytes.NewReader(b))
}

// Unmarshal will read the config file from [github.com/spf13/viper] using
// the appropriate decoding options.
func Unmarshal(v *viper.Viper, cfg *Config) error {
	err := v.Unmarshal(cfg, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		mapstructure.TextUnmarshallerHookFunc(),
		mapstructure.StringToTimeDurationHookFunc(), // default hook
		mapstructure.StringToSliceHookFunc(","),     // default hook
	)))
	cfg.fileUsed = v.ConfigFileUsed()
	return err
}

// JSONSchema returns the JSON schema struct for the [Config] struct.
func JSONSchema() *jsonschema.Schema {
	r := new(jsonschema.Reflector)
	r.KeyNamer = casing.ToCamelCase
	r.Namer = func(t reflect.Type) string {
		return casing.ToCamelCase(t.Name())
	}
	r.RequiredFromJSONSchemaTags = true
	s := r.Reflect(&Config{})
	s.ID = "https://github.com/dinkur/dinkur/raw/main/dinkur.schema.json"
	return s
}
