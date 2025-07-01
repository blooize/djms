<template>
  <v-container class="p-4">
    <v-row no-gutters align="center">
      <v-col cols="" class="d-flex align-start">
        <div class="d-flex align-center">
          <v-img
            class="mr-4"
            height="80"
            width="80"
            rounded="circle"
            :src="AvatarURL"
          />
          <div>
            <p class="text-h4 font-weight-bold">{{ Username }}</p>
            <v-btn
                size="large"
                prepend-icon="mdi-account-circle"
                @click="logout"
                >
                Log out
            </v-btn>
          </div>
        </div>
      </v-col>

      <v-col cols="3" class="d-flex justify-center">
        <club-selection :clubs="clubs" @select-club="($event) => {
          selectedClub = $event
          fetchEvents()
        }" />
      </v-col>
      <v-col cols="3" v-if="selectedClub.id !== -1" class="d-flex justify-center">
        <event-selection :events="events" @select-event="selectedEvent = $event" />
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts" setup>
import { defineProps, onMounted, ref, watch } from 'vue'
import axios from 'axios'

let selectedClub = ref({id: 0, name: ''})
let selectedEvent = ref({id: 0, name: ''})

let clubs = ref([])
let events = ref([])

const props = defineProps([
    'DiscordID',
    'Username',
    'AvatarURL'
])

const logout = () => {
    console.log('Logging out user:', props.Username)
    localStorage.removeItem('jwt')
    window.location.reload()
}
const fetchClubs = async () => {
    const jwtToken = localStorage.getItem('jwt')
    const response = await axios.get('http://localhost:4000/api/clubs', {
      headers: {
        'Authorization': `Bearer ${jwtToken || ''}`
      },
      params: {
        discord_id: props.DiscordID
      }
    })

    clubs.value = response.data.clubs
}

const fetchEvents = async () => {
  const jwtToken = localStorage.getItem('jwt')

  const response = await axios.get(`http://localhost:4000/api/club/events`, {
    headers: {
      'Authorization': `Bearer ${jwtToken || ''}`
    },
    params: {
      club_id: selectedClub.value.id
    }
  })

  events.value = response.data.events
}

watch(selectedClub, (newValue) => {
  if (newValue.id !== 0) {
    fetchEvents()
  }
})

onMounted(() => {
  fetchClubs()
})

</script>