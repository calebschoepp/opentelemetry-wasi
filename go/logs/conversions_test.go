package logs_test

import (
	"encoding/json"
	"testing"

	"github.com/calebschoepp/opentelemetry-wasi/logs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	logApi "go.opentelemetry.io/otel/log"
)

// TestOtelLogValueToJson tests that the logApi.Value serializes into a JSON string correctly.
func TestOtelLogValueToJson(t *testing.T) {
	// Create a map with all different value types
	kvMap := map[string]logApi.Value{
		"key1": logApi.BoolValue(false),
		"key2": logApi.Float64Value(123.456),
		"key3": logApi.Int64Value(41),
		"key4": logApi.BytesValue([]byte("Hello, world!")),
		"key5": logApi.StringValue("This is a string"),
		"key6": logApi.SliceValue(
			logApi.Int64Value(1),
			logApi.Int64Value(2),
			logApi.Int64Value(3),
		),
		"key7": logApi.MapValue(
			logApi.String("nestedkey1", "Hello, from within!"),
		),
	}

	mapValue := logApi.MapValue(convertMapToKeyValues(kvMap)...)
	jsonStr := logs.OtelLogValueToJson(mapValue)

	var actual map[string]any
	err := json.Unmarshal([]byte(jsonStr), &actual)
	require.NoError(t, err, "Failed to unmarshal JSON")

	expected := map[string]any{
		"key1": false,
		"key2": 123.456,
		"key3": float64(41), // JSON numbers are float64
		// 'Hello, world!' encoded to base64
		"key4": "{base64}:SGVsbG8sIHdvcmxkIQ==",
		"key5": "This is a string",
		"key6": []any{float64(1), float64(2), float64(3)},
		"key7": map[string]any{
			"nestedkey1": "Hello, from within!",
		},
	}

	assert.Equal(t, expected, actual, "Serialized JSON does not match expected structure")
}

// convertMapToKeyValues is a helper function to convert a map to KeyValue pairs
func convertMapToKeyValues(m map[string]logApi.Value) []logApi.KeyValue {
	kvs := make([]logApi.KeyValue, 0, len(m))
	for k, v := range m {
		kvs = append(kvs, logApi.KeyValue{
			Key:   k,
			Value: v,
		})
	}
	return kvs
}
