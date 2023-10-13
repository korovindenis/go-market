package handler

type usecase interface {
}

type Handler struct {
	serverUsecase usecase
}

func New(u usecase) (*Handler, error) {
	return &Handler{
		serverUsecase: u,
	}, nil
}
