package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/viru56/go-jwt/intializers"
	"github.com/viru56/go-jwt/models"
	"golang.org/x/crypto/bcrypt"
)

type Body struct {
	Email string
	Password string
}

func Signup(c *gin.Context) {

	// read body
	var body Body;
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	// generate hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to generate password hash",
		})
		return
	}

	// create user in db
	user := models.User{Email: body.Email,Password: string(hash)}
	result := intializers.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
	// return success response

}

func Login(c *gin.Context) {

	// Get the Email and Password from req body
	var body Body;
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Look up requested user
	var user models.User
	intializers.DB.First(&user, "email = ?", body.Email)
	
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not Found",
		})
		return
	} 

	// Validate password
	err:= bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password did not match",
		})
		return
	}

	// Generate a JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // 1 month validity
	})
	
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to generate token",
		})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode);
	c.SetCookie("token",tokenString,3600 * 24 * 30,"/","", false,true)
	// return JWT
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

func GetUser(c *gin.Context) {
	 user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
func Logout(c *gin.Context) {
	c.SetCookie("token","",-1,"/","",false,true)
 c.JSON(http.StatusOK, gin.H{})
}
