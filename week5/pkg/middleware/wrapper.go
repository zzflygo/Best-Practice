package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

//包装middleware

func Wrapper(
	size int,
	total, fail float64,
	duration time.Duration,
) gin.HandlerFunc {
	// 传入参数 获取 滑动窗口对象
	s := NewSlidingWindow(size, total, fail, duration)
	// 启动滑动窗口
	s.Start()
	// 返回滑动窗口状态
	s.Stauts()
	// 判断是否需要熔断
	s.Monitor()
	// 写入req数据
	return func(c *gin.Context) {
		//判断是否熔断
		if s.broken {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "",
			})
			c.Abort()
			return
		}
		c.Next()
		if c.Writer.Status() != http.StatusOK {
			s.RecordReqResult(false)
		} else {
			s.RecordReqResult(true)
		}
	}

}
