package repository

import (
	"context"
	"errors"
	"github.com/Back1ng/wbtech-0/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepo struct {
	conn *pgxpool.Pool
}

func NewOrdersRepo(conn *pgxpool.Pool) *OrderRepo {
	return &OrderRepo{
		conn: conn,
	}
}

func (or *OrderRepo) GetAll(ctx context.Context) ([]entity.Order, error) {
	tx, ok := ctx.Value("tx").(pgx.Tx)
	if !ok {
		return nil, errors.New("no pgx transaction found")
	}

	sql := "SELECT order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shard_key, sm_id, date_created, oof_shard FROM orders"
	rows, err := tx.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entity.Order

	for rows.Next() {
		var o entity.Order
		if err := rows.Scan(
			&o.OrderUID, &o.TrackNumber, &o.Entry, &o.Delivery, &o.Payment, &o.Items,
			&o.Locale, &o.InternalSignature, &o.CustomerId, &o.DeliveryService, &o.ShardKey, &o.SmId,
			&o.DateCreated, &o.OofShard); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	return orders, nil
}

func (or *OrderRepo) Get(ctx context.Context, uid string) (*entity.Order, error) {
	tx, ok := ctx.Value("tx").(pgx.Tx)
	if !ok {
		return nil, errors.New("no pgx transaction found")
	}

	sql := "SELECT order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shard_key, sm_id, date_created, oof_shard FROM orders WHERE order_uid = $1 LIMIT 1"
	row := tx.QueryRow(ctx, sql, uid)

	var o entity.Order
	if err := row.Scan(&o.OrderUID, &o.TrackNumber, &o.Entry, &o.Delivery, &o.Payment, &o.Items,
		&o.Locale, &o.InternalSignature, &o.CustomerId, &o.DeliveryService, &o.ShardKey, &o.SmId,
		&o.DateCreated, &o.OofShard); err != nil {
		return nil, err
	}

	return &o, nil
}

func (or *OrderRepo) Store(ctx context.Context, o entity.Order) error {
	tx, ok := ctx.Value("tx").(pgx.Tx)
	if !ok {
		return errors.New("no pgx transaction found")
	}

	sql := "INSERT INTO orders (order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shard_key, sm_id, date_created, oof_shard) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)"
	commandTag, err := tx.Exec(ctx, sql, o.OrderUID, o.TrackNumber, o.Entry, o.Delivery, o.Payment, o.Items, o.Locale,
		o.InternalSignature, o.CustomerId, o.DeliveryService, o.ShardKey, o.SmId, o.DateCreated, o.OofShard)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows inserted")
	}

	return nil
}
