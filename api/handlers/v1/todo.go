package v1

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"

	pb "github.com/abdullohsattorov/api-gateway/genproto"
	l "github.com/abdullohsattorov/api-gateway/pkg/logger"
	"github.com/abdullohsattorov/api-gateway/pkg/utils"
)

// CreateTodo creates todo
// route /v1/todos [post]
func (h *handlerV1) CreateTodo(c *gin.Context) {
	var (
		body        pb.TodoFunc
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.TodoService().Create(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create todo", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetTodo gets todo by id
// route /v1/todos/{id} [get]
func (h *handlerV1) GetTodo(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	guid := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.TodoService().Get(
		ctx, &pb.ByIdReq{
			Id: guid,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get todo", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// ListTodos returns list of todos
// route /v1/todos/ [get]
func (h *handlerV1) ListTodos(c *gin.Context) {
	queryParams := c.Request.URL.Query()

	params, errStr := utils.ParseQueryParams(queryParams)
	if errStr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errStr[0],
		})
		h.log.Error("failed to parse query params json" + errStr[0])
		return
	}

	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	if params.Time != "0000:00:00" {
		response, err := h.serviceManager.TodoService().ListOverdue(
			ctx, &pb.Time{
				Time:  params.Time,
				Limit: params.Limit,
				Page:  params.Page,
			})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed to list todos", l.Error(err))
			return
		}
		c.JSON(http.StatusOK, response)
	} else {
		fmt.Println("wow")
		response, err := h.serviceManager.TodoService().List(
			ctx, &pb.ListReq{
				Limit: params.Limit,
				Page:  params.Page,
			})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed to list todos", l.Error(err))
			return
		}
		c.JSON(http.StatusOK, response)
	}
}

// UpdateTodo updates todo by id
// route /v1/todos/{id} [put]
func (h *handlerV1) UpdateTodo(c *gin.Context) {
	var (
		body        pb.TodoFunc
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	body.Id = c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.TodoService().Update(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update todo", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteTodo deletes todo by id
// route /v1/todos/{id} [delete]
func (h *handlerV1) DeleteTodo(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	guid := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.TodoService().Delete(
		ctx, &pb.ByIdReq{
			Id: guid,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete todo", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}
