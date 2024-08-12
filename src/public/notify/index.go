package notify

import (
	"Sgrid/src/public"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Notify sends a notification to a specified target.
// The target should be configured centrally. For local development, it's set to log.
func Notify(target string, info string) error {
	if target == "" {
		fmt.Println("notify :: ", info)
		return nil
	}
	idx := os.Getenv(public.ENV_PROCESS_INDEX)
	name := os.Getenv(public.ENV_SGRID_SERVANT_NAME)
	if idx == "" || name == "" || info == "" {
		return fmt.Errorf("notify: env error")
	}
	payload := struct {
		GridID     string `json:"gridId"`
		ServerName string `json:"serverName"`
		Info       string `json:"info"`
	}{
		GridID:     idx,
		ServerName: name,
		Info:       info,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, target, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}

	return nil
}
