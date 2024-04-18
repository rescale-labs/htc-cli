package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	parseArgs()
}

func parseArgs() {
	validateArgs()
	if os.Args[1] == "upload" {
		upload()
	} else {
		download()
	}
}

func upload() {
	uploadCmd := flag.NewFlagSet("upload", flag.ExitOnError)
	src := uploadCmd.String("src", "", "the source path to upload from")
	dest := uploadCmd.String("dest", "", "the destination bucket to upload to")
	region := uploadCmd.String("region", "", "the region to upload to")
	uploadCmd.Parse(os.Args[2:])

	fmt.Printf("Upload command: %s\n", uploadCmd.Args())
	fmt.Printf("src: %s\n", *src)
	fmt.Printf("dest: %s\n", *dest)
	fmt.Printf("region: %s\n", *region)
}

func download() {
	downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)
	src := downloadCmd.String("src", "", "the source bucket to download from")
	dest := downloadCmd.String("dest", "", "the destination path to download to")
	region := downloadCmd.String("region", "", "the region to download from")
	downloadCmd.Parse(os.Args[2:])

	fmt.Printf("Download command: %s\n", downloadCmd.Args())
	fmt.Printf("src: %s\n", *src)
	fmt.Printf("dest: %s\n", *dest)
	fmt.Printf("region: %s\n", *region)
}

func validateArgs() {
	if len(os.Args) < 2 {
		fmt.Println("expected 'upload' or 'download' subcommands")
		os.Exit(1)
	}

	if os.Args[1] != "upload" && os.Args[1] != "download" {
		fmt.Println("expected one of the following sub commands 'upload' or 'download'")
		os.Exit(1)
	}
}
