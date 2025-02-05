package workerpool_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestWorkerpool(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Workerpool Suite")
}
