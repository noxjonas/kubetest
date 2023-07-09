package kubetest

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeT testing.T

var (
	clientset *kubernetes.Clientset
)

// purely for my peace of mind...
func getExpectedContext() string {
	value := os.Getenv("KUBETEST_KUBECONFIG_CONTEXT")

	if value == "" {
		value = "kind-kind"
	}
	return value
}

func (t *KubeT) GetClientset() *kubernetes.Clientset {
	if clientset != nil {
		return clientset
	}

	value := os.Getenv("ENV_VARIABLE_NAME")

	// Check if the environment variable is empty or missing
	if value == "" {
		// Set a default value
		value = "default value"
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("failed to get home directory: %v", err))
	}

	kubeconfigPath := filepath.Join(homeDir, ".kube", "config")

	configLoader := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: "",
		})

	rawConfig, err := configLoader.RawConfig()
	if err != nil {
		panic(fmt.Errorf("failed to load kubeconfig: %v", err))
	}

	restConfig, err := configLoader.ClientConfig()
	if err != nil {
		panic(fmt.Errorf("failed to load kubeconfig: %v", err))
	}

	expectedContext := getExpectedContext()
	if rawConfig.CurrentContext == expectedContext {
		clientset, err = kubernetes.NewForConfig(restConfig)
		if err != nil {
			panic(fmt.Errorf("failed to create Kubernetes clientset: %v", err))
		}

	} else {
		panic(fmt.Errorf("selected context is not '%s'", expectedContext))
	}

	return clientset
}

func (t *KubeT) Deploy(deployment string) {

	//t.C
}
