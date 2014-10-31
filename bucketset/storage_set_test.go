package bucketset_test

import (
	"github.com/mark-rushakoff/ss33/bucketset"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StorageSet", func() {
	Describe("Merge", func() {
		It("Overwrites values in the receiver", func() {
			initial := &bucketset.StorageSet{
				Permanent: &bucketset.Storage{
					Endpoint:        "",
					BucketName:      "",
					AccessKeyId:     "",
					SecretAccessKey: "",
				},
				Cache: &bucketset.Storage{
					Endpoint:        "c.example.com",
					BucketName:      "c-bucket",
					AccessKeyId:     "c-access-key",
					SecretAccessKey: "c-secret",
				},
			}

			downstream := bucketset.StorageSet{
				Permanent: &bucketset.Storage{
					Endpoint:        "p.example.com",
					BucketName:      "p-bucket",
					AccessKeyId:     "p-access-key",
					SecretAccessKey: "p-secret",
				},
				Cache: &bucketset.Storage{
					Endpoint:        "",
					BucketName:      "",
					AccessKeyId:     "",
					SecretAccessKey: "",
				},
			}

			initial.Merge(downstream)

			Expect(initial).To(Equal(
				&bucketset.StorageSet{
					Permanent: &bucketset.Storage{
						Endpoint:        "p.example.com",
						BucketName:      "p-bucket",
						AccessKeyId:     "p-access-key",
						SecretAccessKey: "p-secret",
					},
					Cache: &bucketset.Storage{
						Endpoint:        "c.example.com",
						BucketName:      "c-bucket",
						AccessKeyId:     "c-access-key",
						SecretAccessKey: "c-secret",
					},
				},
			))
		})
	})
})
