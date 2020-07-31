package plugins

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cloud-native-application/rudrx/pkg/cmd/util"
)

func TestLocalSink(t *testing.T) {
	deployment := util.Template{
		Name:  "deployment",
		Type:  util.TypeWorkload,
		Alias: "deployment",
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
		},
		Parameters: []util.Parameter{
			{
				Name:       "image",
				Short:      "i",
				Required:   true,
				FieldPaths: []string{"spec.containers[0].image"},
			},
		},
	}
	statefulset := util.Template{
		Name:  "statefulset",
		Type:  util.TypeWorkload,
		Alias: "stateful",
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Statefulset",
		},
		Parameters: []util.Parameter{
			{
				Name:       "image",
				Short:      "i",
				Required:   true,
				FieldPaths: []string{"spec.containers[0].image"},
			},
		},
	}
	route := util.Template{
		Name:  "route",
		Type:  util.TypeTrait,
		Alias: "route",
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Route",
		},
		Parameters: []util.Parameter{
			{
				Name:       "domain",
				Short:      "d",
				Required:   true,
				FieldPaths: []string{"spec.domain"},
			},
		},
	}

	cases := map[string]struct {
		dir    string
		tmps   []util.Template
		Type   util.DefinitionType
		expDef []util.Template
	}{
		"Test No Templates": {
			dir:  "rudrx-test1",
			tmps: nil,
		},
		"Test Only Workload": {
			dir:    "rudrx-test2",
			tmps:   []util.Template{deployment, statefulset},
			Type:   util.TypeWorkload,
			expDef: []util.Template{deployment, statefulset},
		},
		"Test Only Trait": {
			dir:    "rudrx-test3",
			tmps:   []util.Template{route},
			Type:   util.TypeTrait,
			expDef: []util.Template{route},
		},
		"Test Only Workload But want trait": {
			dir:    "rudrx-test3",
			tmps:   []util.Template{deployment, statefulset},
			Type:   util.TypeTrait,
			expDef: nil,
		},
		"Test Both have Workload and trait But want Workload": {
			dir:    "rudrx-test4",
			tmps:   []util.Template{deployment, route, statefulset},
			Type:   util.TypeWorkload,
			expDef: []util.Template{deployment, statefulset},
		},
		"Test Both have Workload and trait But want Trait": {
			dir:    "rudrx-test5",
			tmps:   []util.Template{deployment, route, statefulset},
			Type:   util.TypeTrait,
			expDef: []util.Template{route},
		},
	}
	for name, c := range cases {
		testInDir(t, name, c.dir, c.tmps, c.expDef, c.Type)
	}
}

func testInDir(t *testing.T, casename, dir string, tmps, defexp []util.Template, Type util.DefinitionType) {
	err := os.MkdirAll(dir, 0755)
	assert.NoError(t, err, casename)
	defer os.RemoveAll(dir)
	err = SinkTemp2Local(tmps, dir)
	assert.NoError(t, err, casename)
	gottmps, err := LoadTempFromLocal(dir)
	assert.NoError(t, err, casename)
	assert.Equal(t, tmps, gottmps, casename)
	if Type != "" {
		gotDef, err := GetDefFromLocal(dir, Type)
		assert.NoError(t, err, casename)
		assert.Equal(t, defexp, gotDef, casename)
	}
}
