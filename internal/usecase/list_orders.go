package usecase

import "github.com/soares-t-o/clean-arch/internal/entity"

type ListOrdersInput struct{}

type OrderOutput struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type ListOrdersOutputDTO struct {
	Orders []OrderOutput `json:"orders"`
}

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrdersUseCase(OrderRepository entity.OrderRepositoryInterface) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: OrderRepository,
	}
}

func (c *ListOrdersUseCase) Execute(input ListOrdersInput) (ListOrdersOutputDTO, error) {
	orders, err := c.OrderRepository.ListOrders()

	if err != nil {
		return ListOrdersOutputDTO{}, err
	}

	var ordersOutput []OrderOutput

	for _, order := range orders {
		ordersOutput = append(ordersOutput, OrderOutput{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.Price + order.Tax,
		})
	}
	return ListOrdersOutputDTO{
		Orders: ordersOutput,
	}, nil

}
