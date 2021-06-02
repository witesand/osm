// Package service models an instance of a service managed by OSM controller and utility routines associated with it.
package service

<<<<<<< HEAD
import (
	"fmt"
	"reflect"
	"strings"
	"strconv"

	"github.com/google/uuid"
)
=======
import "fmt"
>>>>>>> 3d923b3f2d72006f6cdaad056938c492c364196d

const (
	// namespaceNameSeparator used upon marshalling/unmarshalling MeshService to a string
	// or viceversa
	namespaceNameSeparator = "/"
)

// MeshService is the struct defining a service (Kubernetes or otherwise) within a service mesh.
type MeshService struct {
	// If the service resides on a Kubernetes service, this would be the Kubernetes namespace.
	Namespace string

	// The name of the service
	Name string
}

<<<<<<< HEAD
func (ms MeshService) String() string {
	return strings.Join([]string{ms.Namespace, namespaceNameSeparator, ms.Name}, "")
}

// Equals checks if two namespaced services are equal
func (ms MeshService) Equals(service MeshService) bool {
	return reflect.DeepEqual(ms, service)
}

// UnmarshalMeshService unmarshals a NamespaceService type from a string
func UnmarshalMeshService(str string) (*MeshService, error) {
	slices := strings.Split(str, namespaceNameSeparator)
	if len(slices) != 2 {
		return nil, errInvalidMeshServiceFormat
	}

	// Make sure the slices are not empty. Split might actually leave empty slices.
	for _, sep := range slices {
		if len(sep) == 0 {
			return nil, errInvalidMeshServiceFormat
		}
	}

	return &MeshService{
		Namespace: slices[0],
		Name:      slices[1],
	}, nil
}

// ServerName returns the Server Name Identifier (SNI) for TLS connections
func (ms MeshService) ServerName() string {
	return strings.Join([]string{ms.Name, ms.Namespace, "svc", "cluster", "local"}, ".")
}

func (ms MeshService) GetMeshServicePort() MeshServicePort {
	return MeshServicePort{
		Namespace: ms.Namespace,
		Name: ms.Name,
		Port: 0,
	}
}

type MeshServicePort struct {
	// If the service resides on a Kubernetes service, this would be the Kubernetes namespace.
	Namespace string

	// The name of the service
	Name string

	// Service port
	Port int
}

func (ms MeshServicePort) GetMeshService() MeshService {
	return MeshService{
		Namespace: ms.Namespace,
		Name: ms.Name,
	}
}

func (ms MeshServicePort) String() string {
	return strings.Join([]string{ms.Namespace, namespaceNameSeparator, ms.Name, namespaceNameSeparator, strconv.Itoa(ms.Port)}, "")
}

// Equals checks if two namespaced services are equal
func (ms MeshServicePort) Equals(service MeshServicePort) bool {
	return reflect.DeepEqual(ms, service)
}

// UnmarshalMeshServicePort unmarshals a NamespaceService type from a string
func UnmarshalMeshServicePort(str string) (*MeshServicePort, error) {
	slices := strings.Split(str, namespaceNameSeparator)
	if len(slices) != 3 {
		return nil, errInvalidMeshServiceFormat
	}

	// Make sure the slices are not empty. Split might actually leave empty slices.
	for i, sep := range slices {
		if i == 2 {
			// Port can be empty
			continue
		}
		if len(sep) == 0 {
			return nil, errInvalidMeshServiceFormat
		}
	}

	port := 0
	if slices[2] != "" {
		port, _ = strconv.Atoi(slices[2])
	}

	return &MeshServicePort{
		Namespace: slices[0],
		Name:      slices[1],
		Port:      port,
	}, nil
}

=======
>>>>>>> 3d923b3f2d72006f6cdaad056938c492c364196d
// K8sServiceAccount is a type for a namespaced service account
type K8sServiceAccount struct {
	Namespace string
	Name      string
}

// String returns the string representation of the service account object
func (sa K8sServiceAccount) String() string {
	return fmt.Sprintf("%s%s%s", sa.Namespace, namespaceNameSeparator, sa.Name)
}

// IsEmpty returns true if the given service account object is empty
func (sa K8sServiceAccount) IsEmpty() bool {
	return (K8sServiceAccount{}) == sa
}

// ClusterName is a type for a service name
type ClusterName string

// String returns the given ClusterName type as a string
func (c ClusterName) String() string {
	return string(c)
}

// WeightedCluster is a struct of a cluster and is weight that is backing a service
type WeightedCluster struct {
	ClusterName ClusterName `json:"cluster_name:omitempty"`
	Weight      int         `json:"weight:omitempty"`
}
