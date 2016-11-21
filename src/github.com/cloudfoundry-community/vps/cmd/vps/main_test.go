package main_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"net/http"
	"os/exec"
	"time"
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
		vpsConfig, err = gexec.Build("github.com/cloudfoundry-community/vps/cmd/vps")
		Ω(err).ShouldNot(HaveOccurred())
	})

	Context("when starting Virtual Pool Server with given correct arguments", func() {
		It("start virtual pool server and query all vms", func() {
			vpsArgs = Args{
				LogLevel:                 "debug",
				Host:                     "127.0.0.1",
				Port:                     "8889",
				DatabaseDriver:           "postgres",
				DatabaseConnectionString: "postgres://postgres:postgres@localhost/bosh",
			}
			fmt.Println(string(vpsConfig))

			command := exec.Command(string(vpsConfig), vpsArgs.argSlice()...)
			session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Ω(err).ShouldNot(HaveOccurred())

			time.Sleep(20*time.Second)
			resp, err := http.Get(fmt.Sprintf("http://%s:%s/v2/vms", vpsArgs.Host, vpsArgs.Port))
			if err != nil {
				fmt.Println(err.Error())
			}
			defer resp.Body.Close()

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(404))
		})
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
