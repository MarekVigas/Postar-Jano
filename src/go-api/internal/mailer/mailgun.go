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

	return &Client{
		logger:  logger,
		mailgun: mailgun.NewMailgun(domain, key),
	}, nil
}

func (c *Client) ConfirmationMail(ctx context.Context, req *templates.ConfirmationReq) error {
	var b bytes.Buffer
	if err := templates.Confirmation.Execute(&b, req); err != nil {
		return errors.WithStack(err)
	}

	return c.send(ctx, "robot@mailgun.sbb.sk", "Info o letnej akcii v Salezku", b.String(),
		fmt.Sprintf("%s %s <%s>", req.Name, req.Surname, req.Recipient))
}

func (c *Client) send(ctx context.Context, sender string, subject string, body string, recipient string) error {
	msg := c.mailgun.NewMessage(sender, subject, "", recipient)
	msg.SetHtml(body)

	resp, id, err := c.mailgun.Send(ctx, msg)
	if err != nil {
		return errors.Wrap(err, "failed to send a message")
	}
	c.logger.Info("Message sent", zap.String("id", id), zap.String("resp", resp))
	return nil
}
