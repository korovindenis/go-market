package usecases

type storage interface {
}

type Usecases struct {
}

func New(storage storage) (*Usecases, error) {
	return &Usecases{}, nil
}
