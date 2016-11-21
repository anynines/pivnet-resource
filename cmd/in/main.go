package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/go-pivnet/logshim"
	"github.com/anynines/pivnet-resource/concourse"
	"github.com/anynines/pivnet-resource/downloader"
	"github.com/anynines/pivnet-resource/filter"
	"github.com/anynines/pivnet-resource/gp"
	"github.com/anynines/pivnet-resource/in"
	"github.com/anynines/pivnet-resource/in/filesystem"
	"github.com/anynines/pivnet-resource/md5sum"
	"github.com/anynines/pivnet-resource/useragent"
	"github.com/anynines/pivnet-resource/validator"
	"github.com/robdimsdale/sanitizer"
)

var (
	// version is deliberately left uninitialized so it can be set at compile-time
	version string
)

func main() {
	if version == "" {
		version = "dev"
	}

	logger := log.New(os.Stderr, "", log.LstdFlags)

	logger.Printf("PivNet Resource version: %s", version)

	if len(os.Args) < 2 {
		log.Fatalf("not enough args - usage: %s <sources directory>", os.Args[0])
	}

	downloadDir := os.Args[1]

	var input concourse.InRequest
	err := json.NewDecoder(os.Stdin).Decode(&input)
	if err != nil {
		log.Fatalln(err)
	}

	sanitized := concourse.SanitizedSource(input.Source)
	logger.SetOutput(sanitizer.NewSanitizer(sanitized, os.Stderr))

	verbose := false
	ls := logshim.NewLogShim(logger, logger, verbose)

	logger.Printf("Creating download directory: %s", downloadDir)

	err = os.MkdirAll(downloadDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Exiting with error: %s", err)
	}

	err = validator.NewInValidator(input).Validate()
	if err != nil {
		log.Fatalf("Exiting with error: %s", err)
	}

	var endpoint string
	if input.Source.Endpoint != "" {
		endpoint = input.Source.Endpoint
	} else {
		endpoint = pivnet.DefaultHost
	}

	clientConfig := pivnet.ClientConfig{
		Host:      endpoint,
		Token:     input.Source.APIToken,
		UserAgent: useragent.UserAgent(version, "get", input.Source.ProductSlug),
	}

	client := gp.NewClient(
		clientConfig,
		ls,
	)

	d := downloader.NewDownloader(client, downloadDir, ls)
	fs := md5sum.NewFileSummer()

	f := filter.NewFilter(ls)

	fileWriter := filesystem.NewFileWriter(downloadDir, ls)

	response, err := in.NewInCommand(
		ls,
		client,
		f,
		d,
		fs,
		fileWriter,
	).Run(input)
	if err != nil {
		log.Fatalf("Exiting with error: %s", err)
	}

	err = json.NewEncoder(os.Stdout).Encode(response)
	if err != nil {
		log.Fatalf("Exiting with error: %s", err)
	}
}
