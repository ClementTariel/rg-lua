<script setup lang="ts">
  import { computed, ref, watch } from 'vue';

  // TODO: WIP type
  const props = defineProps<{
    game: { turns: any[] },
    metadata: any
  }>();

  const gridRadius: number = 8.5
  const gridSize: number = (gridRadius + 1) * 2;
  const maxTurn = 100;
  const autoRunTurnDuration = 1000;
  enum ActionType {
    Move,
    Attack,
    Guard,
    Suicide
  }
  let gameLoaded = false;
  let animationDone = false;
  setTimeout(() => {
    animationDone = true;
    if (gameLoaded) {
      turn.value = 0;
    }
  }, 1000);


  const turn = ref(-1);

  watch(() => props.game, async (loadedGame, previousGame) => {
    gameLoaded = true;
    if (animationDone) {
      turn.value = 0;
    }
  })

  watch(() => turn.value, async (currentTurn, previousTurn) => {
    if (currentTurn < maxTurn) {
      setTimeout(() => {
        turn.value = currentTurn + 1;
      }, autoRunTurnDuration);
    }
  })


  const grid = computed(() => {
    const computedGrid = Array.from({ length: gridSize }, _ => Array(gridSize));
    if (turn.value >= 0 && props.game?.turns && props.game.turns.length && props.game.turns.length > turn.value) {
      // TODO: WIP type
      Object.values(props.game.turns[turn.value]).forEach((tile: any) => {
        computedGrid[tile.Bot.X][tile.Bot.Y] = tile
      });
    }
    return computedGrid;
  })

  function squaredDistToCenter(x: number, y: number): number {
    const center = 1 + ((gridSize - 1) / 2)
    return (x - center) ** 2 + (y - center) ** 2
  }

  // TODO: WIP type
  function getBot(grid: any[][], x: number, y: number): any {
    return grid?.[x - 1]?.[y - 1]?.Bot
  }

  // TODO: WIP type
  function getAction(grid: any[][], x: number, y: number): any {
    return grid?.[x - 1]?.[y - 1]?.Action
  }

</script>

