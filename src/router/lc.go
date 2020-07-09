package router

import (
	"fmt"
	"trffic_obj"

	"github.com/gin-gonic/gin"
)

// get : /lc
func Lc(c *gin.Context) {
	println("Lc router!")

	fmt.Println("lc")

	c.JSON(200, gin.H{
		"lc": "lc ok",
	})
}

// get : /lc_state_summary
func Lc_state_summary(c *gin.Context) {
	println("group router!")

	fmt.Println("group")

	apiData := trffic_obj.GetLcStateSummary()

	fmt.Println(apiData)
	c.JSON(200, apiData)
}
