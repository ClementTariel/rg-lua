#include "luaBridge.h"

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