<template>
  <div class="game">
    <div class="banners">
      <div class="banner banner-left blue">
        <span class="user-name">{{ metadata.UserName1 }}</span>
        <span class="bot-name">{{ metadata.BotName1 }}</span>
      </div>
      <div class="banner banner-right red">
        <span class="bot-name">{{ metadata.BotName2 }}</span>
        <span class="user-name">{{ metadata.UserName2 }}</span>
      </div>
    </div>
    <div class="turn-count-box">
      <p class="turn-count">{{ turn >= 0 ? turn : 0 }}</p>
    </div>
    <div class="grid">
      <div v-for="x in gridSize" class="row">
        <div v-for="y in gridSize" :style="`--delay: ${(squaredDistToCenter(x, y) / ((gridSize / 2) ** 2))}s`"
          class="cell" :class="{
            'enabled-cell': squaredDistToCenter(x, y) < ((gridRadius) ** 2),
            'blue': getBot(grid, x, y)?.PlayerId === 1,
            'red': getBot(grid, x, y)?.PlayerId === 2,
            'move': getAction(grid, x, y)?.ActionType === ActionType.Move,
            'attack': getAction(grid, x, y)?.ActionType === ActionType.Attack,
            'guard': getAction(grid, x, y)?.ActionType === ActionType.Guard,
            'suicide': getAction(grid, x, y)?.ActionType === ActionType.Suicide
          }">
          <p class="hp">{{ getBot(grid, x, y)?.Hp }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="css" scoped>
  .game {
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 100%;
  }

  @keyframes spawn-banner-left {
    from {
      transform: translate(-100%, 0%);
    }

    to {
      transform: translate(0%, 0%);
    }
  }

  @keyframes spawn-banner-right {
    from {
      transform: translate(100%, 0%);
    }

    to {
      transform: translate(0%, 0%);
    }
  }

  .banners {
    width: 100%;
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    padding-bottom: 16px;
  }

  .banner,
  .banner-left,
  .banner-right {
    --banner-height: 40px;
  }

  .banner {
    display: flex;
    flex-direction: row;
    align-items: end;
    color: white;
    width: 50%;
    height: var(--banner-height);
  }

  .banner-left {
    padding-right: var(--banner-height);
    justify-content: end;

    clip-path: polygon(0% 100%, 100% 0%, 0% 0%);
    -webkit-clip-path: polygon(0% 100%, calc(100% - var(--banner-height)) 100%, 100% 0%, 0% 0%);

    animation-name: spawn-banner-left;
    animation-duration: 0.5s;
    animation-timing-function: ease-in;
  }

  .banner-right {
    padding-left: var(--banner-height);
    justify-content: baseline;

    clip-path: polygon(0% 100%, 100% 100%, 100% 0%, var(--banner-height) 0%);
    -webkit-clip-path: polygon(0% 100%, 100% 100%, 100% 0%, var(--banner-height) 0%);

    animation-name: spawn-banner-right;
    animation-duration: 0.5s;
    animation-timing-function: ease-in;
  }

  .bot-name {
    display: flex;
    height: fit-content;
    font-size: 28px;
    font-weight: 700;
  }

  .user-name {
    display: flex;
    height: fit-content;
    font-size: 12px;
    font-weight: 700;
    padding: 4px 8px
  }

  @keyframes spawn-turn-count {
    from {
      transform: scale(0%);
      visibility: hidden;
      box-shadow: none;
    }

    80% {
      transform: scale(0%);
      visibility: visible;
      box-shadow: none;
    }

    to {
      transform: scale(100%);
      box-shadow: 10px 20px 20px black;
    }
  }

  .turn-count-box {
    width: 48px;
    border: 1px solid black;
    border-radius: 8px;
    padding: 4px;
    margin: 4px;
    background-color: #5a5a5a;
    text-align: center;
    box-shadow: 10px 20px 20px black;

    animation-name: spawn-turn-count;
    animation-duration: 0.5s;
  }

  .turn-count {
    border: 1px solid black;
    border-radius: 4px;
    padding: 4px;
    margin: 0;
    color: white;
    background-color: #404040;
  }

  @keyframes spawn-grid {
    from {
      transform: scale(0%) rotate(-360deg);
      border-radius: 100%;
      box-shadow: none;
    }

    50% {
      box-shadow: none;
    }

    to {
      transform: scale(100%) rotate(0deg);
      border-radius: 20px;
      box-shadow: 10px 20px 20px black;
    }
  }

  .grid {
    padding: 20px;
    border-radius: 20px;
    background-color: #5a5a5a;
    border: 1px solid black;
    display: flex;
    flex-direction: row;
    width: fit-content;
    height: fit-content;
    box-shadow: 10px 20px 20px black;

    animation-name: spawn-grid;
    animation-duration: 0.5s;
    animation-timing-function: ease-out;
  }

  .row {
    display: flex;
    flex-direction: column;
    width: fit-content;
    height: fit-content;
  }

  @keyframes sleepy-cell {
    from {
      visibility: hidden;
    }

    to {
      visibility: hidden;
    }
  }

  @keyframes spawn-cell {
    from {
      transform: scale(0%);
      border-radius: 100%;
    }

    to {
      transform: scale(100%);
      border-radius: 0px;
    }
  }

  .cell {
    border-radius: 1px;
    margin: 1px;
    width: 16px;
    height: 16px;
    border: 1px solid transparent;
    display: flex;
    flex-direction: column;
    justify-items: center;
    align-items: center;
    color: white;
  }

  .enabled-cell {
    background-color: white;

    animation-name: sleepy-cell, spawn-cell;
    animation-duration: calc(0.5s + var(--delay)/2), 0.25s;
    animation-delay: 0s, calc(0.5s + var(--delay)/2);
  }

  .hp {
    font-size: 12px;
    font-weight: 600;
    margin: 0;
    padding: 1px 0px 0px 0px;
  }

  .guard {
    border-color: #00ff00;
  }

  .suicide {
    color: #ff8000;
    border-color: #ff8000;
  }

  .move {
    /* TODO: WIP */
  }

  .attack {
    /* TODO: WIP */
  }

  .red {
    background-color: #800000;
  }

  .blue {
    background-color: #000080;
  }

  :hover.enabled-cell {
    border: 1px solid black;
  }
</style>
