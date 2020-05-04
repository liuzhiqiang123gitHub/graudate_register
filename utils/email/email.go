package email

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"regexp"
	"strconv"
)

func SendMail(mailTo string, subject string, body string) error {
	//定义邮箱服务器连接信息，如果是阿里邮箱 pass填密码，qq邮箱填授权码
	mailConn := map[string]string{
		"user": "1964252235@qq.com",
		"pass": "xourwgilzdopefjh",
		"host": "smtp.qq.com",
		"port": "465",
	}

	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int

	m := gomail.NewMessage()
	//m.SetHeader("From","字节飞舞计算机系统有限公司" + mailConn["user"] )  //这种方式可以添加别名，即“XD Game”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
	m.SetHeader("From", mailConn["user"]) //这种方式可以添加别名，即“XD Game”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
	m.SetHeader("To", mailTo)             //发送给多个用户
	m.SetHeader("Subject", subject)       //设置邮件主题
	m.SetBody("text/html", body)          //设置邮件正文

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

	err := d.DialAndSend(m)
	return err
}
func StartSendEmail(mailTo, topic, body string) error {
	//定义收件人
	//mailTo := "1964252235@qq.com"
	//邮件主题为"Hello"
	//sub := "设置密码"
	// 邮件正文
	//body := "Good"
	err := SendMail(mailTo, topic, body)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

//验证邮箱
func EmailValidate(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

//验证手机
func VerifyMobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}
