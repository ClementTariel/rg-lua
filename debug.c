#include <stdbool.h>
#include <string.h>
#include <stdlib.h>
#include "debug.h"

size_t PRINT_MEMORY_BUDGET = 1 << 15; // 32kb, arbitrary limit

char* copy_string(const char* original_string, bool add_quote, size_t* psize, size_t budget) {
      int len = snprintf(NULL, 0, add_quote ? "\"%s\"" : "%s",original_string);
      size_t allocated_size = sizeof(char)*(len+1);
      char* strVar = NULL;
      if (allocated_size < budget){
            strVar = malloc(allocated_size);
      }
      if (strVar == NULL){
            return NULL;
      }
      snprintf(strVar, len+1, add_quote ? "\"%s\"" : "%s", original_string);
      if (psize != NULL){
            *psize = len;
      }
      return strVar;
}

char* int_to_string(int n, size_t* psize, size_t budget) {
      int len = snprintf(NULL, 0,"%d",n);
      size_t allocated_size = sizeof(char)*(len+1);
      char* strVar = NULL;
      if (allocated_size < budget){
            strVar = malloc(allocated_size);
      }
      if (strVar == NULL){
            return NULL;
      }
      snprintf(strVar, len+1,"%d",n);
      if (psize != NULL){
            *psize = len;
      }
      return strVar;
}

char* float_to_string(float x, size_t* psize, size_t budget) {
      int len = snprintf(NULL, 0,"%f",x);
      size_t allocated_size = sizeof(char)*(len+1);
      char* strVar = NULL;
      if (allocated_size < budget){
            strVar = malloc(allocated_size);
      }
      if (strVar == NULL){
            return NULL;
      }
      snprintf(strVar, len+1,"%f",x);
      if (psize != NULL){
            *psize = len;
      }
      return strVar;
}

char* stringify_list_to_map(char** values_as_str, int values_count, size_t* psize, size_t budget){
      size_t allocated_size = 0;
      int digit_count = 1;
      int lower_bound = 1;
      while (10*lower_bound < values_count+1) {
            digit_count++;
            lower_bound *= 10;
      }
      int keys_as_str_total_size = 0;
      while (digit_count > 0){
            keys_as_str_total_size += digit_count * (values_count + 1 - lower_bound);
            values_count = lower_bound - 1;
            lower_bound /= 10;
            digit_count--;
      }
      int values_as_str_total_size = 0;
      for (int i=0; i<values_count; i++){
            values_as_str_total_size += strlen(values_as_str[i]);
      }
      int total_str_size = keys_as_str_total_size + values_as_str_total_size + 4 * values_count;
      allocated_size = sizeof(char)*(total_str_size + 1);
      char* total_str = NULL;
      if (allocated_size < budget){
            total_str = malloc(allocated_size);
      }
      if (total_str == NULL){
            for (int i=0; i<values_count; i++){
                  free(values_as_str[i]);         
            }
            return NULL;
      }
      budget -= allocated_size;
      for (int i=0; i<values_count; i++){
            size_t key_size = 0;
            char* str_key = int_to_string(i+1, &key_size, budget + allocated_size); // allocated_size already accounted for
            if (str_key == NULL){
                  free(total_str);
                  for (int i=0; i<values_count; i++){
                        free(values_as_str[i]);         
                  }
                  return NULL;
            }
            strcat(total_str, " ");
            strcat(total_str, str_key);
            strcat(total_str, ": ");
            strcat(total_str, values_as_str[i]);
            strcat(total_str, ",");
            free(str_key);
      }
      for (int i=0; i<values_count; i++){
            free(values_as_str[i]);         
      }
      return total_str;
}

