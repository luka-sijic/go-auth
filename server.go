package main

import (
        "app/database"
        "net/http"
        "strings"
        "log"
        "github.com/golang-jwt/jwt/v5"

        "github.com/labstack/echo/v4"
        "github.com/labstack/echo/v4/middleware"

)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
                authHeader := c.Request().Header.Get("Authorization")
                if authHeader == "" {
                        return echo.NewHTTPError(http.StatusUnauthorized, "User must login")
                }

                tokenString := strings.Replace(authHeader, "Bearer ", "", 1)


                // Parse and validate the token
                token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
                        return jwtSecret, nil
                })

                if err != nil {
                        return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
                }

                if claims, ok := token.Claims.(*Claims); ok && token.Valid {
                        if (claims.Status == -1 || claims.Status == 0) {
                                return echo.NewHTTPError(http.StatusUnauthorized, "User is banned")
                        }
                        // Store the claims in the context for later use
                        c.Set("username", claims.Username)
                        c.Set("id", claims.ID)
                        c.Set("role", claims.Role)
                        c.Set("status", claims.Status)
                        return next(c)
                }

                return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
        }
}

func AdminAuth(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
                authHeader := c.Request().Header.Get("Authorization")
                if authHeader == "" {
                        return echo.NewHTTPError(http.StatusUnauthorized, "User must login")
                }

                tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

                // Parse and validate the token
                token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
                        return jwtSecret, nil
                })
                if err != nil || !token.Valid {
                        return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
                }

                // Extract claims
                claims, ok := token.Claims.(*Claims)
                if !ok || token.Valid == false {
                        return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims")
                }
                //log.Println(claims.Role)
                log.Println(claims.Username)
                if claims.Role != 3 {
                        return echo.NewHTTPError(http.StatusForbidden, "Access denied")
                }

                c.Set("username", claims.Username)
                c.Set("id", claims.ID)
                c.Set("role", claims.Role)
                c.Set("status", claims.Status)

                return next(c)
        }
}


func main() {
        e := echo.New()

        database.Connect()
        defer database.Close()

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

        e.GET("/", test2)

        e.Logger.Fatal(e.Start(":8082"))
}

func test2(c echo.Context) error {
        return c.JSON(http.StatusOK, "auth")
}
