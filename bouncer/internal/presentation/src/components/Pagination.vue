<script setup lang="ts">
  // TODO: type
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  defineProps<{ page: any }>();
  defineEmits({
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    updatePaginationSize: (size: number) => true,
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    updatePaginationStart: (start: number) => true,
  });

  const sizes = [10, 20, 50];
</script>

<template>
  <!-- TODO: WIP -->
  <div class="header">

    <div class="sizes">
      <div>rows : </div>
      <button v-for="size in sizes" :key="size" @click="$emit('updatePaginationSize', size)" class="button"
        :class="{ 'current': size === page.size }" v-bind:disabled="size === page.size">{{ size }}</button>
    </div>

    <div v-if="page.last !== -1" class="starts">
      <div v-if="page.start !== 1" class="row">
        <button class="button" @click="$emit('updatePaginationStart', page.start - page.size)">
          previous
        </button>
        <button class="button" @click="$emit('updatePaginationStart', 1)">1</button>
      </div>
      <div v-if="page.start - 3 > 1">...</div>
      <div v-for="offset in 2" :key="offset">
        <button class="button" v-if="page.start - 3 + offset > 1"
          @click="$emit('updatePaginationStart', page.start - 3 + offset)">
          {{ page.start - 3 + offset }}
        </button>
      </div>
      <button class="button current" disabled>{{ page.start }}</button>
      <div v-for="offset in 2" :key="offset" class="row">
        <button class="button" v-if="page.start + offset < page.last"
          @click="$emit('updatePaginationStart', page.start + offset)">
          {{ page.start + offset }}
        </button>
      </div>
      <div v-if="page.last - 3 > page.start">...</div>
      <div v-if="page.start !== page.last" class="row">
        <button class="button" @click="$emit('updatePaginationStart', page.last)">{{ page.last }}</button>
        <button class="button" @click="$emit('updatePaginationStart', page.start + page.size)">
          next
        </button>
      </div>
    </div>
  </div>

  <div v-if="page.elements">
    <slot name="separator"></slot>
    <div v-for="index in page.elements.length ?? 0" :key="index">
      <slot name="data" v-bind="{ index }" />
      <slot name="separator"></slot>
    </div>
    <div v-if="page.elements.length === 0">
      <slot name="empty">
        <div>No element available</div>
      </slot>
      <slot name="separator"></slot>
    </div>
  </div>
</template>

<style lang="css" scoped>
  .button {
    border: none;
    background-color: transparent;
    color: black;
    font-family: inherit;
    font-size: inherit;
  }

  .button.current {
    text-decoration: underline;
  }

  .button:hover {
    text-decoration: underline;

  }

  .header {
    display: flex;
    flex-direction: row;
    padding: 16px 32px;
  }

  .sizes {
    display: flex;
    flex-direction: row;
  }

  .starts {
    display: flex;
    flex-direction: row;
    margin-left: auto;
  }

  .row {
    display: flex;
    flex-direction: row;
  }
</style>
