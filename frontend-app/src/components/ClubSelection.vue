<template>
  <v-container>
    <v-slide-group
      show-arrows
    >
      <v-slide-group-item
        v-for="club in clubNames"
        :key="club"
        class="d-flex align-center justify-center"
      >
        <v-btn
          :class="{ 'selected-club': selectedClub === club }"
          :color="selectedClub === club ? 'primary' : '#000000'"
          @click="selectedClub = club"
        >
          {{ club }}
        </v-btn>
      </v-slide-group-item>
    </v-slide-group>

  </v-container>
</template>

<script lang="ts" setup>
  import { computed, ref, watch } from 'vue'

  const props = defineProps(['clubs'])
  const emit = defineEmits(['select-club'])
  const selectedClub = ref()

  // Create computed array with only club names
  const clubNames = computed(() => {
    return props.clubs ? props.clubs.map((club: any) => club.name) : []
  })

  // Watch for changes in selectedClub and emit to parent
  watch(selectedClub, newValue => {
    let selectedClubObj = props.clubs?.find((club: any) => club.name === newValue)
    if (selectedClubObj === undefined) {
      selectedClubObj = { id: -1, name: newValue }
    }
    emit('select-club', selectedClubObj)
  })

</script>
