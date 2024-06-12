#include <stdbool.h>
#include <string.h>
#include <stdlib.h>
#include "custom_string.h"
#include "string_array.h"
#include "debug.h"

size_t PRINT_MEMORY_BUDGET = 1 << 15; // 32kb, arbitrary limit

String copy_string(const char* original_string, bool add_quote, size_t* pbudget) {
      int len = snprintf(NULL, 0, add_quote ? "\"%s\"" : "%s", original_string);
      String str_var = init_string(pbudget);
      if (is_null_string(str_var)) {
            return (String){.text = NULL, .len = -1};
      }
      expand_string(&str_var, len, pbudget);
      if (is_null_string(str_var)) {
            return (String){.text = NULL, .len = -1};
      }
      snprintf(str_var.text, len+1, add_quote ? "\"%s\"" : "%s", original_string);
      return str_var;
}

String int_to_string(int n, size_t* pbudget) {
      int len = snprintf(NULL, 0,"%d",n);
      String str_var = init_string(pbudget);
      if (is_null_string(str_var)) {
            return (String){.text = NULL, .len = -1};
      }
      expand_string(&str_var, len, pbudget);
      if (is_null_string(str_var)) {
            return (String){.text = NULL, .len = -1};
      }
      snprintf(str_var.text, len+1,"%d",n);
      return str_var;
}

String float_to_string(float x, size_t* pbudget) {
      int len = snprintf(NULL, 0,"%f",x);
      String str_var = init_string(pbudget);
      if (is_null_string(str_var)) {
            return (String){.text = NULL, .len = -1};
      }
      expand_string(&str_var, len, pbudget);
      if (is_null_string(str_var)) {
            return (String){.text = NULL, .len = -1};
      }
      snprintf(str_var.text, len+1,"%f",x);
      return str_var;
}

String stringify_list_to_map(StringArray* pvalues_as_str, size_t* pbudget){
      String total_str = init_string(pbudget);
      if (is_null_string(total_str)){
            free_string_array(pvalues_as_str, pbudget);
            return (String){.text = NULL, .len = -1};
      }
      for (int i=0; i<pvalues_as_str->next_index; i++){
            // lua indexes start at 1
            String str_key = int_to_string(i+1, pbudget);
            if (is_null_string(str_key)){
                  free_string(&total_str, pbudget);
                  free_string_array(pvalues_as_str, pbudget);
                  return (String){.text = NULL, .len = -1};
            }
            if (!is_null_string(total_str)){
                  append_chars(&total_str, " ", pbudget);
            }
            if (!is_null_string(total_str)){
                  append_string(&total_str, &str_key, pbudget);
            }
            free_string(&str_key, pbudget);
            if (!is_null_string(total_str)){
                  append_chars(&total_str, ": ", pbudget);
            }
            if (!is_null_string(total_str)){
                  append_string(&total_str, pvalues_as_str->data[i], pbudget);
            }
            if (!is_null_string(total_str)){
                  append_chars(&total_str, ",", pbudget);
            }
            if (is_null_string(total_str)){
                  free_string_array(pvalues_as_str, pbudget);
                  return (String){.text = NULL, .len = -1};
            }
      }
      free_string_array(pvalues_as_str, pbudget);
      return total_str;
}

