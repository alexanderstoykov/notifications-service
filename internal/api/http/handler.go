package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/alexanderstoykov/notifications-service/internal/storage"
)

type NotificationRepository interface {
	InsertOne(ctx context.Context, notification *storage.Notification) error
}

type Handler struct {
	notificationsRepository NotificationRepository
}

func NewHandler(notificationRepository NotificationRepository) *Handler {
	return &Handler{
		notificationsRepository: notificationRepository,
	}
}

func (h *Handler) HandleSlackNotification(ctx *gin.Context) {
	decoder := json.NewDecoder(ctx.Request.Body)
	defer func() {
		if err := ctx.Request.Body.Close(); err != nil {
			log.Printf("could not close body: %s", err)
		}
	}()

	var message SlackRequest
	if err := decoder.Decode(&message); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})

		return
	}

	if err := validator.New().Struct(&message); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})

		return
	}

	if err := h.notificationsRepository.InsertOne(ctx, &storage.Notification{
		ID:        uuid.New(),
		Type:      storage.NotificationTypeSlack,
		Message:   message.Message,
		Status:    storage.NotificationStatusPending,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func (h *Handler) HandleSMSNotification(ctx *gin.Context) {
	decoder := json.NewDecoder(ctx.Request.Body)
	defer func() {
		if err := ctx.Request.Body.Close(); err != nil {
			log.Printf("could not close body: %s", err)
		}
	}()

	var message SMSRequest
	if err := decoder.Decode(&message); err != nil {
		//http.Error(, `{"error": "Bad request"}`, http.StatusBadRequest)
		//todo gin error
		return
	}

	if err := validator.New().Struct(&message); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})

		return
	}

	if err := h.notificationsRepository.InsertOne(ctx, &storage.Notification{
		ID:        uuid.New(),
		Type:      storage.NotificationTypeSMS,
		Message:   message.Message,
		Receiver:  message.Receiver,
		Status:    storage.NotificationStatusPending,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func (h *Handler) HandleEmailNotification(ctx *gin.Context) {
	decoder := json.NewDecoder(ctx.Request.Body)
	defer func() {
		if err := ctx.Request.Body.Close(); err != nil {
			log.Printf("could not close body: %s", err)
		}
	}()

	var message EmailRequest
	if err := decoder.Decode(&message); err != nil {
		//http.Error(, `{"error": "Bad request"}`, http.StatusBadRequest)
		//todo gin error
		return
	}

	if err := validator.New().Struct(&message); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})

		return
	}

	if err := h.notificationsRepository.InsertOne(ctx, &storage.Notification{
		ID:        uuid.New(),
		Type:      storage.NotificationTypeEmail,
		Message:   message.Message,
		Receiver:  message.Receiver,
		Status:    storage.NotificationStatusPending,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
