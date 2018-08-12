package main

import (
	"flag"
	"log"
	"strings"
	"encoding/json"
	"github.com/mozey/logutil"
	"io/ioutil"
	"fmt"
	"os"
	"path"
	"sort"
)

type ConfigMap map[string]string

// AppDir is the application root
var AppDir string
// Prefix for env vars
var Prefix = "APP"

// Flags
var Env *string
var Update *bool

type ArgMap []string

func (a *ArgMap) String() string {
	return strings.Join(*a, ", ")
}
func (a *ArgMap) Set(value string) error {
	*a = append(*a, value)
	return nil
}

var Keys ArgMap
var Values ArgMap

func Cmd() {
	// If not compiled with ldflags see if AppDir is set on env
	appDirKey := fmt.Sprintf("%v_DIR", Prefix)
	if AppDir == "" {
		AppDir = os.Getenv(appDirKey)
	}

	var config string
	config = path.Join(AppDir, fmt.Sprintf("config.%v.json", *Env))

	b, err := ioutil.ReadFile(config)
	if err != nil {
		logutil.Debugf("Loading config from: %v", config)
		log.Panic(err)
	}

	// The config file must have a flat key value structure
	c := ConfigMap{}
	err = json.Unmarshal(b, &c)
	if err != nil {
		log.Panic(err)
	}

	// Set existing config Keys
	var configKeys []string
	for k := range c {
		configKeys = append(configKeys, k)
	}

	// Sort
	sort.Strings(configKeys)

	if Prefix == "" {
		log.Panicf("Prefix must not be empty")
	}

	if len(Keys) > 0 {
		// Set config key value

		// Validate input
		for i, key := range Keys {
			if !strings.HasPrefix(key, Prefix) {
				log.Panicf("Key must strart with prefix: %v", Prefix)
			}

			if i > len(Values)-1 {
				log.Panicf("Missing value for key: %v", key)
			}
			value := Values[i]

			// Set key value
			c[key] = value
		}

		// Update config
		b, _ := json.MarshalIndent(c, "", "    ")
		if *Update {
			logutil.Debugf("Config updated: %v", config)
			ioutil.WriteFile(config, b, 0)
		} else {
			// Print json
			fmt.Print(string(b))
		}

	} else {
		// Unset env var starting with Prefix
		for _, v := range os.Environ() {
			a := strings.Split(v, "=")
			if len(a) == 2 {
				key := a[0]
				//value := a[1]
				if strings.HasPrefix(key, Prefix) {
					fmt.Println(fmt.Sprintf("unset %v", key))
				}
			}
		}
		// Print commands to set env
		for _, key := range configKeys {
			fmt.Println(fmt.Sprintf("export %v=%v", key, c[key]))
		}
	}
}

func main() {
	log.SetFlags(log.Lshortfile)

	Env = flag.String("env", "dev", "Specify config file to use")
	flag.Var(&Keys, "key", "Set key and print config JSON")
	flag.Var(&Values, "value", "Value for last key specified")
	Update = flag.Bool("update", false, "Update config.json")
	flag.Parse()

	Cmd()
}
