package model

type Pod struct {
	Id int64 `gorm:"primary_key;not_null;auto_increment" json:"id"`
	//pod名称
	PodName string `gorm:"not_null;unique_index;" json:"pod_name"`

	NameSpace string `json:"name_space"`
	//所属团队
	PodTeamId int64 `json:"pod_team_id"`
	//最小CPU
	PodCpuMin float32 `json:"pod_cpu_min"`
	//最大CPU
	PodCpuMax float32 `json:"pod_cpu_max"`
	//副本数量
	PorReplicas int32 `json:"por_replicas"`
	//最小内存
	PodMemoryMin float32 `json:"pod_memory_min"`
	//最大内存
	PodMemoryMax float32 `json:"pod_memory_max"`
	//关联外键
	PodPort []PodPort `gorm:"ForeignKey:PodId" json:"pod_port"`
	//环境变量
	PodEnv []PodEnv `gorm:"ForeignKey:PodId" json:"pod_env"`
	//镜像拉取策略
	/**
	Always :总是拉取
	IfNotPresent:默认值，如果本地存在，则不拉取
	Never:只使用本地镜像，从不拉取
	*/
	PodPullPolicy string `json:"pod_pull_policy"`

	PodRestart string `json:"pod_restart"`
	//pod发布策略 金丝雀发布……
	PodType string `json:"pod_type"`
	//镜像名称
	PodImage string `json:"pod_image"`
	//TODO 挂盘
	// TODO 域名
}
