#ifndef DEBUG_H
#define DEBUG_H

#include <lua.h>
#include <lauxlib.h>
#include <lualib.h>
#include "custom_string.h"
#include "string_array.h"

String copy_string(const char* original_string, bool add_quote, size_t* pbudget);

String int_to_string(int n, size_t* pbudget);

String float_to_string(float x, size_t* pbudget);

String stringify_list_to_map(StringArray* values_as_str, size_t* pbudget);

String luaValueToString(lua_State* pl, int idx, int depth, size_t* pbudget);

int printInLua(lua_State* pl);

#endif