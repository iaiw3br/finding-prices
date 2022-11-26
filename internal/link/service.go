package link

import "context"

type Service struct {
	linkStore Store
}

func NewService(linkStore Store) Service {
	return Service{
		linkStore: linkStore,
	}
}

func (s Service) FindForSearch(ctx context.Context) ([]Search, error) {
	itemsForSearch, err := s.linkStore.ItemsForSearch(ctx)
	if err != nil {
		return nil, err
	}

	return itemsForSearch, nil
}
