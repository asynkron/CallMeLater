<template>
  <q-page class="flex ">
    <div class="q-pa-md full-width">
      <q-table
        :columns="columns"
        :rows="rows"
        row-key="name"
        title="Jobs"
      >
        <template v-slot:body-cell-run="props">
          <q-td :props="props">
            <q-btn color="blue" dense>Run now</q-btn>
            <q-btn color="grey" dense>Logs</q-btn>
          </q-td>
        </template>
      </q-table>
    </div>
  </q-page>
</template>
<script lang="ts">

/* eslint-disable */
import {onMounted, reactive} from "vue";

interface Job {
  id: string;
  status: string;
  url: string;
  scheduledTimestamp: string;
  createdTimestamp: string;
  dataDiscriminator: string;
  parentJobId: string;
}

interface State {
  columns: any[];
  rows: Job[];
}

export default {
  name: "ListJobs",
  setup() {
    let state: State = reactive({
      columns: [
        {name: 'Status', align: 'left', label: 'Status', field: 'status', sortable: true},
        {name: 'Id', align: 'left', label: 'Id', field: 'id', sortable: true},
        {name: 'Url', align: 'left', label: 'Url', field: 'url', sortable: true},
        {name: 'Scheduled at', align: 'left', label: 'Scheduled at', field: 'scheduledTimestamp', sortable: true},
        {name: 'Created at', align: 'left', label: 'Created at', field: 'createdTimestamp', sortable: true},
        {name: 'run', align: 'left', label: 'Run', field: '', sortable: true},
      ],
      rows: [],
    });

    async function getJobs(skip: number, limit: number) {
      const response: any = await fetch(`http://localhost:8080/jobs/${skip}/${limit}`, {
        method: 'get',
        headers: {
          'content-type': 'application/json'
        }
      });

      state.rows = response.jobs
    }

    onMounted(async () => await getJobs(0, 20));

    return state;
  }
}
</script>

<style scoped>

</style>
