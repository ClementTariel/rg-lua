import { shallowRef } from 'vue';
import { defineStore } from 'pinia';
import { useZstdStore } from './zstd';

export const useGameStore = defineStore('game', () => {
  const zstd = useZstdStore();

  // TODO: WIP type
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const game = shallowRef({ turns: [] as any[] });
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const gameMetadata = shallowRef({ Id: undefined } as {[key: string]: any});

  function resetGame(): void {
    game.value = { turns: [] };
    gameMetadata.value = { Id: undefined }
  }

  function loadGame(id: string): void {
    fetch(`http://localhost:5555/match/${id}`, {
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
  }

  return {
    game,
    gameMetadata,
    resetGame,
    loadGame
  };
});
