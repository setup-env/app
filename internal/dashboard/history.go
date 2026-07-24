package dashboard

type History struct {
	limit  int
	values []float64
}

func NewHistory(limit int) History {
	if limit < 1 {
		limit = 1
	}
	return History{limit: limit}
}

func (h History) Add(value float64) History {
	values := make([]float64, 0, min(h.limit, len(h.values)+1))
	start := 0
	if len(h.values) >= h.limit {
		start = len(h.values) - h.limit + 1
	}
	values = append(values, h.values[start:]...)
	values = append(values, value)
	h.values = values
	return h
}

func (h History) Values() []float64 {
	return append([]float64(nil), h.values...)
}

func (h History) Len() int {
	return len(h.values)
}
