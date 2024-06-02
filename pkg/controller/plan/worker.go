/*
Copyright 2021 The Pixiu Authors.

Licensed under the Apache License, Version 2.0 (phe "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package plan

import (
	"context"

	"k8s.io/klog/v2"
	"time"

	"github.com/caoyingjunz/pixiu/pkg/db/model"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/wait"
)

func (p *plan) Run(ctx context.Context, workers int) error {
	klog.Infof("Starting Plan Manager")
	for i := 0; i < workers; i++ {
		go wait.UntilWithContext(ctx, p.worker, time.Second)
	}
	return nil
}

func (p *plan) worker(ctx context.Context) {
	for p.process(ctx) {
	}
}

func (p *plan) process(ctx context.Context) bool {
	key, quit := taskQueue.Get()
	if quit {
		return false
	}
	defer taskQueue.Done(key)

	p.syncHandler(ctx, key.(int64))
	return true
}

type TaskData struct {
	PlanId int64
	Config *model.Config
	Nodes  []model.Node
}

func (t TaskData) validate() error {
	return nil
}

func (p *plan) getTaskData(ctx context.Context, planId int64) (TaskData, error) {
	nodes, err := p.factory.Plan().ListNodes(ctx, planId)
	if err != nil {
		return TaskData{}, err
	}
	cfg, err := p.factory.Plan().GetConfigByPlan(ctx, planId)
	if err != nil {
		return TaskData{}, err
	}

	return TaskData{
		PlanId: planId,
		Config: cfg,
		Nodes:  nodes,
	}, nil
}

// 实际处理函数
// 处理步骤:
// 1. 检查部署参数是否符合要求
// 2. 渲染环境
// 3. 执行部署
// 4. 部署后环境清理
func (p *plan) syncHandler(ctx context.Context, planId int64) {
	klog.Infof("starting plan(%d) task", planId)

	taskData, err := p.getTaskData(ctx, planId)
	if err != nil {
		klog.Errorf("failed to get task data: %v", err)
		return
	}

	handlers := []Handler{
		Check{data: taskData},
	}
	if err = p.syncTasks(handlers...); err != nil {
		klog.Errorf("failed to sync task: %v", err)
	}
}

type Handler interface {
	Name() string                     // 检查项名称
	Step() int                        // 未开始，运行中，异常和完成
	Run() (status string, msg string) // 执行
	GetPlanId() int64
}

type Check struct{ data TaskData }

func (c Check) Name() string {
	return "部署预检查"
}

func (c Check) Step() int {
	return 1
}

func (c Check) Run() (string, string) {
	if err := c.data.validate(); err != nil {
		return "失败", err.Error()
	}
	return "成功", ""
}

func (c Check) GetPlanId() int64 {
	return c.data.PlanId
}

func (p *plan) syncTasks(tasks ...Handler) error {
	for _, task := range tasks {
		planId := task.GetPlanId()
		name := task.Name()

		var (
			object *model.Task
			err    error
		)
		object, err = p.factory.Plan().GetTaskByName(context.TODO(), planId, name)
		if err != nil {
			if errors.IsNotFound(err) {
				object, err = p.factory.Plan().CreatTask(context.TODO(), &model.Task{
					Name:   name,
					PlanId: planId,
					Step:   model.PlanStep(task.Step()),
				})
				if err != nil {
					return err
				}
			}
			return err
		}

		// 执行检查
		status, message := task.Run()

		// 执行完成之后更新状态
		if err = p.factory.Plan().UpdateTask(context.TODO(), object.Id, object.ResourceVersion, map[string]interface{}{
			"status":  status,
			"message": message,
		}); err != nil {
			return err
		}
	}

	return nil
}