char* luaValueToString(lua_State* pl, int idx, int depth, size_t* psize, size_t budget) {
      if (lua_isinteger(pl, idx)) {
            // int
            int n = lua_tointeger(pl, idx);
            return int_to_string(n, psize, budget);
      } else if (lua_isnumber(pl, idx)) {
            // float
            float x = lua_tonumber(pl, idx);
            return float_to_string(x, psize, budget);
      } else if (lua_isboolean(pl, idx)) {
            // bool
            int boolVar = lua_toboolean(pl, idx);
            return copy_string(boolVar ? "true" : "false", false, psize, budget);
      } else if (lua_isstring(pl, idx)) {
            // string (const char*)
            const char* luaStr = lua_tostring(pl, idx);
            return copy_string(luaStr, true, psize, budget);
      } else if (lua_isnoneornil(pl, idx)) {
            // NULL
            return copy_string("<nil>", false, psize, budget);
      } else if (lua_istable(pl, idx)) {
            if (depth > 3){
                  return copy_string("{ ... }", false, psize, budget);
            }
            // list or map
            lua_pushvalue(pl, idx);
            lua_pushnil(pl);
            size_t allocated_size = sizeof(char)*(2 + 1);
            char* total_str = NULL;
            if (allocated_size < budget){
                  total_str = malloc(allocated_size);
            }
            if (total_str == NULL){
                  return NULL;
            }
            budget -= allocated_size;
            *psize = 2;
            strcpy(total_str, "{");
            int array_size = 2;
            allocated_size = sizeof(char*)*array_size;
            char** values_as_str = NULL;
            if (array_size > 0 && allocated_size < budget){
                  values_as_str = malloc(allocated_size);
            }
            if (values_as_str == NULL){
                  free(total_str);
                  return NULL;
            }
            budget -= allocated_size;
            // lua indexes start at 1
            int prev_index = 0;
            while (lua_next(pl, -2) != 0){
                  // stack now contains: -1 => value; -2 => key; -3 => table               
                  // copy the key so that lua_tostring does not modify the original
                  // (avoid issue with lua_tostring that confuses lua_next)
                  lua_pushvalue(pl, -2);
                  // stack now contains: -1 => key; -2 => value; -3 => key; -4 => table
                  if (prev_index >= 0 && lua_isinteger(pl, -1)){
                        int index = lua_tointeger(pl ,-1);
                        if (index != prev_index + 1) {
                              if (prev_index > 0){
                                    size_t beginning_str_size = 0;
                                    char* beginning_str = stringify_list_to_map(values_as_str, prev_index, &beginning_str_size, budget);
                                    if (beginning_str == NULL) {
                                          free(values_as_str);
                                          free(total_str);
                                          return NULL;
                                    }
                                    budget -= sizeof(char)*(beginning_str_size + 1);
                                    size_t old_allocated_size = sizeof(char)*(*psize + 1);
                                    *psize = *psize + beginning_str_size;
                                    allocated_size = sizeof(char)*(*psize + 1);
                                    char* realloc_total_str = NULL;
                                    if (allocated_size < budget){
                                          realloc_total_str = realloc(total_str, allocated_size);
                                    }
                                    if (realloc_total_str == NULL){
                                          free(values_as_str);
                                          free(total_str);
                                          return NULL;
                                    }
                                    budget -= allocated_size;
                                    total_str = realloc_total_str;
                                    budget += old_allocated_size;
                                    strcat(total_str, beginning_str);
                                    free(beginning_str);
                                    budget += sizeof(char)*(beginning_str_size + 1);
                              }
                              prev_index = -1;
                        } else {
                              prev_index = index;
                        }
                  } else {
                        prev_index = -1;
                  }
                  if (prev_index < 0){
                        size_t key_len = 0;
                        char* str_key = luaValueToString(pl, -1, depth+1, &key_len, budget);
                        budget -= sizeof(char)*(key_len + 1);
                        if (str_key == NULL){
                              free(values_as_str);
                              free(total_str);
                              return NULL;
                        }
                        size_t val_len = 0;
                        char* str_val = luaValueToString(pl, -2, depth+1, &val_len, budget);
                        budget -= sizeof(char) * (val_len + 1);
                        if (str_val == NULL){
                              free(str_key);
                              free(values_as_str);
                              free(total_str);
                              return NULL;
                        }
                        // pop value + copy of key, leaving original key
                        lua_pop(pl,2);
                        // stack now contains: -1 => key; -2 => table
                        size_t old_allocated_size = sizeof(char)*(*psize + 1);
                        *psize = *psize + val_len + key_len + 4;
                        allocated_size = sizeof(char)*(*psize + 1);
                        char* realloc_total_str = NULL;
                        if (allocated_size < budget){
                              realloc_total_str = realloc(total_str, allocated_size);
                        }
                        if (realloc_total_str == NULL){
                              free(str_val);
                              free(str_key);
                              free(values_as_str);
                              free(total_str);
                              return NULL;
                        }
                        budget -= allocated_size;
                        total_str = realloc_total_str;
                        budget += old_allocated_size;
                        strcat(total_str, " ");
                        strcat(total_str, str_key);
                        strcat(total_str, ": ");
                        strcat(total_str, str_val);
                        strcat(total_str, ",");
                        free(str_val);
                        budget += sizeof(char) * (val_len + 1);
                        free(str_key);
                        budget += sizeof(char)*(key_len + 1);
                  } else {
                        size_t val_len = 0;
                        char* str_val = luaValueToString(pl, -2, depth+1, &val_len, budget);
                        if (str_val == NULL){
                              free(str_val);
                              for (int i=0; i<prev_index; i++){                  
                                    free(values_as_str[i]);
                              } 
                              free(values_as_str);
                              free(total_str);
                              return NULL;
                        }
                        budget -= sizeof(char) * (val_len + 1);
                        // pop value + copy of key, leaving original key
                        lua_pop(pl,2);
                        // stack now contains: -1 => key; -2 => table
                        if (prev_index-1 >= array_size) { // lua indexes start at 1
                              size_t old_allocated_size = sizeof(char*) * array_size;
                              array_size *= 2;
                              allocated_size = sizeof(char*) * array_size;
                              char** realloc_values_as_str = NULL;
                              if (allocated_size < old_allocated_size || allocated_size < budget){
                                    realloc_values_as_str = realloc(values_as_str, allocated_size);
                              }
                              if (realloc_values_as_str == NULL) {
                                    free(str_val);
                                    for (int i=0; i<prev_index; i++){                  
                                          free(values_as_str[i]);
                                    } 
                                    free(values_as_str);
                                    free(total_str);
                                    return NULL;
                              }
                              budget -= allocated_size;
                              values_as_str = realloc_values_as_str;
                              budget += old_allocated_size;
                        }
                        values_as_str[prev_index-1] = str_val;
                  }
            }
            int total_values_size = 0;
            for (int i=0; i<prev_index; i++){
                  total_values_size += strlen(values_as_str[i]) + 2;
            }
            size_t old_allocated_size = sizeof(char) * (*psize + 1);
            *psize = *psize + total_values_size;
            allocated_size = sizeof(char) * (*psize + 1);
            char* realloc_total_str = NULL;
            if (allocated_size < budget){
                  realloc_total_str = realloc(total_str, allocated_size);
            }
            if (realloc_total_str == NULL) {
                  for (int i=0; i<prev_index; i++){                  
                        free(values_as_str[i]);
                  } 
                  free(values_as_str);
                  free(total_str);
                  return NULL;
            }
            budget -= allocated_size;
            total_str = realloc_total_str;
            budget += old_allocated_size;
            for (int i=0; i<prev_index; i++){
                  strcat(total_str, " ");
                  strcat(total_str, values_as_str[i]);
                  strcat(total_str, ",");
            }
            for (int i=0; i<prev_index; i++){     
                  budget += strlen(values_as_str[i]);
                  free(values_as_str[i]);
            } 
            free(values_as_str);
            budget += sizeof(char*) * array_size;
            lua_pop(pl,1);
            if (*psize >= 2 && total_str[*psize-2] == ','){
                  total_str[*psize-2] = ' ';
            }
            strcat(total_str, "}");
            return total_str;
      } else {
            // Unknown or not handled
            return copy_string("<type not supported>", false, psize, budget);
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
            char* strVar = luaValueToString(pl, i, 0, &size, budget);
            if (strVar == NULL) {
                  printf("%s\n", "<memory error when converting to string>");
                  fflush(stdout);
                  luaL_error(pl, "MemoryError: memory allocation failed or exceeded authorized memory budget");
                  return 1;
            }
            printf("%s", strVar);
            free(strVar);
            if (i != argc) {
                  printf(",");
            }
            fflush(stdout);
      }
      printf("\n");
      fflush(stdout);
      return 0;
}