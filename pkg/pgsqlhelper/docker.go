package pgsqlhelper

import (
	"log"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type PgConfig struct {
	Version  string
	Username string
	Password string
	DB       string
}

type PurgeResourcesFunc func() error

func NewPostgreContainer(conf PgConfig) (*gorm.DB, PurgeResourcesFunc, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Printf("Could not connect to docker: %s", err)
		return nil, nil, err
	}
	// pull image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        conf.Version,
		Env: []string{
			"POSTGRES_USER=" + conf.Username,
			"POSTGRES_PASSWORD=" + conf.Password,
			"POSTGRES_DB=" + conf.DB,
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.Cgroup = "fsync=off synchronous_commit=off archive_mode=off wal_level=minimal shared_buffers=512MB"
	})
	if err != nil {
		log.Printf("Could not start resource: %s", err)
		return nil, nil, err
	}

	var db *gorm.DB
	// exponential backoff-retry, because the application in the container
	// might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		db, err = gorm.Open(
			postgres.Open(
				"host=localhost"+
					" user="+conf.Username+
					" password="+conf.Password+
					" dbname="+conf.DB+
					" port="+resource.GetPort("5432/tcp")+
					" sslmode=disable",
			), &gorm.Config{
				NamingStrategy: schema.NamingStrategy{
					SingularTable: false,
				},
				SkipDefaultTransaction: false,
				PrepareStmt:            true,
			},
		)
		if err != nil {
			log.Printf("failed to connect database :%s", err)
			return err
		}
		sqldb, _ := db.DB()
		return sqldb.Ping()
	}); err != nil {
		log.Printf("Could not connect to docker: %s", err)
		return nil, nil, err
	}
	return db, func() error {
		err := pool.Purge(resource)
		if err != nil {
			log.Printf("Could not purge resource: %s", err)
			return err
		}
		return nil
	}, nil
}
