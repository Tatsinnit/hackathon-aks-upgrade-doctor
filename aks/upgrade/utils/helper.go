package utils

import (
	"context"
	"fmt"
	"reflect"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetPods(clientset kubernetes.Interface, namespace string) (*corev1.PodList, error) {
	// Create a pod interface for the given namespace
	podInterface := clientset.CoreV1().Pods(namespace)

	// List the pods in the given namespace
	podList, err := podInterface.List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, fmt.Errorf("getting pods failed: %w", err)
	}

	return podList, nil
}

// IsMapSubset returns true if mapSubset is a subset of mapSet otherwise false
func IsMapSubset(mapSet interface{}, mapSubset interface{}) bool {
	mapSetValue := reflect.ValueOf(mapSet)
	mapSubsetValue := reflect.ValueOf(mapSubset)
	if fmt.Sprintf("%T", mapSet) != fmt.Sprintf("%T", mapSubset) {
		return false
	}
	if len(mapSetValue.MapKeys()) < len(mapSubsetValue.MapKeys()) {
		return false
	}
	if len(mapSubsetValue.MapKeys()) == 0 {
		return true
	}
	iterMapSubset := mapSubsetValue.MapRange()
	for iterMapSubset.Next() {
		k := iterMapSubset.Key()
		v := iterMapSubset.Value()
		value := mapSetValue.MapIndex(k)
		if !value.IsValid() || v.Interface() != value.Interface() {
			return false
		}
	}
	return true
}
