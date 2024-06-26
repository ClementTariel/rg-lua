#include <pthread.h>
#include <unistd.h>
#include "debug.h"
#include "rg.h"

#define ARENA_RADIUS 8
#define GRID_SIZE 2*ARENA_RADIUS + 3

int get_action(void* pl, void* paction, int bot_id) {
    int clean_stack_size = lua_gettop(pl);
    int err = 0;
    // __RG_CORE_SYSTEM is used as a buffer to simplify data transfer
    lua_getglobal((lua_State*) pl, "__RG_CORE_SYSTEM");
    if (!lua_istable(pl, -1)){
        (pl, lua_gettop(pl) - clean_stack_size);
        return 101;
    }
    lua_getfield(pl, -1, "self");
    if (!lua_istable(pl, -1)){
        lua_pop(pl, lua_gettop(pl) - clean_stack_size);
        return 101;
    }
    lua_geti(pl, -1, bot_id);
    if (!lua_istable(pl, -1)){
        lua_pop(pl, lua_gettop(pl) - clean_stack_size);
        return 101;
    }
    lua_getfield(pl, -3, "act");
    if (!lua_isfunction(pl, -1)){
        lua_pop(pl, lua_gettop(pl) - clean_stack_size);
        return 102;
    }
    lua_pushvalue(pl, -2);
    lua_getfield(pl, -5, "game");
    if (!lua_istable(pl, -1)){
        lua_pop(pl, lua_gettop(pl) - clean_stack_size);
        return 101;
    }

    err = lua_pcall(pl, 2, 1, 0); // 2 arguments, one result
    if (err != 0){
        lua_pop(pl, lua_gettop(pl) - clean_stack_size);
        return err;
    }
    if (!lua_istable(pl, -1)){
        lua_pop(pl, lua_gettop(pl) - clean_stack_size);
        return 103;
    }
    
    lua_getfield(pl, -1, "actionType");
    if (lua_isnumber(pl, -1)){
        int val = lua_tointeger(pl, -1);
        ((struct Action*)paction)->actionType = val;
    } else {
        err = 104;
    }
    lua_pop(pl, 1);

    lua_getfield(pl, -1, "x");
    if (lua_isnumber(pl, -1)){
        int val = lua_tointeger(pl, -1);
        ((struct Action*)paction)->x = val;
    } else if (!lua_isnoneornil(pl, -1)) {
        err = 105;
    }
    lua_pop(pl, 1);

    lua_getfield(pl, -1, "y");
    if (lua_isnumber(pl, -1)){
        int val = lua_tointeger(pl, -1);
        ((struct Action*)paction)->y = val;
    } else if (!lua_isnoneornil(pl, -1)) {
        err = 105;
    }
    lua_pop(pl, 1);
    
    lua_pop(pl, lua_gettop(pl) - clean_stack_size);
    return err;
}

void* get_action_wrapper(void* pparams){
    // allow cancel when stuck in infinite loop
    pthread_setcanceltype(PTHREAD_CANCEL_ASYNCHRONOUS,NULL);
    get_action_thread_params params = *(get_action_thread_params*)pparams;
    void* pl = params.pl;
    void* paction = params.paction;
    int bot_id = params.bot_id;
    int* perr = params.perr;
    bool* pdone = params.pdone;
    pthread_t timeout_thread_id = params.timeout_thread_id;
    int err = get_action(pl, paction, bot_id);
    *perr = err;
    *pdone = true;
    pthread_cancel(timeout_thread_id);
    return NULL;
}

void* timeout_function(void* ptimeout){
    int timeout = *(int*)ptimeout;
    struct timespec ts;
    ts.tv_nsec = (timeout*1000000) % ((int)1e9);
    ts.tv_sec = timeout / 1000;
    nanosleep(&ts, NULL);
    return NULL;
}

int getActionWithTimeoutBridge(void* pl, void* paction, int bot_id, int timeout) {
    bool done = false;
    
    pthread_t timeout_thread_id;

    pthread_create(&timeout_thread_id, NULL, timeout_function, &timeout);
    
    pthread_t action_thread_id;
    int err = 0;
    get_action_thread_params params;
    params.pl = pl;
    params.paction = paction;
    params.bot_id = bot_id;
    params.perr = &err;
    params.pdone = &done;
    params.timeout_thread_id = timeout_thread_id;

    pthread_create(&action_thread_id, NULL, get_action_wrapper, &params);

    pthread_join(timeout_thread_id, NULL);
    if (!done){
        pthread_cancel(action_thread_id);
        if (err == 0){
            err = 106; // timeout
        }
    }
    pthread_join(action_thread_id, NULL);
    return err;
}

int rg_walk_dist_in_lua(lua_State* pl) {
    int argc = lua_gettop(pl);
    if (argc != 2){
        return 0;
    }
    if (!(lua_istable(pl, 1) && lua_istable(pl, 2)) ) {
        return 0;
    }
    lua_getfield(pl, -2, "x");
    if (!lua_isnumber(pl, -1)){
        return 0;
    }
    int x1 = lua_tointeger(pl, -1);
    lua_getfield(pl, -3, "y");
    if (!lua_isnumber(pl, -1)){
        return 0;
    }
    int y1 = lua_tointeger(pl, -1);
    lua_getfield(pl, -3, "x");
    if (!lua_isnumber(pl, -1)){
        return 0;
    }
    int x2 = lua_tointeger(pl, -1);
    lua_getfield(pl, 4, "y");
    if (!lua_isnumber(pl, -1)){
        return 0;
    }
    int y2 = lua_tointeger(pl, -1);
    int d = abs(x1 - x2) + abs(y1 - y2);
    lua_pushinteger(pl, d);
    return 1;
}

