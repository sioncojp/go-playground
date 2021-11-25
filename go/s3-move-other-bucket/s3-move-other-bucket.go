package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type Flag struct {
	AwsEnv    string
	Region    string
	Parallel  int64
	Src       string
	Dest      string
	Dir       string
	ID        string
	BeforeDay int
	Delete    bool
	Check     bool
}

// S3MoveOtherBucket...src/dest bucketのcopy/delete/checkの実施
func S3MoveOtherBucket(f Flag) error {
	c, err := NewClient(f)
	if err != nil {
		return err
	}

	if err := c.BucketCopyDeleteCheck(f.Src); err != nil {
		return err
	}

	// Check flagがtrueの場合、src/destのオブジェクト数出力させる
	if f.Check {
		if err := c.BucketCopyDeleteCheck(f.Dest); err != nil {
			return err
		}

		var srcNumberObject int
		var destNumberObject int

		for _, v := range c.srcObjects {
			srcNumberObject += len(v)
		}

		for _, v := range c.destObjects {
			destNumberObject += len(v)
		}

		fmt.Printf("src number of object: %d\n", srcNumberObject)
		fmt.Printf("dest number of object: %d\n", destNumberObject)
	}

	return nil
}

// BucketCopyDeleteCheck...bucketに対してsrcのcopy/delete/checkの実施
func (c *S3Client) BucketCopyDeleteCheck(bucket string) error {
	var IDs []string

	if c.Flag.ID != "" {
		IDs = []string{fmt.Sprintf("%s/", c.Flag.ID)}
	} else {
		// idを取得するため、バケットのトップディレクトリを取得する
		bucketTopDirs, err := c.ListObjects("", "/", bucket)
		if err != nil {
			return err
		}

		// idの格納。1/ 2/ 3/...のように最後にslash入りがゲットされる
		for _, v := range bucketTopDirs {
			for _, vv := range v.CommonPrefixes {
				IDs = append(IDs, aws.StringValue(vv.Prefix))
			}
		}
	}

	fmt.Printf("%s number of all ID: %d\n", bucket, len(IDs))

	// ここからidをベースに並列で動かす
	sem := semaphore.NewWeighted(c.Flag.Parallel)
	eg := errgroup.Group{}
	var mu sync.RWMutex

	for _, id := range IDs {
		if err := sem.Acquire(context.Background(), 1); err != nil {
			fmt.Printf("failed to acquire semaphore: %v\n", err)
			break
		}
		id := id
		eg.Go(func() error {
			defer sem.Release(1)
			objectList, err := c.CopyDeleteCheck(id, bucket)
			if err != nil {
				return err
			}
			mu.Lock()
			if bucket == c.Flag.Src {
				c.srcObjects[id] = append(c.srcObjects[id], objectList[id]...)
			} else {
				c.destObjects[id] = append(c.destObjects[id], objectList[id]...)
			}
			mu.Unlock()

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

// CopyDeleteCheck...copy / delete / Check　を行う
func (c *S3Client) CopyDeleteCheck(id, bucket string) (map[string][]string, error) {
	// a/aa配下の全キー取得
	copyKey := make(map[string][]string, 0)

	outputImage, err := c.ListObjects(fmt.Sprintf("%s%s/aa/", id, c.Flag.Dir), "/", bucket)
	if err != nil {
		return copyKey, err
	}
	for _, v := range outputImage {
		for _, vv := range v.Contents {
			// x日前のデータを対象とする。ただしdestはcopyしたばかりの可能性があるので無視する
			if isBeforeDay(vv.LastModified, c.Flag.BeforeDay) || bucket == c.Flag.Dest {
				copyKey[id] = append(copyKey[id], aws.StringValue(vv.Key))
			}
		}
	}

	// a/bb配下の全キー取得
	outputMovie, err := c.ListObjects(fmt.Sprintf("%s%s/bb/", id, c.Flag.Dir), "/", bucket)
	if err != nil {

		return copyKey, err
	}
	for _, v := range outputMovie {
		for _, vv := range v.Contents {
			// x日前のデータを対象とする。ただしdestはcopyしたばかりの可能性があるので無視する
			if isBeforeDay(vv.LastModified, c.Flag.BeforeDay) || bucket == c.Flag.Dest {
				copyKey[id] = append(copyKey[id], aws.StringValue(vv.Key))
			}
		}
	}

	// Check flagがtrueの場合、src/destのオブジェクト数出力させるだけなので、delteもcopyもしないので、ここでreturn
	if c.Flag.Check {
		return copyKey, nil
	}

	// Delete flagがtrueの場合は、元のファイルを削除してreturn
	if c.Flag.Delete {
		fmt.Printf("delete now %s...\n", id)
		for _, v := range copyKey[id] {
			if err := c.DeleteObject(v); err != nil {
				return copyKey, err
			}
		}
		// 最後にa/ ディレクトリを消す
		if err := c.DeleteObject(fmt.Sprintf("%s%s/", id, c.Flag.Dir)); err != nil {
			return copyKey, err
		}
		return copyKey, nil
	}

	// copyする
	fmt.Printf("copy now %s...\n", id)
	for _, v := range copyKey[id] {
		if err := c.CopyObject(v); err != nil {
			return copyKey, err
		}
	}
	return copyKey, err
}
