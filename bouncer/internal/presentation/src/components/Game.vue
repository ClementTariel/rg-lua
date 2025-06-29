<script setup lang="ts">
  import { computed, ref, watch } from 'vue';

  // TODO: WIP type
  const props = defineProps<{
    game: { turns: any[] }
  }>();

  const gridRadius: number = 8.5
  const gridSize: number = (gridRadius + 1) * 2;
  const maxTurn = 100;
  const autoRunTurnDuration = 500;
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
      // TODO: type
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


</script>

<template>
  <div class="turn-count-box">
    <p class="turn-count">{{ turn >= 0 ? turn : 0 }}</p>
  </div>
  <div class="grid">
    <div v-for="x in gridSize" class="row">
      <div v-for="y in gridSize" :style="`--delay: ${(squaredDistToCenter(x, y) / ((gridSize / 2) ** 2))}s`"
        class="cell" :class="{
          'enabled-cell': squaredDistToCenter(x, y) < ((gridRadius) ** 2),
          'blue': grid?.[x - 1]?.[y - 1]?.Bot.PlayerId === 1,
          'red': grid?.[x - 1]?.[y - 1]?.Bot.PlayerId === 2,
          'move': grid?.[x - 1]?.[y - 1]?.Action.ActionType === ActionType.Move,
          'attack': grid?.[x - 1]?.[y - 1]?.Action.ActionType === ActionType.Attack,
          'guard': grid?.[x - 1]?.[y - 1]?.Action.ActionType === ActionType.Guard,
          'suicide': grid?.[x - 1]?.[y - 1]?.Action.ActionType === ActionType.Suicide
        }">
        <p class="hp">{{ grid?.[x - 1]?.[y - 1]?.Bot.Hp }}</p>
      </div>
    </div>
  </div>
</template>

<style lang="css" scoped>
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
    animation-name: spawn-grid;
    animation-duration: 0.5s;
    animation-timing-function: ease-out;
    box-shadow: 10px 20px 20px black;
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
    font-weight: bold;
    margin: 0;
    padding: 1px 0px 0px 1px;
  }

  .guard {
    color: #00ff00;
  }

  .suicide {
    color: #ff8000;
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
