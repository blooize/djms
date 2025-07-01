<template>
    <v-container class="">
        <v-select
        v-model="selectedClub"
        label="Select a club"
        :items="clubNames"
        ></v-select>
    </v-container>
</template>

<script lang="ts" setup>
import { computed, ref, watch } from 'vue'

const props = defineProps(['clubs'])
const emit = defineEmits(['select-club'])
const selectedClub = ref('')

// Create computed array with only club names
const clubNames = computed(() => {
  return props.clubs ? props.clubs.map((club: any) => club.name) : []
})

// Watch for changes in selectedClub and emit to parent
watch(selectedClub, (newValue) => {
  if (newValue) {
    // Find the full club object to send back
    const selectedClubObj = props.clubs?.find((club: any) => club.name === newValue)
    emit('select-club', selectedClubObj)
  }
})

</script>