<template>
  <v-container class="p-4">
    <v-row no-gutters>
      <v-col class="d-flex" cols="">
        <div class="d-flex align-center">
          <v-img
            class="mr-4"
            height="80"
            rounded="circle"
            :src="props.avatarURL"
            width="80"
          />
          <div>
            <p class="text-h4 font-weight-bold">{{ props.username }}</p>
            <v-btn
              prepend-icon="mdi-account-circle"
              size="large"
              @click="logout"
            >
              Log out
            </v-btn>
          </div>
        </div>
      </v-col>

      <v-col class="d-flex justify-end" cols="5">
        <club-selection
          :clubs="clubs"
          @select-club="($event) => {
            selectedClub = $event
          }"
        />
        <club-create-dialog class="d-flex align-center mb-4" @new-club="createClub($event)" />
      </v-col>
    </v-row>
    <v-row class="d-flex justify-end" no-gutters>
      <v-col v-if="selectedClub" class="d-flex" cols="5">
        <event-selection :events="events" @select-event="selectedEvent = $event" />
        <event-create-dialog class="d-flex align-center mb-6" @new-event="createEvent($event)" />
      </v-col>
      <v-divider class="mt-3" color="primary" :thickness="3" />
    </v-row>
    <v-row v-if="selectedEvent" no-gutters>
      <v-col cols="5">
        <h1>Talent Slots</h1>
        <event-dash
          class="mt-4"
          :club-i-d="selectedClub.id"
          :discord-i-d="props.discordID"
          :event-i-d="selectedEvent.id"
          :event-name="selectedEvent.name"
          :username="props.username"
        />
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts" setup>
  import type { Club, Event } from '@/types'
  import axios from 'axios'
  import { computed, defineProps, onMounted, ref, watch } from 'vue'

  const selectedClub = ref<Club>()
  const selectedEvent = ref<Event>()

  const clubs = ref<Array<Club>>([])
  const events = ref<Array<Event>>([])

  const props = defineProps<{
    discordID: string
    username: string
    avatarURL: string
  }>()

  const logout = () => {
    console.log('Logging out user:', props.username)
    localStorage.removeItem('jwt')
    window.location.reload()
  }

  const fetchClubs = async () => {
    const jwtToken = localStorage.getItem('jwt')
    const response = await axios.get('http://localhost:4000/api/clubs', {
      headers: {
        Authorization: `Bearer ${jwtToken || ''}`,
      },
    })
    clubs.value = response.data.clubs.length === 0 ? [{ id: 0, name: 'No clubs found' }] : response.data.clubs
  }

  const fetchEvents = async () => {
    const jwtToken = localStorage.getItem('jwt')
    const response = await axios.get(`http://localhost:4000/api/club/events`, {
      headers: {
        Authorization: `Bearer ${jwtToken || ''}`,
      },
      params: {
        club_id: selectedClub.value?.id,
      },
    })

    events.value = response.data.events || []
  }

  const createClub = async (name: string) => {
    const jwtToken = localStorage.getItem('jwt')
    const response = await axios.post('http://localhost:4000/api/club', {
      name: name,
    }, {
      headers: {
        Authorization: `Bearer ${jwtToken || ''}`,
      },
    })

    const newClub = { id: response.data.id, name: response.data.name }
    clubs.value.push(newClub)
  }

  const createEvent = async (name: string) => {
    const jwtToken = localStorage.getItem('jwt')
    const response = await axios.post('http://localhost:4000/api/club/event', {
      name: name,
      club_id: selectedClub.value?.id,
    }, {
      headers: {
        Authorization: `Bearer ${jwtToken || ''}`,
      },
    })
    const newEvent = { id: response.data.id, name: response.data.name, club_id: response.data.club_id }
    events.value.push(newEvent)
  }

  watch(selectedClub, () => {
    fetchEvents()
  })

  onMounted(() => {
    fetchClubs()
  })

</script>
