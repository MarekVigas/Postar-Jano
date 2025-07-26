package registration

import (
	"context"
	"fmt"
	"github.com/MarekVigas/Postar-Jano/internal/model"
	"github.com/MarekVigas/Postar-Jano/internal/repository"
	"github.com/MarekVigas/Postar-Jano/internal/resources"
	"github.com/MarekVigas/Postar-Jano/internal/services/mailer/templates"
	"github.com/MarekVigas/Postar-Jano/pkg/logger"
	"github.com/MarekVigas/Postar-Jano/pkg/payme"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	ErrOverLimit = errors.New("Limit exceeded")
	ErrNotActive = errors.New("Event not active")
)

type PromoValidator interface {
	ValidateTokenWithQueryerContext(ctx context.Context, db sqlx.QueryerContext, token string) (*model.PromoCode, error)
	MarkTokenUsage(ctx context.Context, tx sqlx.QueryerContext, key string) (err error)
}

type EmailSender interface {
	ConfirmationMail(ctx context.Context, req *templates.ConfirmationReq) error
	NotificationMail(ctx context.Context, req *templates.NotificationReq) error
}

type Manager struct {
	postgresDB     *repository.PostgresDB
	promoValidator PromoValidator
	emailSender    EmailSender
}

func NewManager(postgresDB *repository.PostgresDB, promoValidator PromoValidator, emailSender EmailSender) *Manager {
	return &Manager{
		postgresDB:     postgresDB,
		promoValidator: promoValidator,
		emailSender:    emailSender,
	}
}

func (m *Manager) GetAll(ctx context.Context) ([]model.ExtendedRegistration, error) {
	regs, err := repository.ListRegistrations(ctx, m.postgresDB.QueryerContext())
	if err != nil {
		logger.FromCtx(ctx).Error("Failed to list registrations", zap.Error(err))
		return nil, err
	}
	return regs, nil
}

func (m *Manager) GetByToken(ctx context.Context, token string) (*model.ExtendedRegistration, error) {
	reg, err := repository.FindRegistrationByToken(ctx, m.postgresDB.QueryerContext(), token)
	if err != nil {
		logger.UnexpectedError(ctx, err).Error("Failed to get registration", zap.String("token", token))
		return nil, err
	}
	return reg, nil
}

func (m *Manager) GetByID(ctx context.Context, id int) (*model.ExtendedRegistration, error) {
	reg, err := repository.FindRegistrationByID(ctx, m.postgresDB.QueryerContext(), id)
	if err != nil {
		logger.UnexpectedError(ctx, err).Error("Failed to find registration.", zap.Error(err))
		return nil, err
	}
	return reg, nil
}

func (m *Manager) Update(ctx context.Context, registration *model.Registration) error {
	if err := m.postgresDB.WithTxx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		if err := repository.UpdateRegistration(ctx, tx, registration); err != nil {
			logger.UnexpectedError(ctx, err).Error("Failed to update registration.", zap.Error(err))
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (m *Manager) DeleteByID(ctx context.Context, id int) (*model.Registration, error) {
	reg, err := repository.DeleteRegistrationByID(ctx, m.postgresDB.QueryerContext(), id)
	if err != nil {
		logger.UnexpectedError(ctx, err).Error("Failed to delete registration.", zap.Int("reg_id", id), zap.Error(err))
		return nil, err
	}

	return reg, nil
}

func (m *Manager) SendPaymentNotifications(ctx context.Context) (resources.PaymentNotificationResponse, error) {
	var (
		sent        int
		finishedAll bool
	)
	if err := m.postgresDB.WithTxx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		regs, err := repository.ListRegistrationsWithoutNotification(ctx, tx)
		if err != nil {
			logger.FromCtx(ctx).Error("Failed to list registrations", zap.Error(err))
			return err
		}
		for i := range regs {
			notified, err := m.notifyRegistration(ctx, tx, &regs[i])
			if err != nil {
				logger.FromCtx(ctx).Error("Error sending notification", zap.Error(err), zap.Int("id", regs[i].ID))
				// We want to ignore the error and commit the tx
				return nil
			}
			if notified {
				sent++
			}
		}
		finishedAll = true
		return nil
	}); err != nil {
		return resources.PaymentNotificationResponse{}, err
	}
	return resources.PaymentNotificationResponse{
		Sent:        sent,
		FinishedAll: finishedAll,
	}, nil
}

