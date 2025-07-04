<template>
  <v-card v-for="slot in talent_slots.talentslots" :key="slot.date">
    <TalentSlotTable
      class="mt-1"
      :AllTalentNames="allTalentNames"
      :Slot="slot"
    />
  </v-card>
</template>

<script lang="ts" setup>
import axios from "axios";
import { onMounted, defineProps, watch, ref, onUpdated, computed } from "vue";
import type {
  Event,
  TalentSlots,
  TalentSlot,
  Talent,
  DancerSlots,
  DancerSlot,
  Dancer,
} from "@/types";


const props = defineProps<{
  DiscordID: string;
  Username: string;
  EventID: number;
  EventName: string;
  ClubID: number;
}>();

const tslots = ref<TalentSlots>({
  talentslots: [],
  event_id: props.EventID,
  club_id: props.ClubID,
})

const dslots = ref<DancerSlots>({
  dancerslots: [],
  event_id: props.EventID,
  club_id: props.ClubID,
});

const talent_slots = computed(() => tslots.value);
const dancer_slots = computed(() => dslots.value);

const allTalentNames = ref<Array<string>>([]);

const fetchAllTalentNames = async () => {
  const jwtToken = localStorage.getItem("jwt");
  const response = await axios.get(
    `http://localhost:4000/api/talents`,
    {
      headers: {
        Authorization: `Bearer ${jwtToken || ""}`,
      },
    },
  );
  allTalentNames.value = response.data.talents.map((talent: any) => talent.name) as Array<string>;
}

const fetchSlots = async () => {
  const jwtToken = localStorage.getItem("jwt");
  const responseTalentSlots = await axios.get(
    `http://localhost:4000/api/event/talentslots`,
    {
      headers: {
        Authorization: `Bearer ${jwtToken || ""}`,
      },
      params: {
        event_id: props.EventID,
      },
    },
  );

  tslots.value = { 
    talentslots: responseTalentSlots.data.talent_slots as Array<TalentSlot> || [], 
    event_id: props.EventID as number, 
    club_id: props.ClubID as number 
  } as TalentSlots;

  const responseDancerSlots = await axios.get(
    `http://localhost:4000/api/event/dancerslots`,
    {
      headers: {
        Authorization: `Bearer ${jwtToken || ""}`,
      },
      params: {
        event_id: props.EventID,
      },
    },
  );

  dslots.value = { 
    dancerslots: responseDancerSlots.data.dancer_slots as Array<DancerSlot> || [], 
    event_id: props.EventID as number, 
    club_id: props.ClubID as number 
  } as DancerSlots;

};

onMounted(() => {
  fetchAllTalentNames();
  fetchSlots();
});
</script>
