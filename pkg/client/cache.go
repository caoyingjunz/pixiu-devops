/*
Copyright 2021 The Pixiu Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package client

import (
	"sync"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	resourceclient "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
)

type ClusterSet struct {
	Client                *kubernetes.Clientset
	Config                *restclient.Config
	Metric                *resourceclient.MetricsV1beta1Client
	SharedInformerFactory informers.SharedInformerFactory
	StopChan              chan struct{}
	SyncCache             bool //是否开启了集群缓存
}

func (cs *ClusterSet) Complete(cfg []byte) error {
	var err error
	if cs.Config, err = clientcmd.RESTConfigFromKubeConfig(cfg); err != nil {
		return err
	}
	if cs.Client, err = kubernetes.NewForConfig(cs.Config); err != nil {
		return err
	}
	if cs.Metric, err = resourceclient.NewForConfig(cs.Config); err != nil {
		return err
	}

	return nil
}

type store map[string]*ClusterSet

type Cache struct {
	sync.RWMutex
	store
}

func NewClusterCache() *Cache {
	return &Cache{
		store: make(store),
	}
}

func (s *Cache) Get(name string) (*ClusterSet, bool) {
	s.RLock()
	defer s.RUnlock()

	cluster, ok := s.store[name]
	return cluster, ok
}

func (s *Cache) GetConfig(name string) (*restclient.Config, bool) {
	s.RLock()
	defer s.RUnlock()

	clusterSet, ok := s.store[name]
	if !ok {
		return nil, false
	}
	return clusterSet.Config, true
}

func (s *Cache) GetClient(name string) (*kubernetes.Clientset, bool) {
	s.RLock()
	defer s.RUnlock()

	clusterSet, ok := s.store[name]
	if !ok {
		return nil, false
	}

	return clusterSet.Client, true
}

func (s *Cache) Set(name string, cs *ClusterSet) {
	s.Lock()
	defer s.Unlock()

	if s.store == nil {
		s.store = store{}
	}
	s.store[name] = cs
}

func (s *Cache) Delete(name string) {
	s.Lock()
	defer s.Unlock()

	delete(s.store, name)
}

func (s *Cache) List() store {
	s.Lock()
	defer s.Unlock()

	return s.store
}

func (s *Cache) Clear() {
	s.Lock()
	defer s.Unlock()

	s.store = store{}
}
