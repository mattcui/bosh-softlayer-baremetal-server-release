package handlers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.cloudfoundry.org/lager/lagertest"
	"github.com/cloudfoundry-community/vps/models"
	"github.com/cloudfoundry-community/vps/restapi/handlers"
	"github.com/cloudfoundry-community/vps/restapi/handlers/fake_controllers"
	"github.com/cloudfoundry-community/vps/restapi/operations/vm"
	"github.com/go-openapi/runtime/middleware"
)

var _ = Describe("VmHandlerFunc", func() {
	var (
		logger     *lagertest.TestLogger
		controller *fake_controllers.FakeVirtualGuestController

		responseResponder middleware.Responder

		handler *handlers.VMHandler
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("test")
		controller = &fake_controllers.FakeVirtualGuestController{}
		handler = handlers.NewVmHandler(logger, controller)
	})

	Describe("AddVM", func() {
		var (
			vm1  *models.VM
			params vm.AddVMParams
		)

		BeforeEach(func() {
			vm1 = &models.VM{
				Cid: 1234567,
				CPU: 4,
				MemoryMb: 1024,
				IP: "10.0.0.1",
				Hostname: "fake.test.com",
				PrivateVlan: 1234567,
				PublicVlan: 1234568,
			}
			params = vm.NewAddVMParams()
			params.Body = vm1
		})

		JustBeforeEach(func() {
			responseResponder = handler.AddVM(params)
		})

		Context("when the virtual guest is added successful", func() {
			It("added into pool", func() {
				Expect(controller.CreateVMCallCount()).To(Equal(1))
				_, actualVm := controller.CreateVMArgsForCall(0)
				Expect(actualVm).To(Equal(vm1))

				addVmOk, ok := responseResponder.(*vm.AddVMOK)
				Expect(ok).To(BeTrue())
				Expect(addVmOk.GetPayload()).To(Equal("added successfully"))
			})
		})

		Context("when adding virtual guest fails", func() {
			BeforeEach(func() {
				controller.CreateVMReturns(models.ErrUnknownError)
			})


			It("responds with an error", func() {
				addVmDefault, ok := responseResponder.(*vm.AddVMDefault)
				Expect(ok).To(BeTrue())
				Expect(addVmDefault.GetStatusCode()).To(Equal(500))
				Expect(addVmDefault.GetPayload()).To(Equal(models.ErrUnknownError))
			})
		})
	})

	Describe("OrderVmByFilter", func() {
		var (
			vm1 *models.VM
			vmFilter *models.VMFilter
			params vm.OrderVMByFilterParams
		)

		BeforeEach(func() {
			vmFilter = &models.VMFilter{
				CPU: 4,
				MemoryMb: 1024,
				IP: "10.0.0.1",
				PrivateVlan: 1234567,
				PublicVlan: 1234568,
			}
			params.Body = vmFilter
		})

		JustBeforeEach(func() {
			responseResponder = handler.OrderVmByFilter(params)
		})

		Context("when ordering a VM by filter succeeds", func() {
			vm1 = &models.VM{
				Cid: 1234567,
				CPU: 4,
				MemoryMb: 1024,
				IP: "10.0.0.1",
				Hostname: "fake.test.com",
				PrivateVlan: 1234567,
				PublicVlan: 1234568,
			}

			BeforeEach(func() {
				controller.OrderVirtualGuestReturns(vm1, nil)
			})

			It("returns a VM", func() {
				Expect(controller.OrderVirtualGuestCallCount()).To(Equal(1))
				orderVMByFilterOK, ok := responseResponder.(*vm.OrderVMByFilterOK)
				Expect(ok).To(BeTrue())
				Expect(orderVMByFilterOK.GetPayload().VM).To(Equal(vm1))
			})

		})

		Context("when ordering a VM by filter fails due to an unexpected reason", func() {
			BeforeEach(func() {
				controller.OrderVirtualGuestReturns(nil, models.ErrUnknownError)
			})

			It("responds with an unExpected error", func() {
				orderVMByFilterDefault, ok := responseResponder.(*vm.OrderVMByFilterDefault)
				Expect(ok).To(BeTrue())
				Expect(orderVMByFilterDefault.GetStatusCode()).To(Equal(500))
				Expect(orderVMByFilterDefault.GetPayload()).To(Equal(models.ErrUnknownError))
			})

		})

		Context("when ordering a VM by filter fails due to resource not found", func() {
			BeforeEach(func() {
				controller.OrderVirtualGuestReturns(nil, models.ErrResourceNotFound)
			})

			It("responds with a notFound error", func() {
				orderVMByFilterNotFound, ok := responseResponder.(*vm.OrderVMByFilterNotFound)
				Expect(ok).To(BeTrue())
				Expect(orderVMByFilterNotFound.GetStatusCode()).To(Equal(404))
			})

		})

	})

	Describe("UpdateVM", func() {
		var (
			vm1 *models.VM
			params vm.UpdateVMParams
		)

		BeforeEach(func() {
			vm1 = &models.VM{
				Cid: 1234567,
				CPU: 4,
				MemoryMb: 1024,
				IP: "10.0.0.1",
				Hostname: "fake.test.com",
				PrivateVlan: 1234567,
				PublicVlan: 1234568,
			}
			params = vm.NewUpdateVMParams()
			params.Body = vm1
		})

		JustBeforeEach(func() {
			responseResponder = handler.UpdateVM(params)
		})

		Context("when updating virtual guests fails with an unexpected error", func() {
			BeforeEach(func() {
				controller.UpdateVMReturns(models.ErrUnknownError)
			})

			It("responds with an error", func() {
				updateVmDefault, ok := responseResponder.(*vm.UpdateVMDefault)
				Expect(ok).To(BeTrue())
				Expect(updateVmDefault.GetStatusCode()).To(Equal(500))
				Expect(updateVmDefault.GetPayload()).To(Equal(models.ErrUnknownError))
			})
		})

		Context("when updating virtual guests succeeds", func() {
			It("updates virtual guests specifications", func() {
				Expect(controller.UpdateVMCallCount()).To(Equal(1))
				_, actualVm := controller.UpdateVMArgsForCall(0)
				Expect(actualVm).To(Equal(vm1))

				updateVMOK, ok := responseResponder.(*vm.UpdateVMOK)
				Expect(ok).To(BeTrue())
				Expect(updateVMOK.GetPayload()).To(Equal("updated successfully"))
			})
		})
	})

	Describe("DeleteVM", func() {
		Context("when the delete request is normal", func() {
			var (
				params vm.DeleteVMParams
			)

			BeforeEach(func() {
				params = vm.NewDeleteVMParams()
				params.Cid = 1234567
			})

			JustBeforeEach(func() {
				responseResponder = handler.DeleteVM(params)
			})

			Context("when deleting the virtual guest succeeds", func() {
				It("returns no error", func() {
					Expect(controller.DeleteVMCallCount()).To(Equal(1))
					_, cid := controller.DeleteVMArgsForCall(0)
					Expect(cid).To(Equal(params.Cid))

					deleteVmNoContent, ok := responseResponder.(*vm.DeleteVMNoContent)
					Expect(ok).To(BeTrue())
					Expect(deleteVmNoContent.GetPayload()).To(Equal("vm removed"))
				})
			})

			Context("when deleting the vm fails due to resource not found", func() {
				BeforeEach(func() {
					controller.DeleteVMReturns(models.ErrResourceNotFound)
				})

				It("responds with a notFound error", func() {
					deleteVMNotFound, ok := responseResponder.(*vm.DeleteVMNotFound)
					Expect(ok).To(BeTrue())
					Expect(deleteVMNotFound.GetStatusCode()).To(Equal(404))
				})

			})

			Context("when the controller returns an error", func() {
				BeforeEach(func() {
					controller.DeleteVMReturns(models.ErrUnknownError)
				})

				It("provides relevant error information", func() {
					deleteVmDefault, ok := responseResponder.(*vm.DeleteVMDefault)
					Expect(ok).To(BeTrue())
					Expect(deleteVmDefault.GetStatusCode()).To(Equal(500))
					Expect(deleteVmDefault.GetPayload()).To(Equal(models.ErrUnknownError))
				})
			})
		})
	})

	Describe("GetVMByCid", func() {
		var (
			params vm.GetVMByCidParams
		)

		BeforeEach(func() {
			params = vm.NewGetVMByCidParams()
			params.Cid = 1234567
		})

		JustBeforeEach(func() {
			responseResponder = handler.GetVMByCid(params)
		})

		Context("when reading a virtual guest from the controller succeeds", func() {
			var vm1 *models.VM

			BeforeEach(func() {
				vm1 = &models.VM{Cid: params.Cid}
				controller.VirtualGuestByCidReturns(vm1, nil)
			})

			It("fetches virtual guest by cid", func() {
				Expect(controller.VirtualGuestByCidCallCount()).To(Equal(1))
				_, actualCid := controller.VirtualGuestByCidArgsForCall(0)
				Expect(actualCid).To(Equal(params.Cid))
			})

			It("returns the virtual guest", func() {
				getVmByCidOK, ok := responseResponder.(*vm.GetVMByCidOK)
				Expect(ok).To(BeTrue())
				Expect(getVmByCidOK.GetPayload().VM).To(Equal(vm1))
			})
		})

		Context("when the controller returns no virtual guest", func() {
			BeforeEach(func() {
				controller.VirtualGuestByCidReturns(nil, models.ErrResourceNotFound)
			})

			It("returns 404 status code", func() {
				getVmByCidNotFound, ok := responseResponder.(*vm.GetVMByCidNotFound)
				Expect(ok).To(BeTrue())
				Expect(getVmByCidNotFound.GetStatusCode()).To(Equal(404))
			})
		})

		Context("when the controller errors out", func() {
			BeforeEach(func() {
				controller.VirtualGuestByCidReturns(nil, models.ErrUnknownError)
			})

			It("provides relevant error information", func() {
				getVmByCidDefault, ok := responseResponder.(*vm.GetVMByCidDefault)
				Expect(ok).To(BeTrue())
				Expect(getVmByCidDefault.GetStatusCode()).To(Equal(500))
				Expect(getVmByCidDefault.GetPayload()).To(Equal(models.ErrUnknownError))
			})
		})
	})

	Describe("ListVM", func() {
		var (
			vm1 models.VM
			vm2 models.VM
			params vm.ListVMParams
		)

		BeforeEach(func() {
			vm1 = models.VM{Cid: 1234567}
			vm2 = models.VM{Cid: 1234568}
			params = vm.NewListVMParams()
		})

		JustBeforeEach(func() {
			responseResponder = handler.ListVM(params)
		})

		Context("when reading virtual guests from controller succeeds", func() {
			var vms []*models.VM

			BeforeEach(func() {
				vms = []*models.VM{&vm1, &vm2}
				controller.AllVirtualGuestsReturns(vms, nil)
			})

			It("returns a list of virtual guests", func() {
				Expect(controller.AllVirtualGuestsCallCount()).To(Equal(1))
				listVMOK, ok := responseResponder.(*vm.ListVMOK)
				Expect(ok).To(BeTrue())
				Expect(listVMOK.GetPayload().Vms).To(Equal(vms))
			})
		})

		Context("when the controller returns no virtual guest", func() {
			var vms []*models.VM

			BeforeEach(func() {
				vms = []*models.VM{}
				controller.AllVirtualGuestsReturns(vms, nil)
			})

			It("returns 404 status code", func() {
				listVmNotFound, ok := responseResponder.(*vm.ListVMNotFound)
				Expect(ok).To(BeTrue())
				Expect(listVmNotFound.GetStatusCode()).To(Equal(404))
			})
		})

		Context("when the controller errors out", func() {
			BeforeEach(func() {
				controller.AllVirtualGuestsReturns(nil, models.ErrUnknownError)
			})

			It("provides relevant error information", func() {
				listVmDefault, ok := responseResponder.(*vm.ListVMDefault)
				Expect(ok).To(BeTrue())
				Expect(listVmDefault.GetStatusCode()).To(Equal(500))
				Expect(listVmDefault.GetPayload()).To(Equal(models.ErrUnknownError))
			})
		})
	})

	Describe("UpdateVMWithState", func() {
		var (
			vmState *models.VMState
			params vm.UpdateVMWithStateParams
		)

		BeforeEach(func() {
			vmState = &models.VMState{
				State:  models.StateFree,
			}
			params = vm.NewUpdateVMWithStateParams()
			params.Body = vmState
		})

		JustBeforeEach(func() {
			responseResponder = handler.UpdateVMWithState(params)
		})

		Context("when updating VM state fails due to the resource not found", func() {
			BeforeEach(func() {
				controller.UpdateVMWithStateReturns(models.ErrResourceNotFound)
			})

			It("responds with an error", func() {
				updateVMWithStateNotFound, ok := responseResponder.(*vm.UpdateVMWithStateNotFound)
				Expect(ok).To(BeTrue())
				Expect(updateVMWithStateNotFound.GetStatusCode()).To(Equal(404))
			})
		})

		Context("when updating VM state errors out", func() {
			BeforeEach(func() {
				params.Cid = 1234567
				controller.UpdateVMWithStateReturns(models.ErrUnknownError)
			})

			It("responds with an error", func() {
				updateVMWithStateDefault, ok := responseResponder.(*vm.UpdateVMWithStateDefault)
				Expect(ok).To(BeTrue())
				Expect(updateVMWithStateDefault.GetStatusCode()).To(Equal(500))
				Expect(updateVMWithStateDefault.GetPayload()).To(Equal(models.ErrUnknownError))
			})
		})

		Context("when updating VM state succeeds", func() {
			BeforeEach(func() {
				params.Cid = 1234567
			})

			It("updates VM's state", func() {
				Expect(controller.UpdateVMWithStateCallCount()).To(Equal(1))
				_, vmId, actualState := controller.UpdateVMWithStateArgsForCall(0)
				Expect(vmId).To(Equal(int32(1234567)))
				Expect(actualState).To(Equal(&vmState.State))

				updateVMWithStateOK, ok := responseResponder.(*vm.UpdateVMWithStateOK)
				Expect(ok).To(BeTrue())
				Expect(updateVMWithStateOK.GetPayload()).To(Equal("updated successfully"))
			})
		})
	})

	Describe("FindVmsByFilters", func() {
		var (
			vm1 models.VM
			vm2 models.VM
			vmFilter *models.VMFilter
			params vm.FindVmsByFiltersParams
		)

		BeforeEach(func() {
			vm1 = models.VM{CPU: 2}
			vm2 = models.VM{CPU: 4}
			params = vm.NewFindVmsByFiltersParams()
		})

		JustBeforeEach(func() {
			responseResponder = handler.FindVmsByFilters(params)
		})

		Context("when filter is empty and finding VMs by filter succeeds", func() {
			var vms []*models.VM

			BeforeEach(func() {
				vms = []*models.VM{&vm1, &vm2}
				controller.AllVirtualGuestsReturns(vms, nil)
			})

			It("returns all vitual guests", func() {
				Expect(controller.AllVirtualGuestsCallCount()).To(Equal(1))
				findVmsByFiltersOK, ok := responseResponder.(*vm.FindVmsByFiltersOK)
				Expect(ok).To(BeTrue())
				Expect(findVmsByFiltersOK.GetPayload().Vms).To(Equal(vms))
			})

		})

		Context("when filter is not empty finding VMs by filter succeeds", func() {
			var vms []*models.VM

			BeforeEach(func() {
				vms = []*models.VM{&vm1}
				vmFilter = &models.VMFilter{CPU:2}
				params.Body = vmFilter
				controller.VirtualGuestsReturns(vms, nil)
			})

			It("returns virtual guests that match the filter", func() {
				Expect(controller.VirtualGuestsCallCount()).To(Equal(1))
				findVmsByFiltersOK, ok := responseResponder.(*vm.FindVmsByFiltersOK)
				Expect(ok).To(BeTrue())
				Expect(findVmsByFiltersOK.GetPayload().Vms).To(Equal(vms))
			})

		})

		Context("when filter is empty and finding VMs by filter fails", func() {
			BeforeEach(func() {
				controller.AllVirtualGuestsReturns(nil, models.ErrUnknownError)
			})

			It("responds with an error", func() {
				findVmsByFiltersDefault, ok := responseResponder.(*vm.FindVmsByFiltersDefault)
				Expect(ok).To(BeTrue())
				Expect(findVmsByFiltersDefault.GetStatusCode()).To(Equal(500))
				Expect(findVmsByFiltersDefault.GetPayload()).To(Equal(models.ErrUnknownError))
			})

		})

		Context("when filter is not empty and finding VMs by filter fails", func() {
			BeforeEach(func() {
				vmFilter = &models.VMFilter{CPU:2}
				params.Body = vmFilter
				controller.VirtualGuestsReturns(nil, models.ErrUnknownError)
			})

			It("responds with an error", func() {
				findVmsByFiltersDefault, ok := responseResponder.(*vm.FindVmsByFiltersDefault)
				Expect(ok).To(BeTrue())
				Expect(findVmsByFiltersDefault.GetStatusCode()).To(Equal(500))
				Expect(findVmsByFiltersDefault.GetPayload()).To(Equal(models.ErrUnknownError))
			})
		})
	})

	Describe("FindVmsByDeployment", func() {
		var (
			vm1 models.VM
			vm2 models.VM
			params vm.FindVmsByDeploymentParams
		)

		BeforeEach(func() {
			vm1 = models.VM{Cid: 1234567}
			vm2 = models.VM{Cid: 1234568}
			params = vm.NewFindVmsByDeploymentParams()
			params.Deployment = []string{"deployment1", "deployment2"}
		})

		JustBeforeEach(func() {
			responseResponder = handler.FindVmsByDeployment(params)
		})

		Context("when finding VMs by deployments succeeds", func() {
			var vms []*models.VM

			BeforeEach(func() {
				vms = []*models.VM{&vm1, &vm2}
				controller.VirtualGuestsByDeploymentsReturns(vms, nil)
			})

			It("returns VMs which belong to the deployments", func() {
				Expect(controller.VirtualGuestsByDeploymentsCallCount()).To(Equal(1))
				_, deploymentName := controller.VirtualGuestsByDeploymentsArgsForCall(0)
				Expect(deploymentName).To(Equal(params.Deployment))

				findVmsByDeploymentOK, ok := responseResponder.(*vm.FindVmsByDeploymentOK)
				Expect(ok).To(BeTrue())
				Expect(findVmsByDeploymentOK.GetPayload().Vms).To(Equal(vms))
			})
		})

		Context("when finding VMs by deployments fails", func() {
			BeforeEach(func() {
				controller.VirtualGuestsByDeploymentsReturns(nil, models.ErrUnknownError)
			})

			It("responds with an error", func() {
				findVmsByDeploymentDefault, ok := responseResponder.(*vm.FindVmsByDeploymentDefault)
				Expect(ok).To(BeTrue())
				Expect(findVmsByDeploymentDefault.GetStatusCode()).To(Equal(500))
				Expect(findVmsByDeploymentDefault.GetPayload()).To(Equal(models.ErrUnknownError))
			})
		})
	})

	Describe("FindVmsByStates", func() {
		var (
			vm1 models.VM
			vm2 models.VM
			params vm.FindVmsByStatesParams
		)

		BeforeEach(func() {
			vm1 = models.VM{Cid: 1234567, State: "free"}
			vm2 = models.VM{Cid: 1234568, State: "using"}
			params = vm.NewFindVmsByStatesParams()
			params.States = []string{"free", "using"}
		})

		JustBeforeEach(func() {
			responseResponder = handler.FindVmsByStates(params)
		})

		Context("when finding VMs by states succeeds", func() {
			var vms []*models.VM

			BeforeEach(func() {
				vms = []*models.VM{&vm1, &vm2}
				controller.VirtualGuestsByStatesReturns(vms, nil)
			})

			It("returns VMs which match the states", func() {
				Expect(controller.VirtualGuestsByStatesCallCount()).To(Equal(1))
				_, states := controller.VirtualGuestsByStatesArgsForCall(0)
				Expect(states).To(Equal(params.States))

				findVmsByStatesOK, ok := responseResponder.(*vm.FindVmsByStatesOK)
				Expect(ok).To(BeTrue())
				Expect(findVmsByStatesOK.GetPayload().Vms).To(Equal(vms))
			})
		})

		Context("when finding VMs by states fails", func() {
			BeforeEach(func() {
				controller.VirtualGuestsByStatesReturns(nil, models.ErrUnknownError)
			})

			It("responds with an error", func() {
				findVmsByStatesDefault, ok := responseResponder.(*vm.FindVmsByStatesDefault)
				Expect(ok).To(BeTrue())
				Expect(findVmsByStatesDefault.GetStatusCode()).To(Equal(500))
				Expect(findVmsByStatesDefault.GetPayload()).To(Equal(models.ErrUnknownError))
			})
		})
	})
})
