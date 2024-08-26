package usecase

import (
	"context"
	"github.com/Back1ng/wbtech-0/internal/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository interface {
	GetAll(ctx context.Context) ([]entity.Order, error)
	Get(ctx context.Context, uid string) (*entity.Order, error)
	Store(ctx context.Context, o entity.Order) error
}

type OrderCache interface {
	Has(uid string) bool
	Get(uid string) (*entity.Order, error)
	Store(order entity.Order)
	Delete(uid string)
}

type OrderUsecase struct {
	pool  *pgxpool.Pool
	repo  OrderRepository
	cache OrderCache
}

func NewOrderUsecase(pool *pgxpool.Pool, repo OrderRepository, cache OrderCache) *OrderUsecase {
	return &OrderUsecase{
		pool:  pool,
		repo:  repo,
		cache: cache,
	}
}

func (uc *OrderUsecase) GetAllOrders(ctx context.Context) ([]entity.Order, error) {
	tx, err := uc.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	ctxTx := context.WithValue(ctx, "tx", tx)

	orders, err := uc.repo.GetAll(ctxTx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return orders, nil
}

func (uc *OrderUsecase) Get(ctx context.Context, uid string) (*entity.Order, error) {
	tx, err := uc.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	ctxTx := context.WithValue(ctx, "tx", tx)

	var order *entity.Order

	if uc.cache.Has(uid) {
		order, err = uc.cache.Get(uid)
		if err != nil {
			return nil, err
		}

		return order, nil
	}

	order, err = uc.repo.Get(ctxTx, uid)
	if err != nil {
		return nil, err
	}

	uc.cache.Store(*order)

	return order, nil
}

func (uc *OrderUsecase) StoreOrder(ctx context.Context, order entity.Order) error {
	tx, err := uc.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	ctxTx := context.WithValue(ctx, "tx", tx)

	err = uc.repo.Store(ctxTx, order)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	uc.cache.Store(order)

	return nil
}
