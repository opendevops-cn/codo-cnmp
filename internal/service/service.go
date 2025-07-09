package service

import "github.com/google/wire"

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewGreeterService, NewClusterService, NewNameSpaceService, NewNodeService, NewPodService,
	NewDeploymentService, NewEventService, NewWebsocketService, NewPodCommandWebsocketService, NewCloneSetService,
	NewGameServerSetService, NewUserGroupService, NewRoleService, NewStatefulSetService, NewHpaService, NewConfigMapService,
	NewSecretService, NewUserFollowService, NewAuditLogService, NewDaemonSetService, NewGameServerService,
	NewEzRolloutService, NewResourceService, NewLimitRangeService, NewResourceQuotaService, NewSvcService,
	NewIngressService, NewCRDService, NewApiGroupService, NewSideCarSetService, NewScService, NewPVService, NewPvcService,
	NewAgentService, NewCRRService, NewIngressClassService)
