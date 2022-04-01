package mailer

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/MarekVigas/Postar-Jano/internal/mailer/templates"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/pkg/errors"

	"go.uber.org/zap"
)

type Client struct {
	logger  *zap.Logger
	mailgun *mailgun.MailgunImpl

	sender string
}

const (
	mailgunDomain = "MAILGUN_DOMAIN"
	mailgunKey    = "MAILGUN_KEY"
)

func NewClient(logger *zap.Logger) (*Client, error) {
	domain := os.Getenv(mailgunDomain)
	if domain == "" {
		return nil, errors.New("Mailgun domain is not defined.")
	}

	key := os.Getenv(mailgunKey)
	if key == "" {
		return nil, errors.New("Mailgun key is not defined.")
	}

	mg := mailgun.NewMailgun(domain, key)
	mg.SetAPIBase(mailgun.APIBaseEU)

	return &Client{
		logger:  logger,
		mailgun: mg,
		sender:  fmt.Sprintf("robot@%s", domain),
	}, nil
}

func (c *Client) ConfirmationMail(ctx context.Context, req *templates.ConfirmationReq) error {
	var b bytes.Buffer
	if err := templates.Confirmation.Execute(&b, req); err != nil {
		return errors.WithStack(err)
	}

	return c.send(ctx, fmt.Sprintf("Prijatie prihlášky na %s", req.EventName), b.String(),
		fmt.Sprintf("%s %s <%s>", req.ParentName, req.ParentSurname, req.Mail))
}

func (c *Client) send(ctx context.Context, subject string, body string, recipient string) error {
	msg := c.mailgun.NewMessage(c.sender, subject, "", recipient)
	msg.SetHtml(body)

	resp, id, err := c.mailgun.Send(ctx, msg)
	if err != nil {
		return errors.Wrap(err, "failed to send a message")
	}
	c.logger.Info("Message sent", zap.String("id", id), zap.String("resp", resp))
	return nil
}
