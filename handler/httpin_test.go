package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/luca-arch/go-goodies/handler"
	"github.com/stretchr/testify/assert"
)

type StructInt struct {
	IntNum   int   `in:"val"`
	Int16Num int16 `in:"val16"`
	Int32Num int32 `in:"val32"`
	Int64Num int64 `in:"val64"`
}

type StructPtr struct {
	BoolVal  *bool   `in:"bool"`
	IntNum   *int    `in:"val"`
	Int16Num *int16  `in:"val16"`
	Int32Num *int32  `in:"val32"`
	Int64Num *int64  `in:"val64"`
	String   *string `in:"valStr"`
}

type StructRequired struct {
	Param string `in:"sentence,required"`
}

func TestInputFromRequest(t *testing.T) {
	t.Parallel()

	var (
		intNum         = 10
		int16Num int16 = 20
		int32Num int32 = 30
		int64Num int64 = 40
		strVal         = "my string"
		trueVal        = true
	)

	type args struct {
		url string
	}

	type fields struct {
		call func(*http.Request) (any, error)
	}

	type wants struct {
		err string
		out any
	}

	tests := map[string]struct {
		args
		fields
		wants
	}{
		"Struct with numeric types": {
			args{
				url: "https://example.com/?val=10&val16=20&val32=30&val64=40",
			},
			fields{
				call: func(r *http.Request) (any, error) {
					return handler.InputFromRequest[StructInt](r)
				},
			},
			wants{
				out: StructInt{
					IntNum:   intNum,
					Int16Num: int16Num,
					Int32Num: int32Num,
					Int64Num: int64Num,
				},
			},
		},
		"ok - struct with required value": {
			args{
				url: "https://example.com/?sentence=my+string",
			},
			fields{
				call: func(r *http.Request) (any, error) {
					return handler.InputFromRequest[StructRequired](r)
				},
			},
			wants{
				out: StructRequired{
					Param: strVal,
				},
			},
		},
		"ok - struct with pointers": {
			args{
				url: "https://example.com/?bool=1&val=10&val32=30&val64=40&valStr=my+string",
			},
			fields{
				call: func(r *http.Request) (any, error) {
					return handler.InputFromRequest[StructPtr](r)
				},
			},
			wants{
				out: StructPtr{
					BoolVal:  &trueVal,
					IntNum:   &intNum,
					Int16Num: nil,
					Int32Num: &int32Num,
					Int64Num: &int64Num,
					String:   &strVal,
				},
			},
		},
		"error - struct with required value": {
			args{
				url: "https://example.com/",
			},
			fields{
				call: func(r *http.Request) (any, error) {
					return handler.InputFromRequest[StructRequired](r)
				},
			},
			wants{
				err: "invalid input\nmissing required field: sentence",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest(http.MethodGet, test.args.url, nil)

			out, err := test.fields.call(r)

			if test.wants.err != "" {
				assert.EqualError(t, err, test.wants.err)

				return
			}

			assert.Equal(t, test.wants.out, out)
		})
	}
}
