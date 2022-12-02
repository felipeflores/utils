package dbmigrate

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"go.uber.org/zap"

	// postgres import is needed by migrate to connect to database.
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	// ErrDirtyMigration is returned when the migrate status is dirty.
	ErrDirtyMigration = errors.New("migration is dirty, needs a manual fix")
)

// Config is used to receive all parameters to apply a migration.
type Config struct {
	Host, Port   string
	User, Pass   string
	Database     string
	Directory    string
	Logger       *zap.Logger
	ForceVersion int
}

// Up users a PostgreSQL connection to apply all available migrations.
// The argument opts is an array of options.
// The opts format must by in key=value format.
// Example: sslmode=disable connect_timeout=5
func Up(c Config, opts ...string) (err error) {

	var query string
	for i := range opts {
		if i == 0 {
			query = opts[0]
			continue
		}

		query = fmt.Sprintf("%s&%s", query, opts[i])
	}

	pgdsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?%s",
		c.User,
		c.Pass,
		c.Host,
		c.Port,
		c.Database,
		query,
	)

	m, err := migrate.New(
		fmt.Sprintf("file://%s", c.Directory),
		pgdsn,
	)
	if err != nil {
		return err
	}
	defer func() {
		se, de := m.Close()
		if se != nil {
			if err != nil {
				err = fmt.Errorf("%s; %s", err, se)
			} else {
				err = se
			}
		}
		if de != nil {
			if err != nil {
				err = fmt.Errorf("%s; %s", err, se)
			} else {
				err = de
			}
		}
	}()

	lf := map[string]string{
		"migrate_pg_connection_query": query,
		"migrate_pg_host":             c.Host,
	}

	c.Logger.Info(fmt.Sprintf("[dbmigration] ForceVersion is %d", c.ForceVersion))
	if c.ForceVersion > 0 {
		err = m.Force(c.ForceVersion)
		c.Logger.Info(fmt.Sprintf("[dbmigration] Forced Version is %d error detail %v", c.ForceVersion, err))
		if err != nil {
			return err
		}
	}

	mVersion, dirty, err := m.Version()
	if err != nil {
		if err != migrate.ErrNilVersion {
			return err
		}
	}

	if dirty {
		return ErrDirtyMigration
	}

	lf["migration_version"] = fmt.Sprintf("%v", mVersion)
	lf["migration_dirty"] = fmt.Sprintf("%v", dirty)
	c.Logger.Info("current migration status.")

	if err = m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return err
		}
		return nil
	}

	mVersion, dirty, err = m.Version()
	if err != nil {
		return err
	}

	lf["migration_version"] = fmt.Sprintf("%v", mVersion)
	lf["migration_dirty"] = fmt.Sprintf("%v", dirty)
	c.Logger.Info("migration applied")

	return nil
}
