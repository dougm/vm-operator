// Copyright (c) 2019-2023 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package mutation_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/controller-runtime/pkg/client"

	vmopv1 "github.com/vmware-tanzu/vm-operator/api/v1alpha1"
	"github.com/vmware-tanzu/vm-operator/pkg"
	pkgconfig "github.com/vmware-tanzu/vm-operator/pkg/config"
	"github.com/vmware-tanzu/vm-operator/pkg/constants"
	"github.com/vmware-tanzu/vm-operator/pkg/vmprovider/providers/vsphere/network"
	"github.com/vmware-tanzu/vm-operator/test/builder"
)

func intgTests() {
	Describe("Invoking Mutation", intgTestsMutating)
}

type intgMutatingWebhookContext struct {
	builder.IntegrationTestContext
	vm *vmopv1.VirtualMachine
}

func newIntgMutatingWebhookContext() *intgMutatingWebhookContext {
	ctx := &intgMutatingWebhookContext{
		IntegrationTestContext: *suite.NewIntegrationTestContext(),
	}

	ctx.vm = builder.DummyVirtualMachine()
	ctx.vm.Namespace = ctx.Namespace

	return ctx
}

func intgTestsMutating() {
	var (
		ctx *intgMutatingWebhookContext
		vm  *vmopv1.VirtualMachine
	)

	BeforeEach(func() {
		ctx = newIntgMutatingWebhookContext()
		vm = ctx.vm.DeepCopy()
		vm.Spec.NetworkInterfaces = []vmopv1.VirtualMachineNetworkInterface{}
		pkgconfig.SetContext(suite, func(config *pkgconfig.Config) {
			config.NetworkProviderType = pkgconfig.NetworkProviderTypeVDS
		})
	})
	AfterEach(func() {
		pkgconfig.SetContext(suite, func(config *pkgconfig.Config) {
			config.NetworkProviderType = ""
		})
		ctx = nil
	})

	Describe("mutate", func() {
		Context("placeholder", func() {
			BeforeEach(func() {
			})

			AfterEach(func() {
				Expect(ctx.Client.Delete(ctx, vm)).Should(Succeed())
			})

			It("should work", func() {
				err := ctx.Client.Create(ctx, vm)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("Default network interface", func() {
			When("Creating VirtualMachine", func() {
				It("Add default network interface if NetworkInterface is empty and no Annotation", func() {
					err := ctx.Client.Create(ctx, vm)
					Expect(err).ToNot(HaveOccurred())

					modified := &vmopv1.VirtualMachine{}
					Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(vm), modified)).Should(Succeed())
					Expect(modified.Spec.NetworkInterfaces).Should(HaveLen(1))
					Expect(modified.Spec.NetworkInterfaces[0].NetworkType).Should(Equal(network.VdsNetworkType))
					Expect(modified.Spec.NetworkInterfaces[0].NetworkName).Should(Equal(""))
				})
			})

			When("Updating VirtualMachine", func() {
				BeforeEach(func() {
					vm.Annotations[vmopv1.NoDefaultNicAnnotation] = ""
					Expect(ctx.Client.Create(ctx, vm)).To(Succeed())
				})

				AfterEach(func() {
					Expect(ctx.Client.Delete(ctx, vm)).Should(Succeed())
				})

				It("should not add default network interface", func() {
					modified := &vmopv1.VirtualMachine{}
					vmKey := client.ObjectKeyFromObject(vm)
					Expect(ctx.Client.Get(ctx, vmKey, modified)).Should(Succeed())

					delete(modified.Annotations, vmopv1.NoDefaultNicAnnotation)
					Expect(ctx.Client.Update(ctx, modified)).Should(Succeed())
					Expect(ctx.Client.Get(ctx, vmKey, modified)).Should(Succeed())
					Expect(modified.Spec.NetworkInterfaces).Should(BeEmpty())
				})
			})

		})

		Context("SetDefaultPowerState", func() {
			When("Creating VirtualMachine", func() {
				When("When VM PowerState is empty", func() {
					BeforeEach(func() {
						vm.Spec.PowerState = ""
					})
					It("Should set PowerState to PoweredOn", func() {
						Expect(ctx.Client.Create(ctx, vm)).To(Succeed())
						modified := &vmopv1.VirtualMachine{}
						Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(vm), modified)).Should(Succeed())
						Expect(modified.Spec.PowerState).Should(Equal(vmopv1.VirtualMachinePoweredOn))
					})
				})
				When("When VM PowerState is not empty", func() {
					BeforeEach(func() {
						vm.Spec.PowerState = vmopv1.VirtualMachinePoweredOff
					})
					It("Should not mutate PowerState", func() {
						Expect(ctx.Client.Create(ctx, vm)).To(Succeed())
						modified := &vmopv1.VirtualMachine{}
						Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(vm), modified)).Should(Succeed())
						Expect(modified.Spec.PowerState).Should(Equal(vmopv1.VirtualMachinePoweredOff))
					})
				})
			})

			When("Updating VirtualMachine", func() {
				BeforeEach(func() {
					vm.Spec.PowerState = vmopv1.VirtualMachinePoweredOff
					Expect(ctx.Client.Create(ctx, vm)).To(Succeed())
				})
				// This state is not technically possible in production. However,
				// it is used to validate that the power state is not auto-set
				// to poweredOn if empty during an empty. Since the logic for
				// defaulting to poweredOn only works if empty (and on create),
				// it's necessary to replicate the empty state here.
				When("When VM PowerState is empty", func() {
					It("Should not mutate PowerState", func() {
						modified := &vmopv1.VirtualMachine{}
						Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(vm), modified)).Should(Succeed())
						modified.Spec.PowerState = ""
						Expect(ctx.Client.Update(ctx, modified)).To(Succeed())
						Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(vm), modified)).Should(Succeed())
						Expect(modified.Spec.PowerState).Should(BeEmpty())
					})
				})
				When("When VM PowerState is not empty", func() {
					It("Should not mutate PowerState", func() {
						modified := &vmopv1.VirtualMachine{}
						Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(vm), modified)).Should(Succeed())
						modified.Spec.PowerState = vmopv1.VirtualMachinePoweredOff
						Expect(ctx.Client.Update(ctx, modified)).To(Succeed())
						Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(vm), modified)).Should(Succeed())
						Expect(modified.Spec.PowerState).Should(Equal(vmopv1.VirtualMachinePoweredOff))
					})
				})
			})

		})

		Context("ResolveImageName", func() {

			BeforeEach(func() {
				pkgconfig.SetContext(suite, func(config *pkgconfig.Config) {
					config.Features.ImageRegistry = true
				})
			})

			AfterEach(func() {
				pkgconfig.SetContext(suite, func(config *pkgconfig.Config) {
					config.Features.ImageRegistry = false
				})
			})

			When("Creating VirtualMachine", func() {

				When("VM ImageName is already a vmi resource name", func() {

					BeforeEach(func() {
						vm.Spec.ImageName = "vmi-123"
					})

					It("Should not mutate ImageName", func() {
						Expect(ctx.Client.Create(ctx, vm)).To(Succeed())
						modified := &vmopv1.VirtualMachine{}
						Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(vm), modified)).Should(Succeed())
						Expect(modified.Spec.ImageName).Should(Equal("vmi-123"))
					})
				})
			})
		})

		Context("SetNextRestartTime", func() {
			When("create a VM", func() {
				When("spec.nextRestartTime is empty", func() {
					It("should not mutate", func() {
						vm.Spec.NextRestartTime = ""
						Expect(ctx.Client.Create(ctx, vm)).To(Succeed())
						modified := &vmopv1.VirtualMachine{}
						Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(vm), modified)).Should(Succeed())
						Expect(modified.Spec.NextRestartTime).Should(BeEmpty())
					})
				})
				When("spec.nextRestartTime is not empty", func() {
					It("should not mutate", func() {
						vm.Spec.NextRestartTime = "hello"
						Expect(ctx.Client.Create(ctx, vm)).To(Succeed())
						modified := &vmopv1.VirtualMachine{}
						Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(vm), modified)).Should(Succeed())
						Expect(modified.Spec.NextRestartTime).Should(Equal("hello"))
					})
				})
			})

			When("updating VirtualMachine", func() {
				var (
					lastRestartTime time.Time
					modified        *vmopv1.VirtualMachine
				)
				BeforeEach(func() {
					lastRestartTime = time.Now().UTC()
					modified = &vmopv1.VirtualMachine{}
				})
				JustBeforeEach(func() {
					Expect(ctx.Client.Create(ctx, vm)).To(Succeed())
					Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(vm), modified)).Should(Succeed())
				})
				When("spec.nextRestartTime is empty", func() {
					BeforeEach(func() {
						vm.Spec.NextRestartTime = ""
					})
					When("spec.nextRestartTime is empty", func() {
						It("should not mutate", func() {
							modified.Spec.NextRestartTime = ""
							Expect(ctx.Client.Update(ctx, modified)).To(Succeed())
							Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(vm), modified)).Should(Succeed())
							Expect(modified.Spec.NextRestartTime).Should(BeEmpty())
						})
					})
					When("spec.nextRestartTime is now", func() {
						It("should mutate", func() {
							modified.Spec.NextRestartTime = "Now"
							Expect(ctx.Client.Update(ctx, modified)).To(Succeed())
							Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(vm), modified)).Should(Succeed())
							Expect(modified.Spec.NextRestartTime).ShouldNot(BeEmpty())
							_, err := time.Parse(time.RFC3339Nano, modified.Spec.NextRestartTime)
							Expect(err).ToNot(HaveOccurred())
						})
					})
				})
				When("existing spec.nextRestartTime is not empty", func() {
					BeforeEach(func() {
						vm.Spec.NextRestartTime = lastRestartTime.Format(time.RFC3339Nano)
					})
					When("spec.nextRestartTime is empty", func() {
						It("should mutate to original value", func() {
							modified.Spec.NextRestartTime = ""
							Expect(ctx.Client.Update(ctx, modified)).To(Succeed())
							Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(vm), modified)).Should(Succeed())
							Expect(modified.Spec.NextRestartTime).Should(Equal(vm.Spec.NextRestartTime))
						})
					})
					When("spec.nextRestartTime is now", func() {
						It("should mutate", func() {
							modified.Spec.NextRestartTime = "Now"
							Expect(ctx.Client.Update(ctx, modified)).To(Succeed())
							Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(vm), modified)).Should(Succeed())
							Expect(modified.Spec.NextRestartTime).ShouldNot(BeEmpty())
							nextRestartTime, err := time.Parse(time.RFC3339Nano, modified.Spec.NextRestartTime)
							Expect(err).ToNot(HaveOccurred())
							Expect(lastRestartTime.Before(nextRestartTime)).To(BeTrue())
						})
					})
					DescribeTable(
						`spec.nextRestartTime is a non-empty value that is not "now"`,
						append([]any{
							func(nextRestartTime string) {
								modified.Spec.NextRestartTime = nextRestartTime
								err := ctx.Client.Update(ctx, modified)
								Expect(err).To(HaveOccurred())
								expectedErr := field.Invalid(
									field.NewPath("spec", "nextRestartTime"),
									nextRestartTime,
									`may only be set to "now"`)
								Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
							}},
							newInvalidNextRestartTimeTableEntries("should return an invalid field error"))...,
					)
				})
			})
		})
	})

	Context("SetCreatedAtAnnotations", func() {
		When("creating a VM", func() {
			It("should add the created-at annotations", func() {
				Expect(ctx.Client.Create(ctx, ctx.vm)).To(Succeed())
				vm := &vmopv1.VirtualMachine{}
				Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(ctx.vm), vm)).To(Succeed())
				Expect(vm.Annotations).To(HaveKeyWithValue(constants.CreatedAtBuildVersionAnnotationKey, "v1"))
				Expect(vm.Annotations).To(HaveKeyWithValue(constants.CreatedAtSchemaVersionAnnotationKey, vmopv1.SchemeGroupVersion.Version))
			})
		})

		When("updating a VM", func() {
			var oldBuildVersion string
			BeforeEach(func() {
				oldBuildVersion = pkg.BuildVersion
			})
			AfterEach(func() {
				pkg.BuildVersion = oldBuildVersion
			})
			It("should not update the created-at annotations", func() {
				Expect(ctx.Client.Create(ctx, ctx.vm)).To(Succeed())
				vm := &vmopv1.VirtualMachine{}
				Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(ctx.vm), vm)).To(Succeed())
				Expect(vm.Annotations).To(HaveKeyWithValue(constants.CreatedAtBuildVersionAnnotationKey, "v1"))
				Expect(vm.Annotations).To(HaveKeyWithValue(constants.CreatedAtSchemaVersionAnnotationKey, vmopv1.SchemeGroupVersion.Version))

				pkg.BuildVersion = "v2"

				Expect(ctx.Client.Update(ctx, vm)).To(Succeed())
				Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(ctx.vm), vm)).To(Succeed())
				Expect(vm.Annotations).To(HaveKeyWithValue(constants.CreatedAtBuildVersionAnnotationKey, "v1"))
				Expect(vm.Annotations).To(HaveKeyWithValue(constants.CreatedAtSchemaVersionAnnotationKey, vmopv1.SchemeGroupVersion.Version))
			})
		})
	})
}
