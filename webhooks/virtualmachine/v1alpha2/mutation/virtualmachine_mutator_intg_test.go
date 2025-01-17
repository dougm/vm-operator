// Copyright (c) 2019-2023 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package mutation_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/controller-runtime/pkg/client"

	vmopv1 "github.com/vmware-tanzu/vm-operator/api/v1alpha2"
	"github.com/vmware-tanzu/vm-operator/pkg"
	pkgconfig "github.com/vmware-tanzu/vm-operator/pkg/config"
	"github.com/vmware-tanzu/vm-operator/pkg/constants"
	"github.com/vmware-tanzu/vm-operator/pkg/constants/testlabels"
	"github.com/vmware-tanzu/vm-operator/test/builder"
)

func intgTests() {
	Describe(
		"Mutate",
		Label(
			testlabels.Create,
			testlabels.Update,
			testlabels.Delete,
			testlabels.EnvTest,
			testlabels.V1Alpha2,
			testlabels.Mutation,
			testlabels.Webhook,
		),
		intgTestsMutating,
	)
}

type intgMutatingWebhookContext struct {
	builder.IntegrationTestContext
	vm *vmopv1.VirtualMachine
}

func newIntgMutatingWebhookContext() *intgMutatingWebhookContext {
	ctx := &intgMutatingWebhookContext{
		IntegrationTestContext: *suite.NewIntegrationTestContext(),
	}

	ctx.vm = builder.DummyVirtualMachineA2()
	ctx.vm.Namespace = ctx.Namespace

	return ctx
}

