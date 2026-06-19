package repository

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/vivaldi7/golang_code/subscription_service/innternal/models"
	"github.com/vivaldi7/golang_code/subscription_service/logger"
)

const (
	tableName         = "subscriptions"
	idColumn          = "id"
	serviceNameColumn = "service_name"
	priceColumn       = "price"
	userIDColumn      = "user_id"
	startDateColumn   = "start_date"
	endDateColumn     = "end_date"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepositoty(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(ctx context.Context, sub *models.Subscription) (int, error) {
	//Строим DELETE запрос
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(serviceNameColumn, priceColumn, userIDColumn, startDateColumn, endDateColumn).
		Values(sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate).
		Suffix("RETURNING id")

	//Преобразуем в sql и аргументы
	query, args, err := builder.ToSql()
	if err != nil {
		logger.Log.Error("Failed to build create query", "error:", err)
		return 0, err
	}

	var id int

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		logger.Log.Error("Failed to insert create query", "error:", err)
		return 0, err
	}

	logger.Log.Info("Subscription created", "id", id, "user_id", sub.UserID)

	return id, nil
}

func (r *SubscriptionRepository) Get(ctx context.Context, id int) (*models.Subscription, error) {
	//Строим DELETE запрос
	builder := sq.Select(idColumn, serviceNameColumn, priceColumn, startDateColumn, endDateColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{idColumn: id}).
		Limit(1)

	//Преобразуем в sql и аргументы
	query, args, err := builder.ToSql()
	if err != nil {
		logger.Log.Error("Failed to build get query", "error:", err)
		return nil, err
	}

	var sub models.Subscription

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.StartDate, &sub.EndDate)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		logger.Log.Error("Failed to get subscription", "error:", err, "id", id)
		return nil, err
	}

	return &sub, nil
}

func (r *SubscriptionRepository) List(ctx context.Context, page, pageSize int64) ([]models.Subscription, error) {
	//Строим DELETE запрос
	builder := sq.Select(idColumn, serviceNameColumn, priceColumn, startDateColumn, endDateColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Limit(uint64(page)).
		Offset(uint64(pageSize))

	//Преобразуем в sql и аргументы
	query, args, err := builder.ToSql()
	if err != nil {
		logger.Log.Error("Failed to build list query", "error:", err)
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		logger.Log.Error("Failed to list subscription", "error:", err)
		return nil, err
	}
	defer rows.Close()

	var subscriptions []models.Subscription

	for rows.Next() {
		var sub models.Subscription
		err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.StartDate, &sub.EndDate)
		if err != nil {
			logger.Log.Error("Failed to scan subscription", "error:", err)
			return nil, err
		}
		subscriptions = append(subscriptions, sub)
	}

	return subscriptions, nil

}

func (r *SubscriptionRepository) Delete(ctx context.Context, id int) error {
	//Строим DELETE запрос
	builder := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id})

	//Преобразуем в sql и аргументы
	query, args, err := builder.ToSql()
	if err != nil {
		logger.Log.Error("Failed to build delete query", "error:", err)
		return err
	}

	//Выполняем
	rez, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		logger.Log.Error("Failed to execute delete qurty", "error:", err, "id", id)
		return err
	}

	//Проверка на удаление элемента
	rowsAffected, err := rez.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		logger.Log.Warn("No subscription found to delete", "id", id)
		return sql.ErrNoRows
	}

	logger.Log.Info("Subscription deleted", "id", id)
	return nil

}

func (r *SubscriptionRepository) Update(ctx context.Context, id int, updates map[string]interface{}) error {
	//Проверка на наличие того что нужно Апдейтить
	if len(updates) == 0 {
		return nil
	}

	//Строим DELETE запрос
	builder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id})

	for key, value := range updates {
		builder = builder.Set(key, value)
	}

	//Преобразуем в sql и аргументы
	query, args, err := builder.ToSql()
	if err != nil {
		logger.Log.Error("Failed to build update query", "error:", err)
		return err
	}

	//Выполняем
	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		logger.Log.Error("Failed to execute update qurty", "error:", err, "id", id)
		return err
	}

	logger.Log.Info("Subscription updated", "id", id)
	return nil
}

func (r *SubscriptionRepository) GetTotalCost(ctx context.Context, userID uuid.UUID, serviceName, startDate, endDate string) (int, error) {

	builder := sq.Select("COALESCE(SUM(priceColumn), 0)").
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{userIDColumn: userID}).
		Where(sq.Eq{serviceNameColumn: serviceName}).
		Where(sq.GtOrEq{startDateColumn: startDate}).
		Where(sq.Or{
			sq.Eq{endDateColumn: nil},
			sq.LtOrEq{endDateColumn: endDate},
		})

	query, args, err := builder.ToSql()
	if err != nil {
		logger.Log.Error("Failed to build GetTotalCost query", "error:", err)
		return 0, err
	}

	var total int

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&total)
	if err != nil {
		logger.Log.Error("Failed to execute GetTotalCost qurty", "error:", err)
		return 0, err
	}

	return total, nil
}
