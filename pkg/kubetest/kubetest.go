package kubetest

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"os"
	"path/filepath"
	"testing"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeT testing.T

var (
	clientset     *kubernetes.Clientset
	dynamicClient *dynamic.DynamicClient
)

func AsKubeT(t *testing.T) *KubeT {
	return (*KubeT)(t)
}

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

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("failed to get home directory: %v", err))
	}

	kubeconfigPath := filepath.Join(homeDir, ".kube", "config")

	configLoader := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: "", //TODO: decide if to switch context here... or if KUBETEST_KUBECONFIG_CONTEXT is just for safety
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

		dynamicClient, err = dynamic.NewForConfig(restConfig)
		if err != nil {
			panic(fmt.Errorf("failed to create dynamic client: %v", err))
		}

	} else {
		panic(fmt.Errorf("selected context is not '%s'", expectedContext))
	}

	return clientset
}

func (t *KubeT) GetDynamicClientset() *dynamic.DynamicClient {
	if dynamicClient != nil {
		return dynamicClient
	}
	t.GetClientset()
	return dynamicClient
}

func (t *KubeT) WithNamespace(namespace *corev1.Namespace) corev1.Namespace {
	clientset = t.GetClientset()

	result, err := clientset.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
	if err != nil {
		t.Errorf("Failed to create namespace: %s", err)
	}

	t.Cleanup(func() {
		err := clientset.CoreV1().Namespaces().Delete(context.TODO(), result.Name, metav1.DeleteOptions{})
		if err != nil {
			// TODO: it's ok if it is already deleted though?
			t.Errorf("Failed to delete namespace in cleanup: %s", err.Error())
		}
	})

	return *result
}
