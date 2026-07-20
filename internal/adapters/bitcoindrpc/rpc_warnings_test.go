package bitcoindrpc

import (
	"encoding/json"
	"testing"
)

func TestRPCWarningsUnmarshalString(t *testing.T) {
	var w RPCWarnings
	if err := json.Unmarshal([]byte(`"this is a warning"`), &w); err != nil {
		t.Fatalf("unmarshal string: %v", err)
	}
	if len(w) != 1 || w[0] != "this is a warning" {
		t.Fatalf("unexpected: %#v", w)
	}
	if got := w.AsAny(); got != "this is a warning" {
		t.Fatalf("AsAny=%v", got)
	}
}

func TestRPCWarningsUnmarshalEmptyString(t *testing.T) {
	var w RPCWarnings
	if err := json.Unmarshal([]byte(`""`), &w); err != nil {
		t.Fatalf("unmarshal empty string: %v", err)
	}
	if w != nil || w.AsAny() != nil {
		t.Fatalf("empty string should be nil, got %#v", w)
	}
}

func TestRPCWarningsUnmarshalArray(t *testing.T) {
	var w RPCWarnings
	if err := json.Unmarshal([]byte(`["a","b"]`), &w); err != nil {
		t.Fatalf("unmarshal array: %v", err)
	}
	if len(w) != 2 || w[0] != "a" || w[1] != "b" {
		t.Fatalf("unexpected: %#v", w)
	}
	got, ok := w.AsAny().([]string)
	if !ok || len(got) != 2 {
		t.Fatalf("AsAny=%v", w.AsAny())
	}
}

func TestRPCWarningsUnmarshalEmptyArray(t *testing.T) {
	var w RPCWarnings
	if err := json.Unmarshal([]byte(`[]`), &w); err != nil {
		t.Fatalf("unmarshal empty array: %v", err)
	}
	if len(w) != 0 {
		t.Fatalf("unexpected: %#v", w)
	}
}

func TestChainInfoDTOWarningsCore28(t *testing.T) {
	raw := []byte(`{
		"chain":"main",
		"blocks":900000,
		"headers":900000,
		"bestblockhash":"abc",
		"difficulty":1,
		"verificationprogress":1,
		"initialblockdownload":false,
		"chainwork":"00",
		"size_on_disk":1,
		"pruned":false,
		"warnings":["Warning: unknown new rules"]
	}`)
	var dto ChainInfoDTO
	if err := json.Unmarshal(raw, &dto); err != nil {
		t.Fatalf("unmarshal ChainInfoDTO: %v", err)
	}
	if len(dto.Warnings) != 1 || dto.Warnings[0] != "Warning: unknown new rules" {
		t.Fatalf("warnings=%#v", dto.Warnings)
	}
}

func TestNetworkInfoDTOWarningsLegacyString(t *testing.T) {
	raw := []byte(`{
		"version":280000,
		"subversion":"/Satoshi:28.0.0/",
		"protocolversion":70016,
		"localservices":"0000",
		"localservicesnames":[],
		"localrelay":true,
		"timeoffset":0,
		"networkactive":true,
		"connections":8,
		"connections_in":1,
		"connections_out":7,
		"networks":[],
		"relayfee":0.00001,
		"incrementalfee":0.00001,
		"localaddresses":[],
		"warnings":""
	}`)
	var dto NetworkInfoDTO
	if err := json.Unmarshal(raw, &dto); err != nil {
		t.Fatalf("unmarshal NetworkInfoDTO: %v", err)
	}
	if dto.Warnings != nil {
		t.Fatalf("empty warnings should be nil, got %#v", dto.Warnings)
	}
}
