package beth_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBeth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Beth Test Suite")
}
