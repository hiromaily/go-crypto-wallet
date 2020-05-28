package xrp

import (
	"context"

	"github.com/pkg/errors"
)

// https://xrpl.org/server-info-methods.html

// ResponseServerInfo is response data for server_info method
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
					DurationUs  string `json:"duration_us"`
					Transitions int    `json:"transitions"`
				} `json:"connected"`
				Disconnected struct {
					DurationUs  string `json:"duration_us"`
					Transitions int    `json:"transitions"`
				} `json:"disconnected"`
				Full struct {
					DurationUs  string `json:"duration_us"`
					Transitions int    `json:"transitions"`
				} `json:"full"`
				Syncing struct {
					DurationUs  string `json:"duration_us"`
					Transitions int    `json:"transitions"`
				} `json:"syncing"`
				Tracking struct {
					DurationUs  string `json:"duration_us"`
					Transitions int    `json:"transitions"`
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

// ServerInfo calls server_info method
func (r *Ripple) ServerInfo() (*ResponseServerInfo, error) {
	req := RequestCommand{
		ID:      1,
		Command: "server_info",
	}
	var res ResponseServerInfo
	if err := r.wsPublic.Call(context.Background(), &req, &res); err != nil {
		return nil, errors.Wrap(err, "fail to call wsClient.Call()")
	}
	return &res, nil
}
