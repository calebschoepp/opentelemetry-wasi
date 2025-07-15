#include <stdlib.h>
#include <tracing.h>

__attribute__((weak, export_name("canonical_abi_realloc")))
void *canonical_abi_realloc(
void *ptr,
size_t orig_size,
size_t align,
size_t new_size
) {
  if (new_size == 0)
  return (void*) align;
  void *ret = realloc(ptr, new_size);
  if (!ret)
  abort();
  return ret;
}

__attribute__((weak, export_name("canonical_abi_free")))
void canonical_abi_free(
void *ptr,
size_t size,
size_t align
) {
  if (size == 0)
  return;
  free(ptr);
}
#include <string.h>

void tracing_string_set(tracing_string_t *ret, const char *s) {
  ret->ptr = (char*) s;
  ret->len = strlen(s);
}

void tracing_string_dup(tracing_string_t *ret, const char *s) {
  ret->len = strlen(s);
  ret->ptr = canonical_abi_realloc(NULL, 0, 1, ret->len);
  memcpy(ret->ptr, s, ret->len);
}

void tracing_string_free(tracing_string_t *ret) {
  canonical_abi_free(ret->ptr, ret->len, 1);
  ret->ptr = NULL;
  ret->len = 0;
}
void tracing_trace_id_free(tracing_trace_id_t *ptr) {
  tracing_string_free(ptr);
}
void tracing_span_id_free(tracing_span_id_t *ptr) {
  tracing_string_free(ptr);
}
void tracing_tuple2_string_string_free(tracing_tuple2_string_string_t *ptr) {
  tracing_string_free(&ptr->f0);
  tracing_string_free(&ptr->f1);
}
void tracing_trace_state_free(tracing_trace_state_t *ptr) {
  for (size_t i = 0; i < ptr->len; i++) {
    tracing_tuple2_string_string_free(&ptr->ptr[i]);
  }
  canonical_abi_free(ptr->ptr, ptr->len * 16, 4);
}
void tracing_span_context_free(tracing_span_context_t *ptr) {
  tracing_trace_id_free(&ptr->trace_id);
  tracing_span_id_free(&ptr->span_id);
  tracing_trace_state_free(&ptr->trace_state);
}
void tracing_key_free(tracing_key_t *ptr) {
  tracing_string_free(ptr);
}
void tracing_list_string_free(tracing_list_string_t *ptr) {
  for (size_t i = 0; i < ptr->len; i++) {
    tracing_string_free(&ptr->ptr[i]);
  }
  canonical_abi_free(ptr->ptr, ptr->len * 8, 4);
}
void tracing_list_bool_free(tracing_list_bool_t *ptr) {
  canonical_abi_free(ptr->ptr, ptr->len * 1, 1);
}
void tracing_list_float64_free(tracing_list_float64_t *ptr) {
  canonical_abi_free(ptr->ptr, ptr->len * 8, 8);
}
void tracing_list_s64_free(tracing_list_s64_t *ptr) {
  canonical_abi_free(ptr->ptr, ptr->len * 8, 8);
}
void tracing_value_free(tracing_value_t *ptr) {
  switch ((int32_t) ptr->tag) {
    case 0: {
      tracing_string_free(&ptr->val.string);
      break;
    }
    case 4: {
      tracing_list_string_free(&ptr->val.string_array);
      break;
    }
    case 5: {
      tracing_list_bool_free(&ptr->val.bool_array);
      break;
    }
    case 6: {
      tracing_list_float64_free(&ptr->val.f64_array);
      break;
    }
    case 7: {
      tracing_list_s64_free(&ptr->val.s64_array);
      break;
    }
  }
}
void tracing_key_value_free(tracing_key_value_t *ptr) {
  tracing_key_free(&ptr->key);
  tracing_value_free(&ptr->value);
}
void tracing_list_key_value_free(tracing_list_key_value_t *ptr) {
  for (size_t i = 0; i < ptr->len; i++) {
    tracing_key_value_free(&ptr->ptr[i]);
  }
  canonical_abi_free(ptr->ptr, ptr->len * 24, 8);
}
void tracing_event_free(tracing_event_t *ptr) {
  tracing_string_free(&ptr->name);
  tracing_list_key_value_free(&ptr->attributes);
}
void tracing_list_event_free(tracing_list_event_t *ptr) {
  for (size_t i = 0; i < ptr->len; i++) {
    tracing_event_free(&ptr->ptr[i]);
  }
  canonical_abi_free(ptr->ptr, ptr->len * 32, 8);
}
void tracing_link_free(tracing_link_t *ptr) {
  tracing_span_context_free(&ptr->span_context);
  tracing_list_key_value_free(&ptr->attributes);
}
void tracing_list_link_free(tracing_list_link_t *ptr) {
  for (size_t i = 0; i < ptr->len; i++) {
    tracing_link_free(&ptr->ptr[i]);
  }
  canonical_abi_free(ptr->ptr, ptr->len * 36, 4);
}
void tracing_status_free(tracing_status_t *ptr) {
  switch ((int32_t) ptr->tag) {
    case 2: {
      tracing_string_free(&ptr->val.error);
      break;
    }
  }
}
void tracing_option_string_free(tracing_option_string_t *ptr) {
  if (ptr->is_some) {
    tracing_string_free(&ptr->val);
  }
}
void tracing_instrumentation_scope_free(tracing_instrumentation_scope_t *ptr) {
  tracing_string_free(&ptr->name);
  tracing_option_string_free(&ptr->version);
  tracing_option_string_free(&ptr->schema_url);
  tracing_list_key_value_free(&ptr->attributes);
}
void tracing_span_data_free(tracing_span_data_t *ptr) {
  tracing_span_context_free(&ptr->span_context);
  tracing_string_free(&ptr->parent_span_id);
  tracing_string_free(&ptr->name);
  tracing_list_key_value_free(&ptr->attributes);
  tracing_list_event_free(&ptr->events);
  tracing_list_link_free(&ptr->links);
  tracing_status_free(&ptr->status);
  tracing_instrumentation_scope_free(&ptr->instrumentation_scope);
}

