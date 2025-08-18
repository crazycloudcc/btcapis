package tx

import (
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/crazycloudcc/btcapis/types"
)

// ExtractBRC20 scans transaction inputs and returns BRC-20 actions
// found within any ordinal inscriptions embedded in the witness.
// Witnesses that cannot be parsed are silently ignored.
func ExtractBRC20(t *types.Tx) []types.BRC20Action {
	var actions []types.BRC20Action
	for i := range t.Vin {
		info, err := AnalyzeTxInWithIdx(t, i)
		if err != nil || info.Ord == nil || info.Ord.BodyHex == "" {
			continue
		}
		body, err := hex.DecodeString(info.Ord.BodyHex)
		if err != nil {
			continue
		}
		var m map[string]any
		if err := json.Unmarshal(body, &m); err != nil {
			continue
		}
		p, ok := m["p"].(string)
		if !ok || strings.ToLower(p) != "brc-20" {
			continue
		}
		act := types.BRC20Action{}
		if v, ok := m["op"].(string); ok {
			act.Op = v
		}
		if v, ok := m["tick"].(string); ok {
			act.Tick = v
		}
		if v, ok := m["amt"].(string); ok {
			act.Amt = v
		}
		if v, ok := m["max"].(string); ok {
			act.Max = v
		}
		if v, ok := m["lim"].(string); ok {
			act.Lim = v
		}
		actions = append(actions, act)
	}
	return actions
}

// ExtractRunes scans transaction outputs looking for runestone data encoded
// in OP_RETURN scripts. The detection is heuristic and treats any OP_RETURN
// whose payload begins with "rune" (case-insensitive) as a runestone.
// The raw payload hex is returned for further processing by callers.
func ExtractRunes(t *types.Tx) []types.Runestone {
	var stones []types.Runestone
	for _, o := range t.Vout {
		payload, ok := opReturnPayload(o.ScriptPubKey)
		if !ok || len(payload) == 0 {
			continue
		}
		if strings.HasPrefix(strings.ToLower(string(payload)), "rune") {
			stones = append(stones, types.Runestone{BodyHex: hex.EncodeToString(payload)})
		}
	}
	return stones
}

// opReturnPayload concatenates all push-data that follows an OP_RETURN. It is
// a minimal parser sufficient for runestone detection.
func opReturnPayload(script []byte) ([]byte, bool) {
	if len(script) == 0 || script[0] != 0x6a {
		return nil, false
	}
	i := 1
	var out []byte
	for i < len(script) {
		op := script[i]
		i++
		var l int
		switch {
		case op < 0x4c:
			l = int(op)
		case op == 0x4c:
			if i >= len(script) {
				return nil, false
			}
			l = int(script[i])
			i++
		case op == 0x4d:
			if i+1 >= len(script) {
				return nil, false
			}
			l = int(script[i]) | int(script[i+1])<<8
			i += 2
		case op == 0x4e:
			if i+3 >= len(script) {
				return nil, false
			}
			l = int(script[i]) | int(script[i+1])<<8 | int(script[i+2])<<16 | int(script[i+3])<<24
			i += 4
		default:
			return nil, false
		}
		if i+l > len(script) {
			return nil, false
		}
		out = append(out, script[i:i+l]...)
		i += l
	}
	return out, true
}
