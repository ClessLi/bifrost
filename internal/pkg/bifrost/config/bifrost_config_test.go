package config

import (
	"testing"
)

func TestConfig_check(t *testing.T) {
	type fields struct {
		ServiceConfig ServiceConfig
		RAConfig      *RAConfig
		LogConfig     LogConfig
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				ServiceConfig: tt.fields.ServiceConfig,
				RAConfig:      tt.fields.RAConfig,
				LogConfig:     tt.fields.LogConfig,
			}
			if err := c.check(); (err != nil) != tt.wantErr {
				t.Errorf("check() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
