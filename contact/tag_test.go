package contact_test

import (

	"github.com/dfang/wechat-work-go/contact"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("标签管理 API", func() {
	var concat contact.WithApp(app)
	var tagMnagement := concat.NewTagManagement()

	Context("创建", func() {
		tagMnagement.Create()
		tagMnagement.List()
		tagMnagement.xxx()
	})

})
