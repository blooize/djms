<template>
  <v-combobox
    v-model="selectedTalent"
    chips
    clearable
    :items="allTalentNames"
    :label="(talentSlot.date ? new Date(talentSlot.date*1000).toLocaleString() : '')"
    multiple
    variant="outlined"
  />
</template>

<script lang="ts" setup>
  import type { TalentSlot } from '@/types'
  import { defineProps, onMounted, ref } from 'vue'
  const props = defineProps<{
    allTalentNames: Array<string>
    talentSlot: TalentSlot
  }>()

  const selectedTalent = ref<string[]>([])

  onMounted(() => {
    console.log(props.talentSlot.date + '\n' + props.allTalentNames)
  })

</script>
