declare module 'wasi:otel/types@0.2.0-draft' {
  /**
   * The key part of attribute `key-value` pairs.
   */
  export type Key = string;
  /**
   * The value part of attribute `key-value` pairs.
   * 
   * This corresponds with the `AnyValue` type defined in the [attribute spec](https://opentelemetry.io/docs/specs/otel/common/#anyvalue).
   * Because WIT doesn't support recursive types, the data needs to be serialized. JSON is used as the encoding format.
   */
  export type Value = string;
  /**
   * A key-value pair describing an attribute.
   */
  export interface KeyValue {
    /**
     * The attribute name.
     */
    key: Key,
    /**
     * The attribute value.
     */
    value: Value,
  }
  /**
   * An immutable representation of the entity producing telemetry as attributes.
   */
  export interface Resource {
    /**
     * Attributes that identify the resource.
     */
    attributes: Array<KeyValue>,
    /**
     * The schema URL to be recorded in the emitted resource.
     */
    schemaUrl?: string,
  }
  /**
   * Describes the instrumentation scope that produced telemetry.
   */
  export interface InstrumentationScope {
    /**
     * Name of the instrumentation scope.
     */
    name: string,
    /**
     * The library version.
     */
    version?: string,
    /**
     * Schema URL used by this library.
     * https://opentelemetry.io/docs/specs/otel/schemas/#schema-url
     */
    schemaUrl?: string,
    /**
     * Specifies the instrumentation scope attributes to associate with emitted telemetry.
     */
    attributes: Array<KeyValue>,
  }
}
