package mydata

import (
	"encoding/json"
	"fmt"
	"strings"
)

func getDays(days ...int) int {
	daycount := 45

	if len(days) > 0 {
		daycount = days[0]
	}

	return daycount
}

func parseError(b []byte) (error, error) {
	if !strings.Contains(string(b), "errorMessage") {
		return nil, nil
	}

	e := struct {
		Msg string `json:"errorMessage"`
	}{}
	err := json.Unmarshal(b, &e)
	if err != nil {
		return nil, err
	}

	return fmt.Errorf("unexpected error: %s", e.Msg), nil
}
