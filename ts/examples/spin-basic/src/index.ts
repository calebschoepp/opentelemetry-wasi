import { AutoRouter } from "itty-router";
import { BasicTracerProvider } from "@opentelemetry/sdk-trace-base";
import { trace } from "@opentelemetry/api";
import { WasiProcessor } from "opentelemetry-wasi";

// Initialize the tracer provider and configure a simple span processor
const provider = new BasicTracerProvider();
provider.addSpanProcessor(new WasiProcessor());
provider.register();

let router = AutoRouter();

router.get("/", () => {
    let span = trace.getTracer("spin-basic").startSpan("foo");
    span.end();
    return new Response("hello universe");
});

//@ts-ignore
addEventListener("fetch", (event: FetchEvent) => {
    event.respondWith(router.fetch(event.request));
});
