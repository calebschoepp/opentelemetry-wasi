import {
  Datetime as WasiDatetime,
  InstrumentationScope as WasiInstrumentationScope,
} from 'wasi:otel/tracing@0.2.0-draft';
import { Attributes, AttributeValue, HrTime } from '@opentelemetry/api';
import { KeyValue, Value } from 'wasi:otel/types@0.2.0-draft';
import { InstrumentationScope } from '@opentelemetry/core';

/**
 * Converts OpenTelemetry Attributes to WASI Attributes
 */
export function attributesToWasi(attrs: Attributes | undefined): KeyValue[] {
  if (attrs === undefined) {
    return [];
  }
  return Object.entries(attrs).map(([key, value]) => ({
    key,
    value: attributeValueToWasi(value),
  }));
}

/**
 * Converts OpenTelemetry AttributeValue to WASI AttributeValue
 */
function attributeValueToWasi(value: AttributeValue | undefined): Value {
  if (typeof value === 'string') {
    return JSON.stringify(value);
  } else if (typeof value === 'boolean') {
    return JSON.stringify(value);
  } else if (typeof value === 'number') {
    return Number.isInteger(value)
      ? JSON.stringify(BigInt(value))
      : JSON.stringify(value);
  } else if (Array.isArray(value)) {
    const filtered: Array<string | boolean | number> = value.filter(
      (v) => v != null
    );
    if (filtered.length === 0) {
      return JSON.stringify([]);
    }
    const firstType = typeof filtered[0];
    if (firstType === 'string' || firstType == 'boolean') {
      return JSON.stringify(filtered);
    } else if (firstType === 'number') {
      const numbers = filtered as number[];
      if (numbers.every(Number.isInteger)) {
        const bigIntArray = new BigInt64Array(numbers.length);
        numbers.forEach((n, i) => (bigIntArray[i] = BigInt(n)));
        return JSON.stringify(bigIntArray);
      } else {
        return JSON.stringify(new Float64Array(numbers));
      }
    }
  }
  // Default
  return JSON.stringify('');
}

/**
 * Converts OpenTelemetry HrTime to WASI Datetime
 */
export function dateTimeToWasi(time: HrTime): WasiDatetime {
  return {
    seconds: BigInt(time[0]),
    nanoseconds: time[1],
  };
}

/**
 * Converts OpenTelemetry InstrumentationScope to WASI InstrumentationScope.
 */
export function instrumentationScopeToWasi(
  scope: InstrumentationScope
): WasiInstrumentationScope {
  return {
    name: scope.name,
    version: scope.version,
    schemaUrl: scope.schemaUrl,
    // Although other SDKs use the InstrumentationScope.attributes field, the `opentelemetry-js` SDK does not.
    // See https://github.com/open-telemetry/opentelemetry-js/blob/main/packages/opentelemetry-core/src/common/types.ts#L47
    attributes: [],
  };
}
