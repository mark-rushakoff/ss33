package cli

import (
	"os"

	cgcli "github.com/codegangsta/cli"
	"github.com/mark-rushakoff/ss33/bucketset"
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

				bucketSet := storageSetFromContext(c).BucketSet()
				bytesWritten, err := bucketSet.Upload(bucketset.BucketUpload{
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

				bucketSet := storageSetFromContext(c).BucketSet()
				bytesWritten, err := bucketSet.Download(bucketset.BucketDownload{
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

func storageSetFromContext(c *cgcli.Context) *bucketset.StorageSet {
	return &bucketset.StorageSet{
		Cache:     getCacheStorage(c),
		Permanent: getPermanentStorage(c),
	}
}

func getPermanentStorage(c *cgcli.Context) *bucketset.Storage {
	return getStorage("permanent", c)
}

func getCacheStorage(c *cgcli.Context) *bucketset.Storage {
	return getStorage("cache", c)
}

func getStorage(prefix string, c *cgcli.Context) *bucketset.Storage {
	return &bucketset.Storage{
		Endpoint:        c.String(prefix + "-endpoint"),
		BucketName:      c.String(prefix + "-bucket"),
		AccessKeyId:     c.String(prefix + "-access-key-id"),
		SecretAccessKey: c.String(prefix + "-secret-access-key"),
	}
}
