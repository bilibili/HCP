// Package server
/*
 * Copyright 2024-2025 Bilibili Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package server

import (
	"net/http"

	"github.com/bilibili/HCP/app/interface/v1/configs"
	"github.com/bilibili/HCP/app/interface/v1/internal/middleware"
	"github.com/bilibili/HCP/app/interface/v1/internal/service"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
)

// NewHttpServer new http server.
func NewHttpServer(svc *service.Service) (*bm.Engine, error) {
	engine := bm.DefaultServer(configs.Conf.HttpServer)
	var cors []string
	engine.Use(bm.CORS(cors))
	engine.Ping(ping)
	//pb.RegisterCloudDataSyncServiceBMServer(engine, svc)
	initRouter(engine)
	if err := engine.Start(); err != nil {
		panic(err)
	}
	return engine, nil
}

func ping(c *bm.Context) {
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write([]byte("ok"))
}

func initRouter(e *bm.Engine) {
	v1 := e.Group("/api/v1/cloud")
	{
		v1.POST("/login", service.Login)            //用户登录
		v1.POST("/logout", service.Logout)          //用户退出登录
		v1.POST("/sync", service.SyncCloudResource) //同步云资源

		v1db := v1.Group("/db")
		{
			v1db.POST("/query", middleware.CheckLoginMiddleware(func(c *bm.Context, w http.ResponseWriter, r *http.Request) {
				service.Query(c, w, r)
			})) //查询
			v1db.POST("/create", service.Create)                        //创建
			v1db.POST("/update", service.Update)                        //更新
			v1db.POST("/delete", service.Delete)                        //删除
			v1db.POST("/query/formTemplate", service.QueryFormTemplate) // 查询表单模板
		}

		v1job := v1.Group("/job")
		{
			v1job.POST("/create", middleware.CheckLoginMiddleware(func(c *bm.Context, w http.ResponseWriter, r *http.Request) {
				service.CreateJob(c, w, r)
			})) // 创建任务
			v1job.POST("/get", service.GetJob) // 获取任务
		}
		v1user := v1.Group("/user")
		{
			v1user.POST("/overview", middleware.CheckLoginMiddleware(func(c *bm.Context, w http.ResponseWriter, r *http.Request) {
				service.QueryUserOverview(c, w, r)
			})) // 用户概览
			v1user.POST("/getSign", service.GetLoginSign)    // 获取登录签名
			v1user.POST("/getUserInfo", service.GetUserInfo) // 获取用户信息
		}
	}
}
