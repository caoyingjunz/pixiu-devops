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

package httputils

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/caoyingjunz/pixiu/api/server/errors"
	validatorutil "github.com/caoyingjunz/pixiu/api/server/validator"
	"github.com/caoyingjunz/pixiu/pkg/db/model"
)

type Response struct {
	Code    int         `json:"code"`              // 返回的状态码
	Result  interface{} `json:"result,omitempty"`  // 正常返回时的数据，可以为任意数据结构
	Message string      `json:"message,omitempty"` // 异常返回时的错误信息
}

func (r *Response) SetCode(c int) {
	r.Code = c
}

func (r *Response) SetMessage(m interface{}) {
	switch msg := m.(type) {
	case error:
		r.Message = msg.Error()
	case string:
		r.Message = msg
	}
}

func (r *Response) SetMessageWithCode(m interface{}, c int) {
	r.SetCode(c)
	r.SetMessage(m)
}

func (r *Response) Error() string {
	return r.Message
}

func (r *Response) String() string {
	//data, _ := json.Marshal(r)
	//return string(data)
	return ""
}

// NewResponse 构造 http 返回值
// SetSuccess 时设置 code 为 200 并追加 success 的标识
// SetFailed 时设置 code 为 400，也可以自定义设置错误码，并追加报错信息
func NewResponse() *Response {
	return &Response{}
}

// SetSuccess 设置成功返回值
func SetSuccess(c *gin.Context, r *Response) {
	r.SetMessageWithCode("success", http.StatusOK)
	c.JSON(http.StatusOK, r)
}

// SetFailed 设置错误返回值
func SetFailed(c *gin.Context, r *Response, err error) {
	switch e := err.(type) {
	case errors.Error:
		SetFailedWithCode(c, r, e.Code, e)
	case validator.ValidationErrors:
		SetFailedWithValidationError(c, r, validatorutil.TranslateError(e))
	default:
		SetFailedWithCode(c, r, http.StatusBadRequest, err)
	}
}

// SetFailedWithCode 设置错误返回值
func SetFailedWithCode(c *gin.Context, r *Response, code int, err error) {
	r.SetMessageWithCode(err, code)
	c.JSON(http.StatusOK, r)
}

func SetFailedWithValidationError(c *gin.Context, r *Response, e string) {
	r.SetMessageWithCode(e, http.StatusBadRequest)
	c.JSON(http.StatusOK, r)
}

// AbortFailedWithCode 设置错误，code 返回值并终止请求
func AbortFailedWithCode(c *gin.Context, code int, err error) {
	r := NewResponse()
	r.SetMessageWithCode(err, code)
	c.JSON(http.StatusOK, r)
	c.Abort()
}

func ShouldBindAny(c *gin.Context, jsonObject interface{}, uriObject interface{}, queryObject interface{}) error {
	var err error
	if jsonObject != nil {
		if err = c.ShouldBindJSON(jsonObject); err != nil {
			return err
		}
	}
	if uriObject != nil {
		if err = c.ShouldBindUri(uriObject); err != nil {
			return err
		}
	}
	if queryObject != nil {
		if err = c.ShouldBindQuery(queryObject); err != nil {
			return err
		}
	}
	return nil
}

const (
	userKey = "user"
	ip      = "ip"
)

func GetUserFromRequest(ctx context.Context) (*model.User, error) {
	val := ctx.Value(userKey)
	if val == nil {
		return nil, fmt.Errorf("get nil user")
	}

	user, ok := val.(*model.User)
	if !ok {
		return nil, fmt.Errorf("failed to assert user")
	}
	return user, nil
}

func GetUserIdFromRequest(ctx context.Context) (int64, error) {
	user, err := GetUserFromRequest(ctx)
	if err != nil {
		return 0, err
	}
	return user.Id, nil
}

func SetUserToContext(c *gin.Context, user *model.User) {
	c.Set(userKey, user)
}

func GetObjectFromRequest(c *gin.Context) (string, string, bool) {
	return getObjectFromRequest(c.Request.URL.Path)
}

// getObjectFromRequest cuts and returns the object from the request path.
// e.g. /pixiu/clusters/1 -> "clusters" "1" true
func getObjectFromRequest(path string) (obj, sid string, ok bool) {
	// must start with /
	l := len(path)
	if l == 0 || path[0] != '/' {
		return
	}
	subs := strings.Split(path[1:l], "/")
	l = len(subs)
	if l < 2 || subs[0] != "pixiu" {
		return
	}
	if l == 2 {
		// e.g. /pixiu/clusters -> "clusters" "" true
		return subs[1], "", subs[1] != ""
	}
	return subs[1], subs[2], subs[1] != "" && subs[2] != ""
}
