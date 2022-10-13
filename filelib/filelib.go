package filelib

import (
	"fmt"
	"time"
	"log"
	"os"
	"io"
	"bufio"
	"github.com/TwiN/go-color"
	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func S3Download(accessKey, secretAccessKey, bucketName, filename, dest string) {
    file, err := os.Create(dest)
    if err != nil {
        log.Fatalf("Unable to open file %q, %v", dest, err)
    }
    defer file.Close()

	os.Setenv("AWS_ACCESS_KEY_ID", accessKey)
	os.Setenv("AWS_SECRET_ACCESS_KEY", secretAccessKey)

	sess, _ := session.NewSession(&aws.Config{Region: aws.String("us-west-2")},)
	downloader := s3manager.NewDownloader(sess)
	numBytes, err := downloader.Download(file,&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
        Key:    aws.String(filename),
	})
    if err != nil {
        log.Fatalf("Unable to download item %q, %v", filename, err)
    }
    fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
}


func CheckFileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		log.Printf(
			color.Yellow +
            "Failed to read file: %v" +
            color.Reset,
            err,
        )
		return false
	}
	return true
}

func CreateDir(dirName string) {
	err := os.Mkdir(dirName, 0755)
	if err != nil {
		log.Printf(color.Red + "Could not create directory: %s" + color.Reset, dirName)
	}
}

func CopyFile(src, dst string) bool {
	in, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	out,err := os.Create(dst)
	if err != nil {
		log.Printf(color.Red + "Could not create file: %s" + color.Reset, dst)
		return false
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		log.Fatal(err)
	}

	return true
}

func BackupFile(filename string) {
	t := time.Now().Format("2006-01-02T15:04:05")

	in, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	out,err := os.Create(filename + t)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		log.Fatal(err)
	}
	
}

func WriteFile(filename string, content []string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	for i:=0; i<len(content);i++ {
		if (content[i] == "") {
			continue
		}else{
			file.WriteString(content[i] + "\n")
		}
	}
}

func ReadFile(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf(color.Yellow + "Could not find existing file: %s. It will be created\n" + color.Reset, filename)
		}else{
			log.Fatal(err)
		}
	}
	defer file.Close()

	var content []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		content = append(content, line)
	}

	return content
}

