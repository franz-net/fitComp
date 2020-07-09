package main

import "github.com/gin-gonic/contrib/sessions"

func initializeRoutes() {

	store := sessions.NewCookieStore([]byte("sessionSuperSecret"))
	router.Use(sessions.Sessions("sessionName", store))

	api := router.Group("/api/v1")
	// Unauthenticated groups
	{
		api.POST("/login", loginHandler)
		api.POST("/register", registrationHandler) //adds 10 bucks to the pot, needs invite code
	}
	// session termination
	{
		sessionAuth := api.Group("/session")
		sessionAuth.Use(AuthenticationRequired("user", "admin"))
		{
			sessionAuth.GET("/logout", logoutHandler)
		}
	}
	{
		userAuth := api.Group("/user")
		userAuth.Use(AuthenticationRequired("user"))
		{
			userAuth.GET("/my-measurements", myMeasurementsHandler)
			userAuth.POST("/add-measurement", newMeasurementHandler)
		}
	}
	// User | admin auth
	{
		allAuth := api.Group("/all")
		allAuth.Use(AuthenticationRequired("user", "admin"))
		{
			allAuth.GET("/measurements", measurementHistoryHandler)
			allAuth.GET("/prizes", prizeHistoryHandler)
			allAuth.GET("/standings", standingsHandler)
		}
	}

	// Admin auth
	{
		adminAuth := api.Group("/admin")
		adminAuth.Use(AuthenticationRequired("admin"))
		{
			adminAuth.GET("/newInviteCode", newInviteCodeHandler)
			adminAuth.POST("/resetUserPass", resetPasswordHandler)
			adminAuth.GET("/listUsers", listUsersHandler)
			adminAuth.GET("/listInviteCodes", listInviteCodesHandler)
			adminAuth.POST("/deleteUser", deleteUserHandler)
		}
	}
}
