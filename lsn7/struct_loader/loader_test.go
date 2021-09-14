package struct_loader

import (
	"errors"
	"github.com/sirupsen/logrus"
	"reflect"
	"testing"
	"time"
)



func TestLoadSimpleStruct(t *testing.T) {
	type StructSimple struct {
		B  bool
		I  int64
		F  float64
		TD time.Duration
		S  string
	}
	type args struct {
		in              interface{}
		values          map[string]interface{}
		defaultRequired bool
	}
	tests := []struct {
		name string
		args args
		want *StructSimple
	}{
		{"StructSimple1_convert",
			args{&StructSimple{},
				map[string]interface{}{
					"B":  "true",
					"I":  "123",
					"F":  "123.45",
					"TD": "3s",
					"S":  "test",
				},
				false,
			},
			&StructSimple{
				B:  true,
				I:  123,
				F:  123.45,
				TD: 3 * time.Second,
				S:  "test",
			},
		},
		{"StructSimple1",
			args{&StructSimple{},
				map[string]interface{}{
					"B":  true,
					"I":  123,
					"F":  123.45,
					"TD": 3 * time.Second,
					"S":  "test",
				},
				false,
			},
			&StructSimple{
				B:  true,
				I:  123,
				F:  123.45,
				TD: 3 * time.Second,
				S:  "test",
			},
		},
		{"StructSimple2",
			args{&StructSimple{
				B:  true,
				I:  123,
				F:  123.45,
				TD: 3 * time.Second,
				S:  "test",
			},
				map[string]interface{}{},
				false,
			},
			&StructSimple{
				B:  true,
				I:  123,
				F:  123.45,
				TD: 3 * time.Second,
				S:  "test",
			},
		},
		{"StructSimple3",
			args{&StructSimple{
				B:  false,
				I:  1,
				F:  1,
				TD: 1 * time.Second,
				S:  "default",
			},
				map[string]interface{}{
					"B":  true,
					"I":  123,
					"F":  123.45,
					"TD": 3 * time.Second,
					"S":  "test",
				},
				false,
			},
			&StructSimple{
				B:  true,
				I:  123,
				F:  123.45,
				TD: 3 * time.Second,
				S:  "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err:= Load(tt.args.in, tt.args.values, tt.args.defaultRequired); (err != nil) || (!reflect.DeepEqual(tt.args.in, tt.want)) {
				t.Errorf("Load [%v], want [%v] with err: [%v]", tt.args.in, tt.want, err)
			}
		})
	}
}

func TestLoadSimpleStructErrors(t *testing.T) {
	type StructSimple struct {
		B  bool
		I  int64
		F  float64
		TD time.Duration
		S  string
	}
	type args struct {
		in              interface{}
		values          map[string]interface{}
		defaultRequired bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{"ErrorIsNil",
			args{nil,
				map[string]interface{}{},
				false,
			},
			ErrorIsNil,
		},
		{"ErrorInvalidReceiver",
			args{&map[string]string{"": ""},
				map[string]interface{}{},
				false,
			},
			ErrorInvalidReceiver,
		},
		{"ErrorNotConvertable",
			args{&StructSimple{},
				map[string]interface{}{
					"B":  "fail",
					"I":  "fail",
					"F":  "fail",
					"TD": "fail",
					"S":  interface{}(false),
				},
				false,
			},
			ErrorNotConvertable,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Load(tt.args.in, tt.args.values, tt.args.defaultRequired); !errors.Is(err, tt.wantErr) {
				t.Errorf("Load() error = [%v], wantErr [%v]", err, tt.wantErr)
			}
		})
	}
}

func TestLoadStructNamed(t *testing.T) {
	type StructNamed struct {
		B   bool `name:"b"`
		I int64   `name:"i"`
		F float64 `name:"f"`
		TD  time.Duration `name:"td"`
		S   string `name:"s"`
	}
	type args struct {
		in              interface{}
		values          map[string]interface{}
		defaultRequired bool
	}
	tests := []struct {
		name string
		args args
		want *StructNamed
	}{
		{"StructNamed1",
			args{&StructNamed{},
				map[string]interface{}{
					"B":  "true",
					"I":  "123",
					"F":  "123.45",
					"TD": "3s",
					"S":  "test",
				},
				false,
			},
			&StructNamed{},
		},
		{"StructNamed2",
			args{&StructNamed{},
				map[string]interface{}{
					"b":  true,
					"i":  123,
					"f":  123.45,
					"td": 3 * time.Second,
					"s":  "test",
				},
				false,
			},
			&StructNamed{
				B:  true,
				I:  123,
				F:  123.45,
				TD: 3 * time.Second,
				S:  "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err:= Load(tt.args.in, tt.args.values, tt.args.defaultRequired); (err != nil) || (!reflect.DeepEqual(tt.args.in, tt.want)) {
				t.Errorf("Load [%v], want [%v] with err: [%v]", tt.args.in, tt.want, err)
			}
		})
	}
}

