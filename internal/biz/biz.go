package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewGreeterUsecase, NewClusterUseCase, NewIClusterUseCase, NewNodeUseCase, NewINodeUseCase,
	NewOverViewUseCase, NewIOverViewUseCase, NewNameSpaceUseCase, NewINameSpaceUseCase, NewPodUseCase, NewIPodUseCase, NewDeploymentUseCase,
	NewIDeploymentUseCase, NewEventUseCase, NewIEventUseCase, NewCloneSetUseCase, NewICloneSetUseCase, NewGameServerSetUseCase,
	NewIGameServerSetUseCase, NewUserGroupUseCase, NewIGrantedUserGroupRepo, NewIRoleUseCase, NewRoleUseCase, NewIRoleBindingUseCase,
	NewRoleBindingUseCase, NewStatefulSetUseCase, NewIStatefulSetUseCase, NewIUserUseCase, NewUserUseCase,
	NewIUserGroupV2UseCase, NewUserGroupV2UseCase, NewHpaUseCase, NewIHpaUseCase, NewConfigMapUseCase, NewIConfigMapUseCase,
	NewSecretUseCase, NewISecretUseCase, NewUserFollowUseCase, NewIUserFollowUseCase, NewIAuditLogUseCase, NewAuditLogUseCase,
	NewDaemonSetUseCase, NewIDaemonSetUseCase, NewIGameServerUseCase, NewGameServerUseCase, NewIEzRolloutUseCase,
	NewEzRolloutUseCase, NewIResourceUseCase, NewResourceUseCase, NewILimitRangeUseCase, NewLimitRangeUseCase,
	NewResourceQuotaUseCase, NewIResourceQuotaUseCase, NewISvcUseCase, NewSvcUseCase, NewIIngressUseCase,
	NewIngressUseCase, NewICRDUseCase, NewCRDUseCase, NewApiGroupUseCase, NewIApiGroupUseCase,
	NewSideCarSetUseCase, NewISideCarSetUseCase, NewIStorageClassUseCase, NewStorageClassUseCase, NewIPersistentVolumeUseCase,
	NewPersistentVolumeUseCase, NewIPersistentVolumeClaimUseCase, NewPersistentVolumeClaimUseCase, NewIMeshUseCase, NewMeshUseCase,
	NewIAgentUseCase, NewAgentUseCase, NewCRRUseCase, NewICRRUseCase, NewIIngressClassUseCase, NewIngressClassUseCase)
