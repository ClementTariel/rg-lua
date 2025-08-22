<script setup lang="ts">
  import router from '@/router';
  import type { RouteRecordRaw } from 'vue-router';
  defineProps<{
    tabs: RouteRecordRaw[]
  }>()
</script>

<template>
  <nav>
    <div class="tab">
      <RouterLink v-for="tab in tabs" :key="tab.path" :to="tab.path" class="link" :class="{
        'selected-link': router.currentRoute.value.path === tab.path,
        'active-link': router.currentRoute.value.path !== tab.path
      }">
        {{ tab.name }}
      </RouterLink>
    </div>
  </nav>
</template>

<style lang="css" scoped>
  nav {
    width: 100%;
    background-color: white;
  }

  .tab {
    display: flex;
  }

  .link {
    display: flex;
    padding: 24px;
    width: fit-content;
    text-decoration: none;
    color: black;
  }

  .active-link:hover {
    box-shadow: 2px 2px 4px gray;
    z-index: 1;
  }

  .selected-link {
    padding: 12px;
    margin: 12px;
    padding-bottom: 20px;
    margin-bottom: 2px;
    border-bottom: 2px solid black;
  }
</style>
