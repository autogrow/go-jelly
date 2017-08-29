package mydata

import "fmt"

// Query represents a API query that provides chainable method calls ended by
// a call to Points()
type Query struct {
	env     string
	metric  string
	name    string
	days    int
	conn    *Conn
	LastURI string
}

// Metric sets the metric to get data points for
func (q *Query) Metric(metric string) *Query {
	q.metric = metric
	return q
}

// Days sets the number of days ago to get data for
func (q *Query) Days(days int) *Query {
	q.days = days
	return q
}

// Points will retrieve the data points from the API with the given
func (q *Query) Points() (Points, error) {
	if q.metric == "" {
		return Points{}, fmt.Errorf("no metric set")
	}

	uri := fmt.Sprintf("%s/%s/%s?name=%s&days=%d", q.conn.URL, q.env, q.metric, q.name, q.days)
	q.LastURI = uri
	return q.conn.getPoints(uri)
}

// All will populate a given records object with the combined metrics for the given environment
func (q *Query) All(object interface{}) error {
	uri := fmt.Sprintf("%s/%s?name=%s&days=%d", q.conn.URL, q.env, q.name, q.days)
	q.LastURI = uri

	return q.conn.getRecord(uri, object)
}
