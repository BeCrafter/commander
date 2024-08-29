package helper

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortMapValues(t *testing.T) {
	type args struct {
		input map[string]interface{}
	}

	inStr := `{"code":0,"data":{"list":[{"category_id":0,"detail":["https://static0.xx.com/scm/goods-platform/1700118532205_image.png"],"head_picture":["https://static0.xx.com/scm/goods-platform/1700118525742_image.png"],"id":47582,"product_list":[{"RuleDesc":"","code_num":"SP00100001797","desc":"与大师同行礼盒（赠笔记本）","entity":[{"id":2462,"name":"【自采】与大师同行主题礼盒 新"},{"name":"笔记本套装","id":1298},{"id":2461,"name":"【自采】与大师同行主题礼盒袋 新"}],"extra_info":{},"name":"与大师同行礼盒（赠笔记本）","package_id":1296,"package_info":{"package_id":0,"package_type":0},"product_id":47581,"sale_price_cent":490,"sku_entity":{"1091":[{"id":1298,"name":"笔记本套装"}],"2250":[{"id":2461,"name":"与大师同行主题礼盒袋 新"}],"2251":[{"id":2462,"name":"与大师同行主题礼盒 新"}]},"sku_id":2250,"sku_ids":[1091,2251,2250],"source_type":1,"spec_info":[]}],"video":[]}],"total":1},"stat":1,"traceId":"7131e175-2c77-4a43-a30c-0005357ef435"}`
	var input map[string]interface{}
	json.Unmarshal([]byte(inStr), &input)

	outStr := `{"code":0,"data":{"list":[{"category_id":0,"detail":["https://static0.xx.com/scm/goods-platform/1700118532205_image.png"],"head_picture":["https://static0.xx.com/scm/goods-platform/1700118525742_image.png"],"id":47582,"product_list":[{"RuleDesc":"","code_num":"SP00100001797","desc":"与大师同行礼盒（赠笔记本）","entity":[{"id":1298,"name":"笔记本套装"},{"id":2461,"name":"【自采】与大师同行主题礼盒袋 新"},{"id":2462,"name":"【自采】与大师同行主题礼盒 新"}],"extra_info":{},"name":"与大师同行礼盒（赠笔记本）","package_id":1296,"package_info":{"package_id":0,"package_type":0},"product_id":47581,"sale_price_cent":490,"sku_entity":{"1091":[{"id":1298,"name":"笔记本套装"}],"2250":[{"id":2461,"name":"与大师同行主题礼盒袋 新"}],"2251":[{"id":2462,"name":"与大师同行主题礼盒 新"}]},"sku_id":2250,"sku_ids":[1091,2250,2251],"source_type":1,"spec_info":[]}],"video":[]}],"total":1},"stat":1,"traceId":"7131e175-2c77-4a43-a30c-0005357ef435"}
`
	var output map[string]interface{}
	json.Unmarshal([]byte(outStr), &output)

	tests := []struct {
		name     string
		args     args
		expected map[string]interface{}
	}{
		{
			name: "sort string slice",
			args: args{
				input: map[string]interface{}{
					"key1": []string{"c", "a", "b"},
				},
			},
			expected: map[string]interface{}{
				"key1": []string{"a", "b", "c"},
			},
		},
		{
			name: "sort int slice",
			args: args{
				input: map[string]interface{}{
					"key1": []float64{3, 1, 2}, // 此处不能定义为其他类型，受限于json.Unmarshal的限制，只能定义为float64
				},
			},
			expected: map[string]interface{}{
				"key1": []float64{1, 2, 3},
			},
		},
		{
			name: "sort map in slice",
			args: args{
				input: map[string]interface{}{
					"key1": []interface{}{
						map[string]interface{}{"sort_key": "c"},
						map[string]interface{}{"sort_key": "a"},
						map[string]interface{}{"sort_key": "b"},
					},
				},
			},
			expected: map[string]interface{}{
				"key1": []interface{}{
					map[string]interface{}{"sort_key": "a"},
					map[string]interface{}{"sort_key": "b"},
					map[string]interface{}{"sort_key": "c"},
				},
			},
		},
		{
			name: "nested map",
			args: args{
				input: map[string]interface{}{
					"key1": map[string]interface{}{
						"key11": "a",
						"key12": []string{"c", "b", "a"},
						"key13": map[string]interface{}{
							"key131": []float64{2, 3, 1},
						},
					},
				},
			},
			expected: map[string]interface{}{
				"key1": map[string]interface{}{
					"key11": "a",
					"key12": []string{"a", "b", "c"},
					"key13": map[string]interface{}{
						"key131": []float64{1, 2, 3},
					},
				},
			},
		},
		{
			name: "Multi level map sorting",
			args: args{
				input: input,
			},
			expected: output,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.EqualValues(t, tt.expected, SortMapValuesByMap(tt.args.input))
		})
	}
}
