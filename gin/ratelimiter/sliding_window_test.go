package ratelimiter

import (
	"testing"
	"time"
)

func Test_slidingWindow_allow(t *testing.T) {
	type fields struct {
		windowDuration time.Duration
		limit          int64
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "positive",
			fields: fields{
				windowDuration: time.Minute,
				limit:          100,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSlidingWindow(tt.fields.limit, tt.fields.windowDuration)
			if got := s.allow(); got != tt.want {
				t.Errorf("slidingWindow.allow() = %v, want %v", got, tt.want)
			}
		})
	}
}