__attribute__((aligned(8)))
static uint8_t RET_AREA[168];
__attribute__((import_module("tracing"), import_name("on-start")))
void __wasm_import_on_start(int32_t, int32_t, int32_t, int32_t, int32_t, int32_t, int32_t, int32_t);
void on_start(tracing_span_context_t *context) {
  __wasm_import_on_start((int32_t) ((*context).trace_id).ptr, (int32_t) ((*context).trace_id).len, (int32_t) ((*context).span_id).ptr, (int32_t) ((*context).span_id).len, (*context).trace_flags, (*context).is_remote, (int32_t) ((*context).trace_state).ptr, (int32_t) ((*context).trace_state).len);
}
__attribute__((import_module("tracing"), import_name("on-end")))
void __wasm_import_on_end(int32_t);
void on_end(tracing_span_data_t *span) {
  int32_t ptr = (int32_t) &RET_AREA;
  *((int32_t*)(ptr + 4)) = (int32_t) (((*span).span_context).trace_id).len;
  *((int32_t*)(ptr + 0)) = (int32_t) (((*span).span_context).trace_id).ptr;
  *((int32_t*)(ptr + 12)) = (int32_t) (((*span).span_context).span_id).len;
  *((int32_t*)(ptr + 8)) = (int32_t) (((*span).span_context).span_id).ptr;
  *((int8_t*)(ptr + 16)) = ((*span).span_context).trace_flags;
  *((int8_t*)(ptr + 17)) = ((*span).span_context).is_remote;
  *((int32_t*)(ptr + 24)) = (int32_t) (((*span).span_context).trace_state).len;
  *((int32_t*)(ptr + 20)) = (int32_t) (((*span).span_context).trace_state).ptr;
  *((int32_t*)(ptr + 32)) = (int32_t) ((*span).parent_span_id).len;
  *((int32_t*)(ptr + 28)) = (int32_t) ((*span).parent_span_id).ptr;
  *((int8_t*)(ptr + 36)) = (int32_t) (*span).span_kind;
  *((int32_t*)(ptr + 44)) = (int32_t) ((*span).name).len;
  *((int32_t*)(ptr + 40)) = (int32_t) ((*span).name).ptr;
  *((int64_t*)(ptr + 48)) = (int64_t) (((*span).start_time).seconds);
  *((int32_t*)(ptr + 56)) = (int32_t) (((*span).start_time).nanoseconds);
  *((int64_t*)(ptr + 64)) = (int64_t) (((*span).end_time).seconds);
  *((int32_t*)(ptr + 72)) = (int32_t) (((*span).end_time).nanoseconds);
  *((int32_t*)(ptr + 84)) = (int32_t) ((*span).attributes).len;
  *((int32_t*)(ptr + 80)) = (int32_t) ((*span).attributes).ptr;
  *((int32_t*)(ptr + 92)) = (int32_t) ((*span).events).len;
  *((int32_t*)(ptr + 88)) = (int32_t) ((*span).events).ptr;
  *((int32_t*)(ptr + 100)) = (int32_t) ((*span).links).len;
  *((int32_t*)(ptr + 96)) = (int32_t) ((*span).links).ptr;
  switch ((int32_t) ((*span).status).tag) {
    case 0: {
      *((int8_t*)(ptr + 104)) = 0;
      break;
    }
    case 1: {
      *((int8_t*)(ptr + 104)) = 1;
      break;
    }
    case 2: {
      const tracing_string_t *payload25 = &((*span).status).val.error;
      *((int8_t*)(ptr + 104)) = 2;
      *((int32_t*)(ptr + 112)) = (int32_t) (*payload25).len;
      *((int32_t*)(ptr + 108)) = (int32_t) (*payload25).ptr;
      break;
    }
  }
  *((int32_t*)(ptr + 120)) = (int32_t) (((*span).instrumentation_scope).name).len;
  *((int32_t*)(ptr + 116)) = (int32_t) (((*span).instrumentation_scope).name).ptr;
  
  if ((((*span).instrumentation_scope).version).is_some) {
    const tracing_string_t *payload27 = &(((*span).instrumentation_scope).version).val;
    *((int8_t*)(ptr + 124)) = 1;
    *((int32_t*)(ptr + 132)) = (int32_t) (*payload27).len;
    *((int32_t*)(ptr + 128)) = (int32_t) (*payload27).ptr;
    
  } else {
    *((int8_t*)(ptr + 124)) = 0;
    
  }
  
  if ((((*span).instrumentation_scope).schema_url).is_some) {
    const tracing_string_t *payload29 = &(((*span).instrumentation_scope).schema_url).val;
    *((int8_t*)(ptr + 136)) = 1;
    *((int32_t*)(ptr + 144)) = (int32_t) (*payload29).len;
    *((int32_t*)(ptr + 140)) = (int32_t) (*payload29).ptr;
    
  } else {
    *((int8_t*)(ptr + 136)) = 0;
    
  }
  *((int32_t*)(ptr + 152)) = (int32_t) (((*span).instrumentation_scope).attributes).len;
  *((int32_t*)(ptr + 148)) = (int32_t) (((*span).instrumentation_scope).attributes).ptr;
  *((int32_t*)(ptr + 156)) = (int32_t) ((*span).dropped_attributes);
  *((int32_t*)(ptr + 160)) = (int32_t) ((*span).dropped_events);
  *((int32_t*)(ptr + 164)) = (int32_t) ((*span).dropped_links);
  __wasm_import_on_end(ptr);
}
__attribute__((import_module("tracing"), import_name("outer-span-context")))
void __wasm_import_outer_span_context(int32_t);
void outer_span_context(tracing_span_context_t *ret0) {
  int32_t ptr = (int32_t) &RET_AREA;
  __wasm_import_outer_span_context(ptr);
  *ret0 = (tracing_span_context_t) {
    (tracing_string_t) { (char*)(*((int32_t*) (ptr + 0))), (size_t)(*((int32_t*) (ptr + 4))) },
    (tracing_string_t) { (char*)(*((int32_t*) (ptr + 8))), (size_t)(*((int32_t*) (ptr + 12))) },
    (int32_t) (*((uint8_t*) (ptr + 16))),
    (int32_t) (*((uint8_t*) (ptr + 17))),
    (tracing_trace_state_t) { (tracing_tuple2_string_string_t*)(*((int32_t*) (ptr + 20))), (size_t)(*((int32_t*) (ptr + 24))) },
  };
}