func intgTestsMutating() {
	var (
		ctx *intgMutatingWebhookContext
		img *vmopv1.VirtualMachineImage
	)

	BeforeEach(func() {
		ctx = newIntgMutatingWebhookContext()

		img = builder.DummyVirtualMachineImageA2(builder.DummyVMIID)
		img.Namespace = ctx.vm.Namespace
		Expect(ctx.Client.Create(ctx, img)).To(Succeed())
		img.Status.Name = ctx.vm.Spec.ImageName
		Expect(ctx.Client.Status().Update(ctx, img)).To(Succeed())
	})

	AfterEach(func() {
		Expect(ctx.Client.Delete(ctx, img)).To(Succeed())
		img = nil
		ctx = nil
	})

	Describe("mutate", func() {
		Context("placeholder", func() {
			AfterEach(func() {
				Expect(ctx.Client.Delete(ctx, ctx.vm)).To(Succeed())
			})

			It("should work", func() {
				err := ctx.Client.Create(ctx, ctx.vm)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("Default network interface", func() {
			BeforeEach(func() {
				pkgconfig.SetContext(suite, func(config *pkgconfig.Config) {
					config.NetworkProviderType = pkgconfig.NetworkProviderTypeVDS
				})
				ctx.vm.Spec.Network.Interfaces = nil
			})
			AfterEach(func() {
				pkgconfig.SetContext(suite, func(config *pkgconfig.Config) {
					config.NetworkProviderType = ""
				})
			})

			When("Creating VirtualMachine", func() {
				It("Add default network interface if NetworkInterface is empty and no Annotation", func() {
					Expect(ctx.Client.Create(ctx, ctx.vm)).To(Succeed())

					modified := &vmopv1.VirtualMachine{}
					Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(ctx.vm), modified)).To(Succeed())
					Expect(modified.Spec.Network.Interfaces).To(HaveLen(1))
					Expect(modified.Spec.Network.Interfaces[0].Name).To(Equal("eth0"))
					Expect(modified.Spec.Network.Interfaces[0].Network.Kind).To(Equal("Network"))
				})
			})
		})
	})

	Context("SetDefaultPowerState", func() {
		When("Creating VirtualMachine", func() {
			When("When VM PowerState is empty", func() {
				BeforeEach(func() {
					ctx.vm.Spec.PowerState = ""
				})

				It("Should set PowerState to PoweredOn", func() {
					Expect(ctx.Client.Create(ctx, ctx.vm)).To(Succeed())

					modified := &vmopv1.VirtualMachine{}
					Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(ctx.vm), modified)).To(Succeed())
					Expect(modified.Spec.PowerState).To(Equal(vmopv1.VirtualMachinePowerStateOn))
				})
			})
			When("When VM PowerState is not empty", func() {
				BeforeEach(func() {
					ctx.vm.Spec.PowerState = vmopv1.VirtualMachinePowerStateOff
				})

				It("Should not mutate PowerState", func() {
					Expect(ctx.Client.Create(ctx, ctx.vm)).To(Succeed())
					modified := &vmopv1.VirtualMachine{}
					Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(ctx.vm), modified)).To(Succeed())
					Expect(modified.Spec.PowerState).To(Equal(vmopv1.VirtualMachinePowerStateOff))
				})
			})
		})

		When("Updating VirtualMachine", func() {
			BeforeEach(func() {
				ctx.vm.Spec.PowerState = vmopv1.VirtualMachinePowerStateOff
				Expect(ctx.Client.Create(ctx, ctx.vm)).To(Succeed())
			})

			// This state is not technically possible in production. However,
			// it is used to validate that the power state is not auto-set
			// to poweredOn if empty during an empty. Since the logic for
			// defaulting to poweredOn only works if empty (and on create),
			// it's necessary to replicate the empty state here.
			When("When VM PowerState is empty", func() {
				It("Should not mutate PowerState", func() {
					modified := &vmopv1.VirtualMachine{}
					Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(ctx.vm), modified)).To(Succeed())
					modified.Spec.PowerState = ""
					Expect(ctx.Client.Update(ctx, modified)).To(Succeed())
					Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(ctx.vm), modified)).To(Succeed())
					Expect(modified.Spec.PowerState).To(BeEmpty())
				})
			})

			When("When VM PowerState is not empty", func() {
				It("Should not mutate PowerState", func() {
					modified := &vmopv1.VirtualMachine{}
					Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(ctx.vm), modified)).To(Succeed())
					modified.Spec.PowerState = vmopv1.VirtualMachinePowerStateOff
					Expect(ctx.Client.Update(ctx, modified)).To(Succeed())
					Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(ctx.vm), modified)).To(Succeed())
					Expect(modified.Spec.PowerState).To(Equal(vmopv1.VirtualMachinePowerStateOff))
				})
			})
		})
	})

	Context("ResolveImageName", func() {
		When("Creating VirtualMachine", func() {
			When("VM ImageName is already a vmi resource name", func() {

				BeforeEach(func() {
					ctx.vm.Spec.ImageName = "vmi-123"
				})

				It("Should not mutate ImageName", func() {
					Expect(ctx.Client.Create(ctx, ctx.vm)).To(Succeed())
					modified := &vmopv1.VirtualMachine{}
					Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(ctx.vm), modified)).To(Succeed())
					Expect(modified.Spec.ImageName).To(Equal("vmi-123"))
				})
			})
		})
	})

	Context("SetNextRestartTime", func() {
		When("create a VM", func() {
			When("spec.nextRestartTime is empty", func() {
				It("should not mutate", func() {
					ctx.vm.Spec.NextRestartTime = ""
					Expect(ctx.Client.Create(ctx, ctx.vm)).To(Succeed())
					modified := &vmopv1.VirtualMachine{}
					Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(ctx.vm), modified)).To(Succeed())
					Expect(modified.Spec.NextRestartTime).To(BeEmpty())
				})
			})
			When("spec.nextRestartTime is not empty", func() {
				It("should not mutate", func() {
					ctx.vm.Spec.NextRestartTime = "hello"
					Expect(ctx.Client.Create(ctx, ctx.vm)).To(Succeed())
					modified := &vmopv1.VirtualMachine{}
					Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(ctx.vm), modified)).To(Succeed())
					Expect(modified.Spec.NextRestartTime).To(Equal("hello"))
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
				Expect(ctx.Client.Create(ctx, ctx.vm)).To(Succeed())
				Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(ctx.vm), modified)).To(Succeed())
			})
			When("spec.nextRestartTime is empty", func() {
				BeforeEach(func() {
					ctx.vm.Spec.NextRestartTime = ""
				})
				When("spec.nextRestartTime is empty", func() {
					It("should not mutate", func() {
						modified.Spec.NextRestartTime = ""
						Expect(ctx.Client.Update(ctx, modified)).To(Succeed())
						Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(ctx.vm), modified)).To(Succeed())
						Expect(modified.Spec.NextRestartTime).To(BeEmpty())
					})
				})
				When("spec.nextRestartTime is now", func() {
					It("should mutate", func() {
						modified.Spec.NextRestartTime = "Now"
						Expect(ctx.Client.Update(ctx, modified)).To(Succeed())
						Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(ctx.vm), modified)).To(Succeed())
						Expect(modified.Spec.NextRestartTime).ToNot(BeEmpty())
						_, err := time.Parse(time.RFC3339Nano, modified.Spec.NextRestartTime)
						Expect(err).ToNot(HaveOccurred())
					})
				})
			})
			When("existing spec.nextRestartTime is not empty", func() {
				BeforeEach(func() {
					ctx.vm.Spec.NextRestartTime = lastRestartTime.Format(time.RFC3339Nano)
				})
				When("spec.nextRestartTime is empty", func() {
					It("should mutate to original value", func() {
						modified.Spec.NextRestartTime = ""
						Expect(ctx.Client.Update(ctx, modified)).To(Succeed())
						Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(ctx.vm), modified)).To(Succeed())
						Expect(modified.Spec.NextRestartTime).To(Equal(ctx.vm.Spec.NextRestartTime))
					})
				})
				When("spec.nextRestartTime is now", func() {
					It("should mutate", func() {
						modified.Spec.NextRestartTime = "Now"
						Expect(ctx.Client.Update(ctx, modified)).To(Succeed())
						Expect(ctx.Client.Get(ctx, client.ObjectKeyFromObject(ctx.vm), modified)).To(Succeed())
						Expect(modified.Spec.NextRestartTime).ToNot(BeEmpty())
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
