package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"golang.org/x/oauth2/clientcredentials"
)

type Config struct {
	ClientID     string   `toml:"ClientID"`
	ClientSecret string   `toml:"ClientSecret"`
	Scopes       []string `toml:"Scopes"`
}

var (
	// These will get overridden goreleaser.
	buildVersion = "dev"
	buildCommit  = "none"
	buildDate    = "unknown"
)

var configFile = flag.String("config", "config.toml", "Config file name and path")
var showVersion = flag.Bool("version", false, "Show version and exit.")

func getClient() (*http.Client, error) {
	var config Config
	_, err := toml.DecodeFile(*configFile, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %v: %w", *configFile, err)
	}
	ctx := context.Background()
	conf := &clientcredentials.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Scopes:       config.Scopes,
		TokenURL:     "https://auth.globus.org/v2/oauth2/token",
	}

	client := conf.Client(ctx)
	return client, nil
}

func download(source, destination string) error {
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}
	res, err := client.Get(source)
	if err != nil {
		return fmt.Errorf("failed to download %v: %w", source, err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download %v: got status code: %v", source, res.StatusCode)
	}
	out, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("failed to open destination %v for writing: %w", destination, err)
	}
	defer out.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return fmt.Errorf("failed to save %v as %v: %w", source, destination, err)
	}
	return nil
}

func upload(source, destination string) error {
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}
	input, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open source %v for reading: %w", source, err)
	}
	info, err := input.Stat()
	if err != nil {
		return fmt.Errorf("failed to open source %v stats: %w", source, err)
	}
	fileSize := info.Size()

	req, err := http.NewRequest("PUT", destination, input)
	if err != nil {
		return fmt.Errorf("failed to generate request: %w", err)
	}
	req.ContentLength = fileSize

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to upload to %v: %w", destination, err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, err := io.ReadAll(res.Body)
		return fmt.Errorf("failed to upload to %v: got status code %v, response: %v, body read err=%v", destination, res.StatusCode, string(body), err)
	}
	return nil
}

func printVersion() {
	fmt.Printf(`
globus-http %v, git-%v. Built %v
`, buildVersion, buildCommit, buildDate)
}

func help() {
	fmt.Print(`
globus-http allows download and upload throught the Globus HTTP API.

Usage:

	Download
		globus-http [options] download https://g-123456.12345.1234.data.globus.org/filename.txt  filename.txt

	Upload
		globus-http [options] upload filename https://g-123456.12345.1234.data.globus.org/filename.txt

Options: 
`)
	flag.PrintDefaults()
	fmt.Print(`

The configuration file is a TOML file with the following fields:

ClientID = "your-client-id"
ClientSecret = "your-client-secret"
Scopes = ["scope1", "scope2", "scope3"]

`)
	printVersion()

}

func main() {
	flag.Usage = help
	flag.Parse()
	args := flag.Args()
	if *showVersion {
		printVersion()
		os.Exit(0)
	}
	if len(args) < 1 {
		fmt.Println("Missing command.")
		help()
		os.Exit(1)
	}

	switch args[0] {
	case "download":
		if len(args) < 3 {
			fmt.Println("Missing source and/or destination for download")
			help()
			os.Exit(1)
		}
		err := download(args[1], args[2])
		if err != nil {
			fmt.Println("ERROR: ", err)
			os.Exit(1)
		}
		os.Exit(0)

	case "upload":
		if len(args) < 3 {
			fmt.Println("Missing source and/or destination for upload")
			help()
			os.Exit(1)
		}
		err := upload(args[1], args[2])
		if err != nil {
			fmt.Println("ERROR: ", err)
			os.Exit(1)
		}
		os.Exit(0)
	}
}
