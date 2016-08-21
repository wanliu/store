package store

type query struct {
	idx *index
}

type IterFunc func(obj interface{}) error

func (s *store) each(q *query, handle IterFunc) error {

}
