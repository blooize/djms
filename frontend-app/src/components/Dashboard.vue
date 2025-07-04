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
        <club-selection :clubs="clubs" @select-club="($event) => {
          selectedClub = $event
        }"/>
        <club-create-dialog class="d-flex align-center mb-4" @new-club="createClub($event)"/>
      </v-col>
    </v-row>
    <v-row no-gutters class="d-flex justify-end">
      <v-col cols="5" v-if="selectedClub" class="d-flex">
        <event-selection :events="events" @select-event="selectedEvent = $event" />
        <event-create-dialog class="d-flex align-center mb-6" @new-event="createEvent($event)" />
      </v-col>
    <v-divider class="mt-3" color="primary" :thickness="3" />
    </v-row>
    <v-row no-gutters v-if="selectedEvent && selectedClub">
      <v-col cols="5">
        <h1>Talent Slots</h1>
        <event-dash
          :DiscordID="props.DiscordID"
          :Username="props.Username"
          :EventID="selectedEvent.id"
          :EventName="selectedEvent.name"
          :ClubID="selectedClub.id"
          class="mt-4"
          />
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts" setup>
import { defineProps, onMounted, ref, watch, computed } from 'vue'
import type { Club, Event } from '@/types'
import axios from 'axios'

let selectedClub = ref<Club>()
let selectedEvent = ref<Event>()

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
  events.value.push(newEvent)
}

watch(selectedClub, () => {
  fetchEvents()
})

onMounted(() => {
  fetchClubs()
})

</script>