package types

import "fmt"

type ContainerInventory struct {
	DeploymentName       string
	ContainerName        string
	HelmReleaseName      string
	HelmChart            string
	ReplicaCount         int
	Affinity             bool
	PodDisruptionBudget  bool
	UpdateStrategy       string
	ResourceCpuRequested string
	ResourceCpuLimit     string
	ResourceMemRequested string
	ResourceMemLimit     string
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
