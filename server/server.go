package server

import (
        "net/http"

        "github.com/labstack/echo/v4"
        "github.com/labstack/echo/v4/middleware"

)


func Start() {
        e := echo.New()

        e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
                AllowOrigins: []string{"*"}, 
                AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
                AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization}, 
        AllowCredentials: true,
        }))
        e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

        e.Use(middleware.Logger())
        e.Use(middleware.Recover())

        // Authentication Endpoints
        e.POST("/register", registerUser)
        e.POST("/login", loginUser)

        // User Endpoints
        e.GET("/users", getUsers, Auth)
        e.GET("/users/:id", getUser, Auth)
        e.PUT("/users/:id", updateUser, AdminAuth)
        e.DELETE("/users/:id", deleteUser, AdminAuth)

        e.POST("/users/ban", banUser, AdminAuth)
        e.POST("/users/unban", unBanUser, AdminAuth)

        e.GET("/", root)

        e.Logger.Fatal(e.Start(":8082"))
}

func root(c echo.Context) error {
        return c.JSON(http.StatusOK, "auth")
}
