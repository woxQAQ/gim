package mq

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestWorkerpool(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "mq Suite")
}
