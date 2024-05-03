package cli

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func ParseBucket(bucketPath string) (string, string, error) {
	if !strings.HasPrefix(bucketPath, "gs://") {
		return "", "", errors.New("Invalid bucket. Bucket must start with gs://")
	}

	re := regexp.MustCompile("^gs://(.+?)/(.+)")
	match := re.FindStringSubmatch(bucketPath)

	return match[1], strings.TrimRight(match[2], "/"), nil
}

func ParseArgs(args []string) (string, string) {
	help := flag.Bool("h", false, "help message")
	cmd := flag.NewFlagSet("cp", flag.ContinueOnError)
	err := cmd.Parse(args)

	if err != nil {
		log.Fatalf("error parsing args")
	}

	if *help || len(cmd.Args()) != 2 {
		Usage()
	}
	src := cmd.Arg(0)
	dest := cmd.Arg(1)
	return src, dest
}

func Usage() {
	usage := "Usage: htccli cp <src> <dst>\n\n"
	usage += "Available commands:\n"
	usage += "\tcp		uploads or downloads one or more files to a destination URL\n\n"
	usage += "cp arguments:\n"
	usage += "\tsrc		path or cloud storage URI\n"
	usage += "\tdest		path or cloud storage URI\n\n"
	usage += "Environment variables:\n"
	usage += "\tGCP_APPLICATION_CREDENTIALS		JSON string containing GCP application credentials\n"
	fmt.Printf(usage)

	os.Exit(1)
}
