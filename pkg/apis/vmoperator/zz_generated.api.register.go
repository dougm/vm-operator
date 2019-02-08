/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This file was autogenerated by apiregister-gen. Do not edit it manually!

package vmoperator

import (
	"context"
	"fmt"

	"github.com/kubernetes-incubator/apiserver-builder-alpha/pkg/builders"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/rest"
)

var (
	InternalVirtualMachine = builders.NewInternalResource(
		"virtualmachines",
		"VirtualMachine",
		func() runtime.Object { return &VirtualMachine{} },
		func() runtime.Object { return &VirtualMachineList{} },
	)
	InternalVirtualMachineStatus = builders.NewInternalResourceStatus(
		"virtualmachines",
		"VirtualMachineStatus",
		func() runtime.Object { return &VirtualMachine{} },
		func() runtime.Object { return &VirtualMachineList{} },
	)
	InternalVirtualMachineImage = builders.NewInternalResource(
		"virtualmachineimages",
		"VirtualMachineImage",
		func() runtime.Object { return &VirtualMachineImage{} },
		func() runtime.Object { return &VirtualMachineImageList{} },
	)
	InternalVirtualMachineImageStatus = builders.NewInternalResourceStatus(
		"virtualmachineimages",
		"VirtualMachineImageStatus",
		func() runtime.Object { return &VirtualMachineImage{} },
		func() runtime.Object { return &VirtualMachineImageList{} },
	)
	InternalVirtualMachineService = builders.NewInternalResource(
		"virtualmachineservices",
		"VirtualMachineService",
		func() runtime.Object { return &VirtualMachineService{} },
		func() runtime.Object { return &VirtualMachineServiceList{} },
	)
	InternalVirtualMachineServiceStatus = builders.NewInternalResourceStatus(
		"virtualmachineservices",
		"VirtualMachineServiceStatus",
		func() runtime.Object { return &VirtualMachineService{} },
		func() runtime.Object { return &VirtualMachineServiceList{} },
	)
	// Registered resources and subresources
	ApiVersion = builders.NewApiGroup("vmoperator.vmware.com").WithKinds(
		InternalVirtualMachine,
		InternalVirtualMachineStatus,
		InternalVirtualMachineImage,
		InternalVirtualMachineImageStatus,
		InternalVirtualMachineService,
		InternalVirtualMachineServiceStatus,
	)

	// Required by code generated by go2idl
	AddToScheme        = ApiVersion.SchemaBuilder.AddToScheme
	SchemeBuilder      = ApiVersion.SchemaBuilder
	localSchemeBuilder = &SchemeBuilder
	SchemeGroupVersion = ApiVersion.GroupVersion
)

// Required by code generated by go2idl
// Kind takes an unqualified kind and returns a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// Required by code generated by go2idl
// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

// +genclient
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type VirtualMachineImage struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Spec   VirtualMachineImageSpec
	Status VirtualMachineImageStatus
}

// +genclient
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type VirtualMachine struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Spec   VirtualMachineSpec
	Status VirtualMachineStatus
}

type VirtualMachineImageStatus struct {
	Uuid       string
	InternalId string
	PowerState string
}

type VirtualMachineStatus struct {
	Conditions []VirtualMachineCondition
	Host       string
	PowerState string
	Phase      string
	VmIp       string
}

type VirtualMachineSpec struct {
	Image      string
	Resources  VirtualMachineResourcesSpec
	PowerState string
	Env        corev1.EnvVar
	Ports      []VirtualMachinePort
}

type VirtualMachineCondition struct {
	LastProbeTime      metav1.Time
	LastTransitionTime metav1.Time
	Message            string
	Reason             string
	Status             string
	Type               string
}

type VirtualMachinePort struct {
	Port     int
	Ip       string
	Name     string
	Protocol corev1.Protocol
}

type VirtualMachineResourcesSpec struct {
	Capacity VirtualMachineResourceSpec
	Requests VirtualMachineResourceSpec
	Limits   VirtualMachineResourceSpec
}

type VirtualMachineImageSpec struct {
}

type VirtualMachineResourceSpec struct {
	Cpu    int64
	Memory int64
}

// +genclient
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type VirtualMachineService struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Spec   VirtualMachineServiceSpec
	Status VirtualMachineServiceStatus
}

type VirtualMachineServiceSpec struct {
	Type         string
	Ports        []VirtualMachineServicePort
	Selector     map[string]string
	ClusterIP    string
	ExternalName string
}

