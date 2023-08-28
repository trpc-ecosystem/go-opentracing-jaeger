package jaeger

import (
	"context"
	"errors"
	"testing"

	"github.com/opentracing/opentracing-go"
	. "github.com/smartystreets/goconvey/convey"
	"trpc.group/trpc-go/trpc-go"
)

const configInfo = `
plugins:
  tracing:
    jaeger:                               
      serviceName: trpc-ecosystem
      disabled: false
      sampler:
        type: const
        param: 1
      reporter:
        localAgentHostPort: localhost:6831
`

func Test_jaegerPlugin_Type(t *testing.T) {
	tests := []struct {
		name string
		p    *jaegerPlugin
		want string
	}{
		{
			name: "test case",
			p:    new(jaegerPlugin),
			want: pluginType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &jaegerPlugin{}
			if got := p.Type(); got != tt.want {
				t.Errorf("jaegerPlugin.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_metadataTextMap_Set(t *testing.T) {
	type args struct {
		key string
		val string
	}
	tests := []struct {
		name string
		m    metadataTextMap
		args args
	}{
		{
			name: "test case",
			m:    metadataTextMap{},
			args: args{
				key: "key",
				val: "val",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Set(tt.args.key, tt.args.val)
			if string(tt.m[tt.args.key]) != tt.args.val {
				t.Error("err occur")
			}
		})
	}
}

func Test_metadataTextMap_ForeachKey(t *testing.T) {
	type args struct {
		callback func(key, val string) error
	}
	tests := []struct {
		name    string
		m       metadataTextMap
		args    args
		wantErr bool
	}{{
		name: "case fail",
		m:    metadataTextMap{"key": []byte("val")},
		args: args{
			callback: func(key, val string) error {
				return errors.New("error")
			},
		},
		wantErr: true,
	},
		{
			name: "case success",
			m:    metadataTextMap{},
			args: args{
				callback: func(key, val string) error {
					return nil
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.ForeachKey(tt.args.callback); (err != nil) != tt.wantErr {
				t.Errorf("metadataTextMap.ForeachKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServerFilter(t *testing.T) {
	Convey("TestServerFilter", t, func() {
		Convey("case success", func() {
			tracer := opentracing.GlobalTracer()
			f := ServerFilter(tracer)
			_, err := f(trpc.BackgroundContext(), "", func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, nil
			})
			So(err, ShouldBeNil)
		})
	})
}

func TestClientFilter(t *testing.T) {
	Convey("TestClientFilter", t, func() {
		Convey("case success", func() {
			tracer := opentracing.GlobalTracer()
			f := ClientFilter(tracer)
			err := f(trpc.BackgroundContext(), "", "", func(ctx context.Context, req interface{}, rsp interface{}) error {
				return nil
			})
			So(err, ShouldBeNil)
		})
	})
}

type mockDecoder struct {
}

func (d *mockDecoder) Decode(conf interface{}) error {
	return nil
}
