package pgsql_test

import (
	"context"
	"testing"
	"time"

	"github.com/apm-dev/vending-machine/domain"
	"github.com/apm-dev/vending-machine/pkg/pgsqlhelper"
	"github.com/apm-dev/vending-machine/user/data/pgsql"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UserRepoTestSuite struct {
	suite.Suite
	db    *gorm.DB
	purge pgsqlhelper.PurgeResourcesFunc
}

func (s *UserRepoTestSuite) SetupTest() {
	var err error
	s.db, s.purge, err = pgsqlhelper.NewPostgreContainer(pgsqlhelper.PgConfig{
		Version:  "14",
		Username: "admin",
		Password: "root",
		DB:       "vm_db",
	})
	if err != nil {
		panic(err)
	}

	err = s.db.AutoMigrate(&pgsql.User{})
	if err != nil {
		panic(err)
	}
}

func (s *UserRepoTestSuite) TearDownTest() {
	if err := s.purge(); err != nil {
		panic(err)
	}
}

func TestUserRepoTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepoTestSuite))
}

func (s *UserRepoTestSuite) TestInsert() {
	// arrange
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	user, err := domain.NewUser("uname", "passwd", "seller")
	user.Id = 5
	if err != nil {
		panic(err)
	}
	// action
	ur := pgsql.InitUserRepository(s.db)
	id, err := ur.Insert(ctx, *user)
	// assert
	s.NoError(err, "Inserting new user should not return error")
	s.NotEqual(5, id, "prefilled id should skip and generate unique one in db")
}
