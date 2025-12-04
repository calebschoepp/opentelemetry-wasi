import { AutoRouter } from "itty-router";
import { BasicTracerProvider } from "@opentelemetry/sdk-trace-base";
import { WasiSpanProcessor, WasiTraceContextPropagator } from "opentelemetry-wasi";
import { openDefault } from "@spinframework/spin-kv";
import { context, trace } from "@opentelemetry/api";

const propagator = new WasiTraceContextPropagator();
const provider = new BasicTracerProvider({spanProcessors: [new WasiSpanProcessor()]});
provider.register();

let router = AutoRouter();
router
    .get("/", async () => {
        const hostContext = propagator.extract(context.active());
        let tracer = provider.getTracer("basic-spin");
        return tracer.startActiveSpan("main-operation", {}, hostContext, (parentSpan) => {
            const parentContext = trace.setSpan(hostContext, parentSpan);
            parentSpan.setAttribute("my-attribute", "my-value");
            parentSpan.addEvent("Main span event", {"foo": "1"} );
            tracer.startActiveSpan("child-operation", {}, parentContext, (childSpan) => {
                childSpan.addEvent("Sub span event", {"bar": "1"});
                let store = openDefault();
                store.set("foo", "bar");
                childSpan.end();
            });
            parentSpan.end();
            return new Response("Hello, world!");
        });
    });

//@ts-ignore
addEventListener("fetch", (event: FetchEvent) => {
    event.respondWith(router.fetch(event.request));
});
