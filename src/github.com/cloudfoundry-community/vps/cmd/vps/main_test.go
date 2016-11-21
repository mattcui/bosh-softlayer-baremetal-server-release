package main_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"net/http"
	"os/exec"
)

type Args struct {
	LogLevel                 string
	Host                     string
	Port                     string
	DatabaseDriver           string
	DatabaseConnectionString string
}

var (
	vpsConfig string
	vpsArgs   Args
	session   *gexec.Session
	err       error
)

var _ = Describe("Virtual Pool Server", func() {
	BeforeEach(func() {
		//vpsConfig, err = gexec.Build("github.com/cloudfoundry-community/vps/cmd/vps")
		vpsConfig, err = gexec.BuildIn("/home/travis/gopath:/home/travis/gopath/src/github.com/mattcui/bosh-softlayer-baremetal-server-release", "github.com/cloudfoundry-community/vps/cmd/vps")
		Ω(err).ShouldNot(HaveOccurred())
	})

	Context("when starting Virtual Pool Server with given correct arguments", func() {
		vpsArgs = Args{
			LogLevel:                 "debug",
			Host:                     "127.0.0.1",
			Port:                     "8889",
			DatabaseDriver:           "postgres",
			DatabaseConnectionString: "postgres://postgres:postgres@127.0.0.1/bosh",
		}

		command := exec.Command(string(vpsConfig), vpsArgs.argSlice()...)
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())

		resp, err := http.Get(fmt.Sprintf("http://%s:%s/v2/vms", vpsArgs.Host, vpsArgs.Port))
		defer resp.Body.Close()

		Expect(err).To(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(200))
	})

	AfterEach(func() {
		session.Terminate()
		gexec.CleanupBuildArtifacts()
	})
})

func (args Args) argSlice() []string {
	arguments := []string{
		"--logLevel", args.LogLevel,
		"--host", args.Host,
		"--port", args.Port,
		"--databaseDriver", args.DatabaseDriver,
		"--databaseConnectionString", args.DatabaseConnectionString,
	}

	return arguments
}
