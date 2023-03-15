package cryptoprocesser

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

const requestURL = "https://www.binance.com/api/v3/ticker/price"

func requestToBinance(ctx context.Context) ([]Answer, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var ans []Answer
	err = json.Unmarshal(body, &ans)
	if err != nil {
		return nil, err
	}
	return ans, nil
}