func TestLoadStructDefault(t *testing.T) {
	type StructDefault struct {
		B   bool `default:"true"`
		I int64   `default:"12345"`
		F float64 `default:"123.45"`
		TD  time.Duration  `default:"8s"`
		S   string `default:"test"`
		SS []string `default:"[\"test1\",\"test2\",\"test,3\"]"`
		SI []int `default:"[1,2,3]"`
		SB []bool `default:"[\"true\",\"false\"]"`
		SF []float64 `default:"[1.23,2.34,3.45]"`
		SD []time.Duration `default:"[\"1s\",\"2ms\"]"`
		SSS [][]string `default:"[[\"test1\",\"test2\"],[\"test,3\"]]"`
		SSSS [][][]string `default:"[[[\"test1\",\"test2\"]],[[\"test,3\"]]]"`
		SSSI [][][]int `default:"[[[1,2]],[[3]]]"`
		S2S map[string]string `default:"{\"key1\":\"val1\",\"key2\":\"val1\"}"`
		S2I map[string]int `default:"{\"key1\":1,\"key2\":2}"`
		// todo: I2S map[int]string `default:"{\"1\":\"v1\",\"2\":\"v2\"}"`
		S2SI map[string][]int `default:"{\"key1\":[1,2],\"key2\":[3,4,5]}"`

	}
	type args struct {
		in              interface{}
		values          map[string]interface{}
		defaultRequired bool
	}
	tests := []struct {
		name string
		args args
		want *StructDefault
	}{
		{"StructDefault1",
			args{&StructDefault{},
				map[string]interface{}{},
				false,
			},
			&StructDefault{
				B: true,
				I: 12345,
				F: 123.45,
				TD: 8 * time.Second,
				S: "test",
				SS: []string{"test1","test2","test,3"},
				SI: []int{1,2,3},
				SB: []bool{true, false},
				SF: []float64{1.23,2.34,3.45},
				SD: []time.Duration{1*time.Second, 2* time.Millisecond},
				SSS: [][]string{{"test1", "test2"},{"test,3"}},
				SSSS: [][][]string{{{"test1","test2"}},{{"test,3"}}},
				SSSI: [][][]int{{{1,2}},{{3}}},
				S2S: map[string]string{"key1":"val1","key2":"val1"},
				S2I: map[string]int{"key1":1,"key2":2},
				//I2S: map[int]string{1:"v1",2:"v2"},
				S2SI: map[string][]int{"key1":{1,2},"key2":{3,4,5}},

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err:= Load(tt.args.in, tt.args.values, tt.args.defaultRequired); (err != nil) || (!reflect.DeepEqual(tt.args.in, tt.want)) {
				t.Errorf("Load [%v], want [%v] with err: [%v]", tt.args.in, tt.want, err)
			}
		})
	}
}

func TestLoadStructRequiredErrors(t *testing.T) {
	type StructRequired struct {
		B   bool `required:"true"`
		I int64   `required:"true"`
		F float64 `required:"true"`
		TD  time.Duration  `required:"true"`
		S   string `required:"true"`
		SS []string `required:"true"`
		S2S map[string]string `required:"true"`
	}
	type args struct {
		in              interface{}
		values          map[string]interface{}
		defaultRequired bool
	}
	tests := []struct {
		name string
		args args
		wantErr error
	}{
		{"StructRequired_Bool",
			args{&struct{B   bool `required:"true"`}{},
				map[string]interface{}{},
				false,
			},
			ErrorFieldRequired,
		},
		{"StructRequired_Int",
			args{&struct{I  int `required:"true"`}{},
				map[string]interface{}{},
				false,
			},
			ErrorFieldRequired,
		},
		{"StructRequired_Float",
			args{&struct{F   float64 `required:"true"`}{},
				map[string]interface{}{},
				false,
			},
			ErrorFieldRequired,
		},
		{"StructRequired_TD",
			args{&struct{Td time.Duration `required:"true"`}{},
				map[string]interface{}{},
				false,
			},
			ErrorFieldRequired,
		},
		{"StructRequired_Slice",
			args{&struct{SS  []string `required:"true"`}{},
				map[string]interface{}{},
				false,
			},
			ErrorFieldRequired,
		},
		{"StructRequired_Map",
			args{&struct{S2S  map[string]string `required:"true"`}{},
				map[string]interface{}{},
				false,
			},
			ErrorFieldRequired,
		},
		{"StructRequired_Struct",
			args{&struct{Struct struct{S string} `required:"true"`}{},
				map[string]interface{}{},
				false,
			},
			ErrorFieldRequired,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Load(tt.args.in, tt.args.values, tt.args.defaultRequired); !errors.Is(err, tt.wantErr) {
				t.Errorf("Load() error = [%v], wantErr [%v]", err, tt.wantErr)
			}
		})
	}
}


func TestLoadStructNested(t *testing.T) {
	type LogStruct struct {
		Path  string `name:"path" required:"true" default:"~/go"`
		Level logrus.Level `name:"level" default:"info"` // todo: add logrus.Level parsing
	}
	type StructNested struct {
		LogCfg LogStruct `name:"log" required:"true"`  //note: default should be defined in nested struct if needed
	}
	type args struct {
		in              interface{}
		values          map[string]interface{}
		defaultRequired bool
	}
	tests := []struct {
		name string
		args args
		want *StructNested
	}{
		{"StructNested1",
			args{&StructNested{},
				map[string]interface{}{
					"log": LogStruct{Path: "some path", Level: logrus.WarnLevel},
				},
				false,
			},
			&StructNested{
				LogCfg: LogStruct{Path: "some path", Level: logrus.WarnLevel},
			},
		},
		{"StructNested2",
			args{&StructNested{},
				map[string]interface{}{
					"log.path": "some path",
					"log.level": logrus.WarnLevel,
				},
				false,
			},
			&StructNested{
				LogCfg: LogStruct{Path: "some path", Level: logrus.WarnLevel},
			},
		},
		{"StructNested3",
			args{&StructNested{},
				map[string]interface{}{
					"log.level": logrus.WarnLevel,
				},
				false,
			},
			&StructNested{
				LogCfg: LogStruct{Path: "~/go", Level: logrus.WarnLevel},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err:= Load(tt.args.in, tt.args.values, tt.args.defaultRequired); (err != nil) || (!reflect.DeepEqual(tt.args.in, tt.want)) {
				t.Errorf("Load [%v], want [%v] with err: [%v]", tt.args.in, tt.want, err)
			}
		})
	}
}
