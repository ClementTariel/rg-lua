#include "luaBridge.h"

const int CUSTOM_TIMEOUT_ERROR = 408;

static pthread_mutex_t timeout_mutex;
static volatile bool done = false;
static volatile bool timeout_triggered = 0;
static volatile int timeout_count = 0;

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
  return PcallWithTimeoutBridge(pl, nargs, nresults, msgh, -1);
}

void *timeout_function(void *pparams)
{
  timeout_params params = *(timeout_params *)pparams;
  int timeout = params.timeout;
  if (timeout < 0)
  {
    return NULL;
  }

  // TODO maybe take into account a slight shift caused by accumulated errors
  struct timespec ts;
  // TODO: reduce number of chuncks if performance decreases too much
  int chunck_count = 100;
  ts.tv_nsec = (timeout * 1000000 / chunck_count) % ((int)1e9);
  ts.tv_sec = timeout / (1000 * chunck_count);
  bool pending = true;
  while (pending)
  {
    for (int i = 0; i < chunck_count; i++)
    {
      nanosleep(&ts, NULL);
      // Shoud be safe to read without mutex because it is a boolean.
      pending = !done;
      if (!pending)
      {
        // make sure it is not a misreading of done because of mutex bypass
        pthread_mutex_lock(&timeout_mutex);
        pending = !done;
        pthread_mutex_unlock(&timeout_mutex);
        if (!pending)
        {
          return NULL;
        }
      }
    }
    // Should be safe to call without mutex because it is a boolean
    if (!timeout_triggered)
    {
      // Should be safe to modify without mutex because it is a boolean,
      // no other thread can modify it and it can only be set to true
      timeout_triggered = true;
    }
    pthread_mutex_lock(&timeout_mutex);
    timeout_count++;
    pthread_mutex_unlock(&timeout_mutex);
  }
}

void timeout_hook(lua_State *pl, lua_Debug *ar)
{
  // Should be safe to call without mutex because it is a boolean
  if (timeout_triggered)
  {
    pthread_mutex_lock(&timeout_mutex);
    int local_timeout_count = timeout_count;
    pthread_mutex_unlock(&timeout_mutex);
    // make sure that there was no misreading caused by mutex bypass
    if (local_timeout_count > 0)
    {
      // Should be safe to modify without mutex because it is a boolean,
      // no other thread can modify it and it can only be set to true
      done = true;
      luaL_error(pl, "TimeoutError");
      // code here won't be reached
    }
  }
}

int PcallWithTimeoutBridge(void *pl, int nargs, int nresults, int msgh, int timeout)
{
  if (timeout == 0)
  {
    return CUSTOM_TIMEOUT_ERROR;
  }

  int err = 0;
  pthread_t timeout_thread_id;

  int mutex_err = pthread_mutex_init(&timeout_mutex, NULL);
  if (mutex_err != 0)
  {
    return mutex_err;
  }

  pthread_mutex_lock(&timeout_mutex);
  done = false;
  timeout_count = 0;
  pthread_mutex_unlock(&timeout_mutex);

  timeout_params params;
  params.timeout = timeout;
  pthread_create(&timeout_thread_id, NULL, timeout_function, &params);

  if (timeout > 0)
  {
    // ~ same number of hook calls for any timeout
    int num_instruction = timeout * 10;
    lua_sethook(pl, &timeout_hook, LUA_MASKCOUNT, num_instruction);
  }
  err = lua_pcall((lua_State *)pl, nargs, nresults, msgh);
  lua_sethook(pl, &timeout_hook, 0, 0);

  pthread_mutex_lock(&timeout_mutex);
  done = true;
  if (timeout_count > 0)
  {
    err = CUSTOM_TIMEOUT_ERROR; // timeout
  }
  pthread_mutex_unlock(&timeout_mutex);
  pthread_join(timeout_thread_id, NULL);
  mutex_err = pthread_mutex_destroy(&timeout_mutex);
  if (mutex_err != 0)
  {
    return mutex_err;
  }
  return err;
}