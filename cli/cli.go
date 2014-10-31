package cli

import (
	"io"
	"os"

	cgcli "github.com/codegangsta/cli"
	"github.com/rlmcpherson/s3gof3r"
)

func App() *cgcli.App {
	app := cgcli.NewApp()
	app.Name = "ss33"
	app.Usage = "Put and get objects from a permanent and cache S3-compatible storage"

	cacheIoFlags := []cgcli.Flag{
		cgcli.StringFlag{
			Name:  "file",
			Usage: "Path to local file to upload",
		},

		cgcli.StringFlag{
			Name:  "permanent-endpoint",
			Usage: "Endpoint of permanent storage server",
			Value: "s3.amazonaws.com",
		},
		cgcli.StringFlag{
			Name:  "permanent-bucket",
			Usage: "Bucket on permanent storage",
		},
		cgcli.StringFlag{
			Name:  "permanent-access-key-id",
			Usage: "Access Key ID for permanent storage",
		},
		cgcli.StringFlag{
			Name:  "permanent-secret-access-key",
			Usage: "Secret Access Key for permanent storage",
		},
		cgcli.StringFlag{
			Name:  "permanent-key",
			Usage: "Key (path within bucket) for permanent storage",
		},

		cgcli.StringFlag{
			Name:  "cache-endpoint",
			Usage: "Endpoint of cache storage server",
		},
		cgcli.StringFlag{
			Name:  "cache-bucket",
			Usage: "Bucket on cache storage",
		},
		cgcli.StringFlag{
			Name:  "cache-access-key-id",
			Usage: "Access Key ID for cache storage",
		},
		cgcli.StringFlag{
			Name:  "cache-secret-access-key",
			Usage: "Secret Access Key for cache storage",
		},
		cgcli.StringFlag{
			Name:  "cache-key",
			Usage: "Key (path within bucket) for cache storage",
		},
	}

	app.Commands = []cgcli.Command{
		cgcli.Command{
			Name:  "put",
			Usage: "Put a local file onto both the permanent and cache storage",
			Flags: cacheIoFlags,

			Action: func(c *cgcli.Context) {
				localFile, err := os.Open(c.String("file"))
				if err != nil {
					panic(err)
				}

				bucketSet := bucketSetFromContext(c)

				bytesWritten, err := bucketSet.Upload(bucketUpload{
					CacheKey:     c.String("cache-key"),
					PermanentKey: c.String("permanent-key"),
					Source:       localFile,
				})
				if err != nil {
					panic(err)
				}

				stat, err := os.Stat(c.String("file"))
				if err != nil {
					panic(err)
				}

				expectedBytesWritten := stat.Size()
				if bytesWritten != expectedBytesWritten {
					println("Expected bytes written", expectedBytesWritten)
					println("Actual bytes written", bytesWritten)
				}
			},
		},
		cgcli.Command{
			Name:  "get",
			Usage: "Get a file and ensure there is a local copy and a copy in the cache",
			Flags: cacheIoFlags,

			Action: func(c *cgcli.Context) {
				localFile, err := os.Create(c.String("file"))
				if err != nil {
					panic(err)
				}
				defer localFile.Close()

				bucketSet := bucketSetFromContext(c)
				bytesWritten, err := bucketSet.Download(bucketDownload{
					CacheKey:     c.String("cache-key"),
					PermanentKey: c.String("permanent-key"),
					Destination:  localFile,
				})
				if err != nil {
					panic(err)
				}
				localFile.Close()

				stat, err := os.Stat(c.String("file"))
				if err != nil {
					panic(err)
				}

				expectedBytesWritten := stat.Size()
				if bytesWritten != expectedBytesWritten {
					println("Expected bytes written", expectedBytesWritten)
					println("Actual bytes written", bytesWritten)
				}
			},
		},
	}

	return app
}

type bucketSet struct {
	Cache     *s3gof3r.Bucket
	Permanent *s3gof3r.Bucket
}

type bucketUpload struct {
	CacheKey     string
	PermanentKey string
	Source       io.Reader
}

type bucketDownload struct {
	CacheKey     string
	PermanentKey string
	Destination  io.Writer
}

func (b *bucketSet) Upload(settings bucketUpload) (bytesWritten int64, err error) {
	permanentWriter, err := b.Permanent.PutWriter(settings.PermanentKey, nil, b.Permanent.Config)
	if err != nil {
		return 0, err
	}
	defer permanentWriter.Close()

	cacheWriter, err := b.Cache.PutWriter(settings.CacheKey, nil, b.Cache.Config)
	if err != nil {
		return 0, err
	}
	defer cacheWriter.Close()

	everythingWriter := io.MultiWriter(permanentWriter, cacheWriter)

	return io.Copy(everythingWriter, settings.Source)
}

func (b *bucketSet) Download(settings bucketDownload) (bytesWritten int64, err error) {
	cacheReader, _, err := b.Cache.GetReader(settings.CacheKey, b.Cache.Config)
	if err != nil {
		return b.WarmCacheAndDownload(settings)
	}
	defer cacheReader.Close()

	return io.Copy(settings.Destination, cacheReader)
}

func (b *bucketSet) WarmCacheAndDownload(settings bucketDownload) (bytesWritten int64, err error) {
	permanentReader, _, err := b.Permanent.GetReader(settings.PermanentKey, b.Permanent.Config)
	if err != nil {
		return 0, err
	}
	defer permanentReader.Close()

	cacheWriter, err := b.Cache.PutWriter(settings.CacheKey, nil, b.Cache.Config)
	if err != nil {
		return 0, err
	}
	defer cacheWriter.Close()

	everythingWriter := io.MultiWriter(settings.Destination, cacheWriter)

	return io.Copy(everythingWriter, permanentReader)
}

func bucketSetFromContext(c *cgcli.Context) *bucketSet {
	return &bucketSet{
		Cache:     getCacheBucket(c),
		Permanent: getPermanentBucket(c),
	}
}

func getPermanentBucket(c *cgcli.Context) *s3gof3r.Bucket {
	return getBucket("permanent", c)
}

func getCacheBucket(c *cgcli.Context) *s3gof3r.Bucket {
	return getBucket("cache", c)
}

func getBucket(prefix string, c *cgcli.Context) *s3gof3r.Bucket {
	permanentClient := s3gof3r.New(c.String(prefix+"-endpoint"), s3gof3r.Keys{
		AccessKey: c.String(prefix + "-access-key-id"),
		SecretKey: c.String(prefix + "-secret-access-key"),
	})
	return permanentClient.Bucket(c.String(prefix + "-bucket"))
}
