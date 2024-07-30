#ifndef STRING_ARRAY_H
#define STRING_ARRAY_H

#include <stdlib.h>
#include "custom_string.h"

typedef struct StringArray {
    String** data;
    size_t len;
    size_t next_index;
} StringArray;

StringArray init_string_array(size_t* pbudget);

void pop_string_array(StringArray* ps, size_t* pbudget);

void free_string_array(StringArray* ps, size_t* pbudget);

void shrink_string_array(StringArray* ps, size_t* pbudget);

void expand_string_array(StringArray* ps, size_t* pbudget);

void append_to_string_array(StringArray* ps, String* s, size_t* pbudget);

#endif