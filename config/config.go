package config

import (
	"encoding/json"
	"io/ioutil"
)

type StorageSet struct {
	Permanent Storage
	Cache     Storage
}

type Storage struct {
	Endpoint        string
	BucketName      string
	AccessKeyId     string
	SecretAccessKey string
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
