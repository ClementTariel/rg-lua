#ifndef CUSTOM_STRING_H
#define CUSTOM_STRING_H

#include <stdbool.h>

typedef struct String {
    char* text;
    size_t len;
} String;

bool is_null_string(String s);

void free_string(String* s, size_t* pbudget);

String init_string(size_t* pbudget);

void append_preallocated(String* s, char* next);

void expand_string(String* s, size_t size, size_t* pbudget);

void append_string(String* s, String* next, size_t* pbudget);

void append_chars(String* s, char* next, size_t* pbudget);

#endif