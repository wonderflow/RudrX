package plugins

import (
	"context"
	"fmt"
	"io/ioutil"

	util2 "github.com/cloud-native-application/rudrx/pkg/cmd/util"

	"github.com/ghodss/yaml"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/crossplane/oam-kubernetes-runtime/apis/core/v1alpha2"
	"github.com/crossplane/oam-kubernetes-runtime/pkg/oam/util"
)

var _ = Describe("DefinitionFiles", func() {
	ctx := context.Background()
	BeforeEach(func() {

		traitdata, err := ioutil.ReadFile("testdata/traitDef.yaml")
		Expect(err).Should(BeNil())
		var td v1alpha2.TraitDefinition
		Expect(yaml.Unmarshal(traitdata, &td)).Should(BeNil())

		td.Namespace = "default"
		logf.Log.Info("Creating trait definition", "data", td)
		Expect(k8sClient.Create(ctx, &td)).Should(SatisfyAny(BeNil(), &util.AlreadyExistMatcher{}))

		workloaddata, err := ioutil.ReadFile("testdata/workloadDef.yaml")
		Expect(err).Should(BeNil())
		var wd v1alpha2.WorkloadDefinition
		Expect(yaml.Unmarshal(workloaddata, &wd)).Should(BeNil())

		td.Namespace = "default"
		logf.Log.Info("Creating workload definition", "data", wd)
		Expect(k8sClient.Create(ctx, &wd)).Should(SatisfyAny(BeNil(), &util.AlreadyExistMatcher{}))

	})

	It("gettrait", func() {
		traitDefs, err := GetTraitsFromCluster(context.Background(), k8sClient)
		Expect(err).Should(BeNil())
		logf.Log.Info(fmt.Sprintf("%v", traitDefs))
		statefulset := util2.Template{
			Name:  "routes.extend.oam.dev",
			Type:  util2.TypeTrait,
			Alias: "routes",
			Object: map[string]interface{}{
				"apiVersion": "extend.oam.dev/v1alpha2",
				"kind":       "Route",
			},
			Parameters: []util2.Parameter{
				{
					Name:       "domain",
					Short:      "d",
					Required:   true,
					FieldPaths: []string{"spec.domain"},
				},
			},
		}
		Expect(traitDefs).Should(Equal([]util2.Template{statefulset}))
	})
})
