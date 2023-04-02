package model

type PodPort struct {
	ID int64 `gorm:"primary_key;not_null;auto_increment" json:"id"`
	PodID int64 `json:"pod_id"`
	ContainerPort int32 `json:"container_port"`
	Protocol string `json:"protocol"`
}
