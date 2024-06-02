package service

import (
	h "Sgrid/src/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func SystemStatisticsRegisty(ctx *h.SgridServerCtx) {
	GROUP := ctx.Engine.Group(strings.ToLower(ctx.Name))
	GROUP.GET("/system/statistics/getCpuInfo", getCpuInfo)
	GROUP.GET("/system/statistics/getMemoryInfo", getMemoryInfo)
}

// cpu info
type S_Cpu struct {
	CPU         int32  `json:"cpu,omitempty"`
	CacheSize   int32  `json:"cacheSize,omitempty"`
	Cores       int32  `json:"cores,omitempty"`
	Descprition string `json:"descprition,omitempty"`
}

func getCpuInfo(c *gin.Context) {
	cpuInfos, err := cpu.Info()
	if err != nil {
		h.AbortWithError(c, err.Error())
		return
	}
	var resp []*S_Cpu
	for _, v := range cpuInfos {
		resp = append(resp, &S_Cpu{
			CPU:         v.CPU,
			CacheSize:   v.CacheSize,
			Cores:       v.Cores,
			Descprition: v.String(),
		})
	}
	h.AbortWithSucc(c, resp)
}

func getMemoryInfo(c *gin.Context) {
	memInfo, _ := mem.VirtualMemory()
	h.AbortWithSucc(c, memInfo.String())
}
