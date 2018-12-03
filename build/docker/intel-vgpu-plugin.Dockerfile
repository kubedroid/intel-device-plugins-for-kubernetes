FROM golang:1.11-alpine as builder
ARG DIR=/go/src/github.com/intel/intel-device-plugins-for-kubernetes
WORKDIR $DIR
COPY . .
RUN cd cmd/vgpu_plugin; go install
RUN chmod a+x /go/bin/vgpu_plugin

FROM alpine
COPY --from=builder /go/bin/vgpu_plugin /usr/bin/intel_vgpu_device_plugin
CMD ["/usr/bin/intel_vgpu_device_plugin"]
