package main

import (
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Cleint ... S3のクライアント情報を格納
type S3Client struct {
	S3          *s3.S3
	Flag        Flag
	srcObjects  map[string][]string
	destObjects map[string][]string
}

// NewClient ... セッションを持つS3の新しいインスタンス作成
func NewClient(f Flag) (*S3Client, error) {
	c := &S3Client{
		Flag:        f,
		srcObjects:  make(map[string][]string, 0),
		destObjects: make(map[string][]string, 0),
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           f.AwsEnv,
	})

	if err != nil {
		return nil, err
	}

	c.S3 = s3.New(sess, aws.NewConfig().WithMaxRetries(10).WithRegion(f.Region))
	return c, nil
}

// ListObjects...バケットにあるオブジェクトの一覧を取得
func (c *S3Client) ListObjects(prefix, delimiter, bucket string) ([]*s3.ListObjectsV2Output, error) {
	var objects []*s3.ListObjectsV2Output
	pageNum := 0

	input := &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String(delimiter),
	}
	err := c.S3.ListObjectsV2Pages(input,
		func(output *s3.ListObjectsV2Output, lastPage bool) bool {
			pageNum++
			objects = append(objects, output)
			return !lastPage
		})
	if err != nil {
		return objects, err
	}

	return objects, nil
}

// CopyObject...srcからdestにSTANDARD_IAとしてcopyする
// ref: https://docs.aws.amazon.com/code-samples/latest/catalog/go-s3-s3_copy_object.go.html
func (c *S3Client) CopyObject(key string) error {
	input := &s3.CopyObjectInput{
		CopySource:   aws.String(url.PathEscape(fmt.Sprintf("%s/%s", c.Flag.Src, key))),
		Bucket:       aws.String(c.Flag.Dest),
		Key:          aws.String(key),
		StorageClass: aws.String("STANDARD_IA"),
	}
	if _, err := c.S3.CopyObject(input); err != nil {
		return err
	}

	if err := c.S3.WaitUntilObjectExists(&s3.HeadObjectInput{
		Bucket: aws.String(c.Flag.Dest),
		Key:    aws.String(key),
	}); err != nil {
		return err
	}
	return nil
}

// DeleteObjects...srcを削除する
// ref: https://docs.aws.amazon.com/code-samples/latest/catalog/go-s3-s3_delete_object.go.html
func (c *S3Client) DeleteObject(key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(c.Flag.Src),
		Key:    aws.String(key),
	}
	if _, err := c.S3.DeleteObject(input); err != nil {
		return err
	}

	if err := c.S3.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(c.Flag.Src),
		Key:    aws.String(key),
	}); err != nil {
		return err
	}
	return nil
}
