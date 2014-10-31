package cli_test

import (
	"crypto/rand"
	"io/ioutil"
	"os"

	"github.com/mark-rushakoff/ss33/cli"
	"github.com/rlmcpherson/s3gof3r"

	. "github.com/mark-rushakoff/ss33/testhelper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const FILESIZE = 10000

var storageSet = LoadTestStorageSet()

func init() {
	debug := false
	s3gof3r.SetLogger(os.Stderr, "", 15, debug)
}

var _ = Describe("CLI", func() {
	Describe("put", func() {
		It("stores a new file in both the permanent and cache storage", func() {
			randomFileContent := make([]byte, FILESIZE)
			bytesWritten, err := rand.Read(randomFileContent)
			Expect(err).NotTo(HaveOccurred())
			Expect(bytesWritten).To(Equal(FILESIZE))

			file, err := ioutil.TempFile("", "ss33-test-random")
			Expect(err).NotTo(HaveOccurred())
			defer os.Remove(file.Name())
			ioutil.WriteFile(file.Name(), randomFileContent, 0600)

			destName := RandomString()
			defer PurgeFile(storageSet, destName)

			args := []string{
				"placeholder_for_prog_name",
				"put",
				"--file", file.Name(),

				"--permanent-endpoint", storageSet.Permanent.Endpoint,
				"--permanent-bucket", storageSet.Permanent.BucketName,
				"--permanent-key", destName,
				"--permanent-access-key-id", storageSet.Permanent.AccessKeyId,
				"--permanent-secret-access-key", storageSet.Permanent.SecretAccessKey,

				"--cache-endpoint", storageSet.Cache.Endpoint,
				"--cache-bucket", storageSet.Cache.BucketName,
				"--cache-key", destName,
				"--cache-access-key-id", storageSet.Cache.AccessKeyId,
				"--cache-secret-access-key", storageSet.Cache.SecretAccessKey,
			}

			app := cli.App()
			app.Run(args)

			AssertS3FileExistsWithContent(*storageSet.Permanent, destName, randomFileContent)
			AssertS3FileExistsWithContent(*storageSet.Cache, destName, randomFileContent)
		})
	})

	Describe("get", func() {
		It("populates the local file and cache when the file only exists in permanent storage", func() {
			randomFileContent := make([]byte, FILESIZE)
			bytesWritten, err := rand.Read(randomFileContent)
			Expect(err).NotTo(HaveOccurred())
			Expect(bytesWritten).To(Equal(FILESIZE))

			file, err := ioutil.TempFile("", "ss33-test-random")
			Expect(err).NotTo(HaveOccurred())
			defer os.Remove(file.Name())
			ioutil.WriteFile(file.Name(), randomFileContent, 0600)

			destName := RandomString()
			defer PurgeFile(storageSet, destName)

			PutFile(*storageSet.Permanent, destName, file.Name())
			AssertS3FileExistsWithContent(*storageSet.Permanent, destName, randomFileContent)

			os.Remove(file.Name())
			args := []string{
				"placeholder_for_prog_name",
				"get",
				"--file", file.Name(),

				"--permanent-endpoint", storageSet.Permanent.Endpoint,
				"--permanent-bucket", storageSet.Permanent.BucketName,
				"--permanent-key", destName,
				"--permanent-access-key-id", storageSet.Permanent.AccessKeyId,
				"--permanent-secret-access-key", storageSet.Permanent.SecretAccessKey,

				"--cache-endpoint", storageSet.Cache.Endpoint,
				"--cache-bucket", storageSet.Cache.BucketName,
				"--cache-key", destName,
				"--cache-access-key-id", storageSet.Cache.AccessKeyId,
				"--cache-secret-access-key", storageSet.Cache.SecretAccessKey,
			}

			app := cli.App()
			app.Run(args)

			AssertFileExistsWithContent(file.Name(), randomFileContent)
			AssertS3FileExistsWithContent(*storageSet.Cache, destName, randomFileContent)
		})

		It("gets the file from the cache without even checking permanent storage", func() {
			randomFileContent := make([]byte, FILESIZE)
			bytesWritten, err := rand.Read(randomFileContent)
			Expect(err).NotTo(HaveOccurred())
			Expect(bytesWritten).To(Equal(FILESIZE))

			file, err := ioutil.TempFile("", "ss33-test-random")
			Expect(err).NotTo(HaveOccurred())
			defer os.Remove(file.Name())
			ioutil.WriteFile(file.Name(), randomFileContent, 0600)

			destName := RandomString()
			defer PurgeFile(storageSet, destName)

			PutFile(*storageSet.Cache, destName, file.Name())
			AssertS3FileExistsWithContent(*storageSet.Cache, destName, randomFileContent)

			os.Remove(file.Name())
			args := []string{
				"placeholder_for_prog_name",
				"get",
				"--file", file.Name(),

				"--cache-endpoint", storageSet.Cache.Endpoint,
				"--cache-bucket", storageSet.Cache.BucketName,
				"--cache-key", destName,
				"--cache-access-key-id", storageSet.Cache.AccessKeyId,
				"--cache-secret-access-key", storageSet.Cache.SecretAccessKey,
			}

			app := cli.App()
			app.Run(args)

			AssertFileExistsWithContent(file.Name(), randomFileContent)
		})
	})
})