String luaValueToString(lua_State* pl, int idx, int depth, size_t* pbudget) {
      if (lua_isinteger(pl, idx)) {
            // int
            int n = lua_tointeger(pl, idx);
            return int_to_string(n, pbudget);
      } else if (lua_isnumber(pl, idx)) {
            // float
            float x = lua_tonumber(pl, idx);
            return float_to_string(x, pbudget);
      } else if (lua_isboolean(pl, idx)) {
            // bool
            int boolVar = lua_toboolean(pl, idx);
            return copy_string(boolVar ? "true" : "false", false, pbudget);
      } else if (lua_isstring(pl, idx)) {
            // String (const char*)
            const char* luaStr = lua_tostring(pl, idx);
            return copy_string(luaStr, true, pbudget);
      } else if (lua_isnoneornil(pl, idx)) {
            // NULL
            return copy_string("<nil>", false, pbudget);
      } else if (lua_istable(pl, idx)) {
            // list or map
            if (depth > 3){
                  return copy_string("{ ... }", false, pbudget);
            }
            String total_str = init_string(pbudget);
            if (is_null_string(total_str)){
                  return (String){.text = NULL, .len = -1};
            }
            append_chars(&total_str, "{", pbudget);
            if (is_null_string(total_str)){
                  return (String){.text = NULL, .len = -1};
            }
            StringArray values_as_str = init_string_array(pbudget);
            if (values_as_str.data == NULL){
                  free_string(&total_str, pbudget);
                  return (String){.text = NULL, .len = -1};
            }
            // lua indexes start at 1
            int prev_index = 0;
            // Push the value on top of the stack
            lua_pushvalue(pl, idx);
            lua_pushnil(pl);
            while (lua_next(pl, -2) != 0){
                  // stack now contains: -1 => value; -2 => key; -3 => table               
                  // copy the key so that lua_tostring does not modify the original
                  // (avoid issue with lua_tostring that confuses lua_next)
                  lua_pushvalue(pl, -2);
                  // stack now contains: -1 => key; -2 => value; -3 => key; -4 => table
                  int index = -1;
                  if (lua_isinteger(pl, -1)){
                        index = lua_tointeger(pl ,-1);
                  }
                  // lua indexes start at 1
                  if (index > 0 && index == prev_index + 1) {
                        // indexes seem to behave as if it is a list
                        prev_index = index;
                  } else {
                        // inconsistency in indexes, it cannot be a list
                        index = -1;
                        prev_index = -1;
                  }
                  if (index > 0){
                        // list
                        String str_val = luaValueToString(pl, -2, depth+1, pbudget);
                        if (is_null_string(str_val)){
                              free_string_array(&values_as_str, pbudget);
                              free_string(&total_str, pbudget);
                              return (String){.text = NULL, .len = -1};
                        }
                        // pop value + copy of key, leaving original key
                        lua_pop(pl,2);
                        // stack now contains: -1 => key; -2 => table
                        append_to_string_array(&values_as_str, &str_val, pbudget);
                        if (values_as_str.data == NULL){
                              free_string(&total_str, pbudget);
                              return (String){.text = NULL, .len = -1};
                        }
                  } else {
                        // map
                        // add what was previously supposed to be a list
                        String beginning_str = stringify_list_to_map(&values_as_str, pbudget);
                        if (is_null_string(beginning_str)) {
                              free_string_array(&values_as_str, pbudget);
                              free_string(&total_str, pbudget);
                              return (String){.text = NULL, .len = -1};
                        }
                        append_string(&total_str, &beginning_str, pbudget);
                        free_string(&beginning_str, pbudget);
                        if (is_null_string(total_str)){
                              free_string_array(&values_as_str, pbudget);
                              return (String){.text = NULL, .len = -1};
                        }
                        //stringify key
                        String str_key = luaValueToString(pl, -1, depth+1, pbudget);
                        if (is_null_string(str_key)){
                              free_string_array(&values_as_str, pbudget);
                              free_string(&total_str, pbudget);
                              return (String){.text = NULL, .len = -1};
                        }
                        //stringify value
                        String str_val = luaValueToString(pl, -2, depth+1, pbudget);
                        if (is_null_string(str_val)){
                              free_string(&str_key, pbudget);
                              free_string_array(&values_as_str, pbudget);
                              free_string(&total_str, pbudget);
                              return (String){.text = NULL, .len = -1};
                        }
                        // pop value + copy of key, leaving original key
                        lua_pop(pl,2);
                        // stack now contains: -1 => key; -2 => table
                        append_chars(&total_str, " ", pbudget);
                        if (!is_null_string(total_str)){
                              append_string(&total_str, &str_key, pbudget);
                        }
                        free_string(&str_key, pbudget);
                        if (!is_null_string(total_str)){
                              append_chars(&total_str, ": ", pbudget);
                        }
                        if (!is_null_string(total_str)){
                              append_string(&total_str, &str_val, pbudget);
                        }
                        free_string(&str_val, pbudget);
                        if (!is_null_string(total_str)){
                              append_chars(&total_str, ",", pbudget);
                        }
                        if (is_null_string(total_str)){
                              free_string_array(&values_as_str, pbudget);
                              return (String){.text = NULL, .len = -1};
                        }
                  }
            }
            // stack now contains: -1 => table (when lua_next returns 0 it pops the key
            // but does not push anything.)
            // Pop table
            lua_pop(pl,1);
            // Stack is now the same as it was on entry to this function
            // if it was a list, loop over its values
            for (int i=0; i<prev_index; i++){
                  if (!is_null_string(total_str)){
                        append_chars(&total_str, " ", pbudget);
                  }
                  if (!is_null_string(total_str)){
                        append_string(&total_str, values_as_str.data[i], pbudget);
                  }
                  if (!is_null_string(total_str)){
                        append_chars(&total_str, ",", pbudget);
                  }
            }
            free_string_array(&values_as_str, pbudget);
            if (is_null_string(total_str)){
                  return (String){.text = NULL, .len = -1};
            }
            // remove "," after the last element
            if (total_str.len >= 1 && total_str.text[total_str.len-1] == ','){
                  total_str.text[total_str.len-1] = ' ';
            }
            append_chars(&total_str, "}", pbudget);
            return total_str;
      } else {
            // Unknown or not handled
            return copy_string("<type not supported>", false, pbudget);
      }
}

int printInLua(lua_State* pl) {
      if (PRINT_MEMORY_BUDGET < 2){
            return 0;
      }
      size_t budget = PRINT_MEMORY_BUDGET;
      int argc = lua_gettop(pl);
      for (int i=1; i<=argc; i++){
            if (i != 1) {
                  printf(" ");
            }
            size_t size = 0;
            String str_var = luaValueToString(pl, i, 0, &budget);
            if (is_null_string(str_var)) {
                  printf("%s\n", "<memory budget exceeded>");
                  fflush(stdout);
            }
            printf("%s", str_var.text);
            free_string(&str_var, &budget);
            if (i != argc) {
                  printf(",");
            }
            fflush(stdout);
      }
      printf("\n");
      fflush(stdout);
      return 0;
}