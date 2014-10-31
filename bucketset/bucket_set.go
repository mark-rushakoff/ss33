package bucketset

import (
	"io"

	"github.com/codegangsta/cli"
	"github.com/rlmcpherson/s3gof3r"
)

type BucketSet struct {
	cache     *s3gof3r.Bucket
	permanent *s3gof3r.Bucket
}

type BucketUpload struct {
	CacheKey     string
	PermanentKey string
	Source       io.Reader
}

type BucketDownload struct {
	CacheKey     string
	PermanentKey string
	Destination  io.Writer
}

func (b *BucketSet) Upload(settings BucketUpload) (bytesWritten int64, err error) {
	permanentWriter, err := b.permanent.PutWriter(settings.PermanentKey, nil, b.permanent.Config)
	if err != nil {
		return 0, err
	}
	defer permanentWriter.Close()

	cacheWriter, err := b.cache.PutWriter(settings.CacheKey, nil, b.cache.Config)
	if err != nil {
		return 0, err
	}
	defer cacheWriter.Close()

	everythingWriter := io.MultiWriter(permanentWriter, cacheWriter)

	return io.Copy(everythingWriter, settings.Source)
}

func (b *BucketSet) Download(settings BucketDownload) (bytesWritten int64, err error) {
	cacheReader, _, err := b.cache.GetReader(settings.CacheKey, b.cache.Config)
	if err != nil {
		return b.WarmCacheAndDownload(settings)
	}
	defer cacheReader.Close()

	return io.Copy(settings.Destination, cacheReader)
}

func (b *BucketSet) WarmCacheAndDownload(settings BucketDownload) (bytesWritten int64, err error) {
	permanentReader, _, err := b.permanent.GetReader(settings.PermanentKey, b.permanent.Config)
	if err != nil {
		return 0, err
	}
	defer permanentReader.Close()

	cacheWriter, err := b.cache.PutWriter(settings.CacheKey, nil, b.cache.Config)
	if err != nil {
		return 0, err
	}
	defer cacheWriter.Close()

	everythingWriter := io.MultiWriter(settings.Destination, cacheWriter)

	return io.Copy(everythingWriter, permanentReader)
}

func BucketSetFromContext(c *cli.Context) *BucketSet {
	return &BucketSet{
		cache:     getCacheBucket(c),
		permanent: getPermanentBucket(c),
	}
}

func getPermanentBucket(c *cli.Context) *s3gof3r.Bucket {
	return getBucket("permanent", c)
}

func getCacheBucket(c *cli.Context) *s3gof3r.Bucket {
	return getBucket("cache", c)
}

func getBucket(prefix string, c *cli.Context) *s3gof3r.Bucket {
	permanentClient := s3gof3r.New(c.String(prefix+"-endpoint"), s3gof3r.Keys{
		AccessKey: c.String(prefix + "-access-key-id"),
		SecretKey: c.String(prefix + "-secret-access-key"),
	})
	return permanentClient.Bucket(c.String(prefix + "-bucket"))
}
