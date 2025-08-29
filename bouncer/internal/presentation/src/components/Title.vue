<script setup lang="ts">
  defineProps<{ title: string, subtitles?: string[] }>()
</script>

<template>
  <div class="title">
    <h1>{{ title }}</h1>
    <h2 v-for="subtitle in subtitles ?? []" :key="subtitle" class="subtitle"> > {{ subtitle }}</h2>
  </div>
</template>

<style lang="css" scoped>
  .title {
    margin: 16px;
    display: flex;
    flex-direction: row;
  }

  .subtitle {
    display: flex;
    align-self: end;
    padding-left: 8px;
    padding-bottom: 5px;
    height: fit-content;
  }

</style>
