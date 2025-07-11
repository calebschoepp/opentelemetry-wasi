#ifndef __BINDINGS_TRACES_H
#define __BINDINGS_TRACES_H
#ifdef __cplusplus
extern "C"
{
  #endif
  
  #include <stdint.h>
  #include <stdbool.h>
  
  typedef struct {
    char *ptr;
    size_t len;
  } traces_string_t;
  
  void traces_string_set(traces_string_t *ret, const char *s);
  void traces_string_dup(traces_string_t *ret, const char *s);
  void traces_string_free(traces_string_t *ret);
  typedef struct {
    uint64_t seconds;
    uint32_t nanoseconds;
  } traces_datetime_t;
  // The trace that this `span-context` belongs to.
  // 
  // 16 bytes encoded as a hexadecimal string.
  typedef traces_string_t traces_trace_id_t;
  void traces_trace_id_free(traces_trace_id_t *ptr);
  // The id of this `span-context`.
  // 
  // 8 bytes encoded as a hexadecimal string.
  typedef traces_string_t traces_span_id_t;
  void traces_span_id_free(traces_span_id_t *ptr);
  // Flags that can be set on a `span-context`.
  typedef uint8_t traces_trace_flags_t;
  #define TRACES_TRACE_FLAGS_SAMPLED (1 << 0)
  typedef struct {
    traces_string_t f0;
    traces_string_t f1;
  } traces_tuple2_string_string_t;
  void traces_tuple2_string_string_free(traces_tuple2_string_string_t *ptr);
  // Carries system-specific configuration data, represented as a list of key-value pairs. `trace-state` allows multiple tracing systems to participate in the same trace.
  // 
  // If any invalid keys or values are provided then the `trace-state` will be treated as an empty list.
  typedef struct {
    traces_tuple2_string_string_t *ptr;
    size_t len;
  } traces_trace_state_t;
  void traces_trace_state_free(traces_trace_state_t *ptr);
  // Identifying trace information about a span that can be serialized and propagated.
  typedef struct {
    traces_trace_id_t trace_id;
    traces_span_id_t span_id;
    traces_trace_flags_t trace_flags;
    bool is_remote;
    traces_trace_state_t trace_state;
  } traces_span_context_t;
  void traces_span_context_free(traces_span_context_t *ptr);
  // Describes the relationship between the Span, its parents, and its children in a trace.
  typedef uint8_t traces_span_kind_t;
  #define TRACES_SPAN_KIND_CLIENT 0
  #define TRACES_SPAN_KIND_SERVER 1
  #define TRACES_SPAN_KIND_PRODUCER 2
  #define TRACES_SPAN_KIND_CONSUMER 3
  #define TRACES_SPAN_KIND_INTERNAL 4
  // The key part of attribute `key-value` pairs.
  typedef traces_string_t traces_key_t;
  void traces_key_free(traces_key_t *ptr);
  typedef struct {
    traces_string_t *ptr;
    size_t len;
  } traces_list_string_t;
  void traces_list_string_free(traces_list_string_t *ptr);
  typedef struct {
    bool *ptr;
    size_t len;
  } traces_list_bool_t;
  void traces_list_bool_free(traces_list_bool_t *ptr);
  typedef struct {
    double *ptr;
    size_t len;
  } traces_list_float64_t;
  void traces_list_float64_free(traces_list_float64_t *ptr);
  typedef struct {
    int64_t *ptr;
    size_t len;
  } traces_list_s64_t;
  void traces_list_s64_free(traces_list_s64_t *ptr);
  // The value part of attribute `key-value` pairs.
  typedef struct {
    uint8_t tag;
    union {
      traces_string_t string;
      bool boolean;
      double f64;
      int64_t s64;
      traces_list_string_t string_array;
      traces_list_bool_t bool_array;
      traces_list_float64_t f64_array;
      traces_list_s64_t s64_array;
    } val;
  } traces_value_t;
  #define TRACES_VALUE_STRING 0
  #define TRACES_VALUE_BOOL 1
  #define TRACES_VALUE_F64 2
  #define TRACES_VALUE_S64 3
  #define TRACES_VALUE_STRING_ARRAY 4
  #define TRACES_VALUE_BOOL_ARRAY 5
  #define TRACES_VALUE_F64_ARRAY 6
  #define TRACES_VALUE_S64_ARRAY 7
  void traces_value_free(traces_value_t *ptr);
  // A key-value pair describing an attribute.
  typedef struct {
    traces_key_t key;
    traces_value_t value;
  } traces_key_value_t;
  void traces_key_value_free(traces_key_value_t *ptr);
  typedef struct {
    traces_key_value_t *ptr;
    size_t len;
  } traces_list_key_value_t;
  void traces_list_key_value_free(traces_list_key_value_t *ptr);
  // An event describing a specific moment in time on a span and associated attributes.
  typedef struct {
    traces_string_t name;
    traces_datetime_t time;
    traces_list_key_value_t attributes;
  } traces_event_t;
  void traces_event_free(traces_event_t *ptr);
  typedef struct {
    traces_event_t *ptr;
    size_t len;
  } traces_list_event_t;
  void traces_list_event_free(traces_list_event_t *ptr);
  // Describes a relationship to another `span`.
  typedef struct {
    traces_span_context_t span_context;
    traces_list_key_value_t attributes;
  } traces_link_t;
  void traces_link_free(traces_link_t *ptr);
  typedef struct {
    traces_link_t *ptr;
    size_t len;
  } traces_list_link_t;
  void traces_list_link_free(traces_list_link_t *ptr);
  // The `status` of a `span`.
  typedef struct {
    uint8_t tag;
    union {
      traces_string_t error;
    } val;
  } traces_status_t;
  #define TRACES_STATUS_UNSET 0
  #define TRACES_STATUS_OK 1
  #define TRACES_STATUS_ERROR 2
  void traces_status_free(traces_status_t *ptr);
  typedef struct {
    bool is_some;
    traces_string_t val;
  } traces_option_string_t;
  void traces_option_string_free(traces_option_string_t *ptr);
  // Describes the instrumentation scope that produced a span.
  typedef struct {
    traces_string_t name;
    traces_option_string_t version;
    traces_option_string_t schema_url;
    traces_list_key_value_t attributes;
  } traces_instrumentation_scope_t;
  void traces_instrumentation_scope_free(traces_instrumentation_scope_t *ptr);
  // The data associated with a span.
  typedef struct {
    traces_span_context_t span_context;
    traces_string_t parent_span_id;
    traces_span_kind_t span_kind;
    traces_string_t name;
    traces_datetime_t start_time;
    traces_datetime_t end_time;
    traces_list_key_value_t attributes;
    traces_list_event_t events;
    traces_list_link_t links;
    traces_status_t status;
    traces_instrumentation_scope_t instrumentation_scope;
    uint32_t dropped_attributes;
    uint32_t dropped_events;
    uint32_t dropped_links;
  } traces_span_data_t;
  void traces_span_data_free(traces_span_data_t *ptr);
  void traces_on_start(traces_span_context_t *context);
  void traces_on_end(traces_span_data_t *span);
  void traces_outer_span_context(traces_span_context_t *ret0);
  #ifdef __cplusplus
}
#endif
#endif
