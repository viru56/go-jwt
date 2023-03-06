package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/viru56/go-jwt/intializers"
)

type User struct {
	ID        uint
	Email     string
	CreatedAt string
}

func RequireAuth(c *gin.Context) {

	// read cookie
	tokenString, err := c.Cookie("token")
	if err != nil {
		fmt.Println("error in 1", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// validate token
	// Parse takes the token string and a function for looking up the key. The latter is especially
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})

 // No need to check for token expiration
	// if float64(time.Now().Unix()) > claims["exp"].(float64) {
	// 	fmt.Println("error in 3", err)
	// 	c.AbortWithError(http.StatusUnauthorized, errors.New("token expired"))
	// 	return
	// }
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var user User
		intializers.DB.First(&user, claims["sub"])
		fmt.Println("user details", claims["sub"], user)
		if user.ID == 0 {
			fmt.Println("error in 4")
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.Set("user", user)
		c.Next()
	} else {
		c.AbortWithError(http.StatusUnauthorized, err)
	}
}
