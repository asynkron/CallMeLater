<template>
  <q-page class="flex ">
    <div class="q-pa-md full-width">
      <q-table
        :columns="columns"
        :rows="rows"
        row-key="name"
        title="Jobs"
      >
        <template v-slot:body-cell-status="props">
          <q-td :props="props">
            <q-chip :color="statusColor(props.row.status)" square text-color="white">
              {{ status(props.row.status) }}
            </q-chip>

          </q-td>
        </template>

        <template v-slot:body-cell-scheduledTimestamp="props">
          <q-td :props="props">
            {{ formatDate(props.row.scheduledTimestamp) }}
          </q-td>
        </template>

        <template v-slot:body-cell-executedTimestamp="props">
          <q-td :props="props">
            {{ formatDate(props.row.executedTimestamp) }}
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
  status: Function;
  statusColor: Function;
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
      status: function (value: number) {
        switch (value) {
          case 0:
            return "Scheduled";
          case 1:
            return "Succeeded";
          case 2:
            return "Cancelled";
          case 3:
            return "Failed";
          default:
            return "Unknown";
        }
      },
      statusColor: function (value: number) {
        switch (value) {
          case 0:
            return "blue";
          case 1:
            return "green";
          case 2:
            return "gray";
          case 3:
            return "red";
          default:
            return "Unknown";
        }
      },
      columns: [
        {name: 'expander', align: 'left', label: '', field: '', sortable: false},
        {name: 'status', align: 'left', label: 'Status', field: 'status', sortable: true},
        {name: 'cronExpression', align: 'left', label: 'Cron', field: 'cronExpression', sortable: false},
        {name: 'id', align: 'left', label: 'Id', field: 'id', sortable: true},
        {name: 'description', align: 'left', label: 'Description', field: 'description', sortable: false},
        {
          name: 'scheduledTimestamp',
          align: 'left',
          label: 'Next execution',
          field: 'scheduledTimestamp',
          sortable: true
        },
        {name: 'executedTimestamp', align: 'left', label: 'Last execution', field: 'executedTimestamp', sortable: true},
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
