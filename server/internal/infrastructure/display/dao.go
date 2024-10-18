package display

type ArrayXServerRepository struct {
	arr []XServer
}

func NewArrayXServerRepository() *ArrayXServerRepository {
	return &ArrayXServerRepository{}
}

func (r *ArrayXServerRepository) Populate(xServers []XServer) error {
	r.arr = xServers
	return nil
}
func (r *ArrayXServerRepository) GetAll() ([]XServer, error) {
	return r.arr, nil
}
