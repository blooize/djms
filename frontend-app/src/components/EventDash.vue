<template>
  <v-container>
    <v-row>
      <p>{{ props }}</p>
      <p>{{ dancer_slots }}</p>
      <p>{{ talent_slots }}</p>

    </v-row>
  </v-container>
</template>

<script lang="ts" setup>
import axios from "axios";
import { onMounted, defineProps, watch, ref, onUpdated } from "vue";
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

const talent_slots = ref<TalentSlots>();
const dancer_slots = ref<DancerSlots>();

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

  talent_slots.value = { 
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

  dancer_slots.value = { 
    dancerslots: responseDancerSlots.data.dancer_slots as Array<DancerSlot> || [], 
    event_id: props.EventID as number, 
    club_id: props.ClubID as number 
  } as DancerSlots;

  console.log("Talent Slots:", responseTalentSlots.data);
  console.log("Dancer Slots:", dancer_slots.value);
};

onUpdated(() => {
  fetchSlots();
});
</script>
