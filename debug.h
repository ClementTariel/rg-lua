#include <lua.h>
#include <lauxlib.h>
#include <lualib.h>
#include <stdbool.h>

char* copy_string(const char* original_string, bool add_quote, size_t* psize, size_t budget);

char* int_to_string(int n, size_t* psize, size_t budget);

char* float_to_string(float x, size_t* psize, size_t budget);

char* stringify_list_to_map(char** values_as_str, int values_count, size_t* psize, size_t budget);

char* luaValueToString(lua_State* pl, int idx, int depth, size_t* psize, size_t budget);

int printInLua(lua_State* pl);