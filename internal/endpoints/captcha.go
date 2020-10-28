package endpoints

import (
	"bytes"
	"os"
	"time"

	"github.com/dchest/captcha"
	"github.com/gofiber/fiber/v2"
)

var (
	captchaCookieName string
)

func init() {
	captchaCookieName = os.Getenv("CAPTCHA_COOKIE_NAME")
	if captchaCookieName == "" {
		captchaCookieName = "captchaToken"
	}
}

func EndpointGenerateCaptcha(c *fiber.Ctx) error {
	var id string
	var changedId bool

	fromCookie := c.Cookies(captchaCookieName)

	if fromCookie == "" {
		id = captcha.New()
		changedId = true
	} else {
		id = fromCookie
		captcha.Reload(id)
	}

	captchaData, err := generateImageCaptcha(id)
	if err != nil {
		id = captcha.New()
		changedId = true
		captchaData, err = generateImageCaptcha(id)
		if err != nil {
			return err
		}
	}

	if changedId {
		c.Cookie(&fiber.Cookie{
			Name:    captchaCookieName,
			Value:   id,
			Expires: time.Now().Add(time.Minute * 30),
		})
	}

	c.Type("png")
	return c.Send(captchaData)
}

func generateImageCaptcha(id string) ([]byte, error) {
	b := new(bytes.Buffer)
	err := captcha.WriteImage(b, id, 275, 100)
	if err != nil {
		return []byte{}, err
	}

	return b.Bytes(), nil
}
