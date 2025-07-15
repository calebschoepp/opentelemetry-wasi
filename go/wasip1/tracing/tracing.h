#ifndef __BINDINGS_TRACING_H
#define __BINDINGS_TRACING_H
#ifdef __cplusplus
extern "C"
{
  #endif
  
  #include <stdint.h>
  #include <stdbool.h>
  
  typedef struct {
    char *ptr;
    size_t len;
  } tracing_string_t;
  
  void tracing_string_set(tracing_string_t *ret, const char *s);
  void tracing_string_dup(tracing_string_t *ret, const char *s);
  void tracing_string_free(tracing_string_t *ret);
  typedef struct {
    uint64_t seconds;
    uint32_t nanoseconds;
  } tracing_datetime_t;
  // The trace that this `span-context` belongs to.
  // 
  // 16 bytes encoded as a hexadecimal string.
  typedef tracing_string_t tracing_trace_id_t;
  void tracing_trace_id_free(tracing_trace_id_t *ptr);
  // The id of this `span-context`.
  // 
  // 8 bytes encoded as a hexadecimal string.
  typedef tracing_string_t tracing_span_id_t;
  void tracing_span_id_free(tracing_span_id_t *ptr);
  // Flags that can be set on a `span-context`.
  typedef uint8_t tracing_trace_flags_t;
  #define TRACING_TRACE_FLAGS_SAMPLED (1 << 0)
  typedef struct {
    tracing_string_t f0;
    tracing_string_t f1;
  } tracing_tuple2_string_string_t;
  void tracing_tuple2_string_string_free(tracing_tuple2_string_string_t *ptr);
  // Carries system-specific configuration data, represented as a list of key-value pairs. `trace-state` allows multiple tracing systems to participate in the same trace.
  // 
  // If any invalid keys or values are provided then the `trace-state` will be treated as an empty list.
  typedef struct {
    tracing_tuple2_string_string_t *ptr;
    size_t len;
  } tracing_trace_state_t;
  void tracing_trace_state_free(tracing_trace_state_t *ptr);
  // Identifying trace information about a span that can be serialized and propagated.
  typedef struct {
    tracing_trace_id_t trace_id;
    tracing_span_id_t span_id;
    tracing_trace_flags_t trace_flags;
    bool is_remote;
    tracing_trace_state_t trace_state;
  } tracing_span_context_t;
  void tracing_span_context_free(tracing_span_context_t *ptr);
  // Describes the relationship between the Span, its parents, and its children in a trace.
  typedef uint8_t tracing_span_kind_t;
  #define TRACING_SPAN_KIND_CLIENT 0
  #define TRACING_SPAN_KIND_SERVER 1
  #define TRACING_SPAN_KIND_PRODUCER 2
  #define TRACING_SPAN_KIND_CONSUMER 3
  #define TRACING_SPAN_KIND_INTERNAL 4
  // The key part of attribute `key-value` pairs.
  typedef tracing_string_t tracing_key_t;
  void tracing_key_free(tracing_key_t *ptr);
  typedef struct {
    tracing_string_t *ptr;
    size_t len;
  } tracing_list_string_t;
  void tracing_list_string_free(tracing_list_string_t *ptr);
  typedef struct {
    bool *ptr;
    size_t len;
  } tracing_list_bool_t;
  void tracing_list_bool_free(tracing_list_bool_t *ptr);
  typedef struct {
    double *ptr;
    size_t len;
  } tracing_list_float64_t;
  void tracing_list_float64_free(tracing_list_float64_t *ptr);
  typedef struct {
    int64_t *ptr;
    size_t len;
  } tracing_list_s64_t;
  void tracing_list_s64_free(tracing_list_s64_t *ptr);
  // The value part of attribute `key-value` pairs.
  typedef struct {
    uint8_t tag;
    union {
      tracing_string_t string;
      bool boolean;
      double f64;
      int64_t s64;
      tracing_list_string_t string_array;
      tracing_list_bool_t bool_array;
      tracing_list_float64_t f64_array;
      tracing_list_s64_t s64_array;
    } val;
  } tracing_value_t;
  #define TRACING_VALUE_STRING 0
  #define TRACING_VALUE_BOOL 1
  #define TRACING_VALUE_F64 2
  #define TRACING_VALUE_S64 3
  #define TRACING_VALUE_STRING_ARRAY 4
  #define TRACING_VALUE_BOOL_ARRAY 5
  #define TRACING_VALUE_F64_ARRAY 6
  #define TRACING_VALUE_S64_ARRAY 7
  void tracing_value_free(tracing_value_t *ptr);
  // A key-value pair describing an attribute.
  typedef struct {
    tracing_key_t key;
    tracing_value_t value;
  } tracing_key_value_t;
  void tracing_key_value_free(tracing_key_value_t *ptr);
  typedef struct {
    tracing_key_value_t *ptr;
    size_t len;
  } tracing_list_key_value_t;
  void tracing_list_key_value_free(tracing_list_key_value_t *ptr);
  // An event describing a specific moment in time on a span and associated attributes.
  typedef struct {
    tracing_string_t name;
    tracing_datetime_t time;
    tracing_list_key_value_t attributes;
  } tracing_event_t;
  void tracing_event_free(tracing_event_t *ptr);
  typedef struct {
    tracing_event_t *ptr;
    size_t len;
  } tracing_list_event_t;
  void tracing_list_event_free(tracing_list_event_t *ptr);
  // Describes a relationship to another `span`.
  typedef struct {
    tracing_span_context_t span_context;
    tracing_list_key_value_t attributes;
  } tracing_link_t;
  void tracing_link_free(tracing_link_t *ptr);
  typedef struct {
    tracing_link_t *ptr;
    size_t len;
  } tracing_list_link_t;
  void tracing_list_link_free(tracing_list_link_t *ptr);
  // The `status` of a `span`.
  typedef struct {
    uint8_t tag;
    union {
      tracing_string_t error;
    } val;
  } tracing_status_t;
  #define TRACING_STATUS_UNSET 0
  #define TRACING_STATUS_OK 1
  #define TRACING_STATUS_ERROR 2
  void tracing_status_free(tracing_status_t *ptr);
  typedef struct {
    bool is_some;
    tracing_string_t val;
  } tracing_option_string_t;
  void tracing_option_string_free(tracing_option_string_t *ptr);
  // Describes the instrumentation scope that produced a span.
  typedef struct {
    tracing_string_t name;
    tracing_option_string_t version;
    tracing_option_string_t schema_url;
    tracing_list_key_value_t attributes;
  } tracing_instrumentation_scope_t;
  void tracing_instrumentation_scope_free(tracing_instrumentation_scope_t *ptr);
  // The data associated with a span.
  typedef struct {
    tracing_span_context_t span_context;
    tracing_string_t parent_span_id;
    tracing_span_kind_t span_kind;
    tracing_string_t name;
    tracing_datetime_t start_time;
    tracing_datetime_t end_time;
    tracing_list_key_value_t attributes;
    tracing_list_event_t events;
    tracing_list_link_t links;
    tracing_status_t status;
    tracing_instrumentation_scope_t instrumentation_scope;
    uint32_t dropped_attributes;
    uint32_t dropped_events;
    uint32_t dropped_links;
  } tracing_span_data_t;
  void tracing_span_data_free(tracing_span_data_t *ptr);
  void on_start(tracing_span_context_t *context);
  void on_end(tracing_span_data_t *span);
  void outer_span_context(tracing_span_context_t *ret0);
  #ifdef __cplusplus
}
#endif
#endif
