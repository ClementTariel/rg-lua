#ifndef RG_H
#define RG_H

#include <stdbool.h>
#include <pthread.h>
#include <lua.h>
#include <lauxlib.h>
#include <lualib.h>
#include "rgentities/rgentities.h"

int GetActionWithTimeoutBridge(void *pl, void *paction, int bot_id, int timeout);

int rg_walk_dist_in_lua(lua_State *pl);

int rg_dist_in_lua(lua_State *pl);

int locs_equal_in_lua(lua_State *pl);

int create_loc_in_lua(lua_State *pl);

int custom_next(lua_State *pl);

int robots_pairs_in_lua(lua_State *pl);

int loc_map_index_in_lua(lua_State *pl);

int create_loc_table_in_lua(lua_State *pl);

int rg_locs_around_in_lua(lua_State *pl);

int rg_loc_type_in_lua(lua_State *pl);

int rg_toward_in_lua(lua_State *pl);

int luaopen_librobotgame(lua_State *pl);

void LoadRg(void *pl);

#endif