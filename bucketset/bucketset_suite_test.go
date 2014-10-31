package bucketset_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestBucketset(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bucketset Suite")
}
