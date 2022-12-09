package persistence

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"

	"github.com/felipeflores/utils/ferrors"
	"github.com/felipeflores/utils/log"
)

// New creates a new persistence service
func New(config Config, logger log.Logger, serviceName string) (*Service, error) {

	sqltrace.Register(
		"pq",
		&pq.Driver{},
		sqltrace.WithServiceName(serviceName),
	)

	db, err := sqltrace.Open(
		"pq",
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			config.Username,
			config.Password,
			config.Host,
			config.Port,
			config.Schema,
		),
	)
	if err != nil {
		return nil, err
	}

	return &Service{
		DB:     db,
		logger: logger,
	}, nil
}

// Closer to closer an io operation
func (s *Service) Closer(ioc io.Closer) {
	err := ioc.Close()
	if err != nil {
		s.logger.Info(err.Error())
	}
}

// GenerateUUID generate an unique ID
func (s *Service) GenerateUUID() string {
	return uuid.New().String()
}

func (s *Service) GenerateUUIDWithoutPrefix() string {
	return uuid.New().String()
}

func (s *Service) GetJsonObjectFromString(item string) (map[string]interface{}, error) {
	js := make(map[string]interface{})
	err := json.Unmarshal([]byte(item), &js)
	return js, err
}

// HandleError handles an sql error
func (s *Service) HandleError(err error, msg string) error {
	if err == nil {
		return nil
	}

	if err == sql.ErrNoRows {
		return ferrors.NewNotFound(errors.Wrap(err, msg))
	}

	return ferrors.NewInternalServer(errors.Wrap(err, msg))
}

// HandleNoRowsAsNonError handles an sql error ignoring when no rows is found
func (s *Service) HandleNoRowsAsNonError(err error, msg string) error {
	if err == nil || err == sql.ErrNoRows {
		return nil
	}

	return ferrors.NewInternalServer(errors.Wrap(err, msg))
}

// HandleErrorWithTx handles an sql transaction error
func (s *Service) HandleErrorWithTx(tx *sql.Tx, err error, msg string) error {
	rollbackErr := tx.Rollback()
	if rollbackErr != nil {
		return ferrors.NewInternalServer(errors.Wrapf(rollbackErr, "Error rolling back tx of %s", msg))
	}
	return ferrors.NewInternalServer(errors.Wrap(err, msg))
}

// Service holds the database interface
type Service struct {
	DB     *sql.DB
	logger log.Logger
}

// Config holds the database configkids
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Schema   string
}

func toTimestamp(t time.Time) string {
	return t.Format(TimeFormat)
}

// NowUTC is a helper function to generate string-formatted timestamps using time.Now().UTC()
func NowUTC() string {
	return toTimestamp(time.Now().UTC())
}

// TimeFormat is the time format that should be used for timestamps that need to be persisted
const TimeFormat = "2006-01-02 15:04:05"
