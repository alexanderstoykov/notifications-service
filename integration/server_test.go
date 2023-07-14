//go:build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"

	http_server "github.com/alexanderstoykov/notifications-service/internal/api/http"
	"github.com/alexanderstoykov/notifications-service/internal/storage"
)

func TestHandleNotificationRequests(t *testing.T) {
	ctx := context.Background()

	type messageRequest struct {
		Message  string `json:"message"`
		Receiver string `json:"receiver"`
	}

	type test struct {
		name             string
		notificationType storage.NotificationType
		receiver         string
		message          string
		uri              http_server.ServerURI
		expStatusCode    int
		exp              []*storage.Notification
	}

	tests := []test{
		{
			name:             "with slack notification",
			notificationType: storage.NotificationTypeSlack,
			message:          "slack message",
			uri:              http_server.ServerURISlack,
			expStatusCode:    200,

			exp: []*storage.Notification{
				{
					Type:     storage.NotificationTypeSlack,
					Receiver: "",
					Message:  "slack message",
					Status:   storage.NotificationStatusPending,
				},
			},
		}, {
			name:             "with sms notification",
			notificationType: storage.NotificationTypeSMS,
			receiver:         "+359887726152",
			message:          "sms message",
			uri:              http_server.ServerURISms,
			expStatusCode:    200,
			exp: []*storage.Notification{
				{
					Type:     storage.NotificationTypeSMS,
					Receiver: "+359887726152",
					Message:  "sms message",
					Status:   storage.NotificationStatusPending,
				},
			},
		}, {
			name:             "with email notification",
			notificationType: storage.NotificationTypeEmail,
			receiver:         "test@test.com",
			message:          "email message",
			uri:              http_server.ServerURIEmail,
			expStatusCode:    200,
			exp: []*storage.Notification{
				{
					Type:     storage.NotificationTypeEmail,
					Receiver: "test@test.com",
					Message:  "email message",
					Status:   storage.NotificationStatusPending,
				},
			},
		},
		{
			name:             "with failed validation",
			notificationType: storage.NotificationTypeEmail,
			receiver:         "",
			message:          "email message",
			uri:              http_server.ServerURIEmail,
			expStatusCode:    400,
			exp:              nil,
		},
	}

	router := gin.Default()
	handler := http_server.NewHandler(notificationRepo)
	http_server.RegisterRoutes(router, handler)

	for _, tc := range tests {
		purgeTables(ctx)

		message := messageRequest{
			Message:  tc.message,
			Receiver: tc.receiver,
		}

		w := httptest.NewRecorder()
		jsonData, err := json.Marshal(message)

		req, err := http.NewRequest("POST", string(tc.uri), bytes.NewBuffer(jsonData))

		router.ServeHTTP(w, req)

		got, err := listAllNotifications(ctx)

		// assertions
		assert.NoError(t, err)
		assert.Equal(t, tc.expStatusCode, w.Code)

		if diff := deep.Equal(tc.exp, got); diff != nil {
			t.Error(diff)
		}
	}
}
