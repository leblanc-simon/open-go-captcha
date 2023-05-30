package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"leblanc.io/open-go-captcha/config"
	"leblanc.io/open-go-captcha/connection"
	"leblanc.io/open-go-captcha/crypto"
	"leblanc.io/open-go-captcha/log"
	"leblanc.io/open-go-captcha/request"
)

// Args command-line parameters
type Args struct {
	ConfigPath string
}

var cfg config.Config


func ProcessArgs(cfg interface{}) Args {
	var a Args

	f := flag.NewFlagSet("OpenGoCaptcha", 1)

	f.StringVar(&a.ConfigPath, "c", "config.yml", "Path to configuration file")

	fu := f.Usage
	f.Usage = func() {
		fu()
		envHelp, _ := cleanenv.GetDescription(cfg, nil)
		fmt.Fprintln(f.Output())
		fmt.Fprintln(f.Output(), envHelp)
	}

	f.Parse(os.Args[1:])

	return a
}

func main () {
	args := ProcessArgs(&cfg)
	
	// read configuration from the file and environment variables
	if _, err := os.Stat(args.ConfigPath); errors.Is(err, os.ErrNotExist) {
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	} else {
		if err := cleanenv.ReadConfig(args.ConfigPath, &cfg); err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	}

	// Initialize Logger
	log.Initialize(&cfg)

	// Initialize Redis
	connection.Initialize(&cfg)
	defer connection.GetRedisInstance().Close()

	// Initialize Crypt
	crypto.Initialize(&cfg)

	// Initialize request
	request.Initialize(&cfg)

	http.HandleFunc("/", request.GetCaptcha)
	http.HandleFunc("/check", request.CheckCaptcha)
	http.HandleFunc("/confirm", request.CheckCaptcha)

	fmt.Println(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port), nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}