func (m *Manager) CreateNew(ctx context.Context, req *resources.RegisterReq, eventID int) (*resources.RegisterResp, error) {
	var (
		err error
		res model.RegResult
	)

	token, err := uuid.NewUUID()
	if err != nil {
		logger.FromCtx(ctx).Error("Failed to generate uuid", zap.Error(err))
		return nil, errors.WithStack(err)
	}
	res.Token = token.String()

	if err := m.postgresDB.WithTxx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		var err error
		res.Event, err = m.findAndValidateEvent(ctx, tx, eventID, req.PromoCode)
		if err != nil {
			return err
		}

		promoKey, err := m.validatePromoCode(ctx, tx, req.PromoCode)
		if err != nil {
			return err
		}

		if res.RegisteredIDs, err = m.validateEventCapacity(ctx, tx, err, eventID, req); err != nil {
			return err
		}
		res.RegisteredDesc = m.registeredDescriptions(res.RegisteredIDs, res.Event)

		amount, discount := m.computePriceAndDiscount(res.RegisteredIDs, res.Event, promoKey != "")

		newNullString := func(s string) *string {
			if s != "" {
				return &s
			}
			return nil
		}

		newNullInt := func(val int) *int {
			if val != 0 {
				return &val
			}
			return nil
		}

		// Insert into registrations.
		reg, err := (&model.Registration{
			Name:               req.Child.Name,
			Surname:            req.Child.Surname,
			Gender:             req.Child.Gender,
			DateOfBirth:        req.Child.DateOfBirth,
			FinishedSchool:     req.Child.FinishedSchool,
			AttendedPrevious:   req.Child.AttendedPrevious,
			AttendedActivities: req.Membership.AttendedActivities,
			City:               req.Child.City,
			Pills:              req.Medicine.Pills,
			Problems:           req.Health.Problems,
			Notes:              req.Notes,
			ParentName:         req.Parent.Name,
			ParentSurname:      req.Parent.Surname,
			Email:              req.Parent.Email,
			Phone:              req.Parent.Phone,
			Amount:             amount,
			Token:              token.String(),
			PromoCode:          newNullString(promoKey),
			Discount:           newNullInt(discount),
		}).Create(ctx, tx)
		if err != nil {
			return err
		}
		res.Reg = *reg

		// Insert into signups.
		for _, dayID := range req.DayIDs {
			_, err := tx.ExecContext(ctx, `
			INSERT INTO signups(
				day_id,
				registration_id,
				state,
				created_at,
				updated_at
			) VALUES (
				$1,
				$2,
				$3,
				NOW(),
				NOW()
			) RETURNING *
		`, dayID, res.Reg.ID, "init")
			if err != nil {
				return errors.Wrap(err, "failed to create a signup")
			}
		}

		return nil
	}); err != nil {
		if errors.Cause(err) == ErrOverLimit {
			return resources.UnsuccessfulRegisterResponse(res.RegisteredIDs), nil
		}
		return nil, err
	}

	resp, err := m.sendConfirmationMail(ctx, res, err)
	if err != nil {
		return resp, err
	}

	return resources.SuccessRegisterResponse(res.RegisteredIDs, res.Token), nil
}

func (m *Manager) findAndValidateEvent(ctx context.Context, tx *sqlx.Tx, eventID int, promoCode *string) (*model.Event, error) {
	event, err := repository.FindEvent(ctx, tx, eventID)
	if err != nil {
		logger.UnexpectedError(ctx, err).Error("Failed to find event", zap.Int("event_id", eventID))
		return nil, err
	}

	if !event.Active && !(event.PromoRegistration && promoCode != nil) {
		return nil, ErrNotActive
	}
	return event, nil
}

func (m *Manager) validatePromoCode(ctx context.Context, tx *sqlx.Tx, promoCode *string) (string, error) {
	if promoCode == nil {
		return "", nil
	}

	// Validate token
	foundPromoCode, err := m.promoValidator.ValidateTokenWithQueryerContext(ctx, tx, *promoCode)
	if err != nil {
		logger.UnexpectedError(ctx, err).Error("Failed to validate promo code", zap.Error(err))
		return "", err
	}

	// Decrement remaining registrations
	if err := m.promoValidator.MarkTokenUsage(ctx, tx, foundPromoCode.Key); err != nil {
		return "", err
	}
	return foundPromoCode.Key, nil
}

func (m *Manager) validateEventCapacity(ctx context.Context, tx *sqlx.Tx, err error, eventID int, req *resources.RegisterReq) ([]int, error) {
	// List stats.
	stats, err := repository.GetStat(ctx, tx, eventID)
	if err != nil {
		logger.FromCtx(ctx).Error("Failed to get stats", zap.Error(err))
		return nil, err
	}

	var registeredIDs []int
	// Validate limits
	for i := range req.DayIDs {
		for _, stat := range stats {
			if stat.DayID != req.DayIDs[i] {
				continue
			}
			if (stat.GirlsCount + stat.BoysCount) >= stat.Capacity {
				break
			}
			if req.Child.Gender == "male" && stat.LimitBoys != nil {
				if stat.BoysCount >= *stat.LimitBoys {
					break
				}
			}

			if req.Child.Gender == "female" && stat.LimitGirls != nil {
				if stat.GirlsCount >= *stat.LimitGirls {
					break
				}
			}
			registeredIDs = append(registeredIDs, stat.DayID)
		}
	}

	if len(req.DayIDs) != len(registeredIDs) {
		if len(registeredIDs) == 0 {
			registeredIDs = []int{}
		}
		return nil, ErrOverLimit
	}
	return registeredIDs, nil
}

