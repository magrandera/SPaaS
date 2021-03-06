package auth

import (
	"net/http"
	"time"

	"github.com/magrandera/SPaaS/common"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/magrandera/SPaaS/config"
)

// ChangePassword allows for changing the password
func ChangePassword(c echo.Context) error {
	newPassword := c.FormValue("password")
	if len(newPassword) < 8 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "password has to be at leat 8 characters",
		})
	}
	hashedPassword, err := common.HashPassword(newPassword)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	config.Cfg.Config.Set("password", hashedPassword)
	config.Save()
	return c.NoContent(http.StatusOK)
}

// Login is the endpoint for login
func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	if username == config.Cfg.Config.GetString("username") && common.CheckPasswordHash(password, config.Cfg.Config.GetString("password")) {
		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["username"] = username
		claims["admin"] = true
		claims["created"] = time.Now().Unix()
		claims["exp"] = time.Now().Add(time.Hour * 24 * 365).Unix()

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(config.Cfg.Config.GetString("secret")))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}

	return echo.ErrUnauthorized
}

// GetToken generates a token for internal request use
func GetToken() (string, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = "spaas"
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 24 * 365).Unix()
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(config.Cfg.Config.GetString("secret")))
	if err != nil {
		return "", err
	}
	return t, nil
}
