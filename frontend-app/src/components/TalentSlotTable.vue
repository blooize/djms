<template>
  <v-combobox
    v-model="selectedTalent"
    :items="AllTalentNames"
    :label="(Slot.date ? new Date(Slot.date*1000).toLocaleString() : '')"
    clearable
    chips
    multiple
    variant="outlined"
  ></v-combobox>
</template>

<script lang="ts" setup>
import { defineProps, ref } from 'vue'
import type { TalentSlot } from '@/types'
const props = defineProps<{
  AllTalentNames: Array<string>;
  Slot: TalentSlot;
}>();

const selectedTalent = ref<string[]>([]);

</script>