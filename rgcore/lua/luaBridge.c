#include "luaBridge.h"

const int CUSTOM_TIMEOUT_ERROR = 408;

// Copy luaB_next and luaB_pairs from source code to avoid importing an unsafe lib
static int luaB_next(lua_State *L)
{
  luaL_checktype(L, 1, LUA_TTABLE);
  lua_settop(L, 2); /* create a 2nd argument if there isn't one */
  if (lua_next(L, 1))
    return 2;
  else
  {
    lua_pushnil(L);
    return 1;
  }
}

static int pairsmeta(lua_State *L, const char *method, int iszero,
                     lua_CFunction iter)
{
  if (!luaL_getmetafield(L, 1, method))
  {                                   /* no metamethod? */
    luaL_checktype(L, 1, LUA_TTABLE); /* argument must be a table */
    lua_pushcfunction(L, iter);       /* will return generator, */
    lua_pushvalue(L, 1);              /* state, */
    if (iszero)
      lua_pushinteger(L, 0); /* and initial value */
    else
      lua_pushnil(L);
  }
  else
  {
    lua_pushvalue(L, 1); /* argument 'self' to metamethod */
    lua_call(L, 1, 3);   /* get 3 values from metamethod */
  }
  return 3;
}
static int luaB_pairs(lua_State *L)
{
  return pairsmeta(L, "__pairs", 0, luaB_next);
}

int GetStackBridge(void *pl, int depth, void *ptDebug)
{
  return lua_getstack((lua_State *)pl, depth, (lua_Debug *)ptDebug);
}

void GetInfoBridge(void *pl, void *ptDebug)
{
  lua_getinfo((lua_State *)pl, "Snlu", (lua_Debug *)ptDebug);
  printf("%d\n", ((lua_Debug *)ptDebug)->currentline);
}

int DebugCurrentLine(void *ptDebug)
{
  // printf("%d\n",((lua_Debug*)ptDebug)->currentline);
  // return ((lua_Debug*)ptDebug)->currentline;
  return -1;
}

const char *ToStringBridge(void *pl, int idx)
{
  return lua_tostring(pl, (int)idx);
}

void PopStateBridge(void *pl, int idx)
{
  lua_settop((lua_State *)pl, (int)(-(idx + 1)));
}

lua_State *NewStateBridge()
{
  lua_State *pl = luaL_newstate();
  lua_pushcfunction(pl, luaB_next);
  lua_setglobal((lua_State *)pl, "next");
  lua_pushcfunction(pl, luaB_pairs);
  lua_setglobal((lua_State *)pl, "pairs");
  luaL_requiref(pl, "table", luaopen_table, 1);
  lua_pop(pl, 1);
  luaL_requiref(pl, "math", luaopen_math, 1);
  lua_pop(pl, 1);
  return pl;
}

void CloseBridge(void *pl)
{
  lua_close((lua_State *)pl);
}

void PushCFunctionBridge(void *pl, void *fn)
{
  lua_pushcfunction(((lua_State *)pl), (lua_CFunction)fn);
}

void SetGlobalBridge(void *pl, const char *name)
{
  lua_setglobal((lua_State *)pl, name);
}

int LoadStringBridge(void *pl, const char *s)
{
  return luaL_loadstring((lua_State *)pl, s);
}

int PcallBridge(void *pl, int nargs, int nresults, int msgh)
{
  return lua_pcall((lua_State *)pl, nargs, nresults, msgh);
}

void *pcall_wrapper(void *pparams)
{
  // allow cancel when stuck in infinite loop
  // TODO: not safe, might mess with memory allocation
  pthread_setcanceltype(PTHREAD_CANCEL_ASYNCHRONOUS, NULL);
  pcall_thread_params params = *(pcall_thread_params *)pparams;
  void *pl = params.pl;
  int nargs = params.nargs;
  int nresults = params.nresults;
  int msgh = params.msgh;
  int *perr = params.perr;
  bool *pdone = params.pdone;
  pthread_mutex_t *pdone_mutex = params.pdone_mutex;
  pthread_t timeout_thread_id = params.timeout_thread_id;
  int err = lua_pcall((lua_State *)pl, nargs, nresults, msgh);
  *perr = err;
  pthread_mutex_lock(pdone_mutex);
  if (!*pdone)
  {
    *pdone = true;
    pthread_cancel(timeout_thread_id);
  }
  pthread_mutex_unlock(pdone_mutex);
  return NULL;
}

void *timeout_function(void *ptimeout)
{
  // TODO maybe take into account a slight shift caused by accumulated errors
  int timeout = *(int *)ptimeout;
  struct timespec ts;
  int chunck_count = 100;
  ts.tv_nsec = (timeout * 1000000 / chunck_count) % ((int)1e9);
  ts.tv_sec = timeout / (1000 * chunck_count);
  for (int i = 0; i < chunck_count; i++)
  {
    nanosleep(&ts, NULL);
  }
  return NULL;
}

int PcallWithTimeoutBridge(void *pl, int nargs, int nresults, int msgh, int timeout)
{
  bool done = false;
  pthread_mutex_t done_mutex;

  pthread_t timeout_thread_id;
  pthread_t thread_id;

  pthread_mutex_init(&done_mutex, NULL);

  pthread_create(&timeout_thread_id, NULL, timeout_function, &timeout);

  int err = 0;
  pcall_thread_params params;
  params.pl = pl;
  params.perr = &err;
  params.pdone = &done;
  params.pdone_mutex = &done_mutex;
  params.timeout_thread_id = timeout_thread_id;
  params.nargs = nargs;
  params.nresults = nresults;
  params.msgh = msgh;

  pthread_create(&thread_id, NULL, pcall_wrapper, &params);
  pthread_join(timeout_thread_id, NULL);
  pthread_mutex_lock(&done_mutex);
  if (!done)
  {
    done = true;
    pthread_cancel(thread_id);
    if (err == 0)
    {
      err = CUSTOM_TIMEOUT_ERROR; // timeout
    }
  }
  pthread_mutex_unlock(&done_mutex);
  pthread_join(thread_id, NULL);
  return err;
}