package specs

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo"

	"github.com/openshift/openshift-azure/test/clients/azure"
	"github.com/openshift/openshift-azure/test/util/log"
)

var _ = BeforeSuite(func() {
	var err error
	testlogger := log.GetTestLogger()
	fmt.Println("configuring the fake resource provider")
	err = azure.NewClientFromEnvironment(context.Background(), testlogger)
	if err != nil {
		testlogger.Error(err)
	}

	if azure.RPClient == nil {
		testlogger.Error("unable to provision azure client")
		panic("No Azure client")
	}
})
