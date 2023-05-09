package database

import (
	"database/sql"

	"github.com/soares-t-o/clean-arch/internal/entity"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	db.Exec("CREATE TABLE IF NOT EXISTS orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id))")
	return &OrderRepository{Db: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	stmt, err := r.Db.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) GetTotal() (int, error) {
	var total int
	err := r.Db.QueryRow("Select count(*) from orders").Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *OrderRepository) ListOrders() ([]entity.Order, error) {
	var entityOrders []entity.Order
	rows, err := r.Db.Query("SELECT * FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order entity.Order
		err := rows.Scan(&order.ID, &order.Price, &order.Tax, &order.FinalPrice)
		if err != nil {
			return nil, err
		}
		entityOrders = append(entityOrders, order)
	}
	return entityOrders, nil
}
