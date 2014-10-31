package testhelper

import (
	"crypto/rand"
	"encoding/base32"
	"io"
	"io/ioutil"
	"os"

	"github.com/mark-rushakoff/ss33/bucketset"
	"github.com/rlmcpherson/s3gof3r"

	. "github.com/onsi/gomega"
)

func LoadTestStorageSet() *bucketset.StorageSet {
	configPath := os.Getenv("SS33_TEST_CONFIG")
	if configPath == "" {
		panic("Environment variable SS33_TEST_CONFIG was not set; tests cannot continue.")
	}

	storageSet, err := bucketset.StorageSetFromFile(configPath)
	if err != nil {
		panic(err)
	}

	return storageSet
}

func RandomString() string {
	sourceLength := 35
	bytes := make([]byte, sourceLength)
	_, err := rand.Read(bytes)
	Expect(err).NotTo(HaveOccurred())
	return base32.StdEncoding.EncodeToString(bytes)
}

func bucketFromStorage(storage bucketset.Storage) *s3gof3r.Bucket {
	s3Client := s3gof3r.New(storage.Endpoint, s3gof3r.Keys{AccessKey: storage.AccessKeyId, SecretKey: storage.SecretAccessKey})
	return s3Client.Bucket(storage.BucketName)
}

func PurgeFile(storageSet *bucketset.StorageSet, key string) {
	permanentBucket := bucketFromStorage(*storageSet.Permanent)
	permanentBucket.Delete(key)

	cacheBucket := bucketFromStorage(*storageSet.Cache)
	cacheBucket.Delete(key)
}

func AssertS3FileExistsWithContent(storage bucketset.Storage, key string, expectedContent []byte) {
	bucket := bucketFromStorage(storage)
	reader, _, err := bucket.GetReader(key, bucket.Config)
	Expect(err).NotTo(HaveOccurred())
	defer reader.Close()

	content, err := ioutil.ReadAll(reader)
	Expect(err).NotTo(HaveOccurred())

	Expect(content).To(Equal(expectedContent))
}

func AssertFileExistsWithContent(path string, expectedContent []byte) {
	content, err := ioutil.ReadFile(path)
	Expect(err).NotTo(HaveOccurred())

	Expect(content).To(Equal(expectedContent))
}

func PutFile(storage bucketset.Storage, key string, localFilePath string) {
	bucket := bucketFromStorage(storage)

	writer, err := bucket.PutWriter(key, nil, bucket.Config)
	Expect(err).NotTo(HaveOccurred())
	defer writer.Close()

	file, err := os.Open(localFilePath)
	Expect(err).NotTo(HaveOccurred())
	defer file.Close()

	bytesWritten, err := io.Copy(writer, file)
	Expect(err).NotTo(HaveOccurred())

	stat, err := os.Stat(localFilePath)
	Expect(err).NotTo(HaveOccurred())
	Expect(stat.Size()).To(Equal(bytesWritten))
}
