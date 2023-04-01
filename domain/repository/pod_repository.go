package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/ljjgs/pod/domain/model"
)

type IPodRepository interface {
	//初始化表
	InitTable() error
	/**
	Curd
	*/
	FindPodById(int642 int64) (*model.Pod, error)

	CreatPod(pod *model.Pod) (int64, error)

	Delete(int64) error

	UpdatePod(pod *model.Pod) error

	FindAll() ([]*model.Pod, error)
}

//创建一个pod
func NewPodRepository(db *gorm.DB) IPodRepository {
	return &PodRepository{mysqlDb: db}
}

type PodRepository struct {
	mysqlDb *gorm.DB
}

func (p PodRepository) InitTable() error {
	return p.mysqlDb.CreateTable(&model.Pod{}, &model.PodPort{}, &model.PodEnv{}).Error
}

func (p PodRepository) FindPodById(podId int64) (pod *model.Pod, error error) {
	pod = &model.Pod{}
	return pod, p.mysqlDb.Preloads("PodEnv").Preloads("PodPort").First(pod, podId).Error
}

func (p PodRepository) CreatPod(pod *model.Pod) (int64, error) {
	return pod.ID, p.mysqlDb.Create(pod).Error
}

func (p PodRepository) Delete(i int64) error {
	pod := &model.Pod{
		ID: i,
	}
	return p.mysqlDb.Delete(pod).Error
}

func (p PodRepository) UpdatePod(pod *model.Pod) error {
	return p.mysqlDb.Update(pod).Error
}

func (p PodRepository) FindAll() (pods []*model.Pod, error error) {
	return pods, p.mysqlDb.Model(&model.Pod{}).Scan(&pods).Error
}
