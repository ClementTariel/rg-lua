<script setup lang="ts">
  import Title from '@/components/Title.vue';
  import Game from '@/components/Game.vue';
  import { init, decompress } from '@bokuweb/zstd-wasm';
  import { shallowRef } from 'vue';

  // TODO: WIP type
  const game = shallowRef({ turns: [] });
  const gameMetadata = shallowRef({});

  // TODO: WIP
  fetch("http://localhost:5555/highlighted-match", {
    headers: {
      "Accept": "application/json",
      "Content-Type": "application/json",
    }
  }).then(resp => {
    resp.json().then(matchResp => {
      if (matchResp.match) {
        let { CompressedGame, ...metadata } = matchResp.match;
        gameMetadata.value = { ...metadata };
        let bytes = Uint8Array.from(atob(CompressedGame), c => c.charCodeAt(0));
        init().then(() => {
          let utf8Decode = new TextDecoder();
          let payload = utf8Decode.decode(decompress(bytes));
          game.value = { turns: JSON.parse(payload) };
        })
      }
    });
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
