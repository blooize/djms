<template>
  <v-card v-for="talent_slot in talent_slots.talentslots" :key="talent_slot.date">
    <TalentSlotTable
      :all-talent-names="allTalentNames"
      class="mt-1"
      :talent-slot="talent_slot"
    />
  </v-card>
</template>

<script lang="ts" setup>
  import type {
    Dancer,
    DancerSlot,
    DancerSlots,
    Event,
    Talent,
    TalentSlot,
    TalentSlots,
  } from '@/types'
  import axios from 'axios'
  import { computed, defineProps, onMounted, onUpdated, ref, watch } from 'vue'

  const props = defineProps<{
    discordID: string
    username: string
    eventID: number
    eventName: string
    clubID: number
  }>()

  const tslots = ref<TalentSlots>({
    talentslots: [],
    event_id: props.eventID,
    club_id: props.clubID,
  })

  const dslots = ref<DancerSlots>({
    dancerslots: [],
    event_id: props.eventID,
    club_id: props.clubID,
  })

  const talent_slots = computed(() => tslots.value)
  const dancer_slots = computed(() => dslots.value)

  const allTalentNames = ref<Array<string>>([])

  const fetchAllTalentNames = async () => {
    const jwtToken = localStorage.getItem('jwt')
    const response = await axios.get(
      `http://localhost:4000/api/talents`,
      {
        headers: {
          Authorization: `Bearer ${jwtToken || ''}`,
        },
      },
    )
    allTalentNames.value = response.data.talents.map((talent: any) => talent.name) as Array<string>
  }

  const fetchSlots = async () => {
    const jwtToken = localStorage.getItem('jwt')
    const responseTalentSlots = await axios.get(
      `http://localhost:4000/api/event/talentslots`,
      {
        headers: {
          Authorization: `Bearer ${jwtToken || ''}`,
        },
        params: {
          event_id: props.eventID,
        },
      },
    )

    tslots.value = {
      talentslots: responseTalentSlots.data.talent_slots as Array<TalentSlot> || [],
      event_id: props.eventID as number,
      club_id: props.clubID as number,
    } as TalentSlots

    const responseDancerSlots = await axios.get(
      `http://localhost:4000/api/event/dancerslots`,
      {
        headers: {
          Authorization: `Bearer ${jwtToken || ''}`,
        },
        params: {
          event_id: props.eventID,
        },
      },
    )

    dslots.value = {
      dancerslots: responseDancerSlots.data.dancer_slots as Array<DancerSlot> || [],
      event_id: props.eventID as number,
      club_id: props.clubID as number,
    } as DancerSlots
  }

  onMounted(() => {
    fetchAllTalentNames()
    fetchSlots()
  })
</script>
