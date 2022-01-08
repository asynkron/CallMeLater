<template>
  <q-page class="flex ">
    <div class="q-pa-md full-width">
      <q-table
        :columns="columns"
        :rows="rows"
        row-key="name"
        title="Jobs"
        :rows-per-page-options="[0]"
      >
        <template v-slot:body-cell-scheduledTimestamp="props">
          <q-td :props="props">
            {{ formatDate(props.row.scheduledTimestamp) }}
          </q-td>
        </template>

        <template v-slot:body-cell-executedTimestamp="props">
          <q-td :props="props">
            <q-chip :color="statusColor(props.row.executedStatus)" square text-color="white">
              {{ formatDate(props.row.executedTimestamp) }}
            </q-chip>
          </q-td>
        </template>

        <template v-slot:body-cell-expander="props">
          <q-td :props="props">
            <q-btn color="primary" icon="add" outline></q-btn>
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
  statusColor: Function;
  columns: any[];
  rows: Job[];
}

export default {
  name: "ListJobs",
  setup() {
    let state: State = reactive({
      formatDate: formatDate,
      statusColor: statusColor,
      columns: columns(),
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

function formatDate(value: string) {
  if (value) {
    const m = moment(String(value));
    if (m.year() < 2020) {
      return "";
    }
    return m.format('YYYY-MM-DD hh:mm:ss')
  }
}

function statusColor(value: number) {
  switch (value) {
    case 0:
      return "gray";
    case 1:
      return "red";
    case 2:
      return "green";
    default:
      return "Unknown";
  }
}

function columns() {
  return [
    {name: 'expander', align: 'left', label: '', field: '', sortable: false},
    {name: 'cronExpression', align: 'left', label: 'Cron', field: 'cronExpression', sortable: false},
    {name: 'id', align: 'left', label: 'Id', field: 'id', sortable: true},
    {name: 'description', align: 'left', label: 'Description', field: 'description', sortable: false},
    {name: 'scheduledTimestamp', align: 'left', label: 'Next execution', field: 'scheduledTimestamp', sortable: true},
    {name: 'executedTimestamp', align: 'left', label: 'Last execution', field: 'executedTimestamp', sortable: true},
    {name: 'retryCount', align: 'right', label: 'Retries', field: 'retryCount', sortable: true},
  ];
}
</script>

<style scoped>

</style>
