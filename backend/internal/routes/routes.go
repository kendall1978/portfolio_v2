package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kendall1978/project_salary/backend/internal/handlers"
	"github.com/kendall1978/project_salary/backend/internal/middleware"
)

func Setup(r *gin.Engine) {
	api := r.Group("/api")

	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Auth routes — register, login, MFA verify/recover are unauthenticated
	auth := api.Group("/auth")
	{
		auth.GET("/status", handlers.AuthStatus)
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
		auth.POST("/mfa/verify", handlers.MFAVerify)
		auth.POST("/mfa/recover", handlers.MFARecover)
	}

	// Auth routes that require a valid JWT
	authProtected := api.Group("/auth")
	authProtected.Use(middleware.AuthRequired())
	{
		authProtected.POST("/refresh", handlers.RefreshToken)
		authProtected.POST("/mfa/setup", handlers.MFASetup)
	}

	// Protected routes — all require JWT
	profile := api.Group("/profile")
	profile.Use(middleware.AuthRequired())
	{
		profile.GET("", handlers.GetProfile)
		profile.PUT("", handlers.UpdateProfile)
	}

	compensation := api.Group("/compensation")
	compensation.Use(middleware.AuthRequired())
	{
		compensation.GET("", handlers.ListCompensation)
		compensation.POST("", handlers.CreateCompensation)
		compensation.PUT("/:id", handlers.UpdateCompensation)
		compensation.DELETE("/:id", handlers.DeleteCompensation)
	}

	salaries := api.Group("/salaries")
	salaries.Use(middleware.AuthRequired())
	{
		salaries.GET("", handlers.ListSalaries)
		salaries.GET("/trends", handlers.SalaryTrends)
		salaries.GET("/compare", handlers.SalaryCompare)
	}

	api.GET("/regions", middleware.AuthRequired(), handlers.ListRegions)
	api.GET("/occupations", middleware.AuthRequired(), handlers.ListOccupations)

	scrape := api.Group("/scrape")
	scrape.Use(middleware.AuthRequired())
	{
		scrape.POST("/trigger", handlers.ScrapeTrigger)
		scrape.GET("/logs", handlers.ScrapeLogs)
	}

	// Static file serving for uploads
	r.Static("/uploads", "./uploads")

	// Public portfolio routes — no auth required
	api.GET("/blog", handlers.ListBlogPosts)
	api.GET("/blog/:id", handlers.GetBlogPost)
	api.GET("/home-sections", handlers.ListHomeSections)
	api.GET("/social-links", handlers.ListSocialLinks)
	api.GET("/site-settings", handlers.GetSiteSettings)

	// Protected portfolio admin routes
	blogAdmin := api.Group("/blog")
	blogAdmin.Use(middleware.AuthRequired())
	{
		blogAdmin.POST("", handlers.CreateBlogPost)
		blogAdmin.PUT("/:id", handlers.UpdateBlogPost)
		blogAdmin.DELETE("/:id", handlers.DeleteBlogPost)
	}

	homeSections := api.Group("/home-sections")
	homeSections.Use(middleware.AuthRequired())
	{
		homeSections.POST("", handlers.CreateHomeSection)
		homeSections.PUT("/:id", handlers.UpdateHomeSection)
		homeSections.DELETE("/:id", handlers.DeleteHomeSection)
	}

	socialLinks := api.Group("/social-links")
	socialLinks.Use(middleware.AuthRequired())
	{
		socialLinks.POST("", handlers.CreateSocialLink)
		socialLinks.PUT("/:id", handlers.UpdateSocialLink)
		socialLinks.DELETE("/:id", handlers.DeleteSocialLink)
	}

	siteSettings := api.Group("/site-settings")
	siteSettings.Use(middleware.AuthRequired())
	{
		siteSettings.PUT("", handlers.UpdateSiteSettings)
	}

	upload := api.Group("/upload")
	upload.Use(middleware.AuthRequired())
	{
		upload.POST("", handlers.UploadImage)
	}
}
