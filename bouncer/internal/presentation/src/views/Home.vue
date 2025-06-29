<script setup lang="ts">
  import Title from '@/components/Title.vue';
  import Game from '@/components/Game.vue';
  import { init, decompress } from '@bokuweb/zstd-wasm';
  import { shallowRef } from 'vue';

  const game = shallowRef({ turns: [] });

  defineProps<{}>()

  // TODO: WIP
  fetch("http://localhost:5555/highlighted-match", {
    headers: {
      "Accept": "application/json",
      "Content-Type": "application/json",
    }
  }).then(resp => {
    resp.json().then(matchResp => {
      let compressedGame = matchResp.match.CompressedGame;
      // TODO: WIP check compability with chromium
      let bytes = Uint8Array.from(atob(compressedGame), c => c.charCodeAt(0));
      init().then(() => {
        let utf8Decode = new TextDecoder();
        let payload = utf8Decode.decode(decompress(bytes));
        game.value = { turns: JSON.parse(payload) };
      })
    });
  });

</script>

<template>
  <!-- TODO: WIP -->
  <Title title="Welcome to Robot Game" />
  <div class="welcome">
    <Game class="game" :game="game" />
  </div>
</template>

<style lang="css" scoped>
  .welcome {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
  }

  .game {
    margin: auto
  }
</style>
