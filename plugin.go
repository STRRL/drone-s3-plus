package main

import (
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/mattn/go-zglob"
)

// Plugin defines the S3 plugin parameters.
type Plugin struct {
	Endpoint string
	Key      string
	Secret   string
	Bucket   string

	// overwrite object
	Overwrite bool

	// if not "", enable server-side encryption
	// valid values are:
	//     AES256
	//     aws:kms
	Encryption string

	// us-east-1
	// us-west-1
	// us-west-2
	// eu-west-1
	// ap-southeast-1
	// ap-southeast-2
	// ap-northeast-1
	// sa-east-1
	Region string

	// Indicates the files ACL, which should be one
	// of the following:
	//     private
	//     public-read
	//     public-read-write
	//     authenticated-read
	//     bucket-owner-read
	//     bucket-owner-full-control
	Access string

	// Sets the Cache-Control header on each uploaded object
	CacheControl string

	// upload files
	Parallel int

	// Copies the files from the specified directory.
	// Regexp matching will apply to match multiple
	// files
	//
	// Examples:
	//    /path/to/file
	//    /path/to/*.txt
	//    /path/to/*/*.txt
	//    /path/to/**
	Source string
	Target string

	// Strip the prefix from the target path
	StripPrefix string

	// Exclude files matching this pattern.
	Exclude []string

	// Use path style instead of domain style.
	//
	// Should be true for minio and false for AWS.
	PathStyle bool
	// Dry run without uploading/
	DryRun bool
}

func (p *Plugin) Upload(cli *s3.S3, match string) error {
	stat, err := os.Stat(match)
	if err != nil {
		return err // should never happen
	}

	// skip directories
	if stat.IsDir() {
		return nil
	}

	target := filepath.Join(p.Target, strings.TrimPrefix(match, p.StripPrefix))
	if !strings.HasPrefix(target, "/") {
		target = "/" + target
	}

	// amazon S3 has pretty crappy default content-type headers so this pluign
	// attempts to provide a proper content-type.
	content := contentType(match)

	// when executing a dry-run we exit because we don't actually want to
	// upload the file to S3.
	if p.DryRun {
		return nil
	}

	f, err := os.Open(match)
	if err != nil {
		fmt.Printf("opening %s failed, %s\n", match, err)
		return err
	}

	defer f.Close()

	putObjectInput := &s3.PutObjectInput{
		Body:        f,
		Bucket:      &(p.Bucket),
		Key:         &target,
		ACL:         &(p.Access),
		ContentType: &content,
	}

	if p.Encryption != "" {
		putObjectInput.ServerSideEncryption = &(p.Encryption)
	}

	if p.CacheControl != "" {
		putObjectInput.CacheControl = &(p.CacheControl)
	}

	_, err = cli.PutObject(putObjectInput)

	if err != nil {
		fmt.Printf("could not upload file %q to %q, err: %s\n",
			match, p.Bucket+"/"+p.Target, err)
		return err
	}

	return f.Close()
}

// Exec runs the plugin
func (p *Plugin) Exec() error {
	// normalize the target URL
	if strings.HasPrefix(p.Target, "/") {
		p.Target = p.Target[1:]
	}

	// create the client
	conf := &aws.Config{
		Region:           aws.String(p.Region),
		Endpoint:         &p.Endpoint,
		DisableSSL:       aws.Bool(strings.HasPrefix(p.Endpoint, "http://")),
		S3ForcePathStyle: aws.Bool(p.PathStyle),
	}

	if p.Key != "" && p.Secret != "" {
		conf.Credentials = credentials.NewStaticCredentials(p.Key, p.Secret, "")
	} else {
		fmt.Println("AWS Key and/or Secret not provided (falling back to ec2 instance profile)")
	}

	sess := session.Must(session.NewSession(conf))
	client := s3.New(sess)

	matches, err := matches(p.Source, p.Exclude)
	if err != nil {
		fmt.Println("Could not match files")
		return err
	}

	fmt.Println("Attempting to upload files")
	for _, match := range matches {
		fmt.Printf("%-32s --> %s/%s\n", match, p.Bucket, p.Target)
	}

	if p.Parallel == 0 {
		p.Parallel = runtime.NumCPU()
	}

	var wg sync.WaitGroup
	taskChan := make(chan string)

	for i := 0; i < p.Parallel; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for match := range taskChan {
				err := p.Upload(client, match)
				if err != nil {
					fmt.Printf("upload %s failed, %s\n", match, err)
				} else {
					fmt.Printf("upload %s successful\n", match)
				}
			}
		}()
	}

	for _, match := range matches {
		taskChan <- match
	}

	close(taskChan)
	wg.Wait()

	return nil
}

// matches is a helper function that returns a list of all files matching the
// included Glob pattern, while excluding all files that matche the exclusion
// Glob pattners.
func matches(include string, exclude []string) ([]string, error) {
	matches, err := zglob.Glob(include)
	if err != nil {
		return nil, err
	}
	if len(exclude) == 0 {
		return matches, nil
	}

	// find all files that are excluded and load into a map. we can verify
	// each file in the list is not a member of the exclusion list.
	excludem := map[string]bool{}
	for _, pattern := range exclude {
		excludes, err := zglob.Glob(pattern)
		if err != nil {
			return nil, err
		}
		for _, match := range excludes {
			excludem[match] = true
		}
	}

	var included []string
	for _, include := range matches {
		_, ok := excludem[include]
		if ok {
			continue
		}
		included = append(included, include)
	}
	return included, nil
}

// contentType is a helper function that returns the content type for the file
// based on extension. If the file extension is unknown application/octet-stream
// is returned.
func contentType(path string) string {
	ext := filepath.Ext(path)
	typ := mime.TypeByExtension(ext)
	if typ == "" {
		typ = "application/octet-stream"
	}
	return typ
}
