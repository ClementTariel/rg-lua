#include <unistd.h>
#include <math.h>
#include "rgdebug/debug.h"
#include "rg.h"
#include "lua/luaBridge.c"
#include "grid/grid.c"

int get_action(void *pl, void *paction, int bot_id)
{
    int clean_stack_size = lua_gettop(pl);
    int err = 0;
    // __RG_CORE_SYSTEM is used as a buffer to simplify data transfer
    lua_getglobal((lua_State *)pl, "__RG_CORE_SYSTEM");
    if (!lua_istable(pl, -1))
    {
        (pl, lua_gettop(pl) - clean_stack_size);
        return 101;
    }
    lua_getfield(pl, -1, "self");
    if (!lua_istable(pl, -1))
    {
        lua_pop(pl, lua_gettop(pl) - clean_stack_size);
        return 101;
    }
    lua_geti(pl, -1, bot_id);
    if (!lua_istable(pl, -1))
    {
        lua_pop(pl, lua_gettop(pl) - clean_stack_size);
        return 101;
    }
    lua_getfield(pl, -3, "game");
    if (!lua_istable(pl, -1))
    {
        lua_pop(pl, lua_gettop(pl) - clean_stack_size);
        return 101;
    }
    lua_getfield(pl, -1, "robots");
    if (!lua_istable(pl, -1))
    {
        lua_pop(pl, lua_gettop(pl) - clean_stack_size);
        return 101;
    }
    lua_getfield(pl, -5, "act");
    if (!lua_isfunction(pl, -1))
    {
        lua_pop(pl, lua_gettop(pl) - clean_stack_size);
        return 102;
    }
    lua_pushvalue(pl, -4);
    lua_pushvalue(pl, -4);
    err = lua_pcall(pl, 2, 1, 0); // 2 arguments, one result
    if (err != 0)
    {
        lua_pop(pl, lua_gettop(pl) - clean_stack_size);
        return err;
    }
    if (!lua_istable(pl, -1))
    {
        lua_pop(pl, lua_gettop(pl) - clean_stack_size);
        return 103;
    }

    lua_getfield(pl, -1, "actionType");
    if (lua_isinteger(pl, -1))
    {
        int val = lua_tointeger(pl, -1);
        ((Action *)paction)->actionType = val;
    }
    else
    {
        err = 104;
    }
    lua_pop(pl, 1);

    lua_getfield(pl, -1, "x");
    if (lua_isinteger(pl, -1))
    {
        int val = lua_tointeger(pl, -1);
        ((Action *)paction)->x = val;
    }
    else if (!lua_isnoneornil(pl, -1))
    {
        err = 105;
    }
    lua_pop(pl, 1);

    lua_getfield(pl, -1, "y");
    if (lua_isinteger(pl, -1))
    {
        int val = lua_tointeger(pl, -1);
        ((Action *)paction)->y = val;
    }
    else if (!lua_isnoneornil(pl, -1))
    {
        err = 105;
    }
    lua_pop(pl, 1);

    lua_pop(pl, lua_gettop(pl) - clean_stack_size);
    return err;
}

void *get_action_wrapper(void *pparams)
{
    // allow cancel when stuck in infinite loop
    // TODO: not safe, might mess with memory allocation
    pthread_setcanceltype(PTHREAD_CANCEL_ASYNCHRONOUS, NULL);
    get_action_thread_params params = *(get_action_thread_params *)pparams;
    void *pl = params.pl;
    void *paction = params.paction;
    int bot_id = params.bot_id;
    int *perr = params.perr;
    bool *pdone = params.pdone;
    pthread_t timeout_thread_id = params.timeout_thread_id;
    int err = get_action(pl, paction, bot_id);
    *perr = err;
    *pdone = true;
    // // TODO: WIP error here rework timeout cancelation system
    // pthread_cancel(timeout_thread_id);
    return NULL;
}

