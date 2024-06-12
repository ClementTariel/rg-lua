#ifndef RG_H
#define RG_H

#include <lua.h>
#include <lauxlib.h>
#include <lualib.h>

enum ActionType { MOVE = 0, ATTACK= 1, GUARD = 2, SUICIDE = 3};

struct Action {
	int actionType;
	int x;
	int y;
};

typedef struct get_action_thread_params {
    void* pl;
    void* paction;
    int* perr;
    bool* pdone;
    pthread_t timeout_thread_id;
} get_action_thread_params;

int get_action(void* pl, void* paction);

void* get_action_wrapper(void* pparams);

void* timeout_function(void* ptimeout);

int getActionWithTimeoutBridge(void* pl, void* paction, int timeout);

#endif