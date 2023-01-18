package aws

import (
	"fmt"

	"os"

	"path/filepath"

	// uuid
	uuid "github.com/google/uuid"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	s3session *s3.S3
)

func init() {
	// load .env
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// set env vars
	os.Setenv("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID"))
	os.Setenv("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY"))
	os.Setenv("AWS_SESSION_TOKEN", os.Getenv("AWS_SESSION_TOKEN"))

	// use credentials from .env
	// where are the credentials stored?
	// ~/.aws/credentials
	s3session = s3.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
	})))

}

func ListBuckets() (resp *s3.ListBucketsOutput) {
	res, err := s3session.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		fmt.Println(err)
	}
	return res
}

func Gets3(c *gin.Context) {
	fmt.Println(ListBuckets())
}

// gin route
func UploadImage(c *gin.Context, orgId string, name string) (URL string) {
	fileHeader, err := c.FormFile("file")

	if err != nil {
		c.JSON(500, gin.H{
			"message": "error",
		})
	}

	file, err := fileHeader.Open()

	if err != nil {
		c.JSON(500, gin.H{
			"message": "error",
		})
	}

	defer file.Close()
	// get the mimetype
	ending := filepath.Ext(fileHeader.Filename)
	uuid := uuid.New().String()
	endingWithoutDot := ending[1:]
	fmt.Println(endingWithoutDot)
	// upload to s3
	_, err = s3session.PutObject(&s3.PutObjectInput{
		Bucket:             aws.String("rereal-ad-creatives"),
		Key:                aws.String("uuid" + uuid + ending),
		ContentDisposition: aws.String("inline"),
		ContentType:        aws.String("image/jpeg"),
		Body:               file,
	})

	// return the url
	returnURL := "https://rereal-ad-creatives.s3.amazonaws.com/" + "uuid" + uuid + ending
	return returnURL
}

func ListObjects() (resp *s3.ListObjectsOutput) {
	res, err := s3session.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String("rereal-ad-creatives"),
	})
	if err != nil {
		fmt.Println(err)
	}
	return res
}

func GetObjects(c *gin.Context) {
	fmt.Println(ListObjects())
}
