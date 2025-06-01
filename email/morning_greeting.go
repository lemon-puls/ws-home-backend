package email

import (
	"fmt"
	"log"
	"time"

	"ws-home-backend/config"

	"github.com/robfig/cron/v3"
	"gopkg.in/gomail.v2"
)

type MorningGreeting struct {
	fromEmail    string
	fromPassword string
	toEmails     []string
	smtpServer   string
	smtpPort     int
	cron         *cron.Cron
}

func NewMorningGreeting(cfg *config.EmailConfig) *MorningGreeting {
	return &MorningGreeting{
		fromEmail:    cfg.FromEmail,
		fromPassword: cfg.FromPassword,
		toEmails:     cfg.ToEmails,
		smtpServer:   cfg.SmtpServer,
		smtpPort:     cfg.SmtpPort,
		cron:         cron.New(cron.WithSeconds()),
	}
}

func (mg *MorningGreeting) SendGreeting() error {
	// 获取当前时间
	now := time.Now()
	greeting := fmt.Sprintf(`
		<h2>定时问候！</h2>
		<p>现在是 %s</p>
		<p>这是一条定时发送的问候消息！</p>
	`, now.Format("2006-01-02 15:04:05"))

	// 创建邮件发送器
	d := gomail.NewDialer(mg.smtpServer, mg.smtpPort, mg.fromEmail, mg.fromPassword)

	// 为每个收件人单独发送邮件
	for _, toEmail := range mg.toEmails {
		m := gomail.NewMessage()
		m.SetHeader("From", mg.fromEmail)
		m.SetHeader("To", toEmail)
		m.SetHeader("Subject", "定时问候")
		m.SetBody("text/html", greeting)

		if err := d.DialAndSend(m); err != nil {
			return fmt.Errorf("发送问候邮件失败: %v", err)
		}
		log.Printf("成功发送邮件到: %s", toEmail)
	}

	return nil
}

func (mg *MorningGreeting) StartScheduler() {
	// 添加定时任务，每10秒执行一次
	_, err := mg.cron.AddFunc("*/10 * * * * ?", func() {
		if err := mg.SendGreeting(); err != nil {
			log.Printf("发送问候邮件失败: %v", err)
		} else {
			log.Println("所有问候邮件发送成功！")
		}
	})

	if err != nil {
		log.Printf("添加定时任务失败: %v", err)
		return
	}

	// 启动定时任务
	mg.cron.Start()
	log.Printf("定时问候任务已启动，每10秒执行一次，收件人列表: %v", mg.toEmails)
}

func (mg *MorningGreeting) StopScheduler() {
	if mg.cron != nil {
		mg.cron.Stop()
		log.Println("定时问候任务已停止")
	}
}
