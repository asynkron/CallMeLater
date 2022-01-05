<template>
  <router-view/>
</template>
<script>
import {defineComponent} from 'vue';
import moment from "moment";


export default defineComponent({
  name: 'App',
  filters: {
    formatDate: function (value) {
      if (value) {
        return moment(String(value)).format('YYYY-MM-DD hh:mm:ss')
      }
    }
  }
})
</script>
