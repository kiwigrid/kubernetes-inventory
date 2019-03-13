package main

import (
	"flag"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/kiwigrid/kubernetes-inventory/pkg"
	"github.com/kiwigrid/kubernetes-inventory/types"
	"k8s.io/api/extensions/v1beta1"
	v1beta13 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

var VERSION = "latest"
var namespace = ""

func main() {

	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)

	containerList := processResources(clientset)

	pkg.WriteHeader(config)
	for _, c := range containerList {
		pkg.AppendInventoryItem(c)
	}
}

func processResources(clientset *kubernetes.Clientset) []types.ContainerInventory {
	deploymentList, _ := clientset.ExtensionsV1beta1().Deployments(namespace).List(metav1.ListOptions{})
	podDisruptionBudgetList, _ := clientset.PolicyV1beta1().PodDisruptionBudgets(namespace).List(metav1.ListOptions{})

	var containerList []types.ContainerInventory

	for _, deployment := range deploymentList.Items {

		if deployment.Namespace == "kube-system" {
			continue
		}

		for _, container := range deployment.Spec.Template.Spec.Containers {

			chart, rel := getHelmMetadata(&deployment)

			ctn := &types.ContainerInventory{
				ContainerName:        container.Name,
				DeploymentName:       deployment.Name,
				HelmChart:            chart,
				HelmReleaseName:      rel,
				ReplicaCount:         int(*deployment.Spec.Replicas),
				ResourceCpuRequested: container.Resources.Requests.Cpu().String(),
				ResourceMemRequested: container.Resources.Requests.Memory().String(),
				ResourceCpuLimit:     container.Resources.Limits.Cpu().String(),
				ResourceMemLimit:     container.Resources.Limits.Memory().String(),
				UpdateStrategy:       fmt.Sprintf("%v", deployment.Spec.Strategy.Type),
				PodDisruptionBudget:  false,
				Affinity:             deployment.Spec.Template.Spec.Affinity != nil,
				StandardHelmLabels:   hasHelmStandardLabels(&deployment),
			}

			if deployment.Spec.Strategy.RollingUpdate != nil {
				ctn.RollingUpdateMaxSurge = deployment.Spec.Strategy.RollingUpdate.MaxSurge.String()
				ctn.RollingUpdateMaxUnavailable = deployment.Spec.Strategy.RollingUpdate.MaxUnavailable.String()
			}

			checkForPdb(podDisruptionBudgetList, ctn)
			containerList = append(containerList, *ctn)
		}
	}

	return containerList
}

func checkForPdb(pdbList *v1beta13.PodDisruptionBudgetList, ci *types.ContainerInventory) {

	if ci.HelmReleaseName == "" {
		return
	}
	pdbMatches := false
	for _, pdb := range pdbList.Items {

		pdbMatches = pdb.Labels["release"] == ci.HelmReleaseName || pdb.Labels["app.kubernetes.io/instance"] == ci.HelmReleaseName
		if pdbMatches {
			break
		}
	}
	ci.PodDisruptionBudget = pdbMatches
}

func hasHelmStandardLabels(deployment *v1beta1.Deployment) bool {
	return deployment.Labels["app.kubernetes.io/name"] != "" &&
		deployment.Labels["helm.sh/chart"] != "" &&
		deployment.Labels["app.kubernetes.io/managed-by"] != "" &&
		deployment.Labels["app.kubernetes.io/instance"] != ""
}

func getClient(pathToCfg string) (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error
	if pathToCfg == "" {
		// in cluster access
		logrus.Info("Using in cluster config")
		config, err = rest.InClusterConfig()
	} else {
		logrus.Info("Using out of cluster config")
		config, err = clientcmd.BuildConfigFromFlags("", pathToCfg)
	}
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func getHelmMetadata(deployment *v1beta1.Deployment) (string, string) {
	var chart string
	var releaseName string

	if deployment.Labels["helm.sh/chart"] != "" {
		chart = deployment.Labels["helm.sh/chart"]
	} else if deployment.Labels["app"] != "" {
		chart = deployment.Labels["app"]
	}

	if deployment.Labels["app.kubernetes.io/instance"] != "" {
		releaseName = deployment.Labels["app.kubernetes.io/instance"]
	} else if deployment.Labels["release"] != "" {
		releaseName = deployment.Labels["release"]
	}

	return chart, releaseName
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
