package createcluster

import (
	"flag"
	"os"
	"testing"

	"github.com/rancher/distros-test-framework/core/service/customflag"
	"github.com/rancher/distros-test-framework/factory"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMain(m *testing.M) {

	flag.Var(&customflag.ServiceFlag.ClusterConfig.Destroy, "destroy", "Destroy cluster after test")
	flag.Var(&customflag.ServiceFlag.ClusterConfig.Arch, "arch", "Architecture type")

	flag.Parse()
	os.Exit(m.Run())
}

func TestClusterCreateSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Create Cluster Test Suite")
}

var _ = AfterSuite(func() {
	g := GinkgoT()
	if customflag.ServiceFlag.ClusterConfig.Destroy {
		status, err := factory.DestroyCluster(g)
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal("cluster destroyed"))
	}
})
