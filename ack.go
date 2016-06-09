// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
)

type AckRequest struct {
	TxID            string `json:"txid,omitempty"`
	FullTransaction string `json:"fulltransaction,omitempty"`
}

type FactoidTxStatus struct {
	TxID string `json:"txid"`
	GeneralTransactionData
}

type EntryStatus struct {
	CommitTxID string `json:"committxid"`
	EntryHash  string `json:"entryhash"`

	CommitData GeneralTransactionData `json:"commitdata"`
	EntryData  GeneralTransactionData `json:"entrydata"`

	ReserveTransactions          []ReserveInfo `json:"reserveinfo,omitempty"`
	ConflictingRevealEntryHashes []string      `json:"conflictingrevealentryhashes,omitempty"`
}

type ReserveInfo struct {
	TxID    string `json:"txid"`
	Timeout int64  `json:"timeout"` //Unix time
}

type GeneralTransactionData struct {
	TransactionDate int64      `json:"transactiondate"` //Unix time
	Malleated       *Malleated `json:"malleated,omitempty"`
	Status          string     `json:"status"`
}

type Malleated struct {
	MalleatedTxIDs []string `json:"malleatedtxids"`
}

func FactoidACK(txID, fullTransaction string) (*FactoidTxStatus, error) {
	param := AckRequest{TxID: txID, FullTransaction: fullTransaction}
	req := NewJSON2Request("factoid-ack", apiCounter(), param)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	eb := new(FactoidTxStatus)
	if err := json.Unmarshal(resp.JSONResult(), eb); err != nil {
		return nil, err
	}

	return eb, nil
}

func EntryACK(txID, fullTransaction string) (*EntryStatus, error) {
	param := AckRequest{TxID: txID, FullTransaction: fullTransaction}
	req := NewJSON2Request("entry-ack", apiCounter(), param)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	eb := new(EntryStatus)
	if err := json.Unmarshal(resp.JSONResult(), eb); err != nil {
		return nil, err
	}

	return eb, nil
}