int GetActionWithTimeoutBridge(void *pl, void *paction, int bot_id, int timeout)
{
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
    if (!done)
    {
        // // TODO: WIP error here rework timeout cancelation system
        // pthread_cancel(action_thread_id);
        if (err == 0)
        {
            err = CUSTOM_TIMEOUT_ERROR; // timeout
        }
    }
    pthread_join(action_thread_id, NULL);
    return err;
}

int rg_walk_dist_in_lua(lua_State *pl)
{
    int argc = lua_gettop(pl);
    if (argc != 2)
    {
        return 0;
    }
    if (!(lua_istable(pl, 1) && lua_istable(pl, 2)))
    {
        return 0;
    }
    lua_getfield(pl, 1, "x");
    if (!lua_isinteger(pl, -1))
    {
        return 0;
    }
    int x1 = lua_tointeger(pl, -1);
    lua_getfield(pl, 1, "y");
    if (!lua_isinteger(pl, -1))
    {
        return 0;
    }
    int y1 = lua_tointeger(pl, -1);
    lua_getfield(pl, 2, "x");
    if (!lua_isinteger(pl, -1))
    {
        return 0;
    }
    int x2 = lua_tointeger(pl, -1);
    lua_getfield(pl, 2, "y");
    if (!lua_isinteger(pl, -1))
    {
        return 0;
    }
    int y2 = lua_tointeger(pl, -1);
    int d = abs(x1 - x2) + abs(y1 - y2);
    lua_pushinteger(pl, d);
    return 1;
}

int rg_dist_in_lua(lua_State *pl)
{
    int argc = lua_gettop(pl);
    if (argc != 2)
    {
        return 0;
    }
    if (!(lua_istable(pl, 1) && lua_istable(pl, 2)))
    {
        return 0;
    }
    lua_getfield(pl, 1, "x");
    if (!lua_isinteger(pl, -1))
    {
        return 0;
    }
    int x1 = lua_tointeger(pl, -1);
    lua_getfield(pl, 1, "y");
    if (!lua_isinteger(pl, -1))
    {
        return 0;
    }
    int y1 = lua_tointeger(pl, -1);
    lua_getfield(pl, 2, "x");
    if (!lua_isinteger(pl, -1))
    {
        return 0;
    }
    int x2 = lua_tointeger(pl, -1);
    lua_getfield(pl, 2, "y");
    if (!lua_isinteger(pl, -1))
    {
        return 0;
    }
    int y2 = lua_tointeger(pl, -1);
    double d = sqrt((double)((x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2)));
    lua_pushnumber(pl, d);
    return 1;
}

int locs_equal_in_lua(lua_State *pl)
{
    int argc = lua_gettop(pl);
    if (argc != 2)
    {
        return 0;
    }
    lua_pushcfunction(pl, rg_walk_dist_in_lua);
    lua_pushvalue(pl, 1);
    lua_pushvalue(pl, 2);
    // call with 2 arguments, one result
    if (lua_pcall(pl, 2, 1, 0) != 0 || !lua_isinteger(pl, -1))
    {
        return 0;
    }
    int dist = lua_tointeger(pl, -1);
    lua_pushboolean(pl, dist == 0);
    return 1;
}

