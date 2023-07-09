package examples

import (
	"context"
	"github.com/noxjonas/kind-test/pkg/kubetest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
)

func TestMultipleManifests(t *testing.T) {
	err := kubetest.AsKubeT(t).ApplyManifests([]string{"./manifests/*.yaml"})
	if err != nil {
		t.Errorf("failed to create manifests from glob path: %v", err)
	}

	time.Sleep(10 * time.Second)

	clientset := kubetest.AsKubeT(t).GetClientset()
	_, err = clientset.AppsV1().Deployments("manifest-namespace").Get(context.TODO(), "nginx-deployment", metav1.GetOptions{})
	if err != nil {
		t.Errorf("failed to create deployment: %v", err)
	}

	_, err = clientset.CoreV1().Pods("default").Get(context.TODO(), "busybox-pod", metav1.GetOptions{})
	if err != nil {
		t.Errorf("failed to create pod: %v", err)
	}

	// TODO: will fail second time around without cleanup
}
