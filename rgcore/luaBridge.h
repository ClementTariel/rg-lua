#ifndef LUA_BRIDGE_H
#define LUA_BRIDGE_H

#include <stdbool.h>
#include <pthread.h>
#include <lua.h>
#include <lauxlib.h>
#include <lualib.h>

typedef struct pcall_thread_params
{
    void *pl;
    int nargs;
    int nresults;
    int msgh;
    int *perr;
    bool *pdone;
    pthread_mutex_t *pdone_mutex;
    pthread_t timeout_thread_id;
} pcall_thread_params;

int GetStackBridge(void *pl, int depth, void *ptDebug);

void GetInfoBridge(void *pl, void *ptDebug);

int DebugCurrentLine(void *ptDebug);

const char *ToStringBridge(void *pl, int idx);

void PopStateBridge(void *pl, int idx);

lua_State *NewStateBridge();

void CloseBridge(void *pl);

void PushCFunctionBridge(void *pl, void *fn);

void SetGlobalBridge(void *pl, const char *name);

int LoadStringBridge(void *pl, const char *s);

int PcallBridge(void *pl, int nargs, int nresults, int msgh);

void *pcall_wrapper(void *pparams);

void *timeout_function(void *ptimeout);

int PcallWithTimeoutBridge(void *pl, int nargs, int nresults, int msgh, int timeout);

#endif