int create_loc_in_lua(lua_State *pl)
{
    int argc = lua_gettop(pl);
    if (argc != 2)
    {
        return 0;
    }
    if (!lua_isinteger(pl, 1) || !lua_isinteger(pl, 2))
    {
        return 0;
    }
    int x = lua_tointeger(pl, 1);
    int y = lua_tointeger(pl, 2);
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

int custom_next(lua_State *pl)
{
    int argc = lua_gettop(pl);
    if (argc != 2 || !lua_istable(pl, 1) || (!lua_istable(pl, 2) && !lua_isnil(pl, 2)))
    {
        return 0;
    }
    lua_pushvalue(pl, 1);
    if (lua_istable(pl, 2))
    {
        lua_getfield(pl, 2, "x");
    }
    else
    {
        lua_pushnil(pl);
    }
    if (!lua_isinteger(pl, -1))
    {
        // stack : -1 => nil; -2 => game
        if (lua_next(pl, -2) == 0)
        {
            lua_pushnil(pl);
            return 1;
        }
        // stack : -1 => game[x]; -2 => x; -3 => game
    }
    else
    {
        int x = lua_tointeger(pl, -1);
        lua_geti(pl, 1, x);
        // stack : -1 => game[x]; -2 => x; -3 => game
    }
    if (lua_istable(pl, 2))
    {
        lua_getfield(pl, 2, "y");
    }
    else
    {
        lua_pushnil(pl);
    }
    // stack : -1 => y; -2 => game[x]; -3 => x; -4 => game
    int end = false;
    int next_game_x = lua_next(pl, -2);
    while (next_game_x == 0 && !end)
    {
        // stack : -1 => game[x]; -2 => x; -3 => game
        lua_pop(pl, 1);
        // stack : -1 => x; -2 => game
        if (lua_next(pl, -2) != 0)
        {
            // stack : -1 => game[x+1]; -2 => x+1; -3 => game
            lua_pushnil(pl);
            // stack : -1 => nil; -2 => game[x+1]; -3 => x+1; -4 => game
            next_game_x = lua_next(pl, -2);
        }
        else
        {
            end = true;
        }
    }
    if (end)
    {
        lua_pushnil(pl);
        return 1;
    }
    // stack : -1 => game[x][y]; -2 => y; -3 => game[x]; -4 => x; -5 => game
    if (!lua_isinteger(pl, -4) || !lua_isinteger(pl, -2))
    {
        return 0;
    }
    int x = lua_tointeger(pl, -4);
    int y = lua_tointeger(pl, -2);
    lua_pushcfunction(pl, create_loc_in_lua);
    lua_pushinteger(pl, x);
    lua_pushinteger(pl, y);
    if (lua_pcall(pl, 2, 1, 0) != 0)
    {
        return 0;
    }
    // stack : -1 => rg.Loc(x, y); -2 => game[x][y]; ...
    lua_pushvalue(pl, -2);
    // stack : -1 => game[x][y]; -2 => rg.Loc(x, y); -3 => game[x][y]; ...
    return 2;
}

int robots_pairs_in_lua(lua_State *pl)
{
    int argc = lua_gettop(pl);
    if (argc != 1)
    {
        return 0;
    }
    if (!lua_istable(pl, -1))
    {
        return 0;
    }
    lua_pushcfunction(pl, custom_next); // iterator
    lua_pushvalue(pl, 1);               // robots
    lua_pushnil(pl);
    return 3;
}

int loc_map_index_in_lua(lua_State *pl)
{
    int argc = lua_gettop(pl);
    if (argc != 2 || !lua_istable(pl, 1))
    {
        return 0;
    }
    if (!lua_istable(pl, 2))
    {
        lua_pushnil(pl);
        return 1;
    }
    lua_getfield(pl, 1, "x");
    if (!lua_isinteger(pl, -1))
    {
        lua_pushnil(pl);
        return 1;
    }
    int x = lua_tointeger(pl, -1);
    lua_getfield(pl, 1, "y");
    if (!lua_isinteger(pl, -1))
    {
        lua_pushnil(pl);
        return 1;
    }
    int y = lua_tointeger(pl, -1);
    if (x < 0 || x > GRID_SIZE || y < 0 || y > GRID_SIZE)
    {
        lua_pushnil(pl);
        return 1;
    }
    lua_geti(pl, 1, x);
    if (!lua_istable(pl, -1))
    {
        lua_pushnil(pl);
        return 1;
    }
    lua_geti(pl, -1, y);
    return 1;
}

int create_loc_table_in_lua(lua_State *pl)
{
    int argc = lua_gettop(pl);
    if (argc != 0)
    {
        return 0;
    }
    lua_createtable(pl, 0, 0);
    lua_createtable(pl, 0, 1);
    lua_pushcfunction(pl, loc_map_index_in_lua);
    lua_setfield(pl, -2, "__index");
    lua_pushcfunction(pl, robots_pairs_in_lua);
    lua_setfield(pl, -2, "__pairs");
    lua_setmetatable(pl, -2);
    return 1;
}

int rg_locs_around_in_lua(lua_State *pl)
{
    int argc = lua_gettop(pl);
    if (argc < 1 || argc > 2)
    {
        return 0;
    }
    lua_getfield(pl, 1, "x");
    if (!lua_isinteger(pl, -1))
    {
        return 0;
    }
    int x = lua_tointeger(pl, -1);
    lua_getfield(pl, 1, "y");
    if (!lua_isinteger(pl, -1))
    {
        return 0;
    }
    int y = lua_tointeger(pl, -1);
    bool filter_out_normal = false;
    bool filter_out_spawn = false;
    bool filter_out_obstacle = false;
    if (argc == 2)
    {
        if (!lua_istable(pl, 2))
        {
            return 0;
        }
        lua_pushvalue(pl, 2);
        lua_pushnil(pl);
        // stack : -1 => nil; -2 => table;
        while (lua_next(pl, -2) != 0)
        {
            // stack : -1 => value; -2 => key; -3 => table;
            if (!lua_isinteger(pl, -1))
            {
                return 0;
            }
            int loc_type = lua_tointeger(pl, -1);
            lua_pop(pl, 1);
            // stack now contains: -1 => key; -2 => table
            if (loc_type == NORMAL)
            {
                filter_out_normal = true;
            }
            else if (loc_type == SPAWN)
            {
                filter_out_spawn = true;
            }
            else if (loc_type == OBSTACLE)
            {
                filter_out_obstacle = true;
            }
            else
            {
                return 0;
            }
        }
    }
    lua_createtable(pl, 4, 0);
    int loc_count = 0;
    for (int i = 0; i < 4; i++)
    {
        int loc_x = x + (i % 2) * (i - 2);
        int loc_y = y + ((i + 1) % 2) * (i - 1);
        int loc_type = OBSTACLE;
        if (loc_x > 0 && loc_x < GRID_SIZE && loc_y > 0 && loc_y < GRID_SIZE)
        {
            loc_type = GRID[loc_x][loc_y];
        }
        bool filtered_out = false;
        if (loc_type == NORMAL)
        {
            filtered_out = filter_out_normal;
        }
        else if (loc_type == SPAWN)
        {
            filtered_out = filter_out_spawn;
        }
        else if (loc_type == OBSTACLE)
        {
            filtered_out = filter_out_obstacle;
        }
        if (!filtered_out)
        {
            loc_count++;
            lua_pushcfunction(pl, create_loc_in_lua);
            lua_pushinteger(pl, loc_x);
            lua_pushinteger(pl, loc_y);
            if (lua_pcall(pl, 2, 1, 0) != 0)
            {
                return 0;
            }
            lua_seti(pl, -2, loc_count);
        }
    }
    return 1;
}

int rg_loc_type_in_lua(lua_State *pl)
{
    int argc = lua_gettop(pl);
    if (argc != 1)
    {
        return 0;
    }
    lua_getfield(pl, 1, "x");
    if (!lua_isinteger(pl, -1))
    {
        return 0;
    }
    int x = lua_tointeger(pl, -1);
    lua_getfield(pl, 1, "y");
    if (!lua_isinteger(pl, -1))
    {
        return 0;
    }
    int y = lua_tointeger(pl, -1);
    if (x < 0 || y >= GRID_SIZE)
    {
        lua_pushinteger(pl, OBSTACLE);
    }
    else
    {
        lua_pushinteger(pl, GRID[x][y]);
    }
    return 1;
}

int rg_toward_in_lua(lua_State *pl)
{
    int argc = lua_gettop(pl);
    if (argc != 2)
    {
        return 0;
    }
    if (!(lua_istable(pl, 1) && lua_istable(pl, 2)))
    {
        return 0;
    }
    lua_getfield(pl, 1, "x");
    if (!lua_isinteger(pl, -1))
    {
        return 0;
    }
    int x1 = lua_tointeger(pl, -1);
    lua_getfield(pl, 1, "y");
    if (!lua_isinteger(pl, -1))
    {
        return 0;
    }
    int y1 = lua_tointeger(pl, -1);
    lua_getfield(pl, 2, "x");
    if (!lua_isinteger(pl, -1))
    {
        return 0;
    }
    int x2 = lua_tointeger(pl, -1);
    lua_getfield(pl, 2, "y");
    if (!lua_isinteger(pl, -1))
    {
        return 0;
    }
    int y2 = lua_tointeger(pl, -1);
    int loc_x = x1;
    int loc_y = y1;
    if (x1 != x2 || y1 != y2)
    {
        for (int i = 0; i < 4; i++)
        {
            int x = x1 + (i % 2) * (i - 2);
            int y = y1 + ((i + 1) % 2) * (i - 1);
            if (abs(x2 - loc_x) + abs(y2 - loc_y) > abs(x2 - x) + abs(y2 - y))
            {
                loc_x = x;
                loc_y = y;
                continue;
            }
            if (abs(x2 - loc_x) + abs(y2 - loc_y) < abs(x2 - x) + abs(y2 - y))
            {
                continue;
            }
            if (abs(x2 - loc_x) * abs(y2 - loc_y) < abs(x2 - x) * abs(y2 - y))
            {
                loc_x = x;
                loc_y = y;
                continue;
            }
            if (abs(x2 - loc_x) * abs(y2 - loc_y) > abs(x2 - x) * abs(y2 - y))
            {
                continue;
            }
            if ((x2 - x1) * (loc_y - y1) - (y2 - y1) * (loc_x - x1) < (x2 - x1) * (y - y1) - (y2 - y1) * (x - x1))
            {
                loc_x = x;
                loc_y = y;
                continue;
            }
        }
    }
    lua_pushcfunction(pl, create_loc_in_lua);
    lua_pushinteger(pl, loc_x);
    lua_pushinteger(pl, loc_y);
    if (lua_pcall(pl, 2, 1, 0) != 0)
    {
        return 0;
    }
    return 1;
}

int luaopen_librobotgame(lua_State *pl)
{
    static const struct luaL_Reg robotGameLib[] = {
        {"wdist", rg_walk_dist_in_lua},
        {"dist", rg_dist_in_lua},
        {"locs_around", rg_locs_around_in_lua},
        {"loc_type", rg_loc_type_in_lua},
        {"toward", rg_toward_in_lua},
        {"Loc", create_loc_in_lua},
        {"Robots", create_loc_table_in_lua},
        {NULL, NULL}};
    luaL_newlib(pl, robotGameLib);
    lua_pushinteger(pl, GRID_SIZE);
    lua_setfield(pl, -2, "GRID_SIZE");
    lua_pushinteger(pl, ARENA_RADIUS);
    lua_setfield(pl, -2, "ARENA_RADIUS");
    lua_pushcfunction(pl, create_loc_in_lua);
    lua_pushinteger(pl, ARENA_RADIUS + 1);
    lua_pushinteger(pl, ARENA_RADIUS + 1);
    if (lua_pcall(pl, 2, 1, 0) != 0)
    {
        return 0;
    }
    lua_setfield(pl, -2, "CENTER_POINT");
    return 1;
}

void LoadRg(void *pl)
{
    luaL_requiref((lua_State *)pl, "rg", luaopen_librobotgame, 1);
    lua_pop((lua_State *)pl, 1);
}