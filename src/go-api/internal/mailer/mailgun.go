package mailer

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"github.com/MarekVigas/Postar-Jano/internal/config"
	"github.com/MarekVigas/Postar-Jano/internal/mailer/templates"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/pkg/errors"

	"go.uber.org/zap"
)

type Client struct {
	logger  *zap.Logger
	mailgun *mailgun.MailgunImpl

	confirmation *template.Template
	promo        *template.Template
	notification *template.Template
	sender       string
}

func NewClient(cfg *config.Mailer, logger *zap.Logger) (*Client, error) {
	mg := mailgun.NewMailgun(cfg.MailgunDomain, cfg.MailgunKey)
	if cfg.EUBase {
		mg.SetAPIBase(mailgun.APIBaseEU)
	}

	confirmationTemplate, err := templates.LoadFromFile(cfg.ConfirmationMailTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load confirmation template")
	}

	promoTemplate, err := templates.LoadFromFile(cfg.PromoMailTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load promo template")
	}

	var notificationTemplate *template.Template
	if cfg.NotificationMailTemplate != "" {
		var err error
		notificationTemplate, err = templates.LoadFromFile(cfg.NotificationMailTemplate)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load notification template")
		}
	}

	return &Client{
		logger:       logger,
		mailgun:      mg,
		confirmation: confirmationTemplate,
		promo:        promoTemplate,
		notification: notificationTemplate,
		sender:       fmt.Sprintf("robot@%s", cfg.MailgunDomain),
	}, nil
}

func (c *Client) ConfirmationMail(ctx context.Context, req *templates.ConfirmationReq) error {
	var b bytes.Buffer
	if err := c.confirmation.Execute(&b, req); err != nil {
		return errors.WithStack(err)
	}

	return c.send(ctx, fmt.Sprintf("Prijatie prihlášky na %s", req.EventName), b.String(),
		fmt.Sprintf("%s %s <%s>", req.ParentName, req.ParentSurname, req.Mail))
}

func (c *Client) PromoMail(ctx context.Context, req *templates.PromoReq) error {
	var b bytes.Buffer
	if err := c.promo.Execute(&b, req); err != nil {
		return errors.WithStack(err)
	}

	return c.send(ctx, "Prihlasovanie na letne akcie v Salezku", b.String(),
		fmt.Sprintf("<%s>", req.Mail))
}

func (c *Client) NotificationMail(ctx context.Context, req *templates.NotificationReq) error {
	if c.notification == nil {
		return errors.New("notification template not set")
	}
	var b bytes.Buffer
	if err := c.notification.Execute(&b, req); err != nil {
		return errors.WithStack(err)
	}

	return c.send(ctx, fmt.Sprintf("Prijatie platby za %s", req.EventName), b.String(),
		fmt.Sprintf("<%s>", req.Mail))
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
