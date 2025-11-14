package blocklists

import _ "embed"
import "encoding/json"

//go:embed cuuiliste.json
var embedded_cuuiliste []byte

type CUUI struct {
	AddedBy        string `json:"added_by"`
	Domain         string `json:"domain"`
	FirstBlockedOn string `json:"first_blocked_on"`
	// Site string `json:"site"` // always null?
}

func init() {

	var cuuilist []CUUI

	err := json.Unmarshal(embedded_cuuiliste, &cuuilist)

	if err == nil {

		blocklist := make(map[string]bool)

		for _, entry := range cuuilist {

			if entry.Domain != "" && entry.FirstBlockedOn != "" {
				blocklist[entry.Domain] = true
			}

		}

		Blocklists["cuui"] = blocklist

	}

}
