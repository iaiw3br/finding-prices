package price

import "context"

type Service struct {
	priceStore Store
}

func NewService(priceStore Store) Service {
	return Service{
		priceStore: priceStore,
	}
}

func (s Service) Create(ctx context.Context, cp CreatePrice) error {
	err := s.priceStore.Create(ctx, cp)
	if err != nil {
		return err
	}

	return nil
}
