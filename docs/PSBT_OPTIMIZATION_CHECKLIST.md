# PSBT 模块优化 Checklist

## 阻断类问题（必须修复）

### 1. 不支持 P2SH 包裹 SegWit

- [x] 签名阶段识别 P2SH 包裹的 P2WPKH/P2WSH
- [x] 最终化阶段正确处理 FinalScriptSig 和 FinalWitness
- [x] 实现 BIP143 签名逻辑

### 2. P2WSH 多签签名顺序与见证脚本公钥顺序未对齐

- [x] 解析 witnessScript 提取公钥序列
- [x] 按脚本中公钥顺序从 PartialSigs 取签名
- [x] 支持额外数据栈（HASHLOCK 等）

### 3. 未校验 P2WSH/P2WPKH 脚本与 UTXO 一致性

- [x] P2WPKH: hash160(pubkey) == pkScript[2:]
- [x] P2WSH: sha256(witnessScript) == pkScript[2:]
- [x] 校验失败直接报错

### 4. Taproot keypath 最终化忽略 annex

- [x] Keypath 分支正确处理 annex
- [x] 脚本路径支持无签名脚本
- [x] 完善 TapScriptStack 处理

### 5. v2 Combine 未校验 Sequence 一致性

- [x] ensureSameTemplate 加入 Sequence 等价校验
- [x] mergeInput 合并 BIP32 派生信息（去重）

### 6. 签名阶段未覆盖 Legacy 与 P2SH 非见证分支

- [x] 支持 P2PKH 签名
- [x] 支持 P2SH 非见证程序签名
- [x] 统一签名接口

## 风险与规范性（强烈建议优化）

### 7. v0 Finalize 对 segwit 仅接受 WitnessUtxo

- [ ] 允许 segwit 仅带 NonWitnessUtxo
- [ ] 从前序交易取 pkScript/value

### 8. Finalization 后字段清理不一致

- [ ] 统一"final 后尽量最小化"策略
- [ ] 清理 PartialSigs/BIP32/Temp Stacks
- [ ] 配置项控制是否保留 WitnessScript/UTXO

### 9. Taproot 数据建模不够 PSBT-371 友好

- [ ] 区分 tap_key_sig 和 tap_script_sig
- [ ] 引入独立字段避免冲突
- [ ] 支持<xonlyPubKey, leafhash>键值

### 10. P2WPKH 最终化过度依赖 BIP32 条目

- [ ] 回退用 PartialSigs 的 key 直接获取 pubkey
- [ ] 支持缺少 BIP32 时的 finalize

### 11. Key 命名与类型

- [ ] 统一入参为 pubkey []byte
- [ ] 按脚本类型校验长度（33 压缩或 32 x-only）

## 体验与可维护性（加分项）

### 12. 增加 Analyze()自检

- [ ] 输出每个输入缺少的字段
- [ ] 类似 Core 的 analyzepsbt 功能
- [ ] 生成详细报告

### 13. 在 SignInput 存档 SighashType

- [ ] 调用时传入的 sighash 存回 in.SighashType
- [ ] 便于后续合并/审计

### 14. 更细的错误文案

- [ ] 补足上下文信息
- [ ] 包含输入索引、脚本类型判断分支
- [ ] 明确需要提供哪些字段

## 实现优先级

### 高优先级（本周完成）

1. P2SH 包裹 SegWit 支持
2. 脚本一致性校验
3. 签名顺序对齐

### 中优先级（下周完成）

4. Taproot annex 处理
5. v2 合并一致性
6. Legacy 签名支持

### 低优先级（后续迭代）

7. 字段清理优化
8. 数据建模改进
9. 自检功能
10. 错误信息优化

## 测试计划

- [ ] 单元测试覆盖所有脚本类型
- [ ] 集成测试验证 P2SH 包裹 SegWit
- [ ] 压力测试多签场景
- [ ] 兼容性测试 v0/v2
