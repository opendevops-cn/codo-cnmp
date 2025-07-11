// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.27.2
// source: pb/pod.v1.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Pod_ListPod_FullMethodName                 = "/pod.Pod/ListPod"
	Pod_ListPodByNamespace_FullMethodName      = "/pod.Pod/ListPodByNamespace"
	Pod_GetNamespacePodDetail_FullMethodName   = "/pod.Pod/GetNamespacePodDetail"
	Pod_DeletePod_FullMethodName               = "/pod.Pod/DeletePod"
	Pod_BatchDeletePod_FullMethodName          = "/pod.Pod/BatchDeletePod"
	Pod_GetPodCpuMetrics_FullMethodName        = "/pod.Pod/GetPodCpuMetrics"
	Pod_GetPodMemoryMetrics_FullMethodName     = "/pod.Pod/GetPodMemoryMetrics"
	Pod_GetPodContainerMetrics_FullMethodName  = "/pod.Pod/GetPodContainerMetrics"
	Pod_ListControllerPod_FullMethodName       = "/pod.Pod/ListControllerPod"
	Pod_DownloadPodLogs_FullMethodName         = "/pod.Pod/DownloadPodLogs"
	Pod_CreateOrUpdatePodByYaml_FullMethodName = "/pod.Pod/CreateOrUpdatePodByYaml"
	Pod_EvictPod_FullMethodName                = "/pod.Pod/EvictPod"
)

// PodClient is the client API for Pod service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PodClient interface {
	// 查看-云原生管理-Pod-列表
	ListPod(ctx context.Context, in *ListPodRequest, opts ...grpc.CallOption) (*ListPodResponse, error)
	// 查看-云原生管理-命名空间-Pod列表
	ListPodByNamespace(ctx context.Context, in *ListPodByNamespaceRequest, opts ...grpc.CallOption) (*ListPodByNamespaceResponse, error)
	// 查看-云原生管理-命名空间-Pod详情
	GetNamespacePodDetail(ctx context.Context, in *GetNamespacePodDetailRequest, opts ...grpc.CallOption) (*GetNamespacePodDetailResponse, error)
	// 管理-云原生管理-Pod-删除
	DeletePod(ctx context.Context, in *DeletePodRequest, opts ...grpc.CallOption) (*DeletePodResponse, error)
	// 管理-云原生管理-Pod-批量重启
	BatchDeletePod(ctx context.Context, in *BatchDeletePodsRequest, opts ...grpc.CallOption) (*BatchDeletePodsResponse, error)
	// 查看-云原生管理-Pod-CPU指标
	GetPodCpuMetrics(ctx context.Context, in *GetPodMetricsRequest, opts ...grpc.CallOption) (*SidecarMetricResultList, error)
	// 查看-云原生管理-Pod-内存指标
	GetPodMemoryMetrics(ctx context.Context, in *GetPodMetricsRequest, opts ...grpc.CallOption) (*SidecarMetricResultList, error)
	// 查看-云原生管理-Pod-容器指标
	GetPodContainerMetrics(ctx context.Context, in *GetPodContainerMetricsRequest, opts ...grpc.CallOption) (*GetPodContainerMetricsResponse, error)
	// 查看-云原生管理-控制器-Pod列表
	ListControllerPod(ctx context.Context, in *ListControllerPodRequest, opts ...grpc.CallOption) (*ListControllerPodResponse, error)
	// 查看-云原生管理-Pod-下载日志
	DownloadPodLogs(ctx context.Context, in *DownloadPodLogsRequest, opts ...grpc.CallOption) (*DownloadPodLogsResponse, error)
	// 管理-云原生管理-Pod-Yaml创建更新
	CreateOrUpdatePodByYaml(ctx context.Context, in *CreateOrUpdatePodByYamlRequest, opts ...grpc.CallOption) (*CreateOrUpdatePodByYamlResponse, error)
	// 管理-云原生管理-Pod-驱逐
	EvictPod(ctx context.Context, in *EvictPodRequest, opts ...grpc.CallOption) (*EvictPodResponse, error)
}

type podClient struct {
	cc grpc.ClientConnInterface
}

func NewPodClient(cc grpc.ClientConnInterface) PodClient {
	return &podClient{cc}
}

