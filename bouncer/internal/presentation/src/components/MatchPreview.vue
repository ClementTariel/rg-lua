<script setup lang="ts">
  import { computed } from 'vue';
  // TODO: WIP type
  const props = defineProps<{
    preview: any;
    index: number;
  }>();

  function getResultIcon(result: number): string {
    // TODO: use icons
    if (result == 0) {
      return "="
    }
    return result > 0 ? "V" : "X"
  }

  const trueBlueResult = computed(() => {
    return props.preview.Score1 - props.preview.Score2
  })
</script>

<template>
  <!-- TODO: WIP navigate to full match display -->
  <div class="wrapper">
    <div class="count">{{ props.index }}</div>
    <div class="blueName" :class="{ 'winner': trueBlueResult > 0, 'loser': trueBlueResult < 0 }">
      {{ props.preview.BotName1 }}
    </div>
    <div :class="{ 'winner': trueBlueResult > 0, 'loser': trueBlueResult < 0 }">
      {{ getResultIcon(trueBlueResult) }}
    </div>
    <div class="score">
      {{ props.preview.Score1 }} - {{ props.preview.Score2 }}
    </div>
    <div :class="{ 'winner': trueBlueResult < 0, 'loser': trueBlueResult > 0 }">
      {{ getResultIcon(-trueBlueResult) }}
    </div>
    <div class="redName" :class="{ 'winner': trueBlueResult < 0, 'loser': trueBlueResult > 0 }">
      {{ props.preview.BotName2 }}
    </div>
  </div>
</template>

<style lang="css" scoped>
  .wrapper {
    display: flex;
    flex-direction: row;
    padding: 8px;
    column-gap: 32px;
    width: 100%;
  }

  .count {
    width: 50px
  }

  .blueName,
  .redName {
    /* TODO: WIP*/
    width: 150px;
  }

  .blueName {
    text-align: right;
  }

  .redName {
    text-align: left;
  }

  .score {
    width: 70px;
    font-weight: 700;
    text-align: center;
  }

  .winner {
    color: green;
  }

  .loser {
    color: red;
  }
</style>
