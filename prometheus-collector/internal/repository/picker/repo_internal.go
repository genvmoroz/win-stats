package picker

import (
	"encoding/json"
	"fmt"
	"net/http"

	openapi "github.com/genvmoroz/win-stats-prometheus-collector/internal/repository/picker/generated"
	"github.com/samber/lo"
)

func handleResponse[FROM any, TO any](
	resp *http.Response,
	parse func(*http.Response) (*FROM, error),
	transform func(FROM) (TO, error)) (TO, error) {
	if resp == nil {
		return lo.Empty[TO](), fmt.Errorf("response is nil")
	}

	if err := checkError(resp); err != nil {
		return lo.Empty[TO](), err
	}

	from, err := parse(resp)
	if err != nil {
		return lo.Empty[TO](), fmt.Errorf("parse response: %w", err)
	}
	if from == nil {
		return lo.Empty[TO](), fmt.Errorf("parsed response is nil")
	}

	to, err := transform(*from)
	if err != nil {
		return lo.Empty[TO](), fmt.Errorf("transform response: %w", err)
	}

	return to, nil
}

func checkError(resp *http.Response) error {
	switch {
	case resp.StatusCode == http.StatusOK:
		return nil
	default:
		return extractErrorFromResponseBody(resp)
	}
}

func extractErrorFromResponseBody(resp *http.Response) error {
	var em openapi.Error
	if err := json.NewDecoder(resp.Body).Decode(&em); err != nil {
		return fmt.Errorf(
			"unmarshal error message from server response (StatusCode=%d): %w",
			resp.StatusCode, err,
		)
	}

	return fmt.Errorf(
		"server error: (StatusCode=%d)  (Message=%s)",
		lo.FromPtr(em.StatusCode), lo.FromPtr(em.Message),
	)
}
