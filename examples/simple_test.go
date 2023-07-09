package examples

import (
	"context"
	"github.com/noxjonas/kind-test/pkg/kubetest"
	corev1 "k8s.io/api/core/v1"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCheckDefaultNamespaceExists(t *testing.T) {
	clientset := kubetest.AsKubeT(t).GetClientset()

	nsList, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		t.Errorf("failed to fetch namespaces: %v", err)
	}

	for _, ns := range nsList.Items {
		if ns.Name == "default" {
			return
		}
	}

	t.Error("'default' namespace does not exist")
}

func TestWithNamespace(t *testing.T) {
	namespace_ := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-namespace-1",
		},
	}

	namespace := kubetest.AsKubeT(t).WithNamespace(namespace_)

	pod_ := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "busybox-pod",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "busybox",
					Image: "busybox",
					Command: []string{
						"sleep",
						"3600",
					},
				},
			},
		},
	}

	clientset := kubetest.AsKubeT(t).GetClientset()
	_, err := clientset.CoreV1().Pods(namespace.Name).Create(context.TODO(), pod_, metav1.CreateOptions{})
	if err != nil {
		t.Errorf("failed to create pod: %v", err)
	}

	time.Sleep(30 * time.Second)
}
