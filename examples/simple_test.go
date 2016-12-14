package examples

import (
	"fmt"

	"github.com/MasterOfBinary/go-opencl/opencl"
)

func ExampleSimple() {
	const (
		dataSize = 32

		programCode = `
kernel void kern(global float* out)
{
	size_t i = get_global_id(0);
	out[i] = i;
}
`
	)

	platforms, err := opencl.GetPlatforms()
	if err != nil {
		panic(err)
	}

	var cpuDevice *opencl.Device

	var name string
	for _, platform := range platforms {
		err = platform.GetInfo(opencl.PlatformName, &name)
		if err != nil {
			panic(err)
		}
		var devices []*opencl.Device
		devices, err = platform.GetDevices(opencl.DeviceTypeCPU)
		if err != nil {
			panic(err)
		}

		// Use the first device
		if len(devices) > 0 && cpuDevice == nil {
			cpuDevice = devices[0]
		}
	}

	if cpuDevice == nil {
		panic("No device found")
	}

	var context *opencl.Context
	context, err = cpuDevice.CreateContext()
	if err != nil {
		panic(err)
	}
	defer context.Release()

	var commandQueue *opencl.CommandQueue
	commandQueue, err = context.CreateCommandQueue()
	if err != nil {
		panic(err)
	}
	defer commandQueue.Release()

	var program *opencl.Program
	program, err = context.CreateProgramWithSource(programCode)
	if err != nil {
		panic(err)
	}
	defer program.Release()

	var log string
	err = program.Build(&log)
	if err != nil {
		fmt.Println(log)
		panic(err)
	}

	kernel, err := program.CreateKernel("kern")
	if err != nil {
		panic(err)
	}
	defer kernel.Release()

	buffer, err := context.CreateBuffer([]opencl.MemFlags{opencl.MemWriteOnly}, dataSize*4)
	if err != nil {
		panic(err)
	}
	defer buffer.Release()

	err = kernel.SetArg(0, buffer.Size(), buffer)
	if err != nil {
		panic(err)
	}

	err = commandQueue.EnqueueNDRangeKernel(kernel, 1, []uint64{dataSize})
	if err != nil {
		panic(err)
	}

	commandQueue.Flush()
	commandQueue.Finish()

	data := make([]float32, dataSize)

	err = commandQueue.EnqueueReadBuffer(buffer, true, data)
	if err != nil {
		panic(err)
	}

	for _, item := range data {
		fmt.Printf("%v ", item)
	}
	fmt.Println()
	// Output:
	// 0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31 32
}