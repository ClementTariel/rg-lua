<script setup lang="ts">
  import { computed, onUnmounted, ref, watch } from 'vue';

  import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
  import { library } from '@fortawesome/fontawesome-svg-core';

  /* import all the icons in Free Solid, Free Regular, and Brands styles */
  import { fas } from '@fortawesome/free-solid-svg-icons';
  import { far } from '@fortawesome/free-regular-svg-icons';
  import { fab } from '@fortawesome/free-brands-svg-icons';
  library.add(fas, far, fab);

  // TODO: WIP type
  const props = defineProps<{
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    game: { turns: any[] },
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
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
  const animationTimeout = setTimeout(() => {
    animationDone = true;
    if (gameLoaded) {
      turn.value = 0;
    }
  }, 1000);
  let turnTimeout: number | undefined;

  const turn = ref(-1);

  watch(() => props.game, async (currentGame) => {
    if (currentGame?.turns?.length) {
      gameLoaded = true;
      if (animationDone) {
        turn.value = 0;
      }
    }
  },
    { immediate: true })

  watch(turn, async (currentTurn) => {
    if (currentTurn < maxTurn) {
      turnTimeout = setTimeout(() => {
        turn.value = currentTurn + 1;
      }, autoRunTurnDuration);
    }
  })

  const grid = computed(() => {
    const computedGrid = Array.from({ length: gridSize }, () => Array(gridSize));
    if (turn.value >= 0 && props.game?.turns && props.game.turns.length && props.game.turns.length > turn.value) {
      // TODO: WIP type
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
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
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  function getBot(grid: any[][], x: number, y: number): any {
    return grid?.[x - 1]?.[y - 1]?.Bot
  }

  // TODO: WIP type
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  function getAction(grid: any[][], x: number, y: number): any {
    return grid?.[x - 1]?.[y - 1]?.Action
  }

  onUnmounted(() => {
    clearTimeout(animationTimeout);
    clearTimeout(turnTimeout);
  })
</script>

<template>
  <div class="game">
    <div class="turn-count-box">
      <p class="turn-count">{{ turn >= 0 ? turn : 0 }}</p>
    </div>
    <div class="grid">
      <div v-for="x in gridSize" :key="x" class="row">
        <div v-for="y in gridSize" :key="y" :style="`--delay: ${(squaredDistToCenter(x, y) / ((gridSize / 2) ** 2))}s`"
          class="cell" :class="{
            'enabled-cell': squaredDistToCenter(x, y) < ((gridRadius) ** 2),
            'blue': getBot(grid, x, y)?.PlayerId === 1,
            'red': getBot(grid, x, y)?.PlayerId === 2,
            'guard': getAction(grid, x, y)?.ActionType === ActionType.Guard,
            'suicide': getAction(grid, x, y)?.ActionType === ActionType.Suicide
          }">
          <div class="bot-move-wrapper" :class="{
            'move': getAction(grid, x, y)?.ActionType === ActionType.Move,
            'attack': getAction(grid, x, y)?.ActionType === ActionType.Attack
          }">
            <FontAwesomeIcon v-if="getAction(grid, x, y)?.Y + 1 == y - 1" class="bot-move-icon top-icon"
              :icon="['fas', 'caret-up']" />
            <FontAwesomeIcon v-if="getAction(grid, x, y)?.X + 1 == x + 1" class="bot-move-icon right-icon"
              :icon="['fas', 'caret-right']" />
            <FontAwesomeIcon v-if="getAction(grid, x, y)?.Y + 1 == y + 1" class="bot-move-icon bottom-icon"
              :icon="['fas', 'caret-down']" />
            <FontAwesomeIcon v-if="getAction(grid, x, y)?.X + 1 == x - 1" class="bot-move-icon left-icon"
              :icon="['fas', 'caret-left']" />
          </div>
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
    z-index: 1;
    position: relative;
    top: 20px;
    width: 48px;
    border: 1px solid black;
    border-radius: 8px;
    padding: 4px;
    background-color: #5a5a5a;
    text-align: center;

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
    position: relative;
  }

  .enabled-cell {
    background-color: white;

    animation-name: sleepy-cell, spawn-cell;
    animation-duration: calc(0.5s + var(--delay)/2), 0.25s;
    animation-delay: 0s, calc(0.5s + var(--delay)/2);
  }

  .bot-move-wrapper {
    position: absolute;
    width: 16px;
    height: 16px;
  }

  .bot-move-icon {
    position: absolute;
    width: 8px;
    height: 8px;
    font-size: 8px;
  }

  .top-icon {
    top: -4px;
    left: 4px;
  }

  .right-icon {
    top: 4px;
    left: 12px;
  }

  .bottom-icon {
    top: 12px;
    left: 4px;
  }

  .left-icon {
    top: 4px;
    left: -4px;
  }

  .hp {
    align-self: center;
    width: 100%;
    height: 18px;
    text-align: center;
    font-size: 12px;
    font-weight: 700;
    margin: 0;
    padding: 0 0 0 0.5px;
  }

  .guard {
    color: #00ff00;
    border-color: #00ff00;
  }

  .suicide {
    color: #ffff00;
    border-color: #ffff00;
  }

  .move {
    color: #00ffff;
  }

  .attack {
    color: #ff8000;
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
