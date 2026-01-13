export type Key = string;
export type Value = string;

export interface KeyValue {
  key: Key;
  value: Value;
}

export interface Resource {
  attributes: Array<KeyValue>;
  schemaUrl?: string;
}

export interface InstrumentationScope {
  name: string;
  version?: string;
  schemaUrl?: string;
  attributes: Array<KeyValue>;
}
