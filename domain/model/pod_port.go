package model

type PodPort struct {
	Id            int64  `gorm:"auto_increment;not_null;primary_key" json:"id"`
	PodId         int64  `json:"pod_id"`
	ContainerPort int32  `json:"container_port"`
	Protocol      string `json:"protocol"`
	//	TODO host port 等其他的需要添加的
}
