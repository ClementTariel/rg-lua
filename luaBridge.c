#include "luaBridge.h"

// Copy luaB_next and luaB_pairs from source code to avoid importing an unsafe lib
static int luaB_next (lua_State *L) {
  luaL_checktype(L, 1, LUA_TTABLE);
  lua_settop(L, 2);  /* create a 2nd argument if there isn't one */
  if (lua_next(L, 1))
    return 2;
  else {
    lua_pushnil(L);
    return 1;
  }
}
static int pairsmeta (lua_State *L, const char *method, int iszero,
                      lua_CFunction iter) {
  if (!luaL_getmetafield(L, 1, method)) {  /* no metamethod? */
    luaL_checktype(L, 1, LUA_TTABLE);  /* argument must be a table */
    lua_pushcfunction(L, iter);  /* will return generator, */
    lua_pushvalue(L, 1);  /* state, */
    if (iszero) lua_pushinteger(L, 0);  /* and initial value */
    else lua_pushnil(L);
  }
  else {
    lua_pushvalue(L, 1);  /* argument 'self' to metamethod */
    lua_call(L, 1, 3);  /* get 3 values from metamethod */
  }
  return 3;
}
static int luaB_pairs (lua_State *L) {
  return pairsmeta(L, "__pairs", 0, luaB_next);
}

int getStackBridge(void* pl, int depth, void* ptDebug){
	return lua_getstack((lua_State*)pl, depth, (lua_Debug*)ptDebug);
}

void getInfoBridge(void* pl, void* ptDebug){
	lua_getinfo((lua_State*)pl, "Snlu", (lua_Debug*) ptDebug);
	printf("%d\n",((lua_Debug*)ptDebug)->currentline);
}

int debugCurrentLine(void* ptDebug){
	// printf("%d\n",((lua_Debug*)ptDebug)->currentline);
	// return ((lua_Debug*)ptDebug)->currentline;
	return -1;
}

const char* toStringBridge(void* pl, int idx){
	return lua_tostring(pl, (int)idx);
}

void popStateBridge(void* pl, int idx){
	lua_settop((lua_State*)pl, (int)(-(idx+1)));
}

lua_State* newStateBridgeBridge() {
	lua_State* pl = luaL_newstate();
	lua_pushcfunction(pl, luaB_next);
    lua_setglobal((lua_State*)pl, "next");
	lua_pushcfunction(pl, luaB_pairs);
    lua_setglobal((lua_State*)pl, "pairs");
	luaL_requiref(pl, "table", luaopen_table, 1);
	lua_pop(pl, 1);
	luaL_requiref(pl, "math", luaopen_math, 1);
	lua_pop(pl, 1);
    return pl;
}

void closeBridge(void* pl){
	lua_close((lua_State*)pl);
}

void pushCFunctionBridge(void* pl, void* fn){
	lua_pushcfunction(((lua_State*)pl), (lua_CFunction)fn);
}

void setGlobalBridge(void* pl, const char *name){
	lua_setglobal((lua_State*)pl, name);
}

int loadStringBridge(void* pl, const char *s){
	return luaL_loadstring((lua_State*)pl, s);
}

// TODO: handle memory limitation
int pcallBridge(void* pl, int nargs, int nresults, int msgh) {
	return lua_pcall((lua_State*)pl, nargs, nresults, msgh);
}