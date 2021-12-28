<template>
  <q-page class="flex flex-center">
    <img
        alt="Quasar logo"
        src="~assets/quasar-logo-vertical.svg"
        style="width: 200px; height: 200px"
    >
  </q-page>
</template>

<script lang="ts">
import {defineComponent} from 'vue';

export default defineComponent({
  name: 'PageIndex',
  setup() {
    console.log("Hello from PageIndex");
  }
})
</script>
