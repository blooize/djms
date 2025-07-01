<template>
    <v-container class="">
        <v-select
        v-model="selectedEvent"
        label="Select an event"
        :items="eventNames"
        ></v-select>
    </v-container>
</template>

<script lang="ts" setup>
import { computed, ref, watch } from 'vue'

const props = defineProps(['events'])
const emit = defineEmits(['select-event'])
const selectedEvent = ref('')

const eventNames = computed(() => {
  return props.events ? props.events.map((event: any) => event.name) : []
})

watch(selectedEvent, (newValue) => {
  if (newValue) {
    const selectedEventObj = props.events?.find((event: any) => event.name === newValue)
    emit('select-event', selectedEventObj)
  }
})

</script>