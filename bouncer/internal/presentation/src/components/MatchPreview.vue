<script setup lang="ts">
  import { computed } from 'vue';
  import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
  import { library } from '@fortawesome/fontawesome-svg-core';

  /* import all the icons in Free Solid, Free Regular, and Brands styles */
  import { fas } from '@fortawesome/free-solid-svg-icons';
  import { far } from '@fortawesome/free-regular-svg-icons';
  import { fab } from '@fortawesome/free-brands-svg-icons';
  library.add(fas, far, fab);

  // TODO: WIP type
  const props = defineProps<{
    preview: any;
    index: number;
  }>();

  function getResultIcon(result: number): string {
    if (result == 0) {
      return "equals";
    }
    return result > 0 ? "check" : "xmark";
  }

  const trueBlueResult = computed(() => {
    return props.preview.Score1 - props.preview.Score2;
  });
</script>

<template>
  <!-- TODO: WIP navigate to full match display -->
  <div class="wrapper">
    <div class="count">{{ props.index }}</div>
    <div class="blueName" :class="{ 'winner': trueBlueResult > 0, 'loser': trueBlueResult < 0 }">
      {{ props.preview.BotName1 }}
    </div>
    <FontAwesomeIcon :icon="['fas', getResultIcon(trueBlueResult)]"
      :class="{ 'winner': trueBlueResult > 0, 'loser': trueBlueResult < 0 }" />
    <div class="score">
      {{ props.preview.Score1 }} - {{ props.preview.Score2 }}
    </div>
    <FontAwesomeIcon :icon="['fas', getResultIcon(-trueBlueResult)]"
      :class="{ 'winner': trueBlueResult < 0, 'loser': trueBlueResult > 0 }" />
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
