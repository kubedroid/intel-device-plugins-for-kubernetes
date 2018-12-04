// Copyright 2017 Intel Corporation. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
	"path/filepath"

	"github.com/pkg/errors"

	pluginapi "k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1beta1"

	dpapi "github.com/intel/intel-device-plugins-for-kubernetes/internal/deviceplugin"
	"github.com/intel/intel-device-plugins-for-kubernetes/pkg/debug"
)

const (
	sysfsDrmDirectory = "/sys/class/drm"
	devfsDriDirectory = "/dev/dri"
	gpuDeviceRE       = `^card[0-9]+$`
	vgpuDeviceRE      = `^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`
	controlDeviceRE   = `^controlD[0-9]+$`
	vendorString      = "0x8086"

	// Device plugin settings.
	namespace  = "gpu.intel.com"
	deviceType = "i915"
	vgpuDeviceType = "i915-v"
)

type devicePlugin struct {
	sysfsDir string
	devfsDir string

	sharedDevNum int

	gpuDeviceReg     *regexp.Regexp
	vgpuDeviceReg     *regexp.Regexp
	controlDeviceReg *regexp.Regexp
}

func newDevicePlugin(sysfsDir, devfsDir string, sharedDevNum int) *devicePlugin {
	return &devicePlugin{
		sysfsDir:         sysfsDir,
		devfsDir:         devfsDir,
		sharedDevNum:     sharedDevNum,
		gpuDeviceReg:     regexp.MustCompile(gpuDeviceRE),
		vgpuDeviceReg:    regexp.MustCompile(vgpuDeviceRE),
		controlDeviceReg: regexp.MustCompile(controlDeviceRE),
	}
}

func (dp *devicePlugin) Scan(notifier dpapi.Notifier) error {
	for {
		devTree, err := dp.scan()
		if err != nil {
			return err
		}

		notifier.Notify(devTree)

		time.Sleep(5 * time.Second)
	}
}

func (dp *devicePlugin) scan() (dpapi.DeviceTree, error) {
	files, err := ioutil.ReadDir(dp.sysfsDir)
	if err != nil {
		return nil, errors.Wrap(err, "Can't read sysfs folder")
	}

	devTree := dpapi.NewDeviceTree()
	for _, f := range files {
		if dp.gpuDeviceReg.MatchString(f.Name()) {
			dat, err := ioutil.ReadFile(path.Join(dp.sysfsDir, f.Name(), "device/vendor"))
			if err != nil {
				fmt.Println("WARNING: Skipping. Can't read vendor file: ", err)
				continue
			}

			if strings.TrimSpace(string(dat)) == vendorString {
				var nodes []string

				drmFiles, err := ioutil.ReadDir(path.Join(dp.sysfsDir, f.Name(), "device/drm"))
				if err != nil {
					return nil, errors.Wrap(err, "Can't read device/drm folder")
				}

				for _, drmFile := range drmFiles {
					fmt.Println("Evaluating possible GPU device '", drmFile.Name(), "'")
					if dp.controlDeviceReg.MatchString(drmFile.Name()) {
						//Skipping possible drm control node
						continue
					}

					devPath := path.Join(dp.devfsDir, drmFile.Name())
					fmt.Println("Evaluating possible GPU device '", devPath, "'")
					if _, err := os.Stat(devPath); err != nil {
						fmt.Println("Failed to evaluate GPU device ", devPath, "'. ", err)
						continue
					}

					fmt.Println("Found GPU device '", devPath, "'")
					debug.Printf("Adding %s to GPU %s", devPath, f.Name())
					nodes = append(nodes, devPath)
				}

				if len(nodes) > 0 {
					for i := 0; i < dp.sharedDevNum; i++ {
						devID := fmt.Sprintf("%s-%d", f.Name(), i)
						// Currently only one device type (i915) is supported.
						// TODO: check model ID to differentiate device models.
						devTree.AddDevice(deviceType, devID, dpapi.DeviceInfo{
							State: pluginapi.Healthy,
							Nodes: nodes,
						})
					}
				}

				devFiles, err := ioutil.ReadDir(path.Join(dp.sysfsDir, f.Name(), "device"))
				if err != nil {
					return nil, errors.Wrap(err, "Can't read device folder")
				}

				for _, devFile := range devFiles {
					if dp.vgpuDeviceReg.MatchString(devFile.Name()) {
						fmt.Println("Found vGPU device '", devFile.Name(), "'")

						iommuPath, err := filepath.EvalSymlinks(path.Join(dp.sysfsDir, f.Name(), "device", devFile.Name(), "iommu_group"))
						if err != nil {
							return nil, errors.Wrap(err, "Can't read iommu_group")
						}

						iommu, err := os.Stat(iommuPath)
						if err != nil {
							return nil, errors.Wrap(err, "Can't stat iommu_group")
						}

						fmt.Println("The vGPU '", devFile.Name(), "' belongs to iommu group '", iommu.Name(), "'")

						var vgpuNodes []string
						vgpuNodes = append(vgpuNodes, iommuPath)
						vgpuNodes = append(vgpuNodes, "/dev/vfio/vfio")

						devTree.AddDevice(vgpuDeviceType, devFile.Name(), dpapi.DeviceInfo{
							State: pluginapi.Healthy,
							Nodes: vgpuNodes,
						})
					}
				}
			}
		}
	}

	return devTree, nil
}

func main() {
	var sharedDevNum int
	var debugEnabled bool

	flag.IntVar(&sharedDevNum, "shared-dev-num", 1, "number of containers sharing the same GPU device")
	flag.BoolVar(&debugEnabled, "debug", false, "enable debug output")
	flag.Parse()

	if debugEnabled {
		debug.Activate()
	}

	if sharedDevNum < 1 {
		fmt.Println("The number of containers sharing the same GPU must greater than zero")
		os.Exit(1)
	}

	fmt.Println("GPU device plugin started")

	plugin := newDevicePlugin(sysfsDrmDirectory, devfsDriDirectory, sharedDevNum)
	manager := dpapi.NewManager(namespace, plugin)
	manager.Run()
}
