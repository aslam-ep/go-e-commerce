package address

import (
	"context"
	"database/sql"
	"time"
)

type AddressRepository interface {
	Create(ctx context.Context, address *Address) (*Address, error)
	GetCountByUserID(ctx context.Context, userID int) (int, error)
	GetAll(ctx context.Context, userID int) (*[]Address, error)
	GetByID(ctx context.Context, id int, userID int) (*Address, error)
	Update(ctx context.Context, address *Address) (*Address, error)
	SetDefault(ctx context.Context, id int, userID int) error
	Delete(ctx context.Context, id int, userID int) error
}

type addressRepository struct {
	db *sql.DB
}

func NewAddressRepository(db *sql.DB) AddressRepository {
	return &addressRepository{
		db: db,
	}
}

func (r *addressRepository) Create(ctx context.Context, address *Address) (*Address, error) {
	insertQuery := `INSERT INTO addresses(user_id, address_line1, address_line2, postal_code, city, state, country) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, insertQuery,
		address.UserID,
		address.AddressLine1,
		address.AddressLine2,
		address.PostalCode,
		address.City,
		address.State,
		address.Country,
	).Scan(
		&address.ID,
		&address.CreatedAt,
		&address.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return address, nil
}

func (r *addressRepository) GetCountByUserID(ctx context.Context, userID int) (int, error) {
	var count int
	countQuery := `SELECT COUNT(id) as count FROM addresses WHERE user_id = $1;`

	err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *addressRepository) GetAll(ctx context.Context, userID int) (*[]Address, error) {
	selectByUserIDQuery := `SELECT id, user_id, address_line1, address_line2, postal_code, city, state, country, is_default, created_at, updated_at FROM addresses WHERE user_id = $1 ORDER BY id;`

	rows, err := r.db.QueryContext(ctx, selectByUserIDQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []Address
	for rows.Next() {
		var address Address
		if err := rows.Scan(
			&address.ID,
			&address.UserID,
			&address.AddressLine1,
			&address.AddressLine2,
			&address.PostalCode,
			&address.City,
			&address.State,
			&address.Country,
			&address.IsDefault,
			&address.CreatedAt,
			&address.UpdatedAt,
		); err != nil {
			return nil, err
		}

		addresses = append(addresses, address)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &addresses, nil
}

func (r *addressRepository) GetByID(ctx context.Context, id int, userID int) (*Address, error) {
	var address Address
	selectByIDAndUserIDQuery := `SELECT id, user_id, address_line1, address_line2, postal_code, city, state, country, is_default, created_at, updated_at FROM addresses WHERE id=$1 AND user_id = $2;`

	err := r.db.QueryRowContext(ctx, selectByIDAndUserIDQuery, id, userID).Scan(
		&address.ID,
		&address.UserID,
		&address.AddressLine1,
		&address.AddressLine2,
		&address.PostalCode,
		&address.City,
		&address.State,
		&address.Country,
		&address.IsDefault,
		&address.CreatedAt,
		&address.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &address, nil
}

func (r *addressRepository) Update(ctx context.Context, address *Address) (*Address, error) {
	address.UpdatedAt = time.Now()
	updateQuery := `UPDATE addresses SET address_line1 = $1, address_line2 = $2, postal_code = $3, city = $4, state = $5, country = $6, updated_at = $7 WHERE id = $8 AND user_id = $9;`

	_, err := r.db.ExecContext(ctx, updateQuery,
		address.AddressLine1,
		address.AddressLine2,
		address.PostalCode,
		address.City,
		address.State,
		address.Country,
		address.UpdatedAt,
		address.ID,
		address.UserID,
	)

	if err != nil {
		return nil, err
	}

	return address, nil
}

func (r *addressRepository) SetDefault(ctx context.Context, id int, userID int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	unsetQuery := `UPDATE addresses SET is_default = false WHERE user_id = $1;`
	_, err = tx.ExecContext(ctx, unsetQuery, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	setQuery := `UPDATE addresses SET is_default = true WHERE id = $1 AND user_id = $2;`
	_, err = tx.ExecContext(ctx, setQuery, id, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *addressRepository) Delete(ctx context.Context, id int, userID int) error {
	deleteQuery := `DELETE FROM addresses WHERE id = $1 AND user_id = $2;`

	_, err := r.db.ExecContext(ctx, deleteQuery, id, userID)

	return err
}
