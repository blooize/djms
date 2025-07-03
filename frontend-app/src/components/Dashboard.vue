<template>
  <v-container class="p-4">
    <v-row no-gutters>
      <v-col cols="" class="d-flex">
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

      <v-col cols="5" class="d-flex justify-end">
        <club-selection :clubs="clubs" @select-club="selectedClub = $event"/>
        <club-create-dialog class="d-flex align-center mb-4" @new-club="createClub($event)"/>
      </v-col>
    </v-row>
    <v-row no-gutters class="d-flex justify-end">
      <v-col cols="5" v-if="selectedClub" class="d-flex">
        <event-selection :events="events" @select-event="selectedEvent = $event" />
        <event-create-dialog class="d-flex align-center mb-6" @new-event="createEvent($event)" />
      </v-col>
    </v-row>
    <v-row no-gutters>
      <v-col cols="12">
        <event-dash
          v-if="selectedEvent"
          :DiscordID="props.DiscordID"
          :Username="props.Username"
          :EventID="selectedEvent.id"
          :EventName="selectedEvent.name"
          :ClubID="selectedClub?.id as number"
          />
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts" setup>
import { defineProps, onMounted, ref, watch, computed } from 'vue'
import type { Club, Event } from '@/types'
import axios from 'axios'

let selectedClub = ref<Club | undefined>()
let selectedEvent = ref<Event | undefined>()

let clubs = ref<Array<Club>>([])
let events = ref<Array<Event>>([])

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
      }
    })
    if (response.data.clubs.length === 0) {
      clubs.value = [{ id: 0, name: 'No clubs found' }]
    } else {
      clubs.value = response.data.clubs
    }
}

const fetchEvents = async () => {
  const jwtToken = localStorage.getItem('jwt')
  const response = await axios.get(`http://localhost:4000/api/club/events`, {
    headers: {
      'Authorization': `Bearer ${jwtToken || ''}`
    },
    params: {
      club_id: selectedClub.value?.id
    }
  })

  events.value = response.data.events || []
  console.log('Fetched events:', events.value)
}

const createClub = async (name: string) => {
  const jwtToken = localStorage.getItem('jwt')
  const response = await axios.post('http://localhost:4000/api/club', {
    name: name
  }, {
    headers: {
      'Authorization': `Bearer ${jwtToken || ''}`
    }
  })

  let newClub = { id: response.data.id, name: response.data.name }
  console.log('New club created:', newClub)
  clubs.value.push(newClub)
}

const createEvent = async (name: string) => {
  const jwtToken = localStorage.getItem('jwt')
  const response = await axios.post('http://localhost:4000/api/club/event', {
    name: name,
    club_id: selectedClub.value?.id
  }, {
    headers: {
      'Authorization': `Bearer ${jwtToken || ''}`
    }
  })
  let newEvent = { id: response.data.id, name: response.data.name, club_id: response.data.club_id }
  console.log('New event created:', newEvent)
  events.value.push(newEvent)
}

watch(selectedClub, () => {
  fetchEvents()
})

watch(selectedEvent, (newValue) => {
  console.log('Selected event changed:', newValue)
})

onMounted(() => {
  fetchClubs()
})

</script>