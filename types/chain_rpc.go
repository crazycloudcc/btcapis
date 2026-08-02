package types

// ChainNetworkInfo 节点网络信息（getnetworkinfo）
type ChainNetworkInfo struct {
	Version            int64              `json:"version"`
	Subversion         string             `json:"subversion"`
	ProtocolVersion    int64              `json:"protocolversion"`
	LocalServices      string             `json:"localservices"`
	LocalServicesNames []string           `json:"localservicesnames"`
	LocalRelay         bool               `json:"localrelay"`
	TimeOffset         int64              `json:"timeoffset"`
	NetworkActive      bool               `json:"networkactive"`
	Connections        int64              `json:"connections"`
	ConnectionsIn      int64              `json:"connections_in"`
	ConnectionsOut     int64              `json:"connections_out"`
	Networks           []ChainNetworkMeta `json:"networks"`
	RelayFee           float64            `json:"relayfee"`
	IncrementalFee     float64            `json:"incrementalfee"`
	LocalAddresses     []string           `json:"localaddresses"`
	Warnings           any                `json:"warnings"`
}

// ChainNetworkMeta 网络元信息
type ChainNetworkMeta struct {
	Name                      string `json:"name"`
	Limited                   bool   `json:"limited"`
	Reachable                 bool   `json:"reachable"`
	Proxy                     string `json:"proxy"`
	ProxyRandomizeCredentials bool   `json:"proxy_randomize_credentials"`
}

// ChainStats 区块统计（getblockstats）
type ChainStats struct {
	AvgFee             float64   `json:"avgfee"`
	AvgFeeRate         float64   `json:"avgfeerate"`
	AvgTxSize          float64   `json:"avgtxsize"`
	BlockHash          string    `json:"blockhash"`
	FeeratePercentiles []float64 `json:"feerate_percentiles"`
	Height             int64     `json:"height"`
	Ins                int64     `json:"ins"`
	MaxFee             float64   `json:"maxfee"`
	MaxFeeRate         float64   `json:"maxfeerate"`
	MaxTxSize          int64     `json:"maxtxsize"`
	MedianFee          float64   `json:"medianfee"`
	MedianTime         int64     `json:"mediantime"`
	MedianTxSize       float64   `json:"mediantxsize"`
	MinFee             float64   `json:"minfee"`
	MinFeeRate         float64   `json:"minfeerate"`
	MinTxSize          int64     `json:"mintxsize"`
	Outs               int64     `json:"outs"`
	Subsidy            float64   `json:"subsidy"`
	SwtotalSize        int64     `json:"swtotal_size"`
	SwtotalWeight      int64     `json:"swtotal_weight"`
	Swtxs              int64     `json:"swtxs"`
	Time               int64     `json:"time"`
	TotalOut           float64   `json:"total_out"`
	TotalSize          int64     `json:"total_size"`
	TotalWeight        int64     `json:"total_weight"`
	Totalfee           float64   `json:"totalfee"`
	Txs                int64     `json:"txs"`
	UTXOIncrease       int64     `json:"utxo_increase"`
	UTXOSizeInc        int64     `json:"utxo_size_inc"`
	UTXOIncreaseActual int64     `json:"utxo_increase_actual"`
	UTXOSizeIncActual  int64     `json:"utxo_size_inc_actual"`
}

// ChainTip 链分叉信息（getchaintips）
type ChainTip struct {
	Height    int64  `json:"height"`
	Hash      string `json:"hash"`
	BranchLen int64  `json:"branchlen"`
	Status    string `json:"status"`
}

