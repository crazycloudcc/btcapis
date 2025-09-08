package types

import "github.com/shopspring/decimal"

// runes: 用途：每个 Rune 的“元数据/发行规则”，一条记录=一个 Rune。
// 索引建议：
// 		PK(id)，UNIQUE(tx_id)；
// 		BTREE(block_height), BTREE(name), BTREE(number)；
// 		需要以“是否可 mint”检索时，再建复合索引：(terms_height_start,terms_height_end)。
type Runes struct {
	ID               string          `db:"id"`                 // Rune 的唯一 ID（通常是 `block_height-tx_index` 或内部编码）。可作为主键。
	Number           int64           `db:"number"`             // (postgre: int8)Rune 序号（按诞生顺序的编号）。
	Name             string          `db:"name"`               // 原始名称（不含空格/分隔符）。
	SpacedName       string          `db:"spaced_name"`        // 带 spacer 的展示名（协议允许在 name 中插入分隔符以增强可读性）。
	BlockHash        string          `db:"block_hash"`         // 铸造/诞生所在区块哈希。
	BlockHeight      int             `db:"block_height"`       // 诞生区块高度。用于范围/时间查询。
	TxIndex          int64           `db:"tx_index"`           // (postgre: int8)该区块内的交易索引（定位同高内顺序）。
	TxID             string          `db:"tx_id"`              // 诞生交易哈希。
	Divisibility     int16           `db:"divisibility"`       // (postgres: int2)小数位精度（等价 ERC-20 decimals）。
	Premine          decimal.Decimal `db:"premine"`            // (postgres: numberic)预挖/预留发行量（在 terms 之外直接进入供给）。
	Symbol           string          `db:"symbol"`             // 代币符号（截图中默认值显示为 `'r'::text`，不同实现可能为空）。
	TermsAmount      decimal.Decimal `db:"terms_amount"`       // (postgres: numberic)**每次 mint 的单位数量**（按规则铸造时的一次配额）。
	TermsCap         decimal.Decimal `db:"terms_cap"`          // (postgres: numberic)**总的可 mint 次数上限**（“能铸几次”）。达到后不可再 mint。
	TermsHeightStart int             `db:"terms_height_start"` // 允许 mint 的起始高度（闭区间起点）。
	TermsHeightEnd   int             `db:"terms_height_end"`   // 允许 mint 的结束高度（闭区间终点）。
	TermsOffsetStart int             `db:"terms_offset_start"` // 允许 mint 的相对高度起点（相对“诞生高度”的偏移）。实现上常用于 height + offset 的双窗限制。
	TermsOffsetEnd   int             `db:"terms_offset_end"`   // 允许 mint 的相对高度终点。
	Turbo            bool            `db:"turbo"`              // 是否为“turbo”模式（部分实现里代表宽松/快速铸造规则或 gas/费用优化标记，阅读你用的 indexer 定义为准）。
	Cenotaph         bool            `db:"cenotaph"`           // 失效/空墓标记：表示该 Rune 在诞生 tx 上因规则错误被判定为无效（常见于 ordinals 语境中的“Cenotaph”概念）。通常不应被统计进有效供给。
	Timestamp        int64           `db:"timestamp"`          // (postgre: int8)诞生区块时间（Unix 秒）。
}