int locs_equal_in_lua(lua_State* pl) {
    int argc = lua_gettop(pl);
    if (argc != 2){
        return 0;
    }
    lua_pushcfunction(pl, rg_walk_dist_in_lua);
    lua_pushvalue(pl, -3);
    lua_pushvalue(pl, -3);
    // call with 2 arguments, one result
    if (lua_pcall(pl, 2, 1, 0) != 0 || !lua_isnumber(pl, -1)){
        return 0;
    }
    int dist = lua_tointeger(pl, -1);
    lua_pushboolean(pl, dist == 0);
    return 1;
}


int create_loc_in_lua(lua_State* pl) {
    int argc = lua_gettop(pl);
    if (argc != 2){
        return 0;
    }
    if (!lua_isnumber(pl, -1) || !lua_isnumber(pl, -2) ){
        return 0;
    }
    int x = lua_tointeger(pl, -2);
    int y = lua_tointeger(pl, -1);
    lua_createtable(pl, 0, 2);
    lua_pushinteger(pl, x);
    lua_setfield(pl, -2, "x");
    lua_pushinteger(pl, y);
    lua_setfield(pl, -2, "y");
    lua_createtable(pl, 0, 1);
    lua_pushcfunction(pl, locs_equal_in_lua);
    lua_setfield(pl, -2, "__eq");
    lua_setmetatable(pl, -2);
    return 1;
}

int loc_map_index_in_lua(lua_State* pl) {
    int argc = lua_gettop(pl);
    if (argc != 2 || !lua_istable(pl, -2)){
        return 0;
    }
    if (!lua_istable(pl, -1)){
        lua_pushnil(pl);
        return 1;
    }
    lua_getfield(pl, -1, "x");
    if (!lua_isnumber(pl, -1)){
        lua_pushnil(pl);
        return 1;
    }
    int x = lua_tointeger(pl, -1);
    lua_getfield(pl, -2, "y");
    if (!lua_isnumber(pl, -1)){
        lua_pushnil(pl);
        return 1;
    }
    int y = lua_tointeger(pl, -1);
    if (x < 0 || x > GRID_SIZE || y < 0 || y > GRID_SIZE){
        lua_pushnil(pl);
        return 1;
    }
    lua_geti(pl, -4, x);
    if (!lua_istable(pl, -1)){
        lua_pushnil(pl);
        return 1;
    }
    lua_geti(pl, -1, y);
    return 1;
}

int create_loc_table_in_lua(lua_State* pl) {
    int argc = lua_gettop(pl);
    if (argc != 0){
        return 0;
    }
    lua_createtable(pl, 0, 0);
    lua_createtable(pl, 0, 1);
    lua_pushcfunction(pl, loc_map_index_in_lua);
    lua_setfield(pl, -2, "__index");
    lua_setmetatable(pl, -2);
    return 1;
}

int rg_locs_around_in_lua(lua_State* pl) {
    int argc = lua_gettop(pl);
    if (argc != 1){
        return 0;
    }
    lua_getfield(pl, -1, "x");
    if (!lua_isnumber(pl, -1)){
        return 0;
    }
    int x = lua_tointeger(pl, -1);
    lua_getfield(pl, -2, "y");
    if (!lua_isnumber(pl, -1)){
        return 0;
    }
    int y = lua_tointeger(pl, -1);
    lua_createtable(pl, 4, 0);
    for (int i=0; i<4; i++){
        lua_pushcfunction(pl, create_loc_in_lua);
        lua_pushinteger(pl, x + (i%2) * ((i-2)%4));
        lua_pushinteger(pl, y + ((i+1)%2) * ((i-1)%4));
        if (lua_pcall(pl, 2, 1, 0) != 0){
            return 0;
        }
        lua_seti(pl, -2, i + 1);
    }
    return 1;
}



int luaopen_librobotgame(lua_State* pl)
{
    static const struct luaL_Reg robotGameLib [] = {
        {"wdist", rg_walk_dist_in_lua},
        {"locs_around", rg_locs_around_in_lua},
        {"Loc", create_loc_in_lua},
        {"Game", create_loc_table_in_lua},
        {NULL, NULL}
    };
    luaL_newlib(pl, robotGameLib);
    lua_pushinteger(pl, GRID_SIZE);
    lua_setfield(pl, -2, "GRID_SIZE");
    lua_pushinteger(pl, ARENA_RADIUS);
    lua_setfield(pl, -2, "ARENA_RADIUS");
    lua_pushcfunction(pl, create_loc_in_lua);
    lua_pushinteger(pl, ARENA_RADIUS + 1);
    lua_pushinteger(pl, ARENA_RADIUS + 1);
    if (lua_pcall(pl, 2, 1, 0) != 0){
        return 0;
    }
    lua_setfield(pl, -2, "CENTER_POINT");
    return 1;
}

void loadRg(void* pl) {
	luaL_requiref((lua_State*)pl, "rg", luaopen_librobotgame, 1);
	lua_pop((lua_State*)pl, 1);
}