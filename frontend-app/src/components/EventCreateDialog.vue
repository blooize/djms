<template>
  <div>
    <v-btn
      color="primary"
      size="default"
      @click="dialog = true"
    >New Event</v-btn>

    <v-dialog
      v-model="dialog"
      width="500"
    >
      <v-card
        max-width="400"
        title="Enter new Event Name"
      >
        <v-card-text>
          <v-text-field
            v-model="eventName"
            label="Event Name"
            variant="outlined"
            @input="validateInput"
          />
        </v-card-text>
        <template #actions>
          <v-btn
            v-if="validInput"
            append-icon="mdi-check"
            color="primary"
            @click="() => {
              dialog = false
              $emit('new-event', eventName)
            }"
          >
            Accept
          </v-btn>
          <v-btn
            v-if="!validInput"
            append-icon="mdi-close"
            color="primary"
            @click="dialog = false"
          >
            Cancel
          </v-btn>
        </template>
      </v-card>
    </v-dialog>
  </div>

</template>

<script lang="ts" setup>
  import { ref } from 'vue'

  const validInput = ref(false)
  const dialog = ref(false)
  const eventName = ref<string>()

  const validateInput = () => {
    validInput.value = eventName.value && eventName.value.trim() !== '' ? true : false
  }
</script>