// BlockChainInfo 链状态（getblockchaininfo）
type BlockChainInfo struct {
	Chain                string  `json:"chain"`
	Blocks               int64   `json:"blocks"`
	Headers              int64   `json:"headers"`
	BestBlockHash        string  `json:"bestblockhash"`
	Bits                 string  `json:"bits"`
	Target               string  `json:"target"`
	Difficulty           float64 `json:"difficulty"`
	Time                 int64   `json:"time"`
	MedianTime           int64   `json:"mediantime"`
	VerificationProgress float64 `json:"verificationprogress"`
	InitialBlockDownload bool    `json:"initialblockdownload"`
	ChainWork            string  `json:"chainwork"`
	SizeOnDisk           int64   `json:"size_on_disk"`
	Pruned               bool    `json:"pruned"`
	Warnings             any     `json:"warnings"`
}

// BlockHeader 区块头（getblockheader）
type BlockHeader struct {
	Hash              string  `json:"hash"`
	Confirmations     int64   `json:"confirmations"`
	Height            int64   `json:"height"`
	Version           int64   `json:"version"`
	VersionHex        string  `json:"versionHex"`
	MerkleRoot        string  `json:"merkleRoot"`
	Time              int64   `json:"time"`
	MedianTime        int64   `json:"mediantime"`
	Nonce             int64   `json:"nonce"`
	Bits              string  `json:"bits"`
	Target            string  `json:"target"`
	Difficulty        float64 `json:"difficulty"`
	ChainWork         string  `json:"chainwork"`
	NTx               int64   `json:"nTx"`
	PreviousBlockHash string  `json:"previousblockhash"`
	NextBlockHash     string  `json:"nextblockhash"`
}

// Block 区块（getblock verbosity=1）
type Block struct {
	BlockHeader
	StrippedSize int64    `json:"strippedsize"`
	Size         int64    `json:"size"`
	Weight       int64    `json:"weight"`
	Tx           []string `json:"tx"`
}

// BlockTransactions 是 getblock verbosity=2 的紧凑有序交易结果。
type BlockTransactions struct {
	Height       int64              `json:"height"`
	BlockHash    string             `json:"block_hash"`
	Time         int64              `json:"time"`
	Size         int64              `json:"size"`
	Weight       int64              `json:"weight"`
	TxCount      int64              `json:"tx_count"`
	Transactions []BlockTransaction `json:"transactions"`
}

// BlockTransaction 只保留区块聚合需要的交易字段。
type BlockTransaction struct {
	TxID     string `json:"txid"`
	VSize    int64  `json:"vsize"`
	Weight   int64  `json:"weight"`
	FeeSats  int64  `json:"fee_sats"`
	Coinbase bool   `json:"coinbase"`
}

// AddressValidation 地址校验（validateaddress）
type AddressValidation struct {
	IsValid        bool   `json:"isvalid"`
	Address        string `json:"address"`
	ScriptPubKey   string `json:"scriptPubKey"`
	IsScript       bool   `json:"isscript"`
	IsWitness      bool   `json:"iswitness"`
	WitnessVersion int64  `json:"witness_version"`
	WitnessProgram string `json:"witness_program"`
	Error          string `json:"error"`
	ErrorLocations any    `json:"error_locations"`
}

// MempoolTxFees 内存池交易费用
type MempoolTxFees struct {
	Base       float64 `json:"base"`
	Modified   float64 `json:"modified"`
	Ancestor   float64 `json:"ancestor"`
	Descendant float64 `json:"descendant"`
}

// MempoolTxEntry 内存池交易详情（getmempoolentry）
type MempoolTxEntry struct {
	Vsize             int64         `json:"vsize"`
	Weight            int64         `json:"weight"`
	Time              int64         `json:"time"`
	Height            int64         `json:"height"`
	DescendantCount   int64         `json:"descendantcount"`
	DescendantSize    int64         `json:"descendantsize"`
	AncestorCount     int64         `json:"ancestorcount"`
	AncestorSize      int64         `json:"ancestorsize"`
	Wtxid             string        `json:"wtxid"`
	Fees              MempoolTxFees `json:"fees"`
	Depends           []string      `json:"depends"`
	SpentBy           []string      `json:"spentby"`
	Bip125Replaceable bool          `json:"bip125-replaceable"`
	Unbroadcast       bool          `json:"unbroadcast"`
}
