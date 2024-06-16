package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/ccin2p3/go-freeipa/freeipa"
	"github.com/infra-monkey/go-freeipa-importer/importer"
	"github.com/spf13/viper"
)

var (
	host      string
	admin     string
	password  string
	insecure  bool
	configdir string
)

type freeipaConfig struct {
	hostname           string
	principal          string
	principalPassword  string
	insecureConnection bool
}

// Check if a flag "name" is passed as argument
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

// Load configuration from a configuration file.
// This bypasses any parameter given as argument.
func loadConfigFromFile(path string) (freeipaConfig, error) {
	// Check is path exist
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(path)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	// load variables from file
	config := freeipaConfig{viper.GetString("hostname"), viper.GetString("principal"), viper.GetString("principalPassword"), viper.GetBool("insecureConnection")}
	// return the freeipaConfig struct
	return config, nil
}

func main() {
	log.SetPrefix("go-freeipa-importer: ")
	log.SetFlags(0)

	// Read command line arguments
	flag.StringVar(&host, "host", "ipa.ipatest.lan", "The hostname of the FreeIPA host to import resources from.")
	flag.StringVar(&admin, "principal", "admin", "The username of FreeIPA to use to import resources from.")
	flag.StringVar(&password, "password", "P@ssword", "The password of FreeIPA to use to import resources from.")
	flag.BoolVar(&insecure, "insecure", false, "Allow to skip certificate checks.")
	flag.StringVar(&configdir, "configdir", ".", "The hostname of the FreeIPA host to import resources from.")
	flag.Parse()
	var config freeipaConfig
	var err error
	if isFlagPassed("configdir") {
		config, err = loadConfigFromFile(configdir)
		if err != nil {
			log.Fatal(err.Error())
		}
	} else {
		if isFlagPassed("host") && isFlagPassed("admin") && isFlagPassed("password") {
			config = freeipaConfig{host, admin, password, insecure}
		} else {
			log.Fatal("Invalid command line argument combination.")
		}
	}
	//log.Println(config)

	// initialize a freeipa client
	tspt := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.insecureConnection, // WARNING DO NOT USE THIS OPTION IN PRODUCTION
		},
	}
	client, err := freeipa.Connect(config.hostname, tspt, config.principal, config.principalPassword)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Start importing resources
	log.Println("Start import.")

	//Ensure the "output" directory exists.
	_ = os.Mkdir("./output", os.ModePerm)

	//Ensure the "output" directory exists.
	_ = os.Mkdir("./scripts", os.ModePerm)

	// Import users
	err = importer.ImportUsers(client)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("End of FreeIPA import.")
}
