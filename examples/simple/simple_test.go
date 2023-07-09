package simple

import (
	"context"
	"fmt"
	"github.com/noxjonas/kind-test/pkg/kubetest"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSimple(t *testing.T) {
	clientset := (*kubetest.KubeT)(t).GetClientset()
	// Use the custom testing instance to run the test command with additional logic
	nsList, err := clientset.CoreV1().
		Namespaces().
		List(context.Background(), metav1.ListOptions{})
	//checkErr(err)
	fmt.Println(err)

	for _, n := range nsList.Items {
		fmt.Println(n.Name)
	}
}
