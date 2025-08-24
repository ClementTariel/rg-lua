import { defineStore } from "pinia";
import { ref } from "vue";

import { init, compress, decompress } from '@bokuweb/zstd-wasm';

export const useZstdStore = defineStore('zstd', () => {
  const loading = ref(false);
  const initialized = ref(false);

  function initialize(retryLeft = 0): Promise<void> {
    return new Promise((resolve, reject) => {
      if (initialized.value){
        resolve();
        return;
      }
      if (retryLeft < 0) {
        reject();
        loading.value = false;
        return;
      }
      if (!initialized.value && !loading.value) {
        loading.value = true;
        init().then(() => {
          initialized.value = true;
          resolve();
          loading.value = false;
        }).catch(() => {
          initialize(retryLeft - 1)
            .then(() => resolve())
            .catch(() => reject());
        })
      }
    })
  }

  return { loading, initialized, initialize, compress, decompress }
})