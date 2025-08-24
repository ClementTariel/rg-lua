<script setup lang="ts">
  import Title from '@/components/Title.vue';
  import Game from '@/components/Game.vue';
  import { shallowRef } from 'vue';

  import { useZstdStore } from '@/stores/zstd.ts';

  const zstd = useZstdStore();

  // TODO: WIP type
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const game = shallowRef({ turns: [] as any[] });
  const gameMetadata = shallowRef({});

  // TODO: WIP
  fetch("http://localhost:5555/highlighted-match", {
    headers: {
      "Accept": "application/json",
      "Content-Type": "application/json",
    },
  }).then(resp => resp.json())
    .then(matchResp => {
      if (matchResp?.match) {
        const { CompressedGame, ...metadata } = matchResp.match;
        gameMetadata.value = { ...metadata };
        const bytes = Uint8Array.from(atob(CompressedGame), c => c.charCodeAt(0));
        const utf8Decode = new TextDecoder();
        zstd.initialize().then(() => {
          const payload = utf8Decode.decode(zstd.decompress(bytes));
          const turns = JSON.parse(payload);
          game.value = { turns };
        })
      }
    });
</script>

<template>
  <!-- TODO: WIP -->
  <Title title="Welcome to Robot Game LUA" />
  <div class="welcome">
    <Game :game="game" :metadata="gameMetadata" />
  </div>
</template>

<style lang="css" scoped>
  .welcome {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
  }
</style>