// ledger 用途：事件流水表（事实表）。一笔 tx 可能产生多行：mint / send / receive（及可选 burn）。
// 查询口径：
//     “该地址本笔交易净变化” → 对同一 tx_id、address 聚合 receive - send (+ 可选 mint 归属)。
//     “同一 tx 概览” → 按 operation 汇总 mint/receive/send/burn。
// 索引建议：
//     BTREE(tx_id), BTREE(rune_id), BTREE(block_height), BTREE(operation), BTREE(receiver_address), BTREE(address)；
//     事件回放：复合索引 (block_height,tx_index,event_index)。
type RunesLedger struct {
	RuneID          string          `db:"rune_id"`          // 关联的 Rune。
	BlockHash       string          `db:"block_hash"`       // 所在区块哈希。
	BlockHeight     int             `db:"block_height"`     // 区块高度。
	TxIndex         int64           `db:"tx_index"`         // (postgre: int8)区块内交易索引。
	EventIndex      int64           `db:"event_index"`      // (postgre: int8)**同一交易内的事件序号**（确保事件稳定排序，便于回放）。
	TxID            string          `db:"tx_id"`            // 交易哈希。
	Output          int             `db:"output"`           // (postgre: int8)事件关联的 `vout`（收款/分发通常能落到具体输出；`mint` 可能无明确地址时该列可空）。
	Address         string          `db:"address"`          // 地址字段（实现里常作**发送方**或通用地址；搭配 `receiver_address` 使用，见下）。
	ReceiverAddress string          `db:"receiver_address"` // **接收方地址**（`receive` 事件一定有值；`send` 常为空；`mint` 视实现是否能归集到具体 vout）。
	Amount          decimal.Decimal `db:"amount"`           // (postgre: numeric)本事件的数量（正值）。
	Operation       string          `db:"operation"`        // (postgre: public.ledger_operation)**枚举**：常见有 `mint` / `send` / `receive` / `burn`。不同实现也可能仅用三类。
	Timestamp       int64           `db:"timestamp"`        // (postgre: int8)区块时间（Unix 秒）。
}

// balance_changes 用途：地址维度的累计余额快照/变动汇总（比逐事件更适合余额/排行榜查询）。
// 常见用法：
//     取地址最新余额：ORDER BY block_height DESC LIMIT 1；
//     做时间序列：按高度排序取 balance。
// 索引建议：
//     PK/UNIQUE(rune_id,address,block_height)；
//     BTREE(address), BTREE(block_height DESC)。
type RunesBalanceChanges struct {
	RuneID          string          `db:"rune_id"`          // 关联的 Rune。
	BlockHeight     int             `db:"block_height"`     // 统计到的区块高度（该高度时点的余额口径）。
	Address         string          `db:"address"`          // 地址。
	Balance         decimal.Decimal `db:"balance"`          // (postgre: numeric)该地址在该高度的**余额**（注意是累计量，而非当区块的 delta）。
	TotalOperations int64           `db:"total_operations"` // (postgre: int8)该地址到此高度累计发生的相关事件数（便于做增量/校验）。默认 0。
}

// supply_changes 用途：全局供给维度的累计快照/变动统计（每个高度一行或仅在有事件的高度一行）。
// 查询口径：
//     当前总供给 ≈ total_mints - total_burns（或 premine + 累计mint - 累计burn，看实现是否将 premine 计入 total_mints）。
//     供给时间序列：按高度取 total_mints/total_burns。
// 索引建议：
//     PK/UNIQUE(rune_id,block_height)；
//     BTREE(block_height), BTREE(rune_id)。
type RunesSupplyChanges struct {
	RuneID          string          `db:"rune_id"`          // 关联的 Rune。
	BlockHeight     int             `db:"block_height"`     // 统计到的区块高度。
	Minted          decimal.Decimal `db:"minted"`           // (postgre: numeric)**该高度内新增**的铸造量（delta）。
	TotalMints      decimal.Decimal `db:"total_mints"`      // (postgre: numeric)截至该高度的**累计铸造总量**。
	Burned          decimal.Decimal `db:"burned"`           // (postgre: numeric)**该高度内销毁**量（delta，若协议/实现支持）。
	TotalBurns      decimal.Decimal `db:"total_burns"`      // (postgre: numeric)截至该高度的**累计销毁总量**。
	TotalOperations int64           `db:"total_operations"` // (postgre: int8)截至该高度与该 Rune 有关的累计事件数（便于做增量/校验）。默认 0。
}
