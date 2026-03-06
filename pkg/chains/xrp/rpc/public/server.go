package public

import (
	"context"
	"fmt"
)

// https://xrpl.org/server-info-methods.html

// RequestCommand is the minimal request body for commands with no parameters.
type RequestCommand struct {
	ID      int    `json:"id"`
	Command string `json:"command"`
}

// ResponseServerInfo is the wire-format response of the server_info command.
type ResponseServerInfo struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
	Type   string `json:"type"`
	Result struct {
		Info struct {
			BuildVersion    string `json:"build_version"`
			CompleteLedgers string `json:"complete_ledgers"`
			Hostid          string `json:"hostid"`
			IoLatencyMs     int    `json:"io_latency_ms"`
			JqTransOverflow string `json:"jq_trans_overflow"`
			LastClose       struct {
				ConvergeTimeS float64 `json:"converge_time_s"`
				Proposers     int     `json:"proposers"`
			} `json:"last_close"`
			Load struct {
				JobTypes []struct {
					JobType    string `json:"job_type"`
					PeakTime   int    `json:"peak_time,omitempty"`
					PerSecond  int    `json:"per_second"`
					AvgTime    int    `json:"avg_time,omitempty"`
					InProgress int    `json:"in_progress,omitempty"`
				} `json:"job_types"`
				Threads int `json:"threads"`
			} `json:"load"`
			LoadFactor               int    `json:"load_factor"`
			PeerDisconnects          string `json:"peer_disconnects"`
			PeerDisconnectsResources string `json:"peer_disconnects_resources"`
			Peers                    int    `json:"peers"`
			PubkeyNode               string `json:"pubkey_node"`
			PubkeyValidator          string `json:"pubkey_validator"`
			ServerState              string `json:"server_state"`
			ServerStateDurationUs    string `json:"server_state_duration_us"`
			StateAccounting          struct {
				Connected struct {
					DurationUs  string
					Transitions int
				} `json:"connected"`
				Disconnected struct {
					DurationUs  string
					Transitions int
				} `json:"disconnected"`
				Full struct {
					DurationUs  string
					Transitions int
				} `json:"full"`
				Syncing struct {
					DurationUs  string
					Transitions int
				} `json:"syncing"`
				Tracking struct {
					DurationUs  string
					Transitions int
				} `json:"tracking"`
			} `json:"state_accounting"`
			Time            string `json:"time"`
			Uptime          int    `json:"uptime"`
			ValidatedLedger struct {
				Age            int     `json:"age"`
				BaseFeeXrp     float64 `json:"base_fee_xrp"`
				Hash           string  `json:"hash"`
				ReserveBaseXrp int     `json:"reserve_base_xrp"`
				ReserveIncXrp  int     `json:"reserve_inc_xrp"`
				Seq            int     `json:"seq"`
			} `json:"validated_ledger"`
			ValidationQuorum int `json:"validation_quorum"`
			ValidatorList    struct {
				Count      int    `json:"count"`
				Expiration string `json:"expiration"`
				Status     string `json:"status"`
			} `json:"validator_list"`
		} `json:"info"`
	} `json:"result"`
	Error string `json:"error,omitempty"`
}

// ServerInfo calls the server_info WebSocket command and returns the raw wire response.
func (r *PublicRPC) ServerInfo(ctx context.Context) (*ResponseServerInfo, error) {
	req := &RequestCommand{
		ID:      1,
		Command: "server_info",
	}
	var res ResponseServerInfo
	if err := r.caller.Call(ctx, req, &res); err != nil {
		return nil, fmt.Errorf("fail to call wsClient.Call(server_info): %w", err)
	}
	return &res, nil
}
