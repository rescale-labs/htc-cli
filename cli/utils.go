package cli

import (
	"flag"
	"log"
	"os"
	"regexp"
	"strings"
)

func ParseBucket(bucketPath string) (string, string) {
	if !strings.HasPrefix(bucketPath, "gs://") {
		log.Fatal("Invalid bucket. Bucket must start with gs://")
	}

	re := regexp.MustCompile("gs://(.+?)/(.+)")
	match := re.FindStringSubmatch(bucketPath)

	return match[1], match[2]
}

func ParseArgs() (string, string) {
	cmd := flag.NewFlagSet("cp", flag.ContinueOnError)
	err := cmd.Parse(os.Args[2:])
	if err != nil || len(cmd.Args()) != 2 {
		log.Fatalf("error parsing args")
	}
	src := cmd.Arg(0)
	dest := cmd.Arg(1)
	return src, dest
}
