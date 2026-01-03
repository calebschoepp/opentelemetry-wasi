import { AnyValue } from '@opentelemetry/api-logs';
import { logAnyValueToWasi } from './logs';

describe('logAnyValueToWasi', () => {
  /**
   * Tests that the AnyValue serializes into a string correctly.
   */
  test('serialize_otel_log_any_value_to_string', () => {
    const testMap: Record<string, AnyValue> = {
      key1: false,
      key2: 123.456,
      key3: 41,
      key4: new Uint8Array(Buffer.from('Hello, world!', 'utf-8')),
      key5: 'This is a string',
      key6: [1, 2, 3],
      key7: {
        nestedkey1: 'Hello, from within!',
      },
    };

    const result: Record<string, AnyValue> = {};
    for (const [key, value] of Object.entries(testMap)) {
      const serialized = logAnyValueToWasi(value);
      result[key] = JSON.parse(serialized);
    }

    const expected = {
      key1: false,
      key2: 123.456,
      key3: 41,
      // 'Hello, world!' encoded to base64
      key4: '{base64}:SGVsbG8sIHdvcmxkIQ==',
      key5: 'This is a string',
      key6: [1, 2, 3],
      key7: {
        nestedkey1: 'Hello, from within!',
      },
    };

    expect(result).toEqual(expected);
  });
});
