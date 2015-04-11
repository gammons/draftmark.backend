package dropbox_client_test

import (
	dropbox "draftmark/dropbox_client"
	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var dbox = &dropbox.Client{}
var tmpFolder = "/__draftmarktest__"

var _ = BeforeSuite(func() {
	_ = godotenv.Load("test.env")
	setupDropbox()
})

var _ = Describe("Dropbox Client", func() {
	var _ = Describe("GetContent", func() {
		It("can get the content of the file", func() {
			dbox.Dbox.UploadFile("../test.md", tmpFolder+"/test.md", true, "")
			Expect(dbox.GetContent(tmpFolder + "/test.md")).To(ContainSubstring("this is a test file"))
			dbox.Dbox.Delete(tmpFolder + "/test.md")
		})

		Context("When the file cannot be found, or is empty", func() {
			It("returns the error", func() {
				_, err := dbox.GetContent("notfound.md")
				Expect(err.Error()).To(Equal("file does not exist"))
			})
		})
	})
})

func setupDropbox() {
	dbox.AccessToken = os.Getenv("DROPBOX_ACCESS_TOKEN")
	dbox.InitDropbox()
	dbox.Dbox.CreateFolder(tmpFolder)
}
