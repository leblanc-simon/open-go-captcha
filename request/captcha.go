package request

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"leblanc.io/open-go-captcha/captcha"
	"leblanc.io/open-go-captcha/config"
	"leblanc.io/open-go-captcha/crypto"
	"leblanc.io/open-go-captcha/handler"
)

type ResponseIcon struct {
	DataURI string `json:"dataUri"`
	Identifier string `json:"identifier"`
}

type CaptchaHttpResponse struct {
	Token string `json:"token"`
	Icons []ResponseIcon `json:"icons"`
	Lifetime int `json:"lifetime"`
}

func (response *CaptchaHttpResponse) FromCaptcha(captcha captcha.Captcha) {
	response.Token = captcha.Token
	response.Icons = []ResponseIcon{}
	response.Lifetime = cfg.Redis.Expire

	for _, icon := range captcha.Icons {
		response.Icons = append(response.Icons, ResponseIcon{
			DataURI: icon.DataURI,
			Identifier: icon.Identifier,
		})
	}
}

var cfg *config.Config

func Initialize(c *config.Config) {
	cfg = c
}

func GetCaptcha(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	captcha, err := captcha.NewCaptcha(cfg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())

		return
	}

	handler.SetAnswer(captcha.Session, captcha.Answers)

	var captchaHttpResponse CaptchaHttpResponse
	captchaHttpResponse.FromCaptcha(*captcha)

	response, err := json.Marshal(captchaHttpResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())

		return
	}

	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(response))
}

func getSession(w http.ResponseWriter, r *http.Request) string {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return ""
	}

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())

		return ""
	}

	session, err := crypto.Decrypt(r.FormValue("_opc_token"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())

		return ""
	}

	return session
}

func CheckCaptcha(w http.ResponseWriter, r *http.Request) {
	session := getSession(w, r)
	if session == "" {
		return
	}

	answers := strings.Split(r.FormValue("answers"), ",")
	if len(answers) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Not answer found")

		return
	}

	if !handler.CheckAnswer(session, answers) {
		w.WriteHeader(http.StatusForbidden)

		return
	}

	err := handler.StoreValidResult(session)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, "{\"lifetime\": " + strconv.Itoa(cfg.Redis.LongExpire) + "}")
}

func ConfirmCaptcha(w http.ResponseWriter, r *http.Request) {
	session := getSession(w, r)
	if session == "" {
		return
	}

	if !handler.CheckValidResult(session) {
		w.WriteHeader(http.StatusForbidden)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}