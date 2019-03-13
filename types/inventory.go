package types

import "fmt"

type ContainerInventory struct {
	DeploymentName              string
	ContainerName               string
	HelmReleaseName             string
	HelmChart                   string
	ReplicaCount                int
	Affinity                    bool
	PodDisruptionBudget         bool
	UpdateStrategy              string
	RollingUpdateMaxSurge       string
	RollingUpdateMaxUnavailable string
	ResourceCpuRequested        string
	ResourceCpuLimit            string
	ResourceMemRequested        string
	ResourceMemLimit            string
	StandardHelmLabels			bool
}

func (ci ContainerInventory) String() string {
	return fmt.Sprintf("ContainerInventory{"+
		"Name='%s'"+
		", ContainerName='%s'"+
		", HelmChart='%s'"+
		", HelmRelease='%s'"+
		", ResourceCpuRequested='%s'"+
		", ResourceCpuLimit='%s'"+
		", ResourceMemRequested='%s'"+
		", ResourceMemLimit='%s'"+
		"}", ci.DeploymentName, ci.ContainerName, ci.HelmChart, ci.HelmReleaseName, ci.ResourceCpuRequested, ci.ResourceCpuLimit, ci.ResourceMemRequested, ci.ResourceMemLimit)
}

func (ci ContainerInventory) HasResourceConfig() bool {
	return ci.ResourceMemLimit != "0" && ci.ResourceMemRequested != "0" && ci.ResourceCpuLimit != "0" && ci.ResourceCpuRequested != "0"
}
