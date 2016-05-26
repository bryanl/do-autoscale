package autoscale

import (
	"autoscale/gen"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"pkg/backoff"
	"pkg/ctxutil"
	"strings"
	"time"

	"golang.org/x/net/context"

	// import db drivers
	"github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
	_ "github.com/mattes/migrate/driver/postgres"
	"github.com/mattes/migrate/migrate"
)

// NewDB creates a db connection.
func NewDB(ctx context.Context, user, password, addr, database string) (*sql.DB, error) {
	log := ctxutil.LogFromContext(ctx)

	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   addr,
		Path:   database,
	}

	v := url.Values{}
	v.Set("sslmode", "disable")

	u.RawQuery = v.Encode()
	dbURL := u.String()

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.WithError(err).Error("could not open database")
		return nil, err
	}

	var attempt int
	for {
		sleepTime, err := backoff.DefaultPolicy.Duration(attempt)
		if err != nil {
			return nil, err
		}

		time.Sleep(sleepTime)

		log.WithField("attempt", attempt).Info("connecting to db server")
		err = db.Ping()
		if err == nil {
			break
		}

		attempt++
		log.WithError(err).Warn("backing off because db didn't respond to ping")
	}

	log.Info("database is ready")

	env := ctxutil.StringFromContext(ctx, "env")

	if env != "development" {
		if err := migrateDatabase(log, dbURL); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func migrateDatabase(log *logrus.Entry, dbURL string) error {
	migrationDir := "db/migrations"

	tmpDir, err := ioutil.TempDir("", "as-migrations")
	if err != nil {
		log.WithError(err).Error("could not create temp directory")
		return err
	}
	defer os.RemoveAll(tmpDir)

	files, err := gen.AssetDir(migrationDir)
	if err != nil {
		log.WithError(err).Error("unable to find db files")
		return err
	}

	tmpMigrationDir := filepath.Join(tmpDir, migrationDir)
	if err := os.MkdirAll(tmpMigrationDir, 0700); err != nil {
		log.WithError(err).Error("unable to create temp migration directory")
		return err
	}

	for _, file := range files {
		assetPath := filepath.Join(migrationDir, file)
		contents, err := gen.Asset(assetPath)
		if err != nil {
			log.WithError(err).WithField("file-name", assetPath).Error("could not retrieve asset contents")
			return err
		}

		fn := filepath.Join(tmpMigrationDir, file)

		if err := ioutil.WriteFile(fn, contents, 0600); err != nil {
			log.WithError(err).WithField("file-name", assetPath).Error("could not write asset contents")
			return err
		}
	}

	log.WithField("db-url", dbURL).Info("performing migration")
	theErrors, ok := migrate.UpSync(dbURL, tmpMigrationDir)
	if !ok {
		err := fmt.Errorf("db migration failed")

		errorMsgs := []string{}
		for _, err := range theErrors {
			errorMsgs = append(errorMsgs, err.Error())
		}

		log.
			WithError(err).
			WithField("errors", strings.Join(errorMsgs, "|")).
			Error("could not migrate database")
		return err
	}

	return nil
}
