package router

import (
	"outdoor-app-backend/internal/handler"
	"outdoor-app-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Static("/uploads", "./uploads") // 静态资源

	apiV1 := r.Group("/api/v1")
	{
		// ================== 公共接口 (无需登录) ==================
		auth := apiV1.Group("/auth")
		{
			auth.POST("/register", handler.Register)
			auth.POST("/login", handler.Login)
		}

		apiV1.GET("/ws", handler.ConnectWebSocket)
		// ================== 需要登录的接口 (使用 JWT 中间件) ==================
		// 使用 Use() 方法，表示这之下的所有路由都必须先通过 JWTAuth 检查
		authorized := apiV1.Group("/")
		authorized.Use(middleware.JWTAuth())
		{
			// 1. 【活动模块】 (占位)
			activity := authorized.Group("/activity")
			{
				activity.GET("/list", handler.GetActivityList)           //获取活动
				activity.GET("/detail/:id", handler.GetActivityDetail)   //查看具体活动
				activity.POST("/create", handler.CreateActivity)         //发布活动
				activity.DELETE("/delete/:id", handler.DeleteActivity)   //删除活动
				activity.POST("/apply/:id", handler.ApplyActivity)       //用户报名
				activity.POST("/audit", handler.AuditMember)             //发起人审核活动
				activity.GET("/members/:id", handler.GetActivityMembers) //发起人查看审核列表
			}

			// 2. 【路线模块】 (占位)
			route := authorized.Group("/route")
			{
				route.GET("/list", handler.GetRouteList)            //获取线路瀑布流
				route.GET("/detail/:id", handler.GetRouteDetail)    //获取单条路线详情
				route.POST("/favorite/:id", handler.ToggleFavorite) //收藏/取消收藏
				route.POST("/review/:id", handler.CreateReview)     //发布评论
				// 🔐 官方/专家专属接口 (发布路线后门)
				expertRoute := route.Group("/")
				expertRoute.Use(middleware.ExpertRequired())
				{
					// POST /api/v1/route/create
					expertRoute.POST("/create", handler.PublishRoute)
				}
			}
			// 3. 🌟【户外圈子与互动模块】 (朋友圈)
			post := authorized.Group("/post")
			{
				post.GET("/list", handler.GetPostList)         // 动态列表
				post.POST("/create", handler.PublishPost)      // 发布动态
				post.POST("/comment", handler.AddComment)      // 评论
				post.POST("/like/:id", handler.ToggleLike)     // 点赞/取消点赞
				post.DELETE("/delete/:id", handler.DeletePost) //删除动态
			}
			// 4. 【百科模块】 (占位)
			article := authorized.Group("/article")
			{
				// 公共接口 (所有人都能看)
				article.GET("/list", handler.GetArticleList)

				// 🔐 专家专属接口 (加盖 ExpertRequired 中间件！)
				expertOnly := article.Group("/")
				expertOnly.Use(middleware.ExpertRequired())
				{
					expertOnly.POST("/create", handler.PublishArticle)
					expertOnly.PUT("/update/:id", handler.UpdateArticle)
					expertOnly.DELETE("/delete/:id", handler.DeleteArticle)
				}
			}

			// 5. 【个人中心模块】 (已全部实现！)
			profile := authorized.Group("/profile")
			{
				profile.GET("/info", handler.GetProfile)                    // 看资料
				profile.PUT("/update", handler.UpdateProfile)               // 改资料
				profile.GET("/published", handler.GetMyPublishedActivities) // 我发布的
				profile.GET("/joined", handler.GetMyJoinedActivities)       // 我参与的
				profile.GET("/favorites", handler.GetMyFavoriteRoutes)      // 我的收藏
				profile.GET("/messages", handler.GetMyMessages)             // 我的通知
				profile.GET("/posts", handler.GetMyPosts)                   // 我的动态
				profile.POST("/logout", handler.Logout)                     // 退出登录
				profile.POST("/change-password", handler.ChangePassword)    //修改密码
				// 🔐 个人中心-我的百科 (仅专家能进)
				expertProfile := profile.Group("/expert")
				expertProfile.Use(middleware.ExpertRequired())
				{
					expertProfile.GET("/articles", handler.GetMyArticles)
				}
			}

			// 6. 【通用图片上传接口】
			common := authorized.Group("/common")
			{
				common.POST("/upload", handler.UploadImage)
			}

			// 7. 【消息通知模块】
			message := authorized.Group("/messages")
			{
				message.GET("/unread/count", handler.GetUnreadCount) //获取未读消息数量
				message.POST("/read/:id", handler.MarkMessageRead)   //标记消息已读
			}
		}
	}

	return r
}