func (m *Manager) computePriceAndDiscount(registeredIDs []int, event *model.Event, applyPromo bool) (int, int) {
	// Compute price.
	var amount int
	for _, dID := range registeredIDs {
		for _, d := range event.Days {
			if dID != d.ID {
				continue
			}
			amount += d.Price
			break
		}
	}

	// Apply promo discount
	var discount int
	if applyPromo {
		discount = event.PromoDiscount
	}
	return amount, discount
}

func (m *Manager) registeredDescriptions(registeredIDs []int, event *model.Event) []string {
	var descriptions []string
	for _, dID := range registeredIDs {
		for _, d := range event.Days {
			if dID != d.ID {
				continue
			}
			descriptions = append(descriptions, d.Description)
			break
		}
	}
	return descriptions
}

func (m *Manager) sendConfirmationMail(ctx context.Context, res model.RegResult, err error) (*resources.RegisterResp, error) {
	pills := "-"
	if res.Reg.Pills != nil {
		pills = *res.Reg.Pills
	}

	restrictions := "-"
	if res.Reg.Problems != nil {
		restrictions = *res.Reg.Problems
	}
	var info string
	if res.Event.MailInfo != nil {
		info = *res.Event.MailInfo
	}
	var regInfo string
	if res.Event.Info != nil {
		regInfo = *res.Event.Info
	}

	payment, err := m.registrationToPaymentDetails(&res)
	if err != nil {
		logger.FromCtx(ctx).Error("failed to create payment data", zap.Error(err))
		return nil, err
	}

	// Send confirmation mail.
	if err := m.emailSender.ConfirmationMail(ctx, &templates.ConfirmationReq{
		Mail:          res.Reg.Email,
		ParentName:    res.Reg.ParentName,
		ParentSurname: res.Reg.ParentSurname,
		EventName:     res.Event.Title,
		Name:          res.Reg.Name,
		Surname:       res.Reg.Surname,
		Pills:         pills,
		Restrictions:  restrictions,
		Text:          res.Event.OwnerPhone + " " + res.Event.OwnerEmail,
		PhotoURL:      res.Event.OwnerPhoto,
		Sum:           res.Reg.AmountToPay(),
		Owner:         res.Event.OwnerName + " " + res.Event.OwnerSurname,
		Days:          res.RegisteredDesc,
		Info:          info,
		RegInfo:       regInfo,
		Payment:       payment,
	}); err != nil {
		logger.FromCtx(ctx).Error("Failed to send a confirmation mail.", zap.Error(err))
		return nil, err
	}
	return nil, nil
}

func (m *Manager) registrationToPaymentDetails(reg *model.RegResult) (templates.PaymentDetails, error) {
	link, err := payme.NewBuilder().
		IBAN(reg.Event.IBAN).
		Amount(reg.Reg.AmountToPay()).
		PaymentReference(reg.Event.PaymentReference).
		SpecificSymbol(reg.Reg.SpecificSymbol).
		Note(fmt.Sprintf("%s %s %s", reg.Event.Title, reg.Reg.Name, reg.Reg.Surname)).
		Build()
	if err != nil {
		return templates.PaymentDetails{}, err
	}

	details := templates.PaymentDetails{
		IBAN:             reg.Event.IBAN,
		PaymentReference: reg.Event.PaymentReference,
		SpecificSymbol:   reg.Reg.SpecificSymbol,
		Link:             link,
		QRCode:           "", // TODO
	}
	return details, nil
}

func (m *Manager) notifyRegistration(ctx context.Context, tx *sqlx.Tx, reg *model.ExtendedRegistration) (bool, error) {
	if reg.Payed == nil {
		return false, nil
	}
	amountToPay := reg.AmountToPay()
	payed := *reg.Payed
	if payed < amountToPay {
		logger.FromCtx(ctx).Info(
			"Lower amount payed for registration",
			zap.Int("id", reg.ID),
			zap.Int("payed", payed),
			zap.Int("amount_to_pay", amountToPay),
		)
	}
	if err := m.emailSender.NotificationMail(ctx, &templates.NotificationReq{
		Mail:      reg.Email,
		Name:      reg.Name,
		Surname:   reg.Surname,
		EventName: reg.Title,
		Payed:     payed,
	}); err != nil {
		return false, err
	}
	if err := repository.MarkRegistrationAsNotified(ctx, tx, reg.ID); err != nil {
		logger.FromCtx(ctx).Error("Failed to mark registration as notified", zap.Error(err))
		return false, err
	}
	return true, nil

}
