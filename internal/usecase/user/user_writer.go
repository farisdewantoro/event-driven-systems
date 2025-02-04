package user

import (
	"context"
	asynqModels "eventdrivensystem/internal/models/asynq"
	notificationModels "eventdrivensystem/internal/models/notification"
	outboxModels "eventdrivensystem/internal/models/outbox"
	userModels "eventdrivensystem/internal/models/user"

	"time"

	"eventdrivensystem/pkg/errors"
	"eventdrivensystem/pkg/util"
)

type UserUsecaseWriter interface {
	CreateUser(ctx context.Context, param *userModels.CreateUserParam) error
}

func (u *UserUsecase) CreateUser(ctx context.Context, param *userModels.CreateUserParam) error {
	dbTx := u.userDomain.BeginTx(ctx)
	var (
		err error
		now = time.Now()
	)

	defer func() {
		if tmpErr := util.FirstNotNil(recover(), err); tmpErr != nil {
			if tmpErr != err {
				u.log.ErrorWithContext(ctx, tmpErr)
				return
			}

			if errRollback := dbTx.Rollback().Error; errRollback != nil {
				u.log.ErrorWithContext(ctx, errRollback)
				return
			}
		} else {
			if errCommit := dbTx.Commit().Error; errCommit != nil {
				u.log.ErrorWithContext(ctx, errCommit)
				err = errors.ErrSQLTx
				return
			}
		}
	}()

	dbOptions := util.DbOptions{
		Transaction: dbTx,
	}
	user, err := u.userDomain.CreateUser(ctx, param.ToDomain(), dbOptions)

	if err != nil {
		return err
	}

	pNotif := notificationModels.Notification{
		Status:  notificationModels.NotificationStatusPending,
		UserID:  user.ID,
		Type:    notificationModels.NotificationTypeUserRegistration,
		Message: notificationModels.NotificationMessageUserRegistration,
	}

	notif, err := u.notificationDomain.CreateNotification(ctx, &pNotif, dbOptions)
	if err != nil {
		return err
	}

	asynqPayload := asynqModels.AsynqSendNotificationPayload{
		UserID:           user.ID,
		NotificationID:   notif.ID,
		NotificationType: notif.Type,
	}

	asynqJson, err := asynqPayload.ToJSON()
	if err != nil {
		err = errors.ErrParseJsonOutbox
		return err
	}

	outbox := outboxModels.Outbox{
		Payload:         asynqJson,
		DestinationType: outboxModels.OutboxDestinationTypeAsynq,
		EventType:       asynqModels.AsynqTaskSendEmailNotification,
		ExecuteAt:       now,
	}

	err = u.outboxDomain.CreateOutbox(ctx, &outbox, dbOptions)
	if err != nil {
		return err
	}

	return nil
}