func (c *podClient) ListPod(ctx context.Context, in *ListPodRequest, opts ...grpc.CallOption) (*ListPodResponse, error) {
	out := new(ListPodResponse)
	err := c.cc.Invoke(ctx, Pod_ListPod_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *podClient) ListPodByNamespace(ctx context.Context, in *ListPodByNamespaceRequest, opts ...grpc.CallOption) (*ListPodByNamespaceResponse, error) {
	out := new(ListPodByNamespaceResponse)
	err := c.cc.Invoke(ctx, Pod_ListPodByNamespace_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *podClient) GetNamespacePodDetail(ctx context.Context, in *GetNamespacePodDetailRequest, opts ...grpc.CallOption) (*GetNamespacePodDetailResponse, error) {
	out := new(GetNamespacePodDetailResponse)
	err := c.cc.Invoke(ctx, Pod_GetNamespacePodDetail_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *podClient) DeletePod(ctx context.Context, in *DeletePodRequest, opts ...grpc.CallOption) (*DeletePodResponse, error) {
	out := new(DeletePodResponse)
	err := c.cc.Invoke(ctx, Pod_DeletePod_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *podClient) BatchDeletePod(ctx context.Context, in *BatchDeletePodsRequest, opts ...grpc.CallOption) (*BatchDeletePodsResponse, error) {
	out := new(BatchDeletePodsResponse)
	err := c.cc.Invoke(ctx, Pod_BatchDeletePod_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *podClient) GetPodCpuMetrics(ctx context.Context, in *GetPodMetricsRequest, opts ...grpc.CallOption) (*SidecarMetricResultList, error) {
	out := new(SidecarMetricResultList)
	err := c.cc.Invoke(ctx, Pod_GetPodCpuMetrics_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *podClient) GetPodMemoryMetrics(ctx context.Context, in *GetPodMetricsRequest, opts ...grpc.CallOption) (*SidecarMetricResultList, error) {
	out := new(SidecarMetricResultList)
	err := c.cc.Invoke(ctx, Pod_GetPodMemoryMetrics_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *podClient) GetPodContainerMetrics(ctx context.Context, in *GetPodContainerMetricsRequest, opts ...grpc.CallOption) (*GetPodContainerMetricsResponse, error) {
	out := new(GetPodContainerMetricsResponse)
	err := c.cc.Invoke(ctx, Pod_GetPodContainerMetrics_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *podClient) ListControllerPod(ctx context.Context, in *ListControllerPodRequest, opts ...grpc.CallOption) (*ListControllerPodResponse, error) {
	out := new(ListControllerPodResponse)
	err := c.cc.Invoke(ctx, Pod_ListControllerPod_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *podClient) DownloadPodLogs(ctx context.Context, in *DownloadPodLogsRequest, opts ...grpc.CallOption) (*DownloadPodLogsResponse, error) {
	out := new(DownloadPodLogsResponse)
	err := c.cc.Invoke(ctx, Pod_DownloadPodLogs_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *podClient) CreateOrUpdatePodByYaml(ctx context.Context, in *CreateOrUpdatePodByYamlRequest, opts ...grpc.CallOption) (*CreateOrUpdatePodByYamlResponse, error) {
	out := new(CreateOrUpdatePodByYamlResponse)
	err := c.cc.Invoke(ctx, Pod_CreateOrUpdatePodByYaml_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *podClient) EvictPod(ctx context.Context, in *EvictPodRequest, opts ...grpc.CallOption) (*EvictPodResponse, error) {
	out := new(EvictPodResponse)
	err := c.cc.Invoke(ctx, Pod_EvictPod_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PodServer is the server API for Pod service.
// All implementations must embed UnimplementedPodServer
// for forward compatibility
type PodServer interface {
	// 查看-云原生管理-Pod-列表
	ListPod(context.Context, *ListPodRequest) (*ListPodResponse, error)
	// 查看-云原生管理-命名空间-Pod列表
	ListPodByNamespace(context.Context, *ListPodByNamespaceRequest) (*ListPodByNamespaceResponse, error)
	// 查看-云原生管理-命名空间-Pod详情
	GetNamespacePodDetail(context.Context, *GetNamespacePodDetailRequest) (*GetNamespacePodDetailResponse, error)
	// 管理-云原生管理-Pod-删除
	DeletePod(context.Context, *DeletePodRequest) (*DeletePodResponse, error)
	// 管理-云原生管理-Pod-批量重启
	BatchDeletePod(context.Context, *BatchDeletePodsRequest) (*BatchDeletePodsResponse, error)
	// 查看-云原生管理-Pod-CPU指标
	GetPodCpuMetrics(context.Context, *GetPodMetricsRequest) (*SidecarMetricResultList, error)
	// 查看-云原生管理-Pod-内存指标
	GetPodMemoryMetrics(context.Context, *GetPodMetricsRequest) (*SidecarMetricResultList, error)
	// 查看-云原生管理-Pod-容器指标
	GetPodContainerMetrics(context.Context, *GetPodContainerMetricsRequest) (*GetPodContainerMetricsResponse, error)
	// 查看-云原生管理-控制器-Pod列表
	ListControllerPod(context.Context, *ListControllerPodRequest) (*ListControllerPodResponse, error)
	// 查看-云原生管理-Pod-下载日志
	DownloadPodLogs(context.Context, *DownloadPodLogsRequest) (*DownloadPodLogsResponse, error)
	// 管理-云原生管理-Pod-Yaml创建更新
	CreateOrUpdatePodByYaml(context.Context, *CreateOrUpdatePodByYamlRequest) (*CreateOrUpdatePodByYamlResponse, error)
	// 管理-云原生管理-Pod-驱逐
	EvictPod(context.Context, *EvictPodRequest) (*EvictPodResponse, error)
	mustEmbedUnimplementedPodServer()
}

// UnimplementedPodServer must be embedded to have forward compatible implementations.
type UnimplementedPodServer struct {
}

func (UnimplementedPodServer) ListPod(context.Context, *ListPodRequest) (*ListPodResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPod not implemented")
}
func (UnimplementedPodServer) ListPodByNamespace(context.Context, *ListPodByNamespaceRequest) (*ListPodByNamespaceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPodByNamespace not implemented")
}
func (UnimplementedPodServer) GetNamespacePodDetail(context.Context, *GetNamespacePodDetailRequest) (*GetNamespacePodDetailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNamespacePodDetail not implemented")
}
func (UnimplementedPodServer) DeletePod(context.Context, *DeletePodRequest) (*DeletePodResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePod not implemented")
}
func (UnimplementedPodServer) BatchDeletePod(context.Context, *BatchDeletePodsRequest) (*BatchDeletePodsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchDeletePod not implemented")
}
func (UnimplementedPodServer) GetPodCpuMetrics(context.Context, *GetPodMetricsRequest) (*SidecarMetricResultList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPodCpuMetrics not implemented")
}
func (UnimplementedPodServer) GetPodMemoryMetrics(context.Context, *GetPodMetricsRequest) (*SidecarMetricResultList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPodMemoryMetrics not implemented")
}
func (UnimplementedPodServer) GetPodContainerMetrics(context.Context, *GetPodContainerMetricsRequest) (*GetPodContainerMetricsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPodContainerMetrics not implemented")
}
func (UnimplementedPodServer) ListControllerPod(context.Context, *ListControllerPodRequest) (*ListControllerPodResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListControllerPod not implemented")
}
func (UnimplementedPodServer) DownloadPodLogs(context.Context, *DownloadPodLogsRequest) (*DownloadPodLogsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DownloadPodLogs not implemented")
}
func (UnimplementedPodServer) CreateOrUpdatePodByYaml(context.Context, *CreateOrUpdatePodByYamlRequest) (*CreateOrUpdatePodByYamlResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrUpdatePodByYaml not implemented")
}
func (UnimplementedPodServer) EvictPod(context.Context, *EvictPodRequest) (*EvictPodResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EvictPod not implemented")
}
func (UnimplementedPodServer) mustEmbedUnimplementedPodServer() {}

// UnsafePodServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PodServer will
// result in compilation errors.
type UnsafePodServer interface {
	mustEmbedUnimplementedPodServer()
}

func RegisterPodServer(s grpc.ServiceRegistrar, srv PodServer) {
	s.RegisterService(&Pod_ServiceDesc, srv)
}

func _Pod_ListPod_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListPodRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PodServer).ListPod(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Pod_ListPod_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PodServer).ListPod(ctx, req.(*ListPodRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pod_ListPodByNamespace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListPodByNamespaceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PodServer).ListPodByNamespace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Pod_ListPodByNamespace_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PodServer).ListPodByNamespace(ctx, req.(*ListPodByNamespaceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pod_GetNamespacePodDetail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetNamespacePodDetailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PodServer).GetNamespacePodDetail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Pod_GetNamespacePodDetail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PodServer).GetNamespacePodDetail(ctx, req.(*GetNamespacePodDetailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pod_DeletePod_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePodRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PodServer).DeletePod(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Pod_DeletePod_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PodServer).DeletePod(ctx, req.(*DeletePodRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pod_BatchDeletePod_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchDeletePodsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PodServer).BatchDeletePod(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Pod_BatchDeletePod_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PodServer).BatchDeletePod(ctx, req.(*BatchDeletePodsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pod_GetPodCpuMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPodMetricsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PodServer).GetPodCpuMetrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Pod_GetPodCpuMetrics_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PodServer).GetPodCpuMetrics(ctx, req.(*GetPodMetricsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pod_GetPodMemoryMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPodMetricsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PodServer).GetPodMemoryMetrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Pod_GetPodMemoryMetrics_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PodServer).GetPodMemoryMetrics(ctx, req.(*GetPodMetricsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pod_GetPodContainerMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPodContainerMetricsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PodServer).GetPodContainerMetrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Pod_GetPodContainerMetrics_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PodServer).GetPodContainerMetrics(ctx, req.(*GetPodContainerMetricsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pod_ListControllerPod_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListControllerPodRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PodServer).ListControllerPod(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Pod_ListControllerPod_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PodServer).ListControllerPod(ctx, req.(*ListControllerPodRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pod_DownloadPodLogs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DownloadPodLogsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PodServer).DownloadPodLogs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Pod_DownloadPodLogs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PodServer).DownloadPodLogs(ctx, req.(*DownloadPodLogsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pod_CreateOrUpdatePodByYaml_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateOrUpdatePodByYamlRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PodServer).CreateOrUpdatePodByYaml(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Pod_CreateOrUpdatePodByYaml_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PodServer).CreateOrUpdatePodByYaml(ctx, req.(*CreateOrUpdatePodByYamlRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pod_EvictPod_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EvictPodRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PodServer).EvictPod(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Pod_EvictPod_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PodServer).EvictPod(ctx, req.(*EvictPodRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Pod_ServiceDesc is the grpc.ServiceDesc for Pod service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Pod_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pod.Pod",
	HandlerType: (*PodServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListPod",
			Handler:    _Pod_ListPod_Handler,
		},
		{
			MethodName: "ListPodByNamespace",
			Handler:    _Pod_ListPodByNamespace_Handler,
		},
		{
			MethodName: "GetNamespacePodDetail",
			Handler:    _Pod_GetNamespacePodDetail_Handler,
		},
		{
			MethodName: "DeletePod",
			Handler:    _Pod_DeletePod_Handler,
		},
		{
			MethodName: "BatchDeletePod",
			Handler:    _Pod_BatchDeletePod_Handler,
		},
		{
			MethodName: "GetPodCpuMetrics",
			Handler:    _Pod_GetPodCpuMetrics_Handler,
		},
		{
			MethodName: "GetPodMemoryMetrics",
			Handler:    _Pod_GetPodMemoryMetrics_Handler,
		},
		{
			MethodName: "GetPodContainerMetrics",
			Handler:    _Pod_GetPodContainerMetrics_Handler,
		},
		{
			MethodName: "ListControllerPod",
			Handler:    _Pod_ListControllerPod_Handler,
		},
		{
			MethodName: "DownloadPodLogs",
			Handler:    _Pod_DownloadPodLogs_Handler,
		},
		{
			MethodName: "CreateOrUpdatePodByYaml",
			Handler:    _Pod_CreateOrUpdatePodByYaml_Handler,
		},
		{
			MethodName: "EvictPod",
			Handler:    _Pod_EvictPod_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/pod.v1.proto",
}
