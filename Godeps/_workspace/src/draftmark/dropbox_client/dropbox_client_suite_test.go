package dropbox_client_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDraftmark(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Draftmark Suite")
}
