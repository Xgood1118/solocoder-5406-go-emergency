package handler

import (
	"net/http"

	"emergency-dispatch/internal/models"
	"emergency-dispatch/internal/store"

	"github.com/gin-gonic/gin"
)

func ListCallers(c *gin.Context) {
	callers := store.GlobalStore.ListCallers()
	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    callers,
	})
}

func GetCaller(c *gin.Context) {
	id := c.Param("id")
	caller, ok := store.GlobalStore.GetCaller(id)
	if !ok {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "报警人不存在",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    caller,
	})
}

func GetCallerByPhone(c *gin.Context) {
	phone := c.Query("phone")
	if phone == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "phone 参数不能为空",
		})
		return
	}

	caller, ok := store.GlobalStore.GetCallerByPhone(phone)
	if !ok {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "未找到该手机号的报警人",
		})
		return
	}

	var dispatches []*models.Dispatch
	for _, dispatchID := range caller.DispatchIDs {
		d, ok := store.GlobalStore.GetDispatch(dispatchID)
		if ok {
			dispatches = append(dispatches, d)
		}
	}

	result := map[string]interface{}{
		"caller":     caller,
		"dispatches": dispatches,
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    result,
	})
}

type UpdateVIPRequest struct {
	IsVIP    bool   `json:"is_vip"`
	Operator string `json:"operator"`
}

func UpdateCallerVIP(c *gin.Context) {
	id := c.Param("id")
	var req UpdateVIPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	caller, ok := store.GlobalStore.GetCaller(id)
	if !ok {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "报警人不存在",
		})
		return
	}

	caller.IsVIP = req.IsVIP
	store.GlobalStore.UpdateCaller(caller)

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    caller,
	})
}
