package kubetest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
	_ "k8s.io/client-go/kubernetes/scheme"
)

func (t *KubeT) ApplyManifests(manifestFiles []string) error {
	dynamicClient = t.GetDynamicClientset()

	println(manifestFiles)

	for _, pattern := range manifestFiles {
		files, err := filepath.Glob(pattern)
		if err != nil {
			return fmt.Errorf("failed to match files using pattern %s: %v", pattern, err)
		}

		for _, file := range files {
			manifestBytes, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("failed to read manifest file %s: %v", file, err)
			}

			reader := bytes.NewReader(manifestBytes) // Create a single reader for the manifest file

			decoder := yaml.NewYAMLOrJSONDecoder(reader, 4096)
			for {
				obj := &unstructured.Unstructured{}
				err := decoder.Decode(obj)
				if err != nil {
					if err == io.EOF {
						break // Reached the end of the manifest file
					}
					return fmt.Errorf("failed to decode manifest file %s: %v", file, err)
				}

				gvk := obj.GroupVersionKind()
				gvr := schema.GroupVersionResource{
					Group:    gvk.Group,
					Version:  gvk.Version,
					Resource: strings.ToLower(gvk.Kind) + "s",
				}

				_, err = dynamicClient.Resource(gvr).Namespace(obj.GetNamespace()).Create(context.TODO(), obj, metav1.CreateOptions{})
				if err != nil {
					return fmt.Errorf("failed to apply object from manifest file %s: %v", file, err)
				}
			}
		}
	}

	return nil
}
