#include <stdlib.h>
#include "custom_string.h"
#include "string_array.h"

StringArray init_string_array(size_t* pbudget){
      size_t initial_array_size = 2;
      size_t allocated_size = sizeof(String*)*initial_array_size;
      if (*pbudget <= allocated_size){
            return (StringArray){.data = NULL, .len = 0, .next_index = 0};
      }
      String** data = malloc(allocated_size);
      if (data == NULL){
            return (StringArray){.data = NULL, .len = 0, .next_index = 0};
      }
      *pbudget -= allocated_size;
      return (StringArray){.data = data, .len = initial_array_size, .next_index = 0};
}

void pop_string_array(StringArray* ps, size_t* pbudget){
      if (ps->next_index != 0){
            free_string(ps->data[ps->next_index-1], pbudget);
            ps->next_index -= 1;
      }
}

void free_string_array(StringArray* ps, size_t* pbudget){
      if (ps->data == NULL){
            ps->len = 0;
            ps->next_index = 0;
            return;
      }
      while (ps->next_index != 0){
            pop_string_array(ps, pbudget);
      }
      free(ps->data);
      *pbudget += sizeof(String*) * ps->len;
      ps->data = NULL;
      ps->len = 0;
}

void shrink_string_array(StringArray* ps, size_t* pbudget){
      if (ps->len <= 2 || ps->data == NULL){
            return;
      }
      size_t previous_len = ps->len;
      size_t old_allocated_size = sizeof(String*) * previous_len;
      ps->len /= 2;
      size_t allocated_size = sizeof(String*) * ps->len;
      while(ps->next_index > ps->len){
            pop_string_array(ps, pbudget);
      }
      String** realloc_data = realloc(ps->data, allocated_size);
      if (realloc_data == NULL){
            ps->len = previous_len;
            free_string_array(ps, pbudget);
            return;
      }
      ps->data = realloc_data;
      *pbudget += old_allocated_size - allocated_size;
}

void expand_string_array(StringArray* ps, size_t* pbudget){
      size_t previous_len = ps->len;
      size_t old_allocated_size = sizeof(String*) * previous_len;
      ps->len *= 2;
      size_t allocated_size = sizeof(String*) * ps->len;
      if (allocated_size < old_allocated_size || allocated_size >= *pbudget){
            free_string_array(ps, pbudget);
            return ;
      }
      String** realloc_data = realloc(ps->data, allocated_size);
      if (realloc_data == NULL){
            ps->len = previous_len;
            free_string_array(ps, pbudget);
            return;
      }
      ps->data = realloc_data;
      *pbudget += old_allocated_size - allocated_size;
}

void append_to_string_array(StringArray* ps, String* s, size_t* pbudget){
      if (ps->next_index >= ps->len) {
            expand_string_array(ps, pbudget);
            if (ps->data == NULL){
                  return ;
            }
      }
      ps->data[ps->next_index] = s;
      ps->next_index += 1;
}