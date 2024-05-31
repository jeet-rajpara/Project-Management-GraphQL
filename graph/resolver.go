package graph

import (
	"io"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"

	_ "github.com/99designs/gqlgen/graphql/introspection"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{}

// type Timestamp time.Time

// func (t *Timestamp) UnmarshalGQL(v interface{}) error {
// 	valueStr, ok := v.(string)
// 	if !ok {
// 		return errors.New("timestamp must be a string")
// 	}
// 	parsedTime, err := time.Parse(time.RFC3339, valueStr)
// 	if err != nil {
// 		return err
// 	}
// 	*t = Timestamp(parsedTime)
// 	return nil
// }

// func (t Timestamp) MarshalGQL(w io.Writer) {
// 	io.WriteString(w, strconv.Quote(time.Time(t).Format(time.RFC3339)))
// }

func MarshalTimestamp(t time.Time) graphql.Marshaler {
	timestamp := t.Unix() * 1000

	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.FormatInt(timestamp, 10))
	})
}

func UnmarshalTimestamp(v interface{}) (time.Time, error) {
	if tmpStr, ok := v.(int); ok {
		return time.Unix(int64(tmpStr), 0), nil
	}
	return time.Time{}, nil
}
