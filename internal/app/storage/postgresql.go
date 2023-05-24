package storage

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

// PGDB: create struct for DB
type PGDB struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewPostgresDB is a function that returns a new PostgresDB
func NewPostgresDB(param string, logger *zap.Logger) *PGDB {

	db, err := sql.Open("pgx", param)
	if err != nil {
		// If the error is a connection exception, try to reconnect
		//this code wrote for increment #13
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
			logger.Sugar().Errorf("Error opening database connection: %v, trying reconnect", err)
			//time delta for reconnect is 1, 3, 5 seconds
			waitTime := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}
			//try to reconnect
			for _, v := range waitTime {
				time.Sleep(v)
				db, err = sql.Open("pgx", param)
				if err == nil {
					//if connection is successful, break the loop
					break
				}
			}
		}
	}
	return &PGDB{db: db, logger: logger}
	//return &PGDB{db: db}
}

// Close is a function that closes db
func (pgdb *PGDB) Close() {
	err := pgdb.db.Close()
	if err != nil {
		pgdb.logger.Sugar().Errorf("Error closing database connection: %v", err)
	}
}

// PingContext is a function that checks the connection to the database
// used in http httpserver
// realization for increment #11
func (pgdb *PGDB) PingContext(ctx context.Context) error {
	err := pgdb.db.PingContext(ctx)
	pgdb.logger.Sugar().Info("PingContext successfully")
	return err
}

// Create is a function that creates a table in the database
func (pgdb *PGDB) Create(ctx context.Context) (err error) {
	query := `CREATE TABLE IF NOT EXISTS links (
		    id text not null PRIMARY KEY,
		    origin_link text not null
		    );`
	_, err = pgdb.db.ExecContext(ctx, query)
	if err != nil {
		log.Println(err)
		return err
	}
	pgdb.logger.Sugar().Info("Table metrics created successfully")
	//logger.Log.Sugar().Info("Table metrics opened successfully")
	return nil
}

// InsertMetric is a function that inserts a metric into the database
// used in write func
//func (pgdb *PGDB) InsertMetric(ctx context.Context, id string, mtype string, delta *int64, value *float64) (err error) {
//	query := `INSERT INTO metrics (id, mtype, delta, value)
//					VALUES ($1, $2, $3, $4)
//                    ON CONFLICT (id)
//                    DO UPDATE SET
//                    delta=excluded.delta,
//                    value=excluded.value;`
//	_, err = pgdb.db.ExecContext(ctx, query, id, mtype, delta, value)
//	if err != nil {
//		pgdb.logger.Sugar().Errorf("Error inserting metric: %v", err)
//		return err
//	}
//	return nil
//}

