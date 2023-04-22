package scheduler

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"server/config"
	"server/internal/repo"
	"sync"
	"time"
)

type InitSchedulerDep struct {
	Conf               *config.Config
	Logger             logrus.FieldLogger
	RepositoryRegistry *repo.Registry
}

type Scheduler struct {
	Status             status // Статус задач планировщика
	conf               *config.Config
	logger             logrus.FieldLogger
	repositoryRegistry *repo.Registry
}

type status struct {
	WorkerSyncATrucksData workerSyncATrucksData
}

type base struct {
	sync.Mutex
}

func (o *base) MarshalJSON() ([]byte, error) {
	o.Lock()
	defer o.Unlock()
	return json.Marshal(o)
}

type workerSyncATrucksData struct {
	base
	Run  bool      // Состояние метода [true-запущен]
	Last time.Time // Время последнего запуска
}

func NewScheduler(dep InitSchedulerDep) (*Scheduler, error) {

	if dep.Conf == nil {
		return nil, fmt.Errorf("conf")
	}
	if dep.Logger == nil {
		return nil, fmt.Errorf("logger")
	}

	if dep.RepositoryRegistry == nil {
		return nil, fmt.Errorf("repositoryRegistry")
	}

	logger := dep.Logger.WithField("process", "Scheduler")

	return &Scheduler{
		Status:             status{},
		conf:               dep.Conf,
		logger:             logger,
		repositoryRegistry: dep.RepositoryRegistry,
	}, nil
}

func (o *Scheduler) Start() {
	// перед употреблением раскоментить
	//logger := o.logger
	//conf := o.conf.BackGroundWorkers
	//
	//var wg sync.WaitGroup
	//wg.Add(1)

	// Синхронизация данных

}
