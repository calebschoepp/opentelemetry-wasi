declare module 'wasi:otel/types@0.2.0-draft' {
  /**
   * The key part of attribute `key-value` pairs.
   */
  export type Key = string;
  /**
   * The value part of attribute `key-value` pairs.
   */
  export type Value = ValueString | ValueBool | ValueF64 | ValueS64 | ValueStringArray | ValueBoolArray | ValueF64Array | ValueS64Array;
  /**
   * A string value.
   */
  export interface ValueString {
    tag: 'string',
    val: string,
  }
  /**
   * A boolean value.
   */
  export interface ValueBool {
    tag: 'bool',
    val: boolean,
  }
  /**
   * A double precision floating point value.
   */
  export interface ValueF64 {
    tag: 'f64',
    val: number,
  }
  /**
   * A signed 64 bit integer value.
   */
  export interface ValueS64 {
    tag: 's64',
    val: bigint,
  }
  /**
   * A homogeneous array of string values.
   */
  export interface ValueStringArray {
    tag: 'string-array',
    val: Array<string>,
  }
  /**
   * A homogeneous array of boolean values.
   */
  export interface ValueBoolArray {
    tag: 'bool-array',
    val: Array<boolean>,
  }
  /**
   * A homogeneous array of double precision floating point values.
   */
  export interface ValueF64Array {
    tag: 'f64-array',
    val: Float64Array,
  }
  /**
   * A homogeneous array of 64 bit integer values.
   */
  export interface ValueS64Array {
    tag: 's64-array',
    val: BigInt64Array,
  }
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
     * https://github.com/open-telemetry/opentelemetry-specification/blob/v1.9.0/specification/schemas/overview.md#schema-url
     */
    schemaUrl?: string,
    /**
     * Specifies the instrumentation scope attributes to associate with emitted telemetry.
     */
    attributes: Array<KeyValue>,
  }
}
