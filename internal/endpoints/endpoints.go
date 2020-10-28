package endpoints

import (
	"fmt"
	"html"
	"net"
	"net/smtp"
	"os"
	"regexp"
	"strings"

	"github.com/codemicro/contactFormProcessor/internal/helpers"

	"github.com/dchest/captcha"
	"github.com/gofiber/fiber/v2"
	"github.com/jordan-wright/email"
)

var (
	senderAddress    string
	recipientAddress string
	smtpServer       string
	smtpPort         string
	smtpUsername     string
	smtpPassword     string
)

func init() {
	senderAddress = getEnvironmentVar("EMAIL_SENDER_ADDRESS")
	recipientAddress = getEnvironmentVar("EMAIL_RECIPIENT_ADDRESS")
	smtpServer = getEnvironmentVar("EMAIL_SMTP_SERVER")
	smtpPort = getEnvironmentVar("EMAIL_SMTP_PORT")
	smtpUsername = getEnvironmentVar("EMAIL_SMTP_USERNAME")
	smtpPassword = getEnvironmentVar("EMAIL_SMTP_PASSWORD")
}

func getEnvironmentVar(key string) string {
	v := os.Getenv(key)
	if v == "" {
		fmt.Printf("%s envirnoment variable not set.\n", key)
		os.Exit(1)
	}
	return v
}

func EndpointProcessEmail(c *fiber.Ctx) error {

	for _, v := range []string{"subject", "from", "message", "captcha"} {
		if c.FormValue(v) == "" {
			return responseBadRequest(c, fmt.Sprintf("param %s required", v))
		}
	}

	// Gather data
	subject := html.EscapeString(c.FormValue("subject"))
	from := c.FormValue("from")
	if emailOk, err := isEmailValid(from); err != nil {
		return err
	} else if !emailOk {
		return responseBadRequest(c, "invalid email address")
	}
	escapedFrom := html.EscapeString(from)
	message := wordWrap(html.EscapeString(c.FormValue("message")), 100)
	ipAddr := helpers.GetIP(c)

	// Validate CAPTCHA
	captchaCode := c.FormValue("captcha")
	captchaId := c.Cookies(captchaCookieName)
	fmt.Println(captchaCode)
	if !captcha.VerifyString(captchaId, captchaCode) {
		return responseBadRequest(c, "CAPTCHA invalid")
	}

	// Form email
	separator := strings.Repeat("-", 100)
	messageBody := strings.Join([]string{
		"The following content is from a message submitted from your contact form.",
		separator,
		"<b>IP address:</b> " + ipAddr,
		"<b>Sender:</b> " + escapedFrom,
		"<b>Subject: </b>" + subject,
		"<b>Message:</b><br>" + message,
		separator,
		"Message end.",
	}, "<br>\n")

	// Create email
	e := email.NewEmail()
	e.From = fmt.Sprintf("Contact form <%s>", senderAddress)
	e.To = []string{recipientAddress}
	e.ReplyTo = []string{from}
	e.Subject = subject
	e.HTML = []byte(messageBody)

	err := e.Send(fmt.Sprintf("%s:%s", smtpServer, smtpPort), smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer))
	if err != nil {
		return err
	}

	return c.JSON(StatusResponse{"sent"})
}

type StatusResponse struct {
	Status string `json:"status"`
}

func responseBadRequest(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(StatusResponse{message})
}

var (
	emailRegex = regexp.MustCompile("(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])")
)

func isEmailValid(eaddr string) (bool, error) {
	if len(eaddr) < 3 && len(eaddr) > 254 {
		return false, nil
	}
	if !emailRegex.MatchString(eaddr) {
		return false, nil
	}
	// Hopefully anything that would cause this to panic has been filtered out by the above regex thing
	parts := strings.Split(eaddr, "@")

	domain := strings.ToLower(parts[1])

	mx, err := net.LookupMX(domain)
	if err != nil {
		return false, nil // Invalid domains raise an error
	}
	return len(mx) != 0, nil
}

func wordWrap(inputString string, maxLineLen int) string {
	components := strings.Split(inputString, "\n")
	var output []string
	for _, component := range components {
		var line string
		words := strings.Split(component, " ")
		for _, word := range words {
			if len(line+word)+1 > maxLineLen {
				output = append(output, line)
				line = word + " "
			} else {
				line += word + " "
			}
		}
		output = append(output, line+"\n")
	}
	return strings.TrimSpace(strings.Join(output, "<br>"))
}
