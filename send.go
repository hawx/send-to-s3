package main

import (
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"

	"flag"
	"fmt"
	"log"
	"mime"
	"os"
	"path"
)

const helpText = `Usage: send-to-s3 [options]

    --bucket <v>         # Bucket to upload to
    --region <v>         # Region to upload to (default: us-east-1)
    --name <v>           # Path to upload to (remote)
    --src <path>         # File to upload (local)

    --access-key         # AWS Access Key
    --secret-key         # AWS Secret Key

    --help               # Display this message
`

var (
	accessKey = flag.String("access-key", "", "AWS Access Key")
	secretKey = flag.String("secret-key", "", "AWS Secret Key")

	bucket = flag.String("bucket", "", "Bucket name")
	region = flag.String("region", "us-east-1", "Region name (if required)")
	name   = flag.String("name", "", "Path to upload to")
	src    = flag.String("src", "", "File to upload")

	help = flag.Bool("help", false, "Display help message")
)

func main() {
	flag.Parse()

	if *help {
		fmt.Println(helpText)
		return
	}

	if *bucket == "" || *name == "" || *src == "" {
		log.Fatal("Require --bucket, --name and --src")
	}

	file, err := os.Open(*src)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	auth, err := aws.GetAuth(*accessKey, *secretKey)
	if err != nil {
		log.Fatal(err)
	}

	client := s3.New(auth, aws.Regions[*region])
	b := client.Bucket(*bucket)

	typ := mime.TypeByExtension(path.Ext(*src))
	stat, _ := file.Stat()
	if err := b.PutReader(*name, file, stat.Size(), typ, s3.Private); err != nil {
		log.Fatal(err)
	}
}
