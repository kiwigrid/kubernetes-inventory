package pkg

import (
	"fmt"
	"github.com/kiwigrid/kubernetes-inventory/types"
	"k8s.io/client-go/rest"
	"os"
	"strings"
)

const OUTPUT_FILE = "report.md"

func WriteHeader(config *rest.Config) {
	f, err := os.Create(OUTPUT_FILE)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	user := fmt.Sprintf("%v", config.AuthConfigPersister)
	f.WriteString("k8s context: " + strings.Split(user, " ")[1] + "\n\n")
	f.WriteString("| Name | Container | HelmChart | HelmRelease | hasResourceConfig | hasPDB | UpdateStrategy |\n")
	f.WriteString("|---|---|---|---|---|---|---|\n")
}

func AppendInventoryItem(ci types.ContainerInventory) {
	f, err := os.OpenFile(OUTPUT_FILE, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	var line string

	if (ci.UpdateStrategy == "RollingUpdate") {
		line = fmt.Sprintf("| %s | %s | %s | %s | %t | %t | %s (maxSurge: %s, maxUnavailable: %s)|\n",
			ci.DeploymentName,
			ci.ContainerName,
			ci.HelmChart,
			ci.HelmReleaseName,
			ci.HasResourceConfig(),
			ci.PodDisruptionBudget,
			ci.UpdateStrategy,
			ci.RollingUpdateMaxSurge,
			ci.RollingUpdateMaxUnavailable)
	} else {
		line = fmt.Sprintf("| %s | %s | %s | %s | %t | %t | %s |\n",
			ci.DeploymentName,
			ci.ContainerName,
			ci.HelmChart,
			ci.HelmReleaseName,
			ci.HasResourceConfig(),
			ci.PodDisruptionBudget,
			ci.UpdateStrategy)

	}

	if _, err = f.WriteString(line); err != nil {
		panic(err)
	}
}
