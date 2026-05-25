package format

import (
	"testing"
	"time"
)

func TestTimeParseWithDefaultFormat(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "valid date time",
			value:   "15/01/2024 10:30:45",
			want:    time.Date(2024, 1, 15, 10, 30, 45, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "midnight",
			value:   "01/01/2000 00:00:00",
			want:    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "end of day",
			value:   "31/12/2023 23:59:59",
			want:    time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "invalid format - ISO 8601",
			value:   "2024-01-15T10:30:45Z",
			wantErr: true,
		},
		{
			name:    "empty string",
			value:   "",
			wantErr: true,
		},
		{
			name:    "date only without time",
			value:   "15/01/2024",
			wantErr: true,
		},
		{
			name:    "random string",
			value:   "not a date",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TimeParseWithDefaultFormat(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("TimeParseWithDefaultFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.want) {
				t.Errorf("TimeParseWithDefaultFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringParseTmpl(t *testing.T) {
	tests := []struct {
		name    string
		strFmt  string
		values  any
		want    string
		wantErr bool
	}{
		{
			name:   "simple template with struct",
			strFmt: "Hello, {{.Name}}!",
			values: struct{ Name string }{Name: "World"},
			want:   "Hello, World!",
		},
		{
			name:   "template with map",
			strFmt: "{{.greeting}} {{.target}}",
			values: map[string]string{"greeting": "Hi", "target": "there"},
			want:   "Hi there",
		},
		{
			name:   "template with multiple fields",
			strFmt: "{{.First}} {{.Last}} ({{.Age}})",
			values: struct {
				First string
				Last  string
				Age   int
			}{First: "John", Last: "Doe", Age: 30},
			want: "John Doe (30)",
		},
		{
			name:   "empty template",
			strFmt: "",
			values: nil,
			want:   "",
		},
		{
			name:    "invalid template syntax",
			strFmt:  "{{.Name",
			values:  struct{ Name string }{Name: "test"},
			wantErr: true,
		},
		{
			name:    "missing field in values",
			strFmt:  "{{.Missing}}",
			values:  struct{ Name string }{Name: "test"},
			wantErr: true,
		},
		{
			name:   "no placeholders",
			strFmt: "plain text",
			values: nil,
			want:   "plain text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StringParseTmpl(tt.strFmt, tt.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringParseTmpl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("StringParseTmpl() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatConstants(t *testing.T) {
	// Verify the format constants produce expected output
	now := time.Date(2024, 3, 15, 14, 30, 45, 123456789, time.UTC)

	tests := []struct {
		name   string
		format string
		want   string
	}{
		{
			name:   "DayMonthYear",
			format: DayMonthYear,
			want:   "15/03/2024",
		},
		{
			name:   "HourMinSec",
			format: HourMinSec,
			want:   "14:30:45",
		},
		{
			name:   "DayMonthYearHourMinSec",
			format: DayMonthYearHourMinSec,
			want:   "15/03/2024 14:30:45",
		},
		{
			name:   "DayMonthYearHourMinSecMilisec",
			format: DayMonthYearHourMinSecMilisec,
			want:   "15/03/2024 14:30:45.123456789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := now.Format(tt.format)
			if got != tt.want {
				t.Errorf("time.Format(%q) = %q, want %q", tt.format, got, tt.want)
			}
		})
	}
}
