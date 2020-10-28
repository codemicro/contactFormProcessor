# contactFormProcessor

A small service to act as a processor for a simple contact form that includes a sender email address, subject line, message and CAPTCHA. When the form is submitted, an email is sent by this service.

### Usage

```bash
curl --request GET http://127.0.0.1:8000/captcha # view CAPTCHA. Will also set a cookie
curl --request POST \
  --url http://127.0.0.1:8000/send \
  --header 'Content-Type: multipart/form-data' \
  --cookie captchaToken=YbdHoarzgoLg4T7H2in7 \
  --form 'subject=This is a subject' \
  --form from=youremail@domain.net \
  --form 'message=Hey there. This is a message.' \
  --form captcha=396673 # submits form
```

### Prerequisites

* Docker or Go 1.14 or later
* Credentials for a SMTP server

### Setup (Docker)

* Perform one of the following options:
  * Clone this repository and build the Docker image
    ```bash
    git clone https://github.com/codemicro/contactFormProcessor.git
    cd contactFormProcessor
    docker build -t contactformprocessor .
    ```
    * Or, provided it's not been removed by the image retention policy, pull a prebuilt image from the Docker Hub
    ```bash
    docker pull codemicro/contact-form-processor
    ```
    and use the image name `codemicro/contact-form-processor` in the next step.
* Start the Docker container
  ```bash
  docker run -d --restart unless-stopped -p 8000:80 \
    -e EMAIL_SENDER_ADDRESS="noreply@domain.net" \
    -e EMAIL_RECIPIENT_ADDRESS="yourEmail@gmail.com" \
    -e EMAIL_SMTP_SERVER="smtp.gmail.com" \
    -e EMAIL_SMTP_PORT="587" \
    -e EMAIL_SMTP_USERNAME="anotherEmail@gmail.com" \
    -e EMAIL_SMTP_PASSWORD="password!"\
    contactformprocessor
  ```
  This will start the server on port 8000 of your local machine.

### Setup (compilation)

This should run on any platform that can run Go 1.14 or later.

* Clone this repository
  ```bash
  git clone https://github.com/codemicro/contactFormProcessor.git
  cd contactFormProcessor
  ```
  
* Build using Go
  ```bash
  go build github.com/codemicro/contactFormProcessor/cmd/contactFormProcessor
  sudo chmod +x contactFormProcessor
  ```
  
* Run
  ```bash
  EMAIL_SENDER_ADDRESS="noreply@domain.net" 
  EMAIL_RECIPIENT_ADDRESS="yourEmail@gmail.com" 
  EMAIL_SMTP_SERVER="smtp.gmail.com" 
  EMAIL_SMTP_PORT="587" 
  EMAIL_SMTP_USERNAME="anotherEmail@gmail.com" 
  EMAIL_SMTP_PASSWORD="password!"
  ./contactFormProcessor
  ```
  This will start the server on port 80.

### Configuration options

All the following can be used as environment variables to configure the service.

| Key                       | Description                                                        | Default          |
| ------------------------- | ------------------------------------------------------------------ | ---------------- |
| `CAPTCHA_COOKIE_NAME`     | Name of the cookie to use for CAPTCHAs.                            | `"captchaToken"` |
| `EMAIL_SENDER_ADDRESS`    | *Required.* Email sender's address.                                |                  |
| `EMAIL_RECIPIENT_ADDRESS` | *Required.* Where to send any emails.                              |                  |
| `EMAIL_SMTP_SERVER`       | *Required.*                                                        |                  |
| `EMAIL_SMTP_PORT`         | *Required.*                                                        |                  |
| `EMAIL_SMTP_USERNAME`     | *Required.* Username to use for authentication on the SMTP server. |                  |
| `EMAIL_SMTP_PASSWORD`     | *Required.* Password to use for authentication on the SMTP server. |                  |

### Reverse proxying

To reverse proxy this using the Apache2 HTTP server:

```
<Location /a/path>
    ProxyPass http://127.0.0.1:8000/
    ProxyPassReverse http://127.0.0.1:8000/
</Location>
```