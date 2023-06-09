package middleware

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

type CustomRequest struct {
	Payload jwt.MapClaims
}

func AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Load dotenv
		dotenvErr := godotenv.Load(".env")
		if dotenvErr != nil {
			log.Fatal("Error Loading .env File: ", dotenvErr)
		}
		secretKey := []byte(os.Getenv("SECRET_KEY"))

		// Get JWT from request Header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Fatal("Token Error: missing authorization header")
		}

		// Split header to obtain token
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			log.Fatal("Authorization Error: invalid authorization header format")
		}

		// Decode the JWT
		token := headerParts[1]
		decoded, validationErr := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		if validationErr != nil {
			log.Fatal("Validation Error: ", validationErr)
		}

		// Ensure JWT is valid
		if !decoded.Valid {
			log.Fatal("Validation Error: token is invalid")
		}

		// Get claims associated with token
		claims, ok := decoded.Claims.(jwt.MapClaims)
		if !ok {
			log.Fatal("Claims Error: failed to get JWT claims")
		}

		// Attach the payload to the CustomRequest struct
		customReq := &CustomRequest{
			Payload: claims,
		}

		c.Set("custom_request", customReq)
		c.Next()
	}
}
