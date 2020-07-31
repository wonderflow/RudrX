package plugins

import (
	"context"
	"fmt"

	cmdutil "github.com/cloud-native-application/rudrx/pkg/cmd/util"
	corev1alpha2 "github.com/crossplane/oam-kubernetes-runtime/apis/core/v1alpha2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetTemplatesFromCluster(ctx context.Context, c client.Client) ([]cmdutil.Template, error) {
	workloads, err := GetWorkloadsFromCluster(ctx, c)
	if err != nil {
		return nil, err
	}
	traits, err := GetTraitsFromCluster(ctx, c)
	if err != nil {
		return nil, err
	}
	workloads = append(workloads, traits...)
	return workloads, nil
}

func GetWorkloadsFromCluster(ctx context.Context, c client.Client) ([]cmdutil.Template, error) {
	var templates []cmdutil.Template
	var workloadDefs corev1alpha2.WorkloadDefinitionList
	err := c.List(ctx, &workloadDefs)
	if err != nil {
		return nil, fmt.Errorf("list WorkloadDefinition err: %s", err)
	}

	for _, wd := range workloadDefs.Items {
		var tmp cmdutil.Template
		tmp, err := cmdutil.ConvertTemplateJson2Object(wd.Spec.Extension)
		if err != nil {
			fmt.Printf("extract template from workloadDefinition %v err: %v, ignore it\n", wd.Name, err)
			continue
		}
		tmp.Type = cmdutil.TypeWorkload
		tmp.Name = wd.Name
		templates = append(templates, tmp)
	}
	return templates, nil
}

func GetTraitsFromCluster(ctx context.Context, c client.Client) ([]cmdutil.Template, error) {
	var templates []cmdutil.Template
	var traitDefs corev1alpha2.TraitDefinitionList
	err := c.List(ctx, &traitDefs)
	if err != nil {
		return nil, fmt.Errorf("list TraitDefinition err: %s", err)
	}

	for _, td := range traitDefs.Items {
		var tmp cmdutil.Template
		tmp, err := cmdutil.ConvertTemplateJson2Object(td.Spec.Extension)
		if err != nil {
			fmt.Printf("extract template from workloadDefinition %v err: %v, ignore it\n", td.Name, err)
			continue
		}
		tmp.Type = cmdutil.TypeTrait
		tmp.Name = td.Name
		templates = append(templates, tmp)
	}
	return templates, nil
}
