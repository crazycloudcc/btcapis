package btcapis_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/crazycloudcc/btcapis"
)

const (
	testBlockHash  = "000000000000000000006b251a20066090eeec95c803ddf5f8adc937d022cd90"
	testCoinbaseID = "d498583f433a5d63941b5beeeeaceb69e6c233b1685ade3074dec2248b6a687a"
	testTxID       = "1ac03a06e8e915813caabdf0b9f7b40a8851b20dc48567901c4c664a82b3302b"
)

type rpcRequest struct {
	ID     int               `json:"id"`
	Method string            `json:"method"`
	Params []json.RawMessage `json:"params"`
}

func TestGetBlockTransactions(t *testing.T) {
	var calls atomic.Int64
	server := newRPCServer(t, func(t *testing.T, req *http.Request, rpcReq rpcRequest) any {
		calls.Add(1)
		assertGetBlockVerbosityTwo(t, rpcReq)
		return validBlockResult(`"fee":0.00060000`)
	})
	defer server.Close()

	client := newClient(server.URL, 5)
	block, err := client.GetBlockTransactions(context.Background(), testBlockHash)
	if err != nil {
		t.Fatalf("GetBlockTransactions: %v", err)
	}
	if block.Height != 692769 || block.BlockHash != testBlockHash || block.TxCount != 2 || len(block.Transactions) != 2 {
		t.Fatalf("unexpected block: %#v", block)
	}
	if !block.Transactions[0].Coinbase || block.Transactions[0].FeeSats != 0 {
		t.Fatalf("unexpected coinbase: %#v", block.Transactions[0])
	}
	if block.Transactions[1].TxID != testTxID || block.Transactions[1].VSize != 224 || block.Transactions[1].FeeSats != 60000 {
		t.Fatalf("unexpected transaction: %#v", block.Transactions[1])
	}

	var wg sync.WaitGroup
	for i := 0; i < 24; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := client.GetBlockTransactions(context.Background(), testBlockHash); err != nil {
				t.Errorf("concurrent GetBlockTransactions: %v", err)
			}
		}()
	}
	wg.Wait()
	if calls.Load() != 25 {
		t.Fatalf("unexpected call count: %d", calls.Load())
	}
}

func TestGetBlockTransactionsRejectsInvalidData(t *testing.T) {
	tests := []struct {
		name        string
		feeFragment string
	}{
		{name: "missing fee"},
		{name: "fractional sat", feeFragment: `"fee":0.000000001`},
		{name: "negative fee", feeFragment: `"fee":-0.00000001`},
		{name: "overflow", feeFragment: `"fee":92233720368.54775808`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := newRPCServer(t, func(_ *testing.T, _ *http.Request, _ rpcRequest) any {
				return validBlockResult(tt.feeFragment)
			})
			defer server.Close()

			_, err := newClient(server.URL, 5).GetBlockTransactions(context.Background(), testBlockHash)
			if !errors.Is(err, btcapis.ErrInvalidBlockData) {
				t.Fatalf("expected ErrInvalidBlockData, got %v", err)
			}
		})
	}
}

func TestGetBlockTransactionsErrors(t *testing.T) {
	t.Run("block not found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_ = json.NewEncoder(w).Encode(map[string]any{
				"result": nil,
				"error":  map[string]any{"code": -5, "message": "Block not found"},
				"id":     1,
			})
		}))
		defer server.Close()

		_, err := newClient(server.URL, 5).GetBlockTransactions(context.Background(), testBlockHash)
		if !errors.Is(err, btcapis.ErrBlockNotFound) {
			t.Fatalf("expected ErrBlockNotFound, got %v", err)
		}
	})

	t.Run("context timeout", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			time.Sleep(100 * time.Millisecond)
			_ = json.NewEncoder(w).Encode(map[string]any{"result": validBlockResult(`"fee":0.00060000`), "error": nil, "id": 1})
		}))
		defer server.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
		defer cancel()
		_, err := newClient(server.URL, 5).GetBlockTransactions(ctx, testBlockHash)
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("expected context deadline, got %v", err)
		}
	})
}

func newClient(rpcURL string, timeout int) *btcapis.Client {
	return btcapis.New(&btcapis.Config{
		Network: "mainnet",
		Timeout: timeout,
		RPCUrl:  rpcURL,
	})
}

func newRPCServer(t *testing.T, result func(*testing.T, *http.Request, rpcRequest) any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var rpcReq rpcRequest
		if err := json.NewDecoder(req.Body).Decode(&rpcReq); err != nil {
			t.Errorf("decode request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"result": result(t, req, rpcReq),
			"error":  nil,
			"id":     1,
		})
	}))
}

func assertGetBlockVerbosityTwo(t *testing.T, req rpcRequest) {
	t.Helper()
	if req.ID != 1 || req.Method != "getblock" || len(req.Params) != 2 {
		t.Errorf("unexpected RPC request: %#v", req)
		return
	}
	var hash string
	var verbosity int
	if err := json.Unmarshal(req.Params[0], &hash); err != nil {
		t.Errorf("decode hash: %v", err)
	}
	if err := json.Unmarshal(req.Params[1], &verbosity); err != nil {
		t.Errorf("decode verbosity: %v", err)
	}
	if hash != testBlockHash || verbosity != 2 {
		t.Errorf("unexpected params: hash=%q verbosity=%d", hash, verbosity)
	}
}

func validBlockResult(feeFragment string) json.RawMessage {
	if feeFragment != "" {
		feeFragment += ","
	}
	return json.RawMessage(`{
		"hash":"` + testBlockHash + `",
		"height":692769,
		"time":1627613920,
		"size":1000,
		"weight":4000,
		"nTx":2,
		"tx":[
			{"txid":"` + testCoinbaseID + `","vsize":145,"weight":580,"vin":[{"coinbase":"00"}]},
			{"txid":"` + testTxID + `","vsize":224,"weight":896,` + feeFragment + `"vin":[{"txid":"prev"}]}
		]
	}`)
}
