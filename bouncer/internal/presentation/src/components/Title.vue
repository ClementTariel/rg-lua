<script setup lang="ts">
  defineProps<{ title: string }>()
</script>

<template>
  <h1 class="title">{{ title }}</h1>
</template>

<style lang="css" scoped>
  .title {
    margin: 16px;
  }
</style>
