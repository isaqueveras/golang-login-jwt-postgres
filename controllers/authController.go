package controllers

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/isaqueveras/golang-login-jwt-postgres/database"
	"github.com/isaqueveras/golang-login-jwt-postgres/models"
	"golang.org/x/crypto/bcrypt"
)

var SecretKey string = "isaque"

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// Generate password
	password, _ := bcrypt.GenerateFromPassword([]byte(data["passw"]), 14)

	user := models.User{
		Name:  data["name"],
		Email: data["email"],
		Passw: string(password),
	}

	database.DB.Table("users_test_2").Create(&user)
	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	database.DB.
		Table("users_test_2").
		Where("email = ?", data["email"]).
		First(&user)

	if user.Id == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{"message": "user not found"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Passw), []byte(data["passw"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"message": "incorrect password"})
	}

	clams := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 1 day
		Issuer:    strconv.Itoa(int(user.Id)),
	})

	token, err := clams.SignedString([]byte(SecretKey))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{"message": "Could not login"})
	}

	cookie := fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24), // 1 day
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(token)
}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("token")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{"message": "Unauthorized"})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User
	database.DB.
		Table("users_test_2").
		Where("id = ?", claims.Issuer).
		First(&user)

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)
	return c.JSON(fiber.Map{"message": "logout with success"})
}
