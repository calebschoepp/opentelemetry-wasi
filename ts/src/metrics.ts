import {
  AggregationTemporality,
  AggregationType,
  MetricReader,
} from '@opentelemetry/sdk-metrics';
import { export as exportToWasiHost } from 'wasi:otel/metrics@0.2.0-draft';
import { diag } from '@opentelemetry/api';
import { resourceMetricsToWasi } from './types';

export class WasiMetricExporter extends MetricReader {
  constructor() {
    super({
      aggregationSelector: (_instrumentType) => {
        return {
          type: AggregationType.DEFAULT,
        };
      },
      aggregationTemporalitySelector: (_instrumentType) =>
        AggregationTemporality.CUMULATIVE,
    });
  }

  protected override async onForceFlush(): Promise<void> {
    // no-op
  }

  protected override async onShutdown(): Promise<void> {
    // no-op
  }

  /**
   * Exports metric data to a compatible host or component.
   */
  public async export(): Promise<void> {
    const { resourceMetrics, errors } = await this.collect();
    if (errors.length) {
      diag.error('WasiMetricExporter: metrics collection errors', ...errors);
    }
    exportToWasiHost(resourceMetricsToWasi(resourceMetrics));
  }
}
