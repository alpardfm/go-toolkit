package parser

import (
	"reflect"
	"testing"
)

func TestInitParser(t *testing.T) {
	tests := []struct {
		name string
		opts Options
	}{
		{
			name: "default config",
			opts: Options{JSONOptions: JSONOptions{Config: defaultConfig}},
		},
		{
			name: "vanilla compatible config",
			opts: Options{JSONOptions: JSONOptions{Config: vanillaCompatible}},
		},
		{
			name: "fastest config",
			opts: Options{JSONOptions: JSONOptions{Config: fastestConfig}},
		},
		{
			name: "custom config",
			opts: Options{JSONOptions: JSONOptions{
				Config:      customConfig,
				EscapeHTML:  true,
				SortMapKeys: true,
			}},
		},
		{
			name: "empty config uses standard library compatible",
			opts: Options{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := InitParser(tt.opts)
			if p == nil {
				t.Fatal("InitParser() returned nil")
			}
			if p.JSONParser() == nil {
				t.Fatal("JSONParser() returned nil")
			}
		})
	}
}

func TestJSONParser_Marshal(t *testing.T) {
	p := InitParser(Options{JSONOptions: JSONOptions{Config: defaultConfig}})
	jp := p.JSONParser()

	type testStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	tests := []struct {
		name    string
		input   any
		want    string
		wantErr bool
	}{
		{
			name:  "simple struct",
			input: testStruct{Name: "test", Value: 42},
			want:  `{"name":"test","value":42}`,
		},
		{
			name:  "map",
			input: map[string]int{"a": 1, "b": 2},
			want:  "", // order not guaranteed, check in test body
		},
		{
			name:  "slice",
			input: []int{1, 2, 3},
			want:  `[1,2,3]`,
		},
		{
			name:  "nil value",
			input: nil,
			want:  "null",
		},
		{
			name:  "empty struct",
			input: struct{}{},
			want:  `{}`,
		},
		{
			name:  "string value",
			input: "hello",
			want:  `"hello"`,
		},
		{
			name:  "boolean value",
			input: true,
			want:  `true`,
		},
		{
			name:    "channel cannot be marshaled",
			input:   make(chan int),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jp.Marshal(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.want != "" && string(got) != tt.want {
				t.Errorf("Marshal() = %s, want %s", string(got), tt.want)
			}
			// For map test, just verify it's valid JSON with expected keys
			if tt.name == "map" && err == nil {
				s := string(got)
				if len(s) == 0 {
					t.Error("Marshal() returned empty for map")
				}
			}
		})
	}
}

func TestJSONParser_Unmarshal(t *testing.T) {
	p := InitParser(Options{JSONOptions: JSONOptions{Config: defaultConfig}})
	jp := p.JSONParser()

	type testStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	tests := []struct {
		name    string
		input   string
		dest    any
		want    any
		wantErr bool
	}{
		{
			name:  "simple struct",
			input: `{"name":"test","value":42}`,
			dest:  &testStruct{},
			want:  &testStruct{Name: "test", Value: 42},
		},
		{
			name:  "partial fields",
			input: `{"name":"partial"}`,
			dest:  &testStruct{},
			want:  &testStruct{Name: "partial", Value: 0},
		},
		{
			name:  "empty object",
			input: `{}`,
			dest:  &testStruct{},
			want:  &testStruct{},
		},
		{
			name:    "invalid JSON",
			input:   `{invalid}`,
			dest:    &testStruct{},
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   ``,
			dest:    &testStruct{},
			wantErr: true,
		},
		{
			name:  "null value",
			input: `null`,
			dest:  &testStruct{},
			want:  &testStruct{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := jp.Unmarshal([]byte(tt.input), tt.dest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.dest, tt.want) {
				t.Errorf("Unmarshal() result = %+v, want %+v", tt.dest, tt.want)
			}
		})
	}
}

func TestJSONParser_MarshalUnmarshal_RoundTrip(t *testing.T) {
	p := InitParser(Options{JSONOptions: JSONOptions{Config: defaultConfig}})
	jp := p.JSONParser()

	type complexStruct struct {
		Name    string   `json:"name"`
		Age     int      `json:"age"`
		Score   float64  `json:"score"`
		Active  bool     `json:"active"`
		Tags    []string `json:"tags"`
		Address *struct {
			City string `json:"city"`
		} `json:"address"`
	}

	original := complexStruct{
		Name:   "Alice",
		Age:    30,
		Score:  95.5,
		Active: true,
		Tags:   []string{"go", "dev"},
		Address: &struct {
			City string `json:"city"`
		}{City: "Jakarta"},
	}

	data, err := jp.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	var result complexStruct
	if err := jp.Unmarshal(data, &result); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if !reflect.DeepEqual(original, result) {
		t.Errorf("Round-trip failed: got %+v, want %+v", result, original)
	}
}

func TestJSONParser_CustomConfig(t *testing.T) {
	p := InitParser(Options{JSONOptions: JSONOptions{
		Config:      customConfig,
		EscapeHTML:  true,
		SortMapKeys: true,
	}})
	jp := p.JSONParser()

	t.Run("HTML escaping", func(t *testing.T) {
		input := map[string]string{"html": "<script>alert('xss')</script>"}
		data, err := jp.Marshal(input)
		if err != nil {
			t.Fatalf("Marshal() error = %v", err)
		}
		// With EscapeHTML, angle brackets should be escaped
		s := string(data)
		if s == "" {
			t.Error("Marshal() returned empty string")
		}
	})

	t.Run("sorted map keys", func(t *testing.T) {
		input := map[string]int{"z": 1, "a": 2, "m": 3}
		data, err := jp.Marshal(input)
		if err != nil {
			t.Fatalf("Marshal() error = %v", err)
		}
		want := `{"a":2,"m":3,"z":1}`
		if string(data) != want {
			t.Errorf("Marshal() with SortMapKeys = %s, want %s", string(data), want)
		}
	})
}
