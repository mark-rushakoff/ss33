package bucketset

import (
	"encoding/json"
	"io/ioutil"

	"github.com/rlmcpherson/s3gof3r"
)

type StorageSet struct {
	Permanent *Storage
	Cache     *Storage
}

type Storage struct {
	Endpoint        string
	BucketName      string
	AccessKeyId     string
	SecretAccessKey string
}

func (self *StorageSet) Merge(other StorageSet) {
	self.Permanent.Merge(*other.Permanent)
	self.Cache.Merge(*other.Cache)
}

func (self *Storage) Merge(other Storage) {
	if other.Endpoint != "" {
		self.Endpoint = other.Endpoint
	}
	if other.BucketName != "" {
		self.BucketName = other.BucketName
	}
	if other.AccessKeyId != "" {
		self.AccessKeyId = other.AccessKeyId
	}
	if other.SecretAccessKey != "" {
		self.SecretAccessKey = other.SecretAccessKey
	}
}

func (self *StorageSet) BucketSet() *BucketSet {
	return &BucketSet{
		cache:     self.Cache.bucket(),
		permanent: self.Permanent.bucket(),
	}
}

func (self *Storage) bucket() *s3gof3r.Bucket {
	client := s3gof3r.New(self.Endpoint, s3gof3r.Keys{
		AccessKey: self.AccessKeyId,
		SecretKey: self.SecretAccessKey,
	})
	return client.Bucket(self.BucketName)
}

func StorageSetFromFile(path string) (*StorageSet, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var storageSet StorageSet
	if err = json.Unmarshal(content, &storageSet); err != nil {
		return nil, err
	}

	return &storageSet, nil
}
