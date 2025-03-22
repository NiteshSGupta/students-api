package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// struct keyowrd and inner keyowrd have first word capital to automatic export the data
// we have to require of go clean env package for sync the local.yml and this file
// this file is going to make ready to use the configuration of local.yml, and the belowe struct are like boxes which are going to get data from MustLoad function
// this is called struct tags : `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
type HttpServer struct {
	Addr string
}
type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
	Storagepath string `yaml:"storage_path" env-required:"true"`
	HttpServer  `yaml:"http_server"`
}

// ths function name is MustLoad, and here we di't return any error , if this function not run , then there no meaning
func MustLoad() *Config {

	var configPath string

	//we are going to pass environment varirable in configPath

	// os.Getenv is a little helper in Go (from the os package) that lets your program ask the computer, “Hey, do you have a special note for me?” These special notes are called environment variables, and they’re like sticky notes you can leave for your program before it starts running.
	// os: This is a toolbox in Go that helps your program talk to the computer (like opening files or asking about these notes).
	// Getenv: Short for “get environment variable”—it’s the tool that looks for a specific note by its name.
	configPath = os.Getenv("CONFIG_PATH")

	//here we checing that configPath is passed or not
	if configPath == "" {

		//flags are passing like : go run cmd/students-api/main.go -config-path xyz, so we are checking also that CONFIG_PATH are in flag
		// this will work when we run this program like this : go run main.go -config config/local.yaml

		// flag.String: This is a helper from the flag toolbox in Go. It sets up a spot to hold text (a string) that someone might give when they run your program, like -config
		//flags: This is the name of the box you made to hold that text. But here’s the tricky part:
		// flag.String doesn’t give you the text directly—it gives you a pointer to the text.
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()

		//The * is like saying, “Go to the address and get what’s there!”
		// flags is the pointer (the address, like “123 Toy Street”).
		// *flags means “look at the address and take out the actual thing” (the text, like the toy car).
		configPath = *flags

		if configPath == "" {
			log.Fatal("config path is no set")
		}

	}

	//here we go when provided theconfigPath but now we are going to check that provided path have the file or not which we are finding : config/local.yaml
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist :%s", configPath)
	}

	// cfg is name and Config is struct
	// cfg is like a holder (or a box) that’s shaped exactly like your Config struct. It’s ready to hold all the stuff you defined in Config, like Env, Storagepath, and HttpServer.
	var cfg Config

	//so here we done that cleanenv package is serializing the configPath data in cfg which is Config struct
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("can not read config file:%s", err.Error())
	}

	return &cfg

}
