package requestshelper

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/erupshis/effective_mobile/internal/datastructs"
)

func TestFilterPageNumAndPageSize(t *testing.T) {
	type args struct {
		values map[string]interface{}
	}
	tests := []struct {
		name  string
		args  args
		want  map[string]interface{}
		want1 int64
		want2 int64
	}{
		{
			name: "valid base case",
			args: args{
				values: map[string]interface{}{
					"page_num":  4,
					"page_size": 5,
					"name":      "asd",
				},
			},
			want: map[string]interface{}{
				"name": "asd",
			},
			want1: 4,
			want2: 5,
		},
		{
			name: "valid wrong types",
			args: args{
				values: map[string]interface{}{
					"page_num":  "dd",
					"page_size": "ff",
					"name":      "asd",
				},
			},
			want: map[string]interface{}{
				"name": "asd",
			},
			want1: 0,
			want2: 0,
		},
		{
			name: "valid missing",
			args: args{
				values: map[string]interface{}{
					"name": "asd",
				},
			},
			want: map[string]interface{}{
				"name": "asd",
			},
			want1: 0,
			want2: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := FilterPageNumAndPageSize(tt.args.values)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterPageNumAndPageSize() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("FilterPageNumAndPageSize() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("FilterPageNumAndPageSize() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestFilterValues(t *testing.T) {
	type args struct {
		values map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "valid base case",
			args: args{
				values: map[string]interface{}{
					"id":         1,
					"name":       "namev",
					"surname":    "surnamev",
					"patronymic": "patronymicv",
					"age":        12,
					"gender":     "genderv",
					"country":    "countryv",
					"page_num":   1,
					"page_size":  2,
				},
			},
			want: map[string]interface{}{
				"id":         "1",
				"name":       "namev",
				"surname":    "surnamev",
				"patronymic": "patronymicv",
				"age":        "12",
				"gender":     "genderv",
				"country":    "countryv",
				"page_num":   "1",
				"page_size":  "2",
			},
		},
		{
			name: "valid uppercase convert to lower",
			args: args{
				values: map[string]interface{}{
					"name": "nameV",
				},
			},
			want: map[string]interface{}{
				"name": "namev",
			},
		},
		{
			name: "valid removes unknown fields",
			args: args{
				values: map[string]interface{}{
					"name":    "namev",
					"unknown": "hz",
				},
			},
			want: map[string]interface{}{
				"name": "namev",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterValues(tt.args.values); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsPersonDataValid(t *testing.T) {
	type args struct {
		data             *datastructs.PersonData
		allFieldsToCheck bool
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid base case",
			args: args{
				data: &datastructs.PersonData{
					Name:       "nameVal",
					Surname:    "surnameVal",
					Patronymic: "patronymicVal",
					Age:        30,
					Gender:     "genderVal",
					Country:    "countryVal",
				},
				allFieldsToCheck: true,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "valid without patronymic",
			args: args{
				data: &datastructs.PersonData{
					Name:    "nameVal",
					Surname: "surnameVal",
					Age:     30,
					Gender:  "genderVal",
					Country: "countryVal",
				},
				allFieldsToCheck: true,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "valid missing non-critical value with allFieldsToCheck=false",
			args: args{
				data: &datastructs.PersonData{
					Name:       "nameVal",
					Surname:    "surnameVal",
					Patronymic: "patronymicVal",
					Age:        30,
					Gender:     "genderVal",
				},
				allFieldsToCheck: false,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid missing non-critical value with allFieldsToCheck=true",
			args: args{
				data: &datastructs.PersonData{
					Name:       "nameVal",
					Surname:    "surnameVal",
					Patronymic: "patronymicVal",
					Age:        30,
					Gender:     "genderVal",
				},
				allFieldsToCheck: true,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "valid negative/zero age with allFieldsToCheck=false",
			args: args{
				data: &datastructs.PersonData{
					Name:       "nameVal",
					Surname:    "surnameVal",
					Patronymic: "patronymicVal",
					Age:        0,
					Gender:     "genderVal",
					Country:    "countryVal",
				},
				allFieldsToCheck: false,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid negative/zero age with allFieldsToCheck=true",
			args: args{
				data: &datastructs.PersonData{
					Name:       "nameVal",
					Surname:    "surnameVal",
					Patronymic: "patronymicVal",
					Age:        0,
					Gender:     "genderVal",
					Country:    "countryVal",
				},
				allFieldsToCheck: true,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid missing critical value with allFieldsToCheck=true",
			args: args{
				data: &datastructs.PersonData{
					Surname:    "surnameVal",
					Patronymic: "patronymicVal",
					Age:        30,
					Country:    "countryVal",
				},
				allFieldsToCheck: true,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid missing critical value with allFieldsToCheck=false",
			args: args{
				data: &datastructs.PersonData{
					Surname:    "surnameVal",
					Patronymic: "patronymicVal",
					Age:        30,
					Gender:     "genderVal",
					Country:    "countryVal",
				},
				allFieldsToCheck: false,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid all fields empty",
			args: args{
				data:             &datastructs.PersonData{},
				allFieldsToCheck: false,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsPersonDataValid(tt.args.data, tt.args.allFieldsToCheck)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsPersonDataValid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsPersonDataValid() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParsePersonDataFromJSON(t *testing.T) {
	type args struct {
		rawData []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *datastructs.PersonData
		wantErr bool
	}{
		{
			name: "valid base case",
			args: args{
				rawData: []byte(
					`{
						"name": "Name",
						"surname": "Surname",
						"patronymic": "Patronymic",
						"age": 20,
						"gender": "male",
						"country": "rus"
					}`,
				),
			},
			want: &datastructs.PersonData{
				Name:       "Name",
				Surname:    "Surname",
				Patronymic: "Patronymic",
				Age:        20,
				Gender:     "male",
				Country:    "rus",
			},
			wantErr: false,
		},
		{
			name: "invalid broken json structure",
			args: args{
				rawData: []byte(
					`{
						"name": "Name",
						"surname": "Surname",
						"patronymic": "Patronymic",
						"ag`,
				),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePersonDataFromJSON(tt.args.rawData)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePersonDataFromJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParsePersonDataFromJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseQueryValuesIntoMap(t *testing.T) {
	type args struct {
		values url.Values
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid base case",
			args: args{
				values: map[string][]string{
					"name": {"name"},
				},
			},
			want: map[string]interface{}{
				"name": "name",
			},
			wantErr: false,
		},
		{
			name: "valid with uppercase in values",
			args: args{
				values: map[string][]string{
					"name":    {"name"},
					"surname": {"Surname"},
				},
			},
			want: map[string]interface{}{
				"name":    "name",
				"surname": "surname",
			},
			wantErr: false,
		},
		{
			name: "valid with unknown key in values",
			args: args{
				values: map[string][]string{
					"name":    {"name"},
					"unknown": {"unknown"},
				},
			},
			want: map[string]interface{}{
				"name": "name",
			},
			wantErr: false,
		},
		{
			name: "invalid with unknown key in values only",
			args: args{
				values: map[string][]string{
					"unknown": {"unknown"},
				},
			},
			want:    map[string]interface{}{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseQueryValuesIntoMap(tt.args.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseQueryValuesIntoMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseQueryValuesIntoMap() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertQueryValueIntoInt64(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "valid int64 type as input",
			args: args{
				value: 355,
			},
			want: 355,
		},
		{
			name: "valid int type as input",
			args: args{
				value: 355,
			},
			want: 355,
		},
		{
			name: "valid float64 type as input",
			args: args{
				value: 355.5,
			},
			want: 355,
		},
		{
			name: "valid int in string type as input",
			args: args{
				value: "355",
			},
			want: 355,
		},
		{
			name: "valid float64 in string type as input",
			args: args{
				value: "355.5",
			},
			want: 355,
		},
		{
			name: "invalid real string type as input",
			args: args{
				value: "wewffg",
			},
			want: 0,
		},
		{
			name: "invalid bool type as input",
			args: args{
				value: true,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertQueryValueIntoInt64(tt.args.value); got != tt.want {
				t.Errorf("convertQueryValueIntoInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}