type VirtualMachineServiceStatus struct {
}

type VirtualMachineServicePort struct {
	Name       string
	Protocol   string
	Port       int32
	TargetPort int32
}

//
// VirtualMachine Functions and Structs
//
// +k8s:deepcopy-gen=false
type VirtualMachineStrategy struct {
	builders.DefaultStorageStrategy
}

// +k8s:deepcopy-gen=false
type VirtualMachineStatusStrategy struct {
	builders.DefaultStatusStorageStrategy
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type VirtualMachineList struct {
	metav1.TypeMeta
	metav1.ListMeta
	Items []VirtualMachine
}

func (VirtualMachine) NewStatus() interface{} {
	return VirtualMachineStatus{}
}

func (pc *VirtualMachine) GetStatus() interface{} {
	return pc.Status
}

func (pc *VirtualMachine) SetStatus(s interface{}) {
	pc.Status = s.(VirtualMachineStatus)
}

func (pc *VirtualMachine) GetSpec() interface{} {
	return pc.Spec
}

func (pc *VirtualMachine) SetSpec(s interface{}) {
	pc.Spec = s.(VirtualMachineSpec)
}

func (pc *VirtualMachine) GetObjectMeta() *metav1.ObjectMeta {
	return &pc.ObjectMeta
}

func (pc *VirtualMachine) SetGeneration(generation int64) {
	pc.ObjectMeta.Generation = generation
}

func (pc VirtualMachine) GetGeneration() int64 {
	return pc.ObjectMeta.Generation
}

// Registry is an interface for things that know how to store VirtualMachine.
// +k8s:deepcopy-gen=false
type VirtualMachineRegistry interface {
	ListVirtualMachines(ctx context.Context, options *internalversion.ListOptions) (*VirtualMachineList, error)
	GetVirtualMachine(ctx context.Context, id string, options *metav1.GetOptions) (*VirtualMachine, error)
	CreateVirtualMachine(ctx context.Context, id *VirtualMachine) (*VirtualMachine, error)
	UpdateVirtualMachine(ctx context.Context, id *VirtualMachine) (*VirtualMachine, error)
	DeleteVirtualMachine(ctx context.Context, id string) (bool, error)
}

// NewRegistry returns a new Registry interface for the given Storage. Any mismatched types will panic.
func NewVirtualMachineRegistry(sp builders.StandardStorageProvider) VirtualMachineRegistry {
	return &storageVirtualMachine{sp}
}

// Implement Registry
// storage puts strong typing around storage calls
// +k8s:deepcopy-gen=false
type storageVirtualMachine struct {
	builders.StandardStorageProvider
}

func (s *storageVirtualMachine) ListVirtualMachines(ctx context.Context, options *internalversion.ListOptions) (*VirtualMachineList, error) {
	if options != nil && options.FieldSelector != nil && !options.FieldSelector.Empty() {
		return nil, fmt.Errorf("field selector not supported yet")
	}
	st := s.GetStandardStorage()
	obj, err := st.List(ctx, options)
	if err != nil {
		return nil, err
	}
	return obj.(*VirtualMachineList), err
}

func (s *storageVirtualMachine) GetVirtualMachine(ctx context.Context, id string, options *metav1.GetOptions) (*VirtualMachine, error) {
	st := s.GetStandardStorage()
	obj, err := st.Get(ctx, id, options)
	if err != nil {
		return nil, err
	}
	return obj.(*VirtualMachine), nil
}

func (s *storageVirtualMachine) CreateVirtualMachine(ctx context.Context, object *VirtualMachine) (*VirtualMachine, error) {
	st := s.GetStandardStorage()
	obj, err := st.Create(ctx, object, nil, &metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return obj.(*VirtualMachine), nil
}

func (s *storageVirtualMachine) UpdateVirtualMachine(ctx context.Context, object *VirtualMachine) (*VirtualMachine, error) {
	st := s.GetStandardStorage()
	obj, _, err := st.Update(ctx, object.Name, rest.DefaultUpdatedObjectInfo(object), nil, nil, false, &metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	return obj.(*VirtualMachine), nil
}

func (s *storageVirtualMachine) DeleteVirtualMachine(ctx context.Context, id string) (bool, error) {
	st := s.GetStandardStorage()
	_, sync, err := st.Delete(ctx, id, &metav1.DeleteOptions{})
	return sync, err
}

//
// VirtualMachineImage Functions and Structs
//
// +k8s:deepcopy-gen=false
type VirtualMachineImageStrategy struct {
	builders.DefaultStorageStrategy
}

// +k8s:deepcopy-gen=false
type VirtualMachineImageStatusStrategy struct {
	builders.DefaultStatusStorageStrategy
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type VirtualMachineImageList struct {
	metav1.TypeMeta
	metav1.ListMeta
	Items []VirtualMachineImage
}

func (VirtualMachineImage) NewStatus() interface{} {
	return VirtualMachineImageStatus{}
}

func (pc *VirtualMachineImage) GetStatus() interface{} {
	return pc.Status
}

func (pc *VirtualMachineImage) SetStatus(s interface{}) {
	pc.Status = s.(VirtualMachineImageStatus)
}

func (pc *VirtualMachineImage) GetSpec() interface{} {
	return pc.Spec
}

func (pc *VirtualMachineImage) SetSpec(s interface{}) {
	pc.Spec = s.(VirtualMachineImageSpec)
}

func (pc *VirtualMachineImage) GetObjectMeta() *metav1.ObjectMeta {
	return &pc.ObjectMeta
}

func (pc *VirtualMachineImage) SetGeneration(generation int64) {
	pc.ObjectMeta.Generation = generation
}

func (pc VirtualMachineImage) GetGeneration() int64 {
	return pc.ObjectMeta.Generation
}

// Registry is an interface for things that know how to store VirtualMachineImage.
// +k8s:deepcopy-gen=false
type VirtualMachineImageRegistry interface {
	ListVirtualMachineImages(ctx context.Context, options *internalversion.ListOptions) (*VirtualMachineImageList, error)
	GetVirtualMachineImage(ctx context.Context, id string, options *metav1.GetOptions) (*VirtualMachineImage, error)
	CreateVirtualMachineImage(ctx context.Context, id *VirtualMachineImage) (*VirtualMachineImage, error)
	UpdateVirtualMachineImage(ctx context.Context, id *VirtualMachineImage) (*VirtualMachineImage, error)
	DeleteVirtualMachineImage(ctx context.Context, id string) (bool, error)
}

// NewRegistry returns a new Registry interface for the given Storage. Any mismatched types will panic.
func NewVirtualMachineImageRegistry(sp builders.StandardStorageProvider) VirtualMachineImageRegistry {
	return &storageVirtualMachineImage{sp}
}

// Implement Registry
// storage puts strong typing around storage calls
// +k8s:deepcopy-gen=false
type storageVirtualMachineImage struct {
	builders.StandardStorageProvider
}

func (s *storageVirtualMachineImage) ListVirtualMachineImages(ctx context.Context, options *internalversion.ListOptions) (*VirtualMachineImageList, error) {
	if options != nil && options.FieldSelector != nil && !options.FieldSelector.Empty() {
		return nil, fmt.Errorf("field selector not supported yet")
	}
	st := s.GetStandardStorage()
	obj, err := st.List(ctx, options)
	if err != nil {
		return nil, err
	}
	return obj.(*VirtualMachineImageList), err
}

func (s *storageVirtualMachineImage) GetVirtualMachineImage(ctx context.Context, id string, options *metav1.GetOptions) (*VirtualMachineImage, error) {
	st := s.GetStandardStorage()
	obj, err := st.Get(ctx, id, options)
	if err != nil {
		return nil, err
	}
	return obj.(*VirtualMachineImage), nil
}

func (s *storageVirtualMachineImage) CreateVirtualMachineImage(ctx context.Context, object *VirtualMachineImage) (*VirtualMachineImage, error) {
	st := s.GetStandardStorage()
	obj, err := st.Create(ctx, object, nil, &metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return obj.(*VirtualMachineImage), nil
}

func (s *storageVirtualMachineImage) UpdateVirtualMachineImage(ctx context.Context, object *VirtualMachineImage) (*VirtualMachineImage, error) {
	st := s.GetStandardStorage()
	obj, _, err := st.Update(ctx, object.Name, rest.DefaultUpdatedObjectInfo(object), nil, nil, false, &metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	return obj.(*VirtualMachineImage), nil
}

func (s *storageVirtualMachineImage) DeleteVirtualMachineImage(ctx context.Context, id string) (bool, error) {
	st := s.GetStandardStorage()
	_, sync, err := st.Delete(ctx, id, &metav1.DeleteOptions{})
	return sync, err
}

//
// VirtualMachineService Functions and Structs
//
// +k8s:deepcopy-gen=false
type VirtualMachineServiceStrategy struct {
	builders.DefaultStorageStrategy
}

// +k8s:deepcopy-gen=false
type VirtualMachineServiceStatusStrategy struct {
	builders.DefaultStatusStorageStrategy
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type VirtualMachineServiceList struct {
	metav1.TypeMeta
	metav1.ListMeta
	Items []VirtualMachineService
}

func (VirtualMachineService) NewStatus() interface{} {
	return VirtualMachineServiceStatus{}
}

func (pc *VirtualMachineService) GetStatus() interface{} {
	return pc.Status
}

func (pc *VirtualMachineService) SetStatus(s interface{}) {
	pc.Status = s.(VirtualMachineServiceStatus)
}

func (pc *VirtualMachineService) GetSpec() interface{} {
	return pc.Spec
}

func (pc *VirtualMachineService) SetSpec(s interface{}) {
	pc.Spec = s.(VirtualMachineServiceSpec)
}

func (pc *VirtualMachineService) GetObjectMeta() *metav1.ObjectMeta {
	return &pc.ObjectMeta
}

func (pc *VirtualMachineService) SetGeneration(generation int64) {
	pc.ObjectMeta.Generation = generation
}

func (pc VirtualMachineService) GetGeneration() int64 {
	return pc.ObjectMeta.Generation
}

// Registry is an interface for things that know how to store VirtualMachineService.
// +k8s:deepcopy-gen=false
type VirtualMachineServiceRegistry interface {
	ListVirtualMachineServices(ctx context.Context, options *internalversion.ListOptions) (*VirtualMachineServiceList, error)
	GetVirtualMachineService(ctx context.Context, id string, options *metav1.GetOptions) (*VirtualMachineService, error)
	CreateVirtualMachineService(ctx context.Context, id *VirtualMachineService) (*VirtualMachineService, error)
	UpdateVirtualMachineService(ctx context.Context, id *VirtualMachineService) (*VirtualMachineService, error)
	DeleteVirtualMachineService(ctx context.Context, id string) (bool, error)
}

// NewRegistry returns a new Registry interface for the given Storage. Any mismatched types will panic.
func NewVirtualMachineServiceRegistry(sp builders.StandardStorageProvider) VirtualMachineServiceRegistry {
	return &storageVirtualMachineService{sp}
}

// Implement Registry
// storage puts strong typing around storage calls
// +k8s:deepcopy-gen=false
type storageVirtualMachineService struct {
	builders.StandardStorageProvider
}

func (s *storageVirtualMachineService) ListVirtualMachineServices(ctx context.Context, options *internalversion.ListOptions) (*VirtualMachineServiceList, error) {
	if options != nil && options.FieldSelector != nil && !options.FieldSelector.Empty() {
		return nil, fmt.Errorf("field selector not supported yet")
	}
	st := s.GetStandardStorage()
	obj, err := st.List(ctx, options)
	if err != nil {
		return nil, err
	}
	return obj.(*VirtualMachineServiceList), err
}

func (s *storageVirtualMachineService) GetVirtualMachineService(ctx context.Context, id string, options *metav1.GetOptions) (*VirtualMachineService, error) {
	st := s.GetStandardStorage()
	obj, err := st.Get(ctx, id, options)
	if err != nil {
		return nil, err
	}
	return obj.(*VirtualMachineService), nil
}

func (s *storageVirtualMachineService) CreateVirtualMachineService(ctx context.Context, object *VirtualMachineService) (*VirtualMachineService, error) {
	st := s.GetStandardStorage()
	obj, err := st.Create(ctx, object, nil, &metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return obj.(*VirtualMachineService), nil
}

func (s *storageVirtualMachineService) UpdateVirtualMachineService(ctx context.Context, object *VirtualMachineService) (*VirtualMachineService, error) {
	st := s.GetStandardStorage()
	obj, _, err := st.Update(ctx, object.Name, rest.DefaultUpdatedObjectInfo(object), nil, nil, false, &metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	return obj.(*VirtualMachineService), nil
}

func (s *storageVirtualMachineService) DeleteVirtualMachineService(ctx context.Context, id string) (bool, error) {
	st := s.GetStandardStorage()
	_, sync, err := st.Delete(ctx, id, &metav1.DeleteOptions{})
	return sync, err
}
