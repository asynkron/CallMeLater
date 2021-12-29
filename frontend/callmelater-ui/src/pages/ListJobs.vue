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
<script>

class Job {

}

export default {
  name: "ListJobs",
  setup() {

    return {
      columns: [
        {name: 'Status', align: 'left', label: 'Status', field: 'status', sortable: true},
        {name: 'Id', align: 'left', label: 'Id', field: 'id', sortable: true},
        {name: 'Url', align: 'left', label: 'Url', field: 'url', sortable: true},
        {name: 'Scheduled at', align: 'left', label: 'Scheduled at', field: 'scheduledTimestamp', sortable: true},
        {name: 'Created at', align: 'left', label: 'Created at', field: 'createdTimestamp', sortable: true},
        {name: 'run', align: 'left', label: 'Run', field: '', sortable: true},
      ],
      rows: [
        {
          status: "Future",
          id: "1232",
          url: "http://www.google.com",
          scheduledTimestamp: "2021-12-12T12:12:12.000Z",
          completed: false
        },
        {
          status: "Retries: 3",
          id: "4232",
          url: "https://localhost:8080/api/v1/treats/1",
          scheduledTimestamp: "2021-12-19T13:12:12.000Z",
          completed: false
        },
        {
          status: "Completed",
          id: "534534",
          url: "http://linkedin.com",
          scheduledTimestamp: "2021-12-31T12:01:01.000Z",
          completed: false
        },
      ],
    };
  }
}
</script>

<style scoped>

</style>
