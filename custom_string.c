#include <stdbool.h>
#include <string.h>
#include <stdlib.h>
#include "custom_string.h"

bool is_null_string(String s){
      return s.text == NULL;
}

void free_string(String* s, size_t* pbudget) {
      if (!is_null_string(*s)){
            free(s->text);
            *pbudget += s->len;
      }
      s->text = NULL;
      s->len = -1;
}

String init_string(size_t* pbudget){
      if (*pbudget <= 1){
            return (String){.text = NULL, .len = -1}; 
      }
      char* text = malloc(1);
      if (text == NULL){
            return (String){.text = NULL, .len = -1}; 
      }
      text[0] = '\0';
      *pbudget -= 1;
      return (String){.text = text, .len = 0};
}

void append_preallocated(String* s, char* next) {
      strcat(s->text, next);
}

void expand_string(String* s, size_t size, size_t* pbudget){
      size_t allocated_size = s->len + size + 1;
      if (allocated_size >= *pbudget || allocated_size <= s->len){
            free_string(s, pbudget);
            return;
      }
      char* realloc_text = realloc(s->text, allocated_size);
      if (realloc_text == NULL){
            free_string(s, pbudget);
            return;
      }
      *pbudget -= size;
      s->text = realloc_text;
      s->len += size;
}

void append_string(String* s, String* next, size_t* pbudget) {
      expand_string(s, next->len, pbudget);
      if (is_null_string(*s)){
            free_string(s, pbudget);
            return;
      }
      append_preallocated(s, next->text);
}

void append_chars(String* s, char* next, size_t* pbudget) {
      expand_string(s, strlen(next), pbudget);
      if (is_null_string(*s)){
            free_string(s, pbudget);
            return;
      }
      append_preallocated(s, next);
}

