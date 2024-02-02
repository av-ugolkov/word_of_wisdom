package server

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"testing"

	pkgCache "word_of_wisdom/internal/pkg/cache"
	"word_of_wisdom/internal/pkg/data"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestServer_processRequest(t *testing.T) {
	testData := `{"Version":1,"ZerosCount":5,"Date":1706815068,"Resource":"192.168.112.3:40488","Rand":"NmExOGMxNzktMmUzMi00YWU4LTg4ZmMtY2Y5MjczNDA1Nzcw","Counter":827471}`
	testUUID, _ := uuid.Parse("6a18c179-2e32-4ae8-88fc-cf9273405770")
	type fields struct {
		cache     cache
		listener  net.Listener
		firstZero int
	}
	type args struct {
		dataStr    string
		clientInfo string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *data.Data
		wantErr error
	}{
		{
			name: "correct challenge",
			fields: fields{
				cache:     pkgCache.NewCache(),
				firstZero: 4,
			},
			args: args{
				dataStr:    "challenge::word of wisdom",
				clientInfo: "client_ip",
			},
			wantErr: nil,
		},
		{
			name: "invalid data",
			fields: fields{
				cache:     pkgCache.NewCache(),
				firstZero: 4,
			},
			args: args{
				dataStr:    "challenge:word of wisdom",
				clientInfo: "client_ip",
			},
			wantErr: errors.New("invalid data"),
		},
		{
			name: "invalid key",
			fields: fields{
				cache:     pkgCache.NewCache(),
				firstZero: 4,
			},
			args: args{
				dataStr:    "chalenge::word of wisdom",
				clientInfo: "client_ip",
			},
			wantErr: fmt.Errorf("unknown key"),
		},
		{
			name: "correct response",
			fields: fields{
				cache:     pkgCache.NewCache(),
				firstZero: 4,
			},
			args: args{
				dataStr:    fmt.Sprintf("response::%s", testData),
				clientInfo: "192.168.112.3:40488",
			},
			wantErr: nil,
		},
		{
			name: "correct response",
			fields: fields{
				cache:     pkgCache.NewCache(),
				firstZero: 4,
			},
			args: args{
				dataStr:    fmt.Sprintf("response::%s", testData),
				clientInfo: "191.168.112.3:40488",
			},
			wantErr: fmt.Errorf("invalid hashcash resource"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				cache:           tt.fields.cache,
				firstZerosCount: tt.fields.firstZero,
			}
			s.cache.Add(testUUID)
			got, err := s.processRequest(tt.args.dataStr, tt.args.clientInfo)
			if tt.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tt.wantErr.Error()))
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}
