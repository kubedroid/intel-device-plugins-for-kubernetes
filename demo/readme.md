# Intel Device Plugin Demo for Kubernetes

## Demo overview

Acceleration of compute and data processing of workloads like video
transcoding, compression, and encryption is enabled in Kubernetes with
the [Device Plugin Framework](https://kubernetes.io/docs/concepts/extend-kubernetes/compute-storage-net/device-plugins/). This repository
contains a set of Kubernetes plugins and instructions to enable Intel devices
for the acceleration of your workloads orchestrated by Kubernetes.

The current list of supported Intel Device Plugins includes:

- [GPU Device Plugin](https://github.com/intel/intel-device-plugins-for-kubernetes/blob/master/cmd/gpu_plugin/README.md) with support for [Intel® Graphics Technology](https://www.intel.com/content/www/us/en/architecture-and-technology/visual-technology/graphics-overview.html)

- [Intel® FPGA Device Plugin](https://github.com/intel/intel-device-plugins-for-kubernetes/blob/master/cmd/fpga_plugin/README.md)

- [Intel® QuickAssist Technology (QAT) Device Plugin](https://github.com/intel/intel-device-plugins-for-kubernetes/blob/master/cmd/qat_plugin/README.md)

We have included an example demo and configuration information for the Intel
Device Plugins for Kubernetes below. Please join us on the sig-node-rtk channel
on [kubernetes.slack.com](https://kubernetes.slack.com/) to ask questions,
contribute to the community, and learn about the work we are doing with
Kubernetes and the Device Plugin Framework.

## Intel® FPGA Device Plugin demo video

The screencast demonstrates the deployment of the Intel® FPGA Device Plugin for
Kubernetes and executes a sample GZIP compression workload. The demo begins
with a fully [configured Kubernetes cluster](https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/)
with the [Go runtime](https://golang.org/doc/install).

### Demo platform configuration

- 1-node, 2x Intel® Xeon® E5-2680 CPU @ 2.4 GHz
- Total memory 251 GB
- Intel® Arria® 10 GX FPGA
- Software
- OpenSUSE 15 (Kernel: 4.18.9-200.fc28.x86_64)
- Kubernetes 1.11
- CRI-O 1.11.1
- Intel® FPGA Device Plugin built from master branch

### Demo steps

1. Validate the status of the [Kubernetes cluster](https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/).
2. Clone the [Intel Device Plugins for Kubernetes source](https://github.com/intel/intel-device-plugins-for-kubernetes).
3. Provision the [admission controller webhook](https://github.com/intel/intel-device-plugins-for-kubernetes/blob/master/cmd/fpga_admissionwebhook/README.md).
4. Provision the [Intel® FPGA Device Plugin](https://github.com/intel/intel-device-plugins-for-kubernetes/blob/master/cmd/fpga_plugin/README.md).
5. Create bitstream storage for the Intel® FPGA.
6. Run the sample GZIP compression workload.

### Screencast

Intel® FPGA Device Plugin deployment

[![Intel® FPGA Device Plugin deployment](https://asciinema.org/a/mRn15bkRRUzTG4kp2UeHQX6gk.png)](https://asciinema.org/a/mRn15bkRRUzTG4kp2UeHQX6gk)

## Intel® QuickAssist Technology Device Plugin OpenSSL demo video

The screencast demonstrates the deployment of the [Intel® QAT Device Plugin](https://github.com/intel/intel-device-plugins-for-kubernetes/blob/master/cmd/qat_plugin/README.md) for
Kubernetes and executes a sample QAT accelerated OpenSSL workload with the OCI
compatible Kata Containers runtime, a lightweight Virtual Machine (VMs) that feel
and perform like traditional containers, but provides the workload isolation
and security advantages of VMs. The demo begins with a fully [configured Kubernetes cluster](https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/)
and [Kata Containers](https://github.com/kata-containers/documentation/tree/master/install) runtime for workloads.

### Demo platform configuration

- 1-node, 2x Intel® Xeon® CPU E5-2687W v4 @ 3.00GHz
- Total memory 251 GB DDR4
- Intel® QAT C62x chipset
- Software
- OpenSUSE 15 (Kernel:4.12.14-lp150.12.22-default)
- Kubernetes 1.12
- Containerd 1.2
- Kata Containers 1.3.0
- Intel® QAT Device Plugin built from master
- QAT 1.7 L.4.3.0-00033

### Demo steps

1. Load the host drivers and prepare the virtual function (VF) devices
2. Check the Kubernetes cluster is in good shape
3. Deploy the Intel® QAT device plugin for Kubernetes
4. Deploy an Intel® QAT Accelerated OpenSSL workload
5. Testing!

### Screencast
Intel® QAT Device Plugin deployment

[![Intel® QAT Device Plugin deployment](https://asciinema.org/a/2N7wF3c9oeCuB9sFqTNm2gmOr.png)](https://asciinema.org/a/2N7wF3c9oeCuB9sFqTNm2gmOr)
