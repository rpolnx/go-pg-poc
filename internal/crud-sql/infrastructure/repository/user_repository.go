package repository

import (
	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"rpolnx.com.br/crud-sql/internal/crud-sql/domain/model"
	port "rpolnx.com.br/crud-sql/internal/crud-sql/domain/ports"
	"time"
)

type userRepository struct {
	Db *pg.DB
}

func (svc *userRepository) GetAllUsersOut() ([]*model.User, error) {
	users := make([]*model.User, 0)

	err := svc.Db.Model(&users).
		Where("deleted_at IS NULL").
		Select()

	return users, err
}

func (svc *userRepository) GetOneUserOut(id *int64) (*model.User, error) {
	user := &model.User{
		Id: id,
	}

	err := svc.Db.Model(user).
		WherePK().
		Select()

	return user, err
}

func (svc *userRepository) CreateUserOut(user *model.User) (*int64, error) {
	user.Id = nil

	insert, err := svc.Db.
		Model(user).
		Insert()

	if err != nil {
		return nil, err
	}

	affected := insert.RowsAffected()
	logrus.Info(affected)

	return user.Id, err
}

func (svc *userRepository) UpdateUserOut(id *int64, user *model.User) (*int64, error) {
	now := time.Now()
	user.Id = id
	user.UpdatedAt = &now

	updated, err := svc.Db.
		Model(user).
		WherePK().
		UpdateNotZero()

	if err != nil {
		return nil, err
	}

	affected := updated.RowsAffected()

	logrus.Infof("Afftected %d", affected)

	if affected == 0 {
		return nil, errors.New("not found")
	}

	return user.Id, err
}

func (svc *userRepository) DeleteUserOut(id *int64) error {
	now := time.Now()
	user := &model.User{
		Id:        id,
		UpdatedAt: &now,
		DeletedAt: &now,
	}

	deleted, err := svc.Db.
		Model(user).
		WherePK().
		UpdateNotZero()

	if err != nil {
		return err
	}

	affected := deleted.RowsAffected()
	logrus.Info(affected)

	return err
}

func NewUserRepository(db *pg.DB) port.UserPort {
	return &userRepository{
		Db: db,
	}
}
