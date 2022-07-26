package docker_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestDocker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Docker Suite")
}

var Db *gorm.DB
var cleanupDocker func()

var _ = BeforeSuite(func() {
	// setup *gorm.Db with docker
	Db, cleanupDocker = setupGormWithDocker()
})

var _ = AfterSuite(func() {
	// cleanup resource
	cleanupDocker()
})

var _ = BeforeEach(func() {
	// clear db tables before each test
	err := Db.Exec(`DROP SCHEMA public CASCADE;CREATE SCHEMA public;`).Error
	Ω(err).To(Succeed())
})

const (
	dbName = "test"
	passwd = "test"
)

func setupGormWithDocker() (*gorm.DB, func()) {
	pool, err := dockertest.NewPool("")
	chk(err)

	runOpt := &dockertest.RunOptions{
		Repository: "postgres", // set image name
		Tag:        "14",       // set version
		Env: []string{
			"POSTGRES_PASSWORD=" + passwd, // default passwd
			"POSTGRES_DB=" + dbName,
		},
	}

	fnConfig := func(config *docker.HostConfig) {
		config.AutoRemove = true // set AutoRemove to true so that stopped container goes away by itself

		config.RestartPolicy = docker.NeverRestart() // don't restart container
	}

	resource, err := pool.RunWithOptions(runOpt, fnConfig)
	chk(err)
	// call clean up function to release resource
	fnCleanup := func() {
		err := resource.Close()
		chk(err)
	}

	conStr := fmt.Sprintf("host=localhost port=%s user=postgres dbname=%s password=%s sslmode=disable",
		resource.GetPort("5432/tcp"),
		dbName,
		passwd,
	)

	dialector := postgres.Open(conStr) // create dialector

	fnRetry := func() error {
		// no log output since there will be several errors during setup
		cfg := &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)}
		gdb, err := gorm.Open(dialector, cfg)
		if err != nil {
			return err
		}
		db, err := gdb.DB()
		if err != nil {
			return err
		}
		return db.Ping()
	}

	// retry until container is ready
	err = pool.Retry(fnRetry)
	chk(err)

	// container is ready, create and return *gorm.Db for testing
	db, err := gorm.Open(dialector, &gorm.Config{})
	chk(err)

	return db, fnCleanup
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
