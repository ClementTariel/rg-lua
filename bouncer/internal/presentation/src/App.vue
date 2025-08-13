<script setup lang="ts">
  import { routes } from './router';
  import NavBar from './views/NavBar.vue';
</script>

<template>
  <NavBar :tabs="routes" />
  <RouterView />
</template>

<style>

  html,
  body,
  #app {
    width: 100vw;
    height: 100vh;
    background-color: #f2f2f2;
    padding: 0;
    margin: 0;
  }

  #app {
    display: flex;
    flex-direction: column;
    font-family: roboto;
  }
</style>
