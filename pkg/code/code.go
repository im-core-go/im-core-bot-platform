package code

import (
	"context"
	"crypto/rand"
	_ "embed"
	"errors"
	"fmt"
	"github.com/im-core-go/im-core-bot-platform/pkg/logger"
	"github.com/im-core-go/im-core-bot-platform/pkg/mail"
	"github.com/im-core-go/im-core-bot-platform/pkg/sms"
	"math/big"

	"github.com/redis/go-redis/v9"
)

var (
	ErrSendTooFrequent      = errors.New("send code too frequent")
	ErrInvalidCode          = errors.New("invalid code")
	ErrTooManyVerifications = errors.New("too many verifications")
	ErrWrongCode            = errors.New("wrong code")
)

//go:embed lua/set_code.lua
var setCodeScript string

//go:embed lua/verify_code.lua
var varifyCodeScript string

type Manager struct {
	sms  sms.SMS
	mail mail.Manager
	cmd  redis.Cmdable
}

func NewManager(cmd redis.Cmdable) *Manager {
	qqMail, err := mail.NewQQMail()
	if err != nil {
		panic(err)
	}
	return &Manager{
		sms:  nil,
		mail: qqMail,
		cmd:  cmd,
	}
}

func (m *Manager) SendEmailCode(ctx context.Context, target string) error {
	if m.cmd == nil {
		return errors.New("redis cmd is nil")
	}
	code, err := m.genCode()
	if err != nil {
		return err
	}
	err = m.mail.Send(mail.Message{
		Title:    "验证码",
		Content:  code,
		Appendix: "",
	}, []string{target})
	if err != nil {
		logger.L().Errorf("send email code failed: %v", err)
		return err
	}
	res, err := m.cmd.Eval(ctx, setCodeScript, []string{m.genRedisKey(target)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		return nil
	case -2:
		return ErrSendTooFrequent
	default:
		logger.L().Errorf("send code failed: invalid redis ret code :%d", res)
		return errors.New("system err")
	}
}
func (m *Manager) VerifyCode(ctx context.Context, code string, target string) (bool, error) {
	if m.cmd == nil {
		return false, errors.New("redis cmd is nil")
	}
	res, err := m.cmd.Eval(ctx, varifyCodeScript, []string{m.genRedisKey(target)}, code).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		return true, nil
	case -1:
		return false, ErrTooManyVerifications
	case -2:
		return false, ErrWrongCode
	case -3:
		return false, ErrInvalidCode
	default:
		logger.L().Errorf("verify code failed: invalid redis ret code :%d", res)
		return false, errors.New("system err")
	}
}
func (m *Manager) genRedisKey(target string) string {
	return "code:" + target
}

func (m *Manager) genCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
