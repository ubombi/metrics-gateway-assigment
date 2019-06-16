package clickhouse

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/ubombi/timeseries/api"

	_ "github.com/kshvakov/clickhouse"
)

var (
	CommitInterval = flag.Duration("clickhouseCommitInterval", 10*time.Second, "Events will be comitted to Clickhouse each time interval, or if batch is already bigger than x")
	Workers        = flag.Int("clickhouseWorkers", 4, "Insert workers")
	DSN            = flag.String("clickhouseDSN", "native://127.0.0.1:9000?debug=true&block_size=100000", "clickhouse connection string")
)

// NewStorage uses command line args instead of parameters
func NewStorage(ctx context.Context) *Storage {
	var err error

	s := Storage{}
	s.ctx, s.cancel = context.WithCancel(ctx)
	s.input = make(chan api.Event, 10000)

	s.db, err = connect(*DSN)
	if err != nil {
		log.Fatal(err)
	}

	if _, err = s.db.Exec(ddlQuery); err != nil {
		log.Fatal(err)
	}
	s.workers = *Workers
	s.errors = make(chan error)

	return &s
}

type Storage struct {
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
	db      *sql.DB
	input   chan api.Event
	workers int
	errors  chan error
}

// Shutdown stops receiving events to store, inserts existing quque and exits gracefuly.
func (s *Storage) Shutdown() {
	select {
	case <-s.ctx.Done():
		return
	default:
	}

	s.cancel()
	s.wg.Wait()
}

// Store enquque event to be processed by workers. Returns error if called after Shutdown
func (s *Storage) Store(e api.Event) error {
	select {
	case <-s.ctx.Done():
		return errors.New("storage is shutted down")
	default:
	}
	s.input <- e
	return nil
}

// Start workers. Blocka untill shutdown is called.
func (s *Storage) Start() error {
	s.wg.Add(s.workers)
	defer s.wg.Wait()

	go func() {
		<-s.ctx.Done()
		close(s.input)
	}()

	for i := 0; i < s.workers; i++ {
		go s.work()
	}

	for {
		select {
		case <-s.ctx.Done():
			return nil
		case e := <-s.errors:
			// Some kind of errors can be handled here and/or ignored with continue

			s.Shutdown()
			return e
		}
	}
}

// checkErr is a helper function that redirect errors to errors chanel. After context cancelation errors are ignored to avoid deadlocks
func (s *Storage) checkErr(err error) {
	if err == nil {
		return
	}
	select {
	case <-s.ctx.Done():
	case s.errors <- err:
	}
}

// work serialize and send events from input quque with own chickhouse connection and commit them with defined intervals
func (s *Storage) work() {
	defer s.wg.Done()
	defer log.Print("worker done")

	tx, stmt, err := makeInsert(s.db)
	s.checkErr(err)
	defer tx.Commit()

	commitTicker := time.NewTicker(*CommitInterval)
	defer commitTicker.Stop()

	uncomitted := 0
	for {
		select {
		case <-commitTicker.C:
			if uncomitted == 0 {
				continue
			}
			s.checkErr(tx.Commit())
			tx, stmt, err = makeInsert(s.db)
			s.checkErr(err)

		case Event, ok := <-s.input:
			if !ok {
				// exit ony after input quque is drained
				return
			}

			e := convertToInternal(Event)
			err = execEventInsert(stmt, &e)
			s.checkErr(err)

			uncomitted++
		}
	}
}

// execEventInsert wraps exec query to make select in `work` more readable
func execEventInsert(stmt *sql.Stmt, e *event) (err error) {
	_, err = stmt.Exec(
		e.EventType,
		e.Ts,
		e.StringParamNames,
		e.StringParamValues,
		e.IntParamNames,
		e.IntParamValues,
		e.FloatParamNames,
		e.FloatParamValues,
		e.UID,
	)
	return
}
