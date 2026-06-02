package security

import (
	"strings"

	"github.com/Ajulll22/payment-ai-assistant/pkg/constant"
)

func MaskSensitive(data map[string]any) map[string]any {
	for k, v := range data {
		// cek apakah key termasuk sensitif
		for _, sKey := range constant.SensitiveKeys {
			if strings.EqualFold(k, sKey) {
				data[k] = "*****" // masking
			}
		}

		// kalau value berupa nested object (map)
		if nested, ok := v.(map[string]any); ok {
			data[k] = MaskSensitive(nested)
		}

		// kalau value berupa array
		if arr, ok := v.([]any); ok {
			for i, item := range arr {
				if itemMap, ok := item.(map[string]any); ok {
					arr[i] = MaskSensitive(itemMap)
				}
			}
			data[k] = arr
		}
	}
	return data
}
