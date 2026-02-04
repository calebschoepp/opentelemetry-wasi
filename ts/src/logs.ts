import { LogRecordProcessor, SdkLogRecord } from '@opentelemetry/sdk-logs';
import { onEmit as emitToWasi } from 'wasi:otel/logs@0.2.0-draft';
import { LogRecord as WasiLogRecord } from 'wasi:otel/logs@0.2.0-draft';
import { KeyValue as WasiKeyValue } from 'wasi:otel/types@0.2.0-draft';
import { dateTimeToWasi, instrumentationScopeToWasi } from './types';
import { AnyValue, AnyValueMap } from '@opentelemetry/api-logs';

export class WasiLogProcessor implements LogRecordProcessor {
  onEmit(logRecord: SdkLogRecord): void {
    emitToWasi(logRecordToWasi(logRecord));
  }

  async forceFlush(): Promise<void> {
    // no-op
  }

  async shutdown(): Promise<void> {
    // no-op
  }
}

/**
 * Converts an OpenTelemetry log record to a WASI log record.
 */
function logRecordToWasi(r: SdkLogRecord): WasiLogRecord {
  return {
    timestamp: r.hrTime ? dateTimeToWasi(r.hrTime) : undefined,
    observedTimestamp: r.hrTimeObserved
      ? dateTimeToWasi(r.hrTimeObserved)
      : undefined,
    severityText: r.severityText,
    severityNumber: r.severityNumber,
    body: r.body ? logAnyValueToWasi(r.body) : undefined,
    attributes: r.attributes ? logAttributesToWasi(r.attributes) : undefined,
    eventName: r.eventName,
    resource: {
      attributes: logAttributesToWasi(r.resource.attributes),
      schemaUrl: r.resource.schemaUrl,
    },
    instrumentationScope: instrumentationScopeToWasi(r.instrumentationScope),
    traceId: undefined,
    spanId: undefined,
    traceFlags: undefined,
  };
}

/**
 * Converts OpenTelemetry log attributes to WASI attributes.
 */
function logAttributesToWasi(attrs: AnyValueMap): WasiKeyValue[] {
  const result: WasiKeyValue[] = [];
  for (const [k, v] of Object.entries(attrs)) {
    result.push({ key: k, value: logAnyValueToWasi(v) });
  }
  return result;
}

export function logAnyValueToWasi(v: AnyValue): string {
  if (v instanceof Uint8Array) {
    return JSON.stringify(
      'data:application/octet-stream;base64,' +
        Buffer.from(v).toString('base64')
    );
  } else if (v === null || v === undefined) {
    return 'null';
  } else if (typeof v === 'string') {
    return JSON.stringify(v);
  } else if (typeof v === 'number' || typeof v === 'boolean') {
    return String(v);
  } else if (Array.isArray(v)) {
    const items = v.map((item) => logAnyValueToWasi(item)).join(',');
    return `[${items}]`;
  } else {
    const pairs = Object.entries(v)
      .map(
        ([key, value]) => `${JSON.stringify(key)}:${logAnyValueToWasi(value)}`
      )
      .join(',');
    return `{${pairs}}`;
  }
}
