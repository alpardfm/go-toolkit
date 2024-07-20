package sorter

import (
	"fmt"
	"sort"
	"testing"

	"github.com/alpardfm/go-toolkit/files"
)

func Test_Sorter(t *testing.T) {
	type (
		user struct {
			id   int
			name string
		}

		params struct {
			items    []user
			lessFunc func(items []user, i, j int) bool
		}

		wantResult struct {
			resultFunc func(result []user) bool
		}

		test struct {
			name       string
			params     params
			wantResult wantResult
		}
	)
	tests := []test{
		{
			name: "sort by id asc",
			params: params{
				items: []user{
					{id: 3, name: "c"},
					{id: 2, name: "b"},
					{id: 1, name: "a"},
				},
				lessFunc: func(items []user, i int, j int) bool {
					return items[i].id < items[j].id
				},
			},
			wantResult: wantResult{
				resultFunc: func(result []user) bool {
					return result[0].id == 1
				},
			},
		},
	}

	f := files.GetCurrentMethodName()
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v:%v", f, tt.name), func(t *testing.T) {
			srt := NewSorter(&tt.params.items, tt.params.lessFunc)
			sort.Sort(srt)

			if isSorted := tt.wantResult.resultFunc(tt.params.items); !isSorted {
				t.Fatalf("want is sorted but not")
			}
		})
	}
}
