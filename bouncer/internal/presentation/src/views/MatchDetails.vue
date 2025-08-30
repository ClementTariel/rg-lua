<!-- TODO: WIP -->
<script setup lang="ts">
  import Title from '@/components/Title.vue';
  import Game from '@/components/Game.vue';
  import { useGameStore } from '@/stores/game.ts';
  import { computed, ref, watch } from 'vue';
  import { useRoute } from 'vue-router';

  const gameStore = useGameStore();

  const route = useRoute();

  const gameId = ref("");

  watch(
    () => route.params.id,
    (newId, oldId) => {
      if (Array.isArray(newId)) {
        return;
      }
      if (newId !== oldId || !gameId.value.length) {
        gameId.value = newId;
      }
      if (gameStore.gameMetadata?.Id !== gameId.value) {
        gameStore.resetGame();
        gameStore.loadGame(gameId.value)
      }
    },
    { immediate: true }
  )

  const subtitle = computed(() => `${gameStore.gameMetadata?.BotName1} VS ${gameStore.gameMetadata?.BotName2}`)
</script>

<template>
  <Title title="Matchs" :subtitles="[subtitle]" />
  <Game :game="gameStore.game" :metadata="gameStore.gameMetadata" />

</template>

<style lang="css" scoped></style>
