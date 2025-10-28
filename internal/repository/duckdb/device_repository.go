package duckdb

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/leandronowras/device-api/internal/device"
	"github.com/leandronowras/device-api/internal/repository"
)

type deviceRepo struct {
	db *sql.DB
}

func NewDeviceRepository(db *sql.DB) repository.DeviceRepository {
	repo := &deviceRepo{db: db}
	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS devices (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		brand TEXT NOT NULL,
		state TEXT NOT NULL,
		creation_time TIMESTAMP NOT NULL
	)`)
	return repo
}

func (r *deviceRepo) Save(ctx context.Context, d *device.Device) (*device.Device, error) {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO devices (id, name, brand, state, creation_time) VALUES (?, ?, ?, ?, ?)`,
		d.ID(), d.Name(), d.Brand(), d.State(), d.CreationTime())
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (r *deviceRepo) FindByID(ctx context.Context, id string) (*device.Device, error) {
	var idVal, name, brand, state string
	var creationTime string

	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, brand, state, creation_time FROM devices WHERE id = ?`, id).
		Scan(&idVal, &name, &brand, &state, &creationTime)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		return nil, err
	}

	ct, err := parseTime(creationTime)
	if err != nil {
		return nil, err
	}

	return device.NewWithID(idVal, name, brand, state, ct)
}

func (r *deviceRepo) FindAll(ctx context.Context, brand, state *string) ([]*device.Device, error) {
	query := `SELECT id, name, brand, state, creation_time FROM devices`
	clauses := []string{}
	args := []any{}

	if brand != nil && *brand != "" {
		clauses = append(clauses, "LOWER(brand) = LOWER(?)")
		args = append(args, *brand)
	}
	if state != nil && *state != "" {
		clauses = append(clauses, "LOWER(state) = LOWER(?)")
		args = append(args, *state)
	}
	if len(clauses) > 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}
	query += " ORDER BY creation_time DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*device.Device
	for rows.Next() {
		var id, name, brand, state, creationTime string
		if err := rows.Scan(&id, &name, &brand, &state, &creationTime); err != nil {
			return nil, err
		}
		ct, err := parseTime(creationTime)
		if err != nil {
			return nil, err
		}
		d, err := device.NewWithID(id, name, brand, state, ct)
		if err != nil {
			return nil, err
		}
		list = append(list, d)
	}
	return list, rows.Err()
}

func (r *deviceRepo) Update(ctx context.Context, d *device.Device) (*device.Device, error) {
	res, err := r.db.ExecContext(ctx,
		`UPDATE devices SET name = ?, brand = ?, state = ? WHERE id = ?`,
		d.Name(), d.Brand(), d.State(), d.ID())
	if err != nil {
		return nil, err
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		return nil, sql.ErrNoRows
	}
	return d, nil
}

func (r *deviceRepo) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM devices WHERE id = ?`, id)
	if err != nil {
		return err
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func parseTime(s string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02 15:04:05",
		"2006-01-02 15:04:05.999999999",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("unable to parse time: " + s)
}
