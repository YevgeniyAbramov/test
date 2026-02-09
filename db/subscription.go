package db

import (
	"database/sql"
	"errors"
	"strconv"
	"test/models"
)

func (db *DB) CreateSubscription(subscription *models.Subscription) error {
	query := `INSERT INTO subscriptions.subscription (service_name, price, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return db.conn.QueryRow(query, subscription.ServiceName, subscription.Price, subscription.UserID, subscription.StartDate, subscription.EndDate).Scan(&subscription.ID)
}

func (db *DB) GetSubscription(id int) (models.Subscription, error) {
	var subscription models.Subscription
	query := `SELECT * FROM subscriptions.subscription WHERE id = $1 AND deleted_at IS NULL`
	err := db.conn.Get(&subscription, query, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return subscription, errors.New("subscription not found")
		}
		return subscription, err
	}

	return subscription, nil
}

func (db *DB) DeleteSubscription(id int) error {
	query := `UPDATE subscriptions.subscription SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	result, err := db.conn.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("subscription not found")
	}

	return nil
}

func (db *DB) ListSubscriptions(page, limit int) ([]models.Subscription, int, error) {
	var subscriptions []models.Subscription

	var total int
	countQuery := `SELECT COUNT(*) FROM subscriptions.subscription WHERE deleted_at IS NULL`
	err := db.conn.Get(&total, countQuery)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	query := `SELECT * FROM subscriptions.subscription WHERE deleted_at IS NULL ORDER BY id asc LIMIT $1 OFFSET $2`
	err = db.conn.Select(&subscriptions, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return subscriptions, total, nil
}

func (db *DB) UpdateSubscription(id int, req *models.UpdateSubscriptionRequest) (models.Subscription, error) {
	query := `
	UPDATE subscriptions.subscription
	SET service_name = COALESCE($1, service_name),
	    price = COALESCE($2, price),
	    start_date = COALESCE($3, start_date),
	    end_date = COALESCE($4, end_date),
	    updated_at = NOW()
	WHERE id = $5 AND deleted_at IS NULL
	RETURNING *
	`

	var subscription models.Subscription

	err := db.conn.QueryRowx(query,
		req.ServiceName,
		req.Price,
		req.StartDate,
		req.EndDate,
		id,
	).StructScan(&subscription)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Subscription{}, errors.New("subscription not found")
		}
		return models.Subscription{}, err
	}

	return subscription, nil
}

func (db *DB) GetTotalCost(req *models.TotalCostRequest) (int, error) {
	query := `
		SELECT COALESCE(SUM(price), 0) 
		FROM subscriptions.subscription 
		WHERE deleted_at IS NULL 
		  AND start_date <= $1 
		  AND (end_date IS NULL OR end_date >= $2)`

	args := []interface{}{req.PeriodEnd, req.PeriodStart}

	if req.UserID != nil {
		query += " AND user_id = $3"
		args = append(args, *req.UserID)
	}
	if req.ServiceName != nil {
		query += " AND service_name = $" + strconv.Itoa(len(args)+1)
		args = append(args, *req.ServiceName)
	}

	var total int
	err := db.conn.Get(&total, query, args...)
	if err != nil {
		return 0, err
	}
	return total, nil
}
