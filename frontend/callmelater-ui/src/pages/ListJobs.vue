<template>
  <q-page class="flex ">
    <div class="q-pa-md full-width">
      <q-table
        :columns="columns"
        :rows="rows"
        row-key="name"
        title="Jobs"
      >

        <template v-slot:body-cell-scheduledTimestamp="props">
          <q-td :props="props">
            {{ formatDate(props.row.scheduledTimestamp) }}
          </q-td>
        </template>

        <template v-slot:body-cell-createdTimestamp="props">
          <q-td :props="props">
            {{ formatDate(props.row.createdTimestamp) }}
          </q-td>
        </template>

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
import moment from "moment";

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
  formatDate: Function;
  columns: any[];
  rows: Job[];
}

export default {
  name: "ListJobs",
  setup() {
    let state: State = reactive({
      formatDate: function (value: string) {
        if (value) {
          return moment(String(value)).format('YYYY-MM-DD hh:mm:ss')
        }
      },
      columns: [
        {name: 'Status', align: 'left', label: 'Status', field: 'status', sortable: true},
        {name: 'Id', align: 'left', label: 'Id', field: 'id', sortable: true},
        {name: 'Url', align: 'left', label: 'Url', field: 'url', sortable: true},
        {name: 'scheduledTimestamp', align: 'left', label: 'Scheduled at', field: 'scheduledTimestamp', sortable: true},
        {name: 'createdTimestamp', align: 'left', label: 'Created at', field: 'createdTimestamp', sortable: true},
        {name: 'run', align: 'left', label: 'Run', field: '', sortable: true},
      ],
      rows: [],
    });

    async function getJobs(skip: number, limit: number) {
      const response = await fetch(`http://localhost:8080/jobs/${skip}/${limit}`, {
        method: 'get',
        headers: {
          'content-type': 'application/json'
        }
      });

      const json: any = await response.json();

      console.log(response);

      state.rows = json.jobs
    }

    onMounted(async () => await getJobs(0, 20));

    return state;
  }
}
</script>

<style scoped>

</style>
