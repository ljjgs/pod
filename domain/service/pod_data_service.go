package service

import (
	"context"
	"github.com/ljjgs/pod/domain/model"
	"github.com/ljjgs/pod/domain/repository"
	"github.com/ljjgs/pod/proto/pod"
	v1 "k8s.io/api/apps/v1"
	v13 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strconv"
)

type IPodDataService interface {
	AddPod(*model.Pod) (int64, error)
	Delete(int64) error
	Update(pod *model.Pod) error
	FindPodById(int64) (*model.Pod, error)
	FindAllPod() ([]*model.Pod, error)
	CreatToK8s(info *pod.PodInfo) error
	DeleteFromK8s(pod *model.Pod) error
	UpdateToK8s(info *pod.PodInfo) error
}
type PodDataService struct {
	PodRepository repository.IPodRepository
	K8sClientSet  *kubernetes.Clientset
	deployment    *v1.Deployment
}

func (p PodDataService) AddPod(m *model.Pod) (int64, error) {
	return p.PodRepository.CreatPod(m)
}

func (p PodDataService) Delete(i int64) error {
	return p.PodRepository.Delete(i)
}

func (p PodDataService) Update(pod *model.Pod) error {
	return p.PodRepository.UpdatePod(pod)
}

func (p PodDataService) FindPodById(i int64) (*model.Pod, error) {
	return p.PodRepository.FindPodById(i)
}

func (p PodDataService) FindAllPod() ([]*model.Pod, error) {
	return p.PodRepository.FindAll()
}

func (p PodDataService) CreatToK8s(pod *pod.PodInfo) (error error) {
	p.SetDeployment(pod)
	//	去k8s里面去请求看看namespace是否存在
	if _, error = p.K8sClientSet.AppsV1().Deployments(pod.PodNamespace).Get(
		context.TODO(), pod.PodName, v12.GetOptions{}); error != nil {

		if _, error = p.K8sClientSet.AppsV1().Deployments(pod.PodNamespace).Create(
			context.TODO(), p.deployment, v12.CreateOptions{}); error != nil {
			panic(error)
			return error
		}
	}
	return nil
}

func (p *PodDataService) DeleteFromK8s(pod *model.Pod) error {
	if err := p.K8sClientSet.AppsV1().Deployments(pod.NameSpace).Delete(context.TODO(), pod.PodName, v12.DeleteOptions{}); err != nil {
		return err
	} else {
		if err := p.Delete(pod.Id); err != nil {
			return err
		}
		return nil
	}

}

func (p PodDataService) UpdateToK8s(pod *pod.PodInfo) (error error) {
	p.SetDeployment(pod)
	if _, err := p.K8sClientSet.AppsV1().Deployments(pod.PodNamespace).Get(
		context.TODO(), pod.PodName, v12.GetOptions{}); err != nil {
		//			如果不存在
		return err
	} else {
		if _, error := p.K8sClientSet.AppsV1().Deployments(pod.PodNamespace).Update(
			context.TODO(), p.deployment, v12.UpdateOptions{}); error != nil {
			return error
		}
		return nil
	}
}

func NewPodDataService(repository repository.IPodRepository, clientSet *kubernetes.Clientset) IPodDataService {
	return &PodDataService{
		PodRepository: repository,
		K8sClientSet:  clientSet,
		deployment:    &v1.Deployment{},
	}
}

func (p *PodDataService) getContainerPort(podInfo *pod.PodInfo) (containerPort []v13.ContainerPort) {
	for _, v := range podInfo.PodPort {
		containerPort = append(containerPort, v13.ContainerPort{
			Name:          "port" + strconv.FormatInt(int64(v.ContainerPort), 10),
			ContainerPort: v.ContainerPort,
			Protocol:      p.getProtocol(v.Protocol),
		})
	}
	return containerPort
}

func (p *PodDataService) getProtocol(protocol string) v13.Protocol {
	switch protocol {
	case "TCP":
		return "TCP"
	case "UDP":
		return "UDP"
	case "SCTP":
		return "SCTP"
	default:
		return "TCP"
	}
}
func (p *PodDataService) getEnv(pod *pod.PodInfo) (env []v13.EnvVar) {
	for _, v := range pod.PodEnv {
		env = append(env, v13.EnvVar{
			Name:      v.EnvKey,
			Value:     v.EnvValue,
			ValueFrom: nil,
		})
	}
	return env
}

func (p *PodDataService) getResource(pod *pod.PodInfo) (source v13.ResourceRequirements) {

	source.Limits = v13.ResourceList{
		"cpu":    resource.MustParse(strconv.FormatFloat(float64(pod.PodCpuMax), 'f', 6, 64)),
		"memory": resource.MustParse(strconv.FormatFloat(float64(pod.PodMemoryMax), 'f', 6, 64)),
	}

	source.Requests = v13.ResourceList{
		"cpu":    resource.MustParse(strconv.FormatFloat(float64(pod.PodCpuMax), 'f', 6, 64)),
		"memory": resource.MustParse(strconv.FormatFloat(float64(pod.PodMemoryMax), 'f', 6, 64)),
	}
	return source
}

func (p *PodDataService) SetDeployment(pod *pod.PodInfo) {
	deployment := &v1.Deployment{}
	deployment.TypeMeta = v12.TypeMeta{
		Kind:       "deployment",
		APIVersion: "v1",
	}
	deployment.ObjectMeta = v12.ObjectMeta{
		Name:         pod.PodName,
		GenerateName: "",
		Namespace:    pod.PodNamespace,
		Labels: map[string]string{
			"app-name": pod.PodName,
		},
	}
	deployment.Name = pod.PodName
	deployment.Spec = v1.DeploymentSpec{
		Replicas: &pod.PodReplicas,
		Selector: &v12.LabelSelector{
			MatchLabels: map[string]string{
				"app-name": pod.PodName,
			},
			MatchExpressions: nil,
		},
		Template: v13.PodTemplateSpec{

			ObjectMeta: v12.ObjectMeta{
				Labels: map[string]string{
					"app-name": pod.PodName,
				},
			},
			Spec: v13.PodSpec{
				Containers: []v13.Container{
					{
						Name:            pod.PodName,
						Image:           pod.PodImages,
						Ports:           p.getContainerPort(pod),
						Env:             p.getEnv(pod),
						Resources:       p.getResource(pod),
						ImagePullPolicy: p.getImagePullPolicy(pod),
					},
				},
			},
		},
		Strategy:                v1.DeploymentStrategy{},
		MinReadySeconds:         0,
		RevisionHistoryLimit:    nil,
		Paused:                  false,
		ProgressDeadlineSeconds: nil,
	}

}

func (p PodDataService) getImagePullPolicy(pod *pod.PodInfo) v13.PullPolicy {
	switch pod.PodPullPolicy {
	case "Always":
		return "Always"
	case "Never":
		return "Never"
	case "IfNotPresent":
		return "IfNotPresent"
	default:
		return "Always"
	}

}
