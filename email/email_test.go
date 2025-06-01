package email_test

import (
	"gopkg.in/gomail.v2"
	"log"
	"testing"
	"time"
)

func TestEmail(t *testing.T) {

	t.Run("TestEmail", func(t *testing.T) {
		// 创建新的邮件消息
		m := gomail.NewMessage()

		// 设置邮件头部信息
		m.SetHeader("From", "13434615275@163.com") // 发送方
		m.SetHeader("To", "2814869489@qq.com")     // 接收方
		m.SetHeader("Subject", "测试邮件")             // 邮件主题

		// 使用这种方式发送图片会被退信 报被系统标记为垃圾邮件
		//// 嵌入图片并获取 CID
		//m.Embed("../assets/ws.jpg", gomail.Rename("ws.jpg"), gomail.SetHeader(map[string][]string{
		//	"Content-Type": {"image/jpeg"},
		//}))
		// 邮件内容，包括嵌入的图片
		//m.SetBody("text/html", `<h2>这是嵌入图片的测试邮件！</h2><img src="cid:ws.jpg">`)

		// 添加附件（图片）
		//m.Attach("../assets/ws.jpg") // 替换成你的图片路径
		//m.SetBody("text/html", "<h2>我是小哆啦呀！</h2><p style='color:red'>这是一封测试邮件。</p>") // 邮件内容，支持HTML格式

		m.Embed("../assets/ws.jpg")
		m.SetBody("text/html", `<img src="cid:ws.jpg" alt="My image" />`)

		// 设置邮件服务器信息
		d := gomail.NewDialer(
			"smtp.163.com",        // SMTP服务器地址
			25,                    // 端口号
			"13434615275@163.com", // 发件人邮箱账号
			"EH2hEmJnzpcYPNjv",    // 发件人邮箱密码
		)

		// 发送邮件
		if err := d.DialAndSend(m); err != nil {
			log.Fatalf("发送邮件时出错: %v", err) // 错误处理
		} else {
			log.Println("邮件发送成功！")
		}

		time.Sleep(120 * time.Second)
	})
}
