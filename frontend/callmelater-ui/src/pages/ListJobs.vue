<template>
  <q-page class="flex ">
    <div class="q-pa-md full-width">

      <q-toolbar class="q-mb-sm">

        <q-btn color="primary" icon="play_circle_outline">Trigger now</q-btn>
        <q-btn class="q-ml-md" icon="clear">Remove</q-btn>

        <q-input v-model="search" class="bg-white col q-ml-md" dense placeholder="Search" standout="bg-primary">
          <template v-slot:prepend>
            <q-icon v-if="search === ''" name="search"/>
            <q-icon v-else class="cursor-pointer" name="clear" @click="search = ''"/>
          </template>
          <template v-slot:after>
            <q-btn color="primary" icon="search" label="Search" outline square v-on:click="fetch"></q-btn>
          </template>
        </q-input>

      </q-toolbar>
      <q-table
        v-model:selected="selected"
        color="primary"
        row-key="id"
        :columns="columns"
        :rows="rows"
        selection="multiple"
        title="Jobs"

      >
        <template v-slot:body-cell-id="props">
          <q-td key="link" :props="props">
            <div class="row justify-center">
              <div style="width:80px;text-overflow: ellipsis;overflow:hidden">{{ props.row.id }}</div>
            </div>
          </q-td>
        </template>

        <template v-slot:body-cell-jobType="props">
          <q-td :props="props">
            <div>HTTP</div>
          </q-td>
        </template>

        <template v-slot:body-cell-scheduleTimestamp="props">
          <q-td :props="props">
            {{ formatDate(props.row.scheduleTimestamp) }}
          </q-td>
        </template>

        <template v-slot:body-cell-executedTimestamp="props">
          <q-td :props="props">
            <q-chip
              :color="statusColor(props.row.executedStatus)"
              :icon="statusIcon(props.row.executedStatus)"
              square text-color="white">
              {{ formatDate(props.row.executedTimestamp) }}
            </q-chip>
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
  statusIcon: Function;
  fetch: Function;
  columns: any[];
  rows: Job[];
  selected: string[];
  search: string;
  skip: number;
  limit: number;
}

export default {
  name: "ListJobs",
  setup() {
    let state: State = reactive({
      formatDate: formatDate,
      statusColor: statusColor,
      statusIcon: statusIcon,
      columns: columns(),
      rows: [],
      selected: [],
      search: "",
      skip: 0,
      limit: 20,
      fetch: () => getJobs(state.skip, state.limit, state.search)
    });

    async function getJobs(skip: number, limit: number, search: string) {
      const response = await fetch(`http://localhost:8080/jobs?skip=${skip}&limit=${limit}&search=${search}`, {
        method: 'get',
        headers: {
          'content-type': 'application/json'
        }
      });

      const json: any = await response.json();

      console.log(response);

      state.rows = json.jobs
    }

    onMounted(async () => await state.fetch());

    return state;
  }
}

function formatDate(value: string) {
  if (value) {
    const m = moment(String(value));

    //hack null time
    if (m.year() < 2020) {
      return "";
    }

    if (moment.now() - m.valueOf() < 1000 * 60 * 60) {
      return m.fromNow();
    }

    return m.format('YYYY-MM-DD hh:mm:ss')
  }
}

function statusColor(value: number) {
  switch (value) {
    case 0:
      return "";
    case 1:
      return "positive";
    case 2:
      return "negative";
    case 3:
      return "deep-orange";
    default:
      return "Unknown";
  }
}

function statusIcon(value: number) {
  switch (value) {
    case 0:
      return "";
    case 1:
      return "check_circle_outline";
    case 2:
      return "error_outline";
    case 3:
      return "warning_amber";
    default:
      return "Unknown";
  }
}

function columns() {
  return [
    {name: 'jobType', align: 'left', label: 'Job Type', field: 'jobType', sortable: true},
    {name: 'id', align: 'left', label: 'Id', field: 'id', sortable: true},
    {name: 'scheduleCronExpression', align: 'left', label: 'Cron', field: 'scheduleCronExpression', sortable: false},
    {name: 'description', align: 'left', label: 'Description', field: 'description', sortable: false},
    {
      name: 'scheduleTimestamp',
      align: 'left',
      label: 'Next execution',
      field: 'scheduleTimestamp',
      sortable: true,
      sort: 'asc'
    },
    {name: 'executedTimestamp', align: 'left', label: 'Last execution', field: 'executedTimestamp', sortable: true},
    {name: 'executedCount', align: 'right', label: 'Retries', field: 'executedCount', sortable: true},
  ];
}
</script>

<style scoped>

</style>
