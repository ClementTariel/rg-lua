#include <lua.h>
#include <lauxlib.h>
#include <lualib.h>

int getStackBridge(void* pl, int depth, void* ptDebug);

void getInfoBridge(void* pl, void* ptDebug);

int debugCurrentLine(void* ptDebug);

const char* toStringBridge(void* pl, int idx);

void popStateBridge(void* pl, int idx);

lua_State* newStateBridgeBridge();

void closeBridge(void* pl);

void pushCFunctionBridge(void* pl, void* fn);

void setGlobalBridge(void* pl, const char *name);

int loadStringBridge(void* pl, const char *s);

int pcallBridge(void* pl, int nargs, int nresults, int msgh);