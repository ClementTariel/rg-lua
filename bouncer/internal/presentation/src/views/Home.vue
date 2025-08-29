<script setup lang="ts">
  import Title from '@/components/Title.vue';
  import Game from '@/components/Game.vue';
  import { shallowRef } from 'vue';

  import { useZstdStore } from '@/stores/zstd.ts';

  const zstd = useZstdStore();

  // TODO: WIP type
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const game = shallowRef({ turns: [] as any[] });
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const gameMetadata = shallowRef({} as { [key: string]: any });

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
    <div class="banners">
      <div class="banner banner-left blue">
        <span class="user-name">{{ gameMetadata.UserName1 }}</span>
        <span class="bot-name">{{ gameMetadata.BotName1 }}</span>
      </div>
      <div class="banner banner-right red">
        <span class="bot-name">{{ gameMetadata.BotName2 }}</span>
        <span class="user-name">{{ gameMetadata.UserName2 }}</span>
      </div>
    </div>
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

    clip-path: polygon(0% 100%, calc(100% - var(--banner-height)) 100%, 100% 0%, 0% 0%);
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

  .red {
    background-color: #800000;
  }

  .blue {
    background-color: #000080;
  }
</style>
