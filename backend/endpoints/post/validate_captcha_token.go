package post

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// ValidateCAPTCHAToken handles the endpoint to verify a CAPTCHA token.
// Parameters:
//   - c: the gin context, the request and response http.
//
// It sends a request to Google's reCAPTCHA API and returns whether the token is valid.
func ValidateCAPTCHAToken(c *gin.Context) {
	var captchaResponseBody struct {
		Success bool `json:"success"`
	}

	token := c.Query("token")

	resp, err := http.Get("https://www.google.com/recaptcha/api/siteverify?secret=" + os.Getenv("RECAPTCHA_SECRET_KEY") + "&response=" + token)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&captchaResponseBody)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"human": captchaResponseBody.Success})
}
