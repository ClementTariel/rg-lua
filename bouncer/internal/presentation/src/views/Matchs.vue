<script setup lang="ts">
  import { ref } from 'vue';
  import Title from '@/components/Title.vue';
  import Pagination from '@/components/Pagination.vue';
  import MatchPreview from '@/components/MatchPreview.vue';

  // TODO: WIP type
  const page = ref({ start: 1, size: 10, last: -1, elements: undefined });

  function fetchPagination(): void {
    // TODO: WIP
    fetch(`http://localhost:5555/matchs?start=${page.value.start}&size=${page.value.size}`, {
      headers: {
        "Accept": "application/json",
        "Content-Type": "application/json",
      }
    }).then(resp => {
      resp.json().then(matchsResp => {
        if (matchsResp.summaries) {
          // TODO: WIP
          page.value = { ...page.value, elements: matchsResp.summaries, last: 1 };
        }
      });
    });
  }

  // TODO: type
  function updatePage(updatedPage: any): void {
    page.value = { ...updatedPage };
    fetchPagination()
  }

  function updatePaginationStart(start: number): void {
    updatePage({ ...page.value, start })
  }

  function updatePaginationSize(size: number): void {
    updatePage({ ...page.value, size })
  }

  fetchPagination()
</script>

<template>
  <Title title="Matchs" />

  <Pagination :page="page" @updatePaginationSize="updatePaginationSize" @updatePaginationStart="updatePaginationStart">
    <template #data="{ index }">
      <MatchPreview :index="page.start + index - 1" :preview="page.elements?.[index - 1]"
        :class="{ 'contrast': !!(index % 2) }" />
    </template>
    <template v-slot:separator>
      <div class="separator"></div>
    </template>
    <template v-slot:empty>
      <div>No match available</div>
    </template>
  </Pagination>

</template>

<style lang="css" scoped>
  .preview {
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .contrast {
    background-color: white;
  }

  .separator {
    width: 100%;
    height: 1px;
    background-color: black;
  }
</style>
