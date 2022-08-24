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

package demo

import "github.com/gin-gonic/gin"

// demoRouter is a router to talk with the cicd controller
type demoRouter struct{}

// NewRouter initializes a new container router
func NewRouter(ginEngine *gin.Engine) {
	s := &demoRouter{}
	s.initRoutes(ginEngine)
}

func (s *demoRouter) initRoutes(ginEngine *gin.Engine) {
	demoRoute := ginEngine.Group("/demo")
	{
		demoRoute.POST("/create", s.createDemo)
		demoRoute.GET("/detail", s.getDemo)
	}
}
