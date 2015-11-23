package helper

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/astaxie/beego"
	"net/smtp"
	"strconv"
	"strings"
)

// test -> /interface/test
func UserSpaceUrl(username string) string {
	return "/space/" + username
}

// User.Login -> /user/login
func UrlFor(str string) string {
	return strings.ToLower(beego.UrlFor(str))
}

func MD5(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// 断言 interface -> string
func GetStringFromInterface(i interface{}) string {
	var data string

	if str, ok := i.(string); ok {
		data = str
	}

	return data
}

func EncodeEmail(email string) string {
	index := strings.Index(email, "@")
	return email[:index-2] + "**" + email[index:]
}

// 发送邮件
// from : 发送者名字
// emailAccount : 发送者邮箱账号
// emailPassword : 发送者邮箱密码
// smtpEmailServer : smtp 服务器地址
// to : 接受者邮箱地址
// subject : 发送邮箱主题
// body : 发送邮箱内容
// mailtype : 邮件类型
func SendEmail(from, emailAccount, emailPassword, smtpEmailServer, to, subject, body, mailtype string) error {
	var contentType string

	address := strings.Split(smtpEmailServer, ":")

	auth := smtp.PlainAuth("", emailAccount, emailPassword, address[0])

	if mailtype == "html" {
		contentType = "Content-Type: text/html; charset=UTF-8"
	} else {
		contentType = "Content-Type: text/plain; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + from + "<" + emailAccount + ">" + "\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	sendTo := strings.Split(to, ";")
	err := smtp.SendMail(smtpEmailServer, auth, emailAccount, sendTo, msg)
	if err != nil {
		return err
	}

	return nil
}

// 对中文进行unicode编码
func HtmlUnicode(cn string) string {
	rs := []rune(cn)
	html := ""

	for _, r := range rs {
		rint := int(r)

		if rint < 128 {
			html += string(r)
		} else {
			html += "&#" + strconv.Itoa(int(r)) + ";"
		}
	}

	return html
}

// 对中文进行unicode编码
func JsonUnicode(cn string) string {
	rs := []rune(cn)
	json := ""

	for _, r := range rs {
		rint := int(r)

		if rint < 128 {
			json += string(r)
		} else {
			json += "\\u" + strconv.FormatInt(int64(rint), 16)
		}
	}

	return json
}
