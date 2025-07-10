package sysConfig

import (
	"sync"
)

type JobSched struct {
	MaxRetryCazFailPerSched uint32 `json:"max_retry_caz_fail_per_sched"`
	MaxRunSec               uint32 `json:"max_run_sec"`
}

type JobSchedConfig struct {
	SysSubConfig[JobSched]
}

func InstJobSchedConfig(install ...interface{}) *JobSchedConfig {
	if len(install) > 0 && install[0].(bool) {
		sub := &JobSchedConfig{}
		sub.Cache = &sync.Map{}
		sub.SysKey = KeyJobSchedConfig
		RegistSubConfig(KeyJobSchedConfig, sub)
		return sub
	}
	inst := GetSubConfig(KeyJobSchedConfig)
	return inst.(*JobSchedConfig)
}
