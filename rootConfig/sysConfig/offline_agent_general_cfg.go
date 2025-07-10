package sysConfig

type SysOfflineAgentGeneralCfg struct {
	Agentgenerals    []int32         `json:"agent_generals"`     // 代理总代列表
	ProfitRadio      int32           `json:"profit_radio"`       // 占成比例
	GameRebateRadios map[int32]int32 `json:"game_rebate_radios"` // 游戏返水比例
	MaxAgentLevel    int32           `json:"max_agent_level"`    // 最多多少层级
}
