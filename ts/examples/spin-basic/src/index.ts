import { AutoRouter } from 'itty-router';
import { BasicTracerProvider } from '@opentelemetry/sdk-trace-base';
import { context, trace, metrics } from '@opentelemetry/api';
import { MeterProvider } from '@opentelemetry/sdk-metrics';
import {
  WasiSpanProcessor,
  WasiTraceContextPropagator,
  WasiMetricExporter,
} from 'opentelemetry-wasi';
import { openDefault } from '@spinframework/spin-kv';
import { resourceFromAttributes } from '@opentelemetry/resources';
import { ATTR_SERVICE_NAME } from '@opentelemetry/semantic-conventions';

const metricExporter = new WasiMetricExporter();
metrics.setGlobalMeterProvider(
  new MeterProvider({
    resource: resourceFromAttributes({
      [ATTR_SERVICE_NAME]: 'spin-metrics',
    }),
    readers: [metricExporter],
  })
);

const propagator = new WasiTraceContextPropagator();
const provider = new BasicTracerProvider({
  spanProcessors: [new WasiSpanProcessor()],
});

const router = AutoRouter();
router.get('/', async () => {
  // Metrics
  const attrs = { spinKey1: 'spinValue1', spinKey2: 'spinValue2' };
  const meter = metrics.getMeter('spin-meter');

  const counter = meter.createCounter('spin-counter');
  counter.add(10, attrs);

  const upDownCounter = meter.createUpDownCounter('spin-up-down-counter');
  upDownCounter.add(-10, attrs);
  upDownCounter.add(5, attrs);

  const histogram = meter.createHistogram('spin-histogram');
  histogram.record(10, attrs);
  histogram.record(23, attrs);

  const gauge = meter.createGauge('spin-gauge');
  gauge.record(15);

  await metricExporter.export();

  // Tracing
  const hostContext = propagator.extract(context.active());
  const tracer = provider.getTracer('basic-spin');
  return tracer.startActiveSpan(
    'main-operation',
    {},
    hostContext,
    (parentSpan) => {
      const parentContext = trace.setSpan(hostContext, parentSpan);
      parentSpan.setAttribute('my-attribute', 'my-value');
      parentSpan.addEvent('Main span event', { foo: '1' });
      tracer.startActiveSpan(
        'child-operation',
        {},
        parentContext,
        (childSpan) => {
          childSpan.addEvent('Sub span event', { bar: '1' });
          const store = openDefault();
          store.set('foo', 'bar');
          childSpan.end();
        }
      );
      parentSpan.end();
      return new Response('Hello, world!');
    }
  );
});

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
//@ts-ignore
addEventListener('fetch', (event: FetchEvent) => {
  event.respondWith(router.fetch(event.request));
});
