package smtp

import (
	"fmt"

	"github.com/alpardfm/go-toolkit/codes"
	"github.com/alpardfm/go-toolkit/errors"
	mail "github.com/go-mail/gomail"
)

type GoMailInterface interface {
	DialAndSend(messages ...*mail.Message) error
}

type ConfigGoMail struct {
	Username string
	Password string
	Host     string
	Port     int64
}

type gomail struct {
	dialer *mail.Dialer
}

func InitGoMail(cfg ConfigGoMail) GoMailInterface {
	result := gomail{
		dialer: mail.NewDialer(cfg.Host, int(cfg.Port), cfg.Username, cfg.Password),
	}

	return &result
}

func (g *gomail) DialAndSend(messages ...*mail.Message) error {
	if err := g.dialer.DialAndSend(messages...); err != nil {
		return errors.NewWithCode(codes.CodeSMTPError, fmt.Sprintf("failed to send messages, %v", err))
	}
	return nil
}
