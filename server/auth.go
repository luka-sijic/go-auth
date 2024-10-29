package server 

import (
        "app/database"
        "app/models"

        "net/http"
        "context"
        "strconv"
        "time"
        "log"

        "github.com/google/uuid"
        "github.com/labstack/echo/v4"
        "github.com/golang-jwt/jwt/v5"
)


func loginUser(c echo.Context) error {
        login := new(models.LoginRequest)
        if err := c.Bind(login); err != nil {
                log.Println(err)
                return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
        }

        // Database query and password validation
        var storedHash string 
        err := database.DB.QueryRow(context.Background(), "SELECT id, password, role, status FROM users WHERE username=$1", login.Username).Scan(&login.ID, &storedHash, &login.Role, &login.Status)
        if err != nil || !checkPasswordHash(login.Password, storedHash) {
            log.Println(err)
            return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid username or password"})
        }

        if (login.Status == 0) {
            log.Println(err)
            return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User is banned from entering"})
        }

        // Create jwt token
        claims := &Claims{
                Username: login.Username,
                ID: login.ID,
                Role: login.Role,
                Status: login.Status,
                RegisteredClaims: jwt.RegisteredClaims{
                        ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
                },
        }

        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

        t, err := token.SignedString(jwtSecret)
        if err != nil {
                echo.NewHTTPError(http.StatusInternalServerError, "could not generate token")
        }

        return c.JSON(http.StatusOK, map[string]string{
                "token": t,
        })
}

func registerUser(c echo.Context) error {
        user := new(models.User)
        if err := c.Bind(user); err != nil {
                log.Println(err)
                return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to bind to user"})
        }

        // Hash the provided password
        hashedPassword, err := hashPassword(user.Password)
        if err != nil {
                log.Println(err)
                return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
        }
        // Create UUID
        id, _ := uuid.NewUUID()

        // Prepare database statement
        _, err = database.DB.Exec(context.Background(), "INSERT INTO users (id, username, password) VALUES ($1,$2,$3)", id, user.Username, hashedPassword)
        if err != nil {
                log.Println("Failed to create user:", err)
                return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
        }

        return c.JSON(http.StatusCreated, map[string]string{"message": "User created successfully"})
}

func getUser(c echo.Context) error {
    identifier := c.Param("id")

    var user models.User
    if isUUID(identifier) {
        err := database.DB.QueryRow(context.Background(), "SELECT id, username, credits, role, status, country, rating, avatar, TO_CHAR(creationdate, 'MM-DD-YY') AS formatted_date FROM users WHERE id=$1", identifier).
        Scan(&user.ID, &user.Username, &user.Credits, &user.Role, &user.Status, &user.Country, &user.Rating, &user.Avatar, &user.CreationDate)
        if err != nil {
            log.Println(err)
            return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
        }
    } else {
        err := database.DB.QueryRow(context.Background(), "SELECT id, username, credits, role, status, country, rating, avatar, TO_CHAR(creationdate, 'MM-DD-YY') AS formatted_date FROM users WHERE username=$1", identifier).
        Scan(&user.ID, &user.Username, &user.Credits, &user.Role, &user.Status, &user.Country, &user.Rating, &user.Avatar, &user.CreationDate)
        if err != nil {
            log.Println(err)
            return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
        }
    }
    
    return c.JSON(http.StatusOK, user)
}

func isUUID(s string) bool {
        _, err := uuid.Parse(s)
        return err == nil
}

func getUsers(c echo.Context) error {
        var users []models.User

        // Query the database to get all users
        rows, err := database.DB.Query(context.Background(), "SELECT id, username, credits, role, status, country, rating, avatar, TO_CHAR(creationdate, 'MM-DD-YY') AS formatted_date FROM users")
        if err != nil {
                return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve users"})
        }
        defer rows.Close()

        // Loop through the result set and scan each user
        for rows.Next() {
                var user models.User
                err := rows.Scan(&user.ID, &user.Username, &user.Credits, &user.Role, &user.Status, &user.Country, &user.Rating, &user.Avatar, &user.CreationDate)
                if err != nil {
                        log.Printf("Error scanning user: %v", err)
                        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error scanning user data"})
                }
                users = append(users, user)
        }

        // Check for errors after iterating through the rows
        if err = rows.Err(); err != nil {
                return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error processing user data"})
        }

        // Return the users information as JSON
        return c.JSON(http.StatusOK, users)
}



func updateUser(c echo.Context) error {
    id, _ := strconv.Atoi(c.Param("id"))
    user := new(models.User)
    if err := c.Bind(user); err != nil {
        return err
    }
    _, err := database.DB.Exec(context.Background(), "UPDATE users SET username=$1, password=$2 WHERE id=$3", user.Username, user.Password, id)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }
    return c.JSON(http.StatusOK, user)
}

func deleteUser(c echo.Context) error {
    id, _ := strconv.Atoi(c.Param("id"))
    _, err := database.DB.Exec(context.Background(), "DELETE FROM users WHERE id=$1", id)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }
    return c.NoContent(http.StatusNoContent)
}

func banUser(c echo.Context) error {
        username := c.Get("username").(string)
        ban := new(models.BanRequest)
        if err := c.Bind(ban); err != nil {
                log.Println(err)
                return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error scanning data"})
        }

        query := "UPDATE users SET status=0 WHERE username=$1"
        _, err := database.DB.Exec(context.Background(), query, ban.Username)
        if err != nil {
                return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
        }

        query = "INSERT INTO bans (userbanned, username, reason) VALUES ($1,$2,$3)"
        _, err = database.DB.Exec(context.Background(), query, ban.Username, username, ban.Reason)
        if err != nil {
                return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
        }

        return c.JSON(http.StatusOK, "user banned")
}

func unBanUser(c echo.Context) error {
        ban := new(models.BanRequest)
        if err := c.Bind(ban); err != nil {
                log.Println(err)
                return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error scanning data"})
        }

        query := "UPDATE users SET status=1 WHERE username=$1"
        _, err := database.DB.Exec(context.Background(), query, ban.Username)
        if err != nil {
                return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
        }

        query = "UPDATE bans SET active=0 WHERE userbanned=$1"
        _, err = database.DB.Exec(context.Background(), query, ban.Username)
        if err != nil {
                return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
        }

        return c.JSON(http.StatusOK, "user unbanned")
}
