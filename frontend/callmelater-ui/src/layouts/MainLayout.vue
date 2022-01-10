<template>
  <q-layout class="shadow-2 rounded-borders" view="hHh Lpr lff">
    <q-header>
      <q-toolbar>
        <q-toolbar-title>
          Call Me Later
        </q-toolbar-title>
        <div>blabla</div>
      </q-toolbar>
    </q-header>
    <q-page-container>

      <q-tabs
        class="bg-grey-3"
        align="left"
        indicator-color="primary"
        inline-label>
        <q-route-tab
          exact
          icon="alarm"
          label="Jobs"
          to="/jobs"
        />

      </q-tabs>

      <router-view/>
    </q-page-container>
  </q-layout>
</template>

<script>

import {defineComponent, ref} from 'vue'

const linksList = [
  {
    title: 'Docs',
    caption: 'quasar.dev',
    icon: 'school',
    link: 'https://quasar.dev'
  },
  {
    title: 'Github',
    caption: 'github.com/quasarframework',
    icon: 'code',
    link: 'https://github.com/quasarframework'
  },
  {
    title: 'Discord Chat Channel',
    caption: 'chat.quasar.dev',
    icon: 'chat',
    link: 'https://chat.quasar.dev'
  },
  {
    title: 'Forum',
    caption: 'forum.quasar.dev',
    icon: 'record_voice_over',
    link: 'https://forum.quasar.dev'
  },
  {
    title: 'Twitter',
    caption: '@quasarframework',
    icon: 'rss_feed',
    link: 'https://twitter.quasar.dev'
  },
  {
    title: 'Facebook',
    caption: '@QuasarFramework',
    icon: 'public',
    link: 'https://facebook.quasar.dev'
  },
  {
    title: 'Quasar Awesome',
    caption: 'Community Quasar projects',
    icon: 'favorite',
    link: 'https://awesome.quasar.dev'
  }
];

export default defineComponent({
  name: 'MainLayout',


  setup() {
    const leftDrawerOpen = ref(false)

    return {
      essentialLinks: linksList,
      leftDrawerOpen,
      toggleLeftDrawer() {
        leftDrawerOpen.value = !leftDrawerOpen.value
      }
    }
  }
})
</script>