// UpdateMetric is a function that updates a metric in the database
// used in write func
//func (pgdb *PGDB) UpdateMetric(ctx context.Context, id string, mtype string, delta *int64, value *float64) (err error) {
//	query := `UPDATE metrics SET delta = $1, value = $2 WHERE id = $3 AND mtype = $4;`
//	_, err = pgdb.db.ExecContext(ctx, query, delta, value, id, mtype)
//	if err != nil {
//		pgdb.logger.Sugar().Errorf("Error updating metric: %v", err)
//		return err
//	}
//	return nil
//}
//
//// GetAllID is a function that returns all id from the database
//// used in write func
//func (pgdb *PGDB) GetAllID(ctx context.Context) (ids []string, err error) {
//	query := `SELECT id FROM metrics;`
//	rows, err := pgdb.db.QueryContext(ctx, query)
//	if err != nil {
//		pgdb.logger.Sugar().Errorf("Error getting all id: %v", err)
//		return nil, err
//	}
//	if rows.Err() != nil {
//		pgdb.logger.Sugar().Errorf("Error getting all id: %v", err)
//		return nil, err
//	}
//	defer func() {
//		err := rows.Close()
//		if err != nil {
//			pgdb.logger.Sugar().Errorf("Error closing rows: %v", err)
//		}
//	}()
//	for rows.Next() {
//		var id string
//		err = rows.Scan(&id)
//		if err != nil {
//			pgdb.logger.Sugar().Errorf("Error scanning rows: %v", err)
//			return nil, err
//		}
//		ids = append(ids, id)
//	}
//	return ids, nil
//}
//
//// GetValuesByID is a function that returns value by id from the database
//// used in autotests
//func (pgdb *PGDB) GetValueByID(ctx context.Context, req []byte) (res []byte, err error) {
//
//	var met metrics
//	err = json.Unmarshal(req, &met)
//	if err != nil {
//		pgdb.logger.Sugar().Errorf("Error unmarshaling request: %v", err)
//		return nil, err
//	}
//	query := `SELECT delta, value FROM metrics WHERE id = $1 AND mtype = $2;`
//	row := pgdb.db.QueryRowContext(ctx, query, met.ID, met.MType)
//	err = row.Scan(&met.Delta, &met.Value)
//	if err != nil {
//		pgdb.logger.Sugar().Errorf("Error scanning row: %v", err)
//		return nil, err
//	}
//	res, err = json.Marshal(met)
//	if err != nil {
//		pgdb.logger.Sugar().Errorf("Error marshaling response: %v", err)
//		return nil, err
//	}
//	return res, nil
//}
//
//// Read is a function that reads all metrics from the database
//// used in main module through interface
//func (pgdb *PGDB) Read(ctx context.Context) ([][]byte, error) {
//	var result [][]byte
//	rows, err := pgdb.db.QueryContext(ctx, "SELECT * FROM metrics;")
//	if err != nil {
//		pgdb.logger.Sugar().Errorf("Error reading from database: %v", err)
//		return nil, err
//	}
//	if rows.Err() != nil {
//		pgdb.logger.Sugar().Errorf("Error reading from database: %v", err)
//		return nil, err
//	}
//	var met metrics
//	for rows.Next() {
//		err = rows.Scan(&met.ID, &met.MType, &met.Delta, &met.Value)
//		if err != nil {
//			pgdb.logger.Sugar().Errorf("Error scan data rows from database: %v", err)
//			return nil, err
//		}
//		metJSON, err := json.Marshal(met)
//		if err != nil {
//			pgdb.logger.Sugar().Errorf("Error marshal data rows from database: %v", err)
//			return nil, err
//		}
//		result = append(result, metJSON)
//	}
//	pgdb.logger.Sugar().Info("Read from DB successfully")
//	return result, nil
//}
//
//// Write is a function that writes metrics to the database
//// used in main module through interface
//func (pgdb *PGDB) Write(ctx context.Context, met [][]byte) error {
//	pgdb.logger.Sugar().Info("Saving to DB...")
//	for _, metric := range met {
//		var met metrics
//		err := json.Unmarshal(metric, &met)
//		if err != nil {
//			pgdb.logger.Sugar().Errorf("Error unmarshal data: %v", err)
//			return err
//		}
//
//		tx, err := pgdb.db.BeginTx(ctx, nil)
//		if err != nil {
//			return err
//		}
//
//		err = pgdb.InsertMetric(ctx, met.ID, met.MType, met.Delta, met.Value)
//		if err != nil {
//			pgdb.logger.Sugar().Errorf("Error inserting data: %v, rollback transaction", err)
//			err2 := tx.Rollback()
//			if err2 != nil {
//				pgdb.logger.Sugar().Errorf("Error rollback transaction: %v", err2)
//			}
//			return err
//		}
//		err = tx.Commit()
//		if err != nil {
//			pgdb.logger.Sugar().Errorf("Error commit transaction: %v", err)
//		}
//
//	}
//	pgdb.logger.Sugar().Infof("Commited all transactions, data saved to DB successfully")
//	return nil
//
//}
