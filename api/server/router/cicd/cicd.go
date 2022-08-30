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

package cicd

import (
	"github.com/gin-gonic/gin"
)

// cicdRouter is a router to talk with the cicd controller
type cicdRouter struct{}

// NewRouter initializes a new container router
func NewRouter(ginEngine *gin.Engine) {
	s := &cicdRouter{}
	s.initRoutes(ginEngine)
}

func (s *cicdRouter) initRoutes(ginEngine *gin.Engine) {
	cicdRoute := ginEngine.Group("/cicd")
	{
		cicdRoute.POST("/jobs/:name/run", s.runJob)
		cicdRoute.DELETE("/jobs/:name", s.deleteJob)
		cicdRoute.GET("/jobs", s.getAllJobs)
		cicdRoute.POST("/jobs", s.createJob)
		cicdRoute.POST("/jobs/copy", s.copyJob)
		cicdRoute.POST("/jobs/rename", s.renameJob)
		cicdRoute.POST("/view/:add_view_job/view_name", s.addViewJob)
		cicdRoute.POST("/safeRestart/:safeRestart", s.safeRestart)
	}
}
