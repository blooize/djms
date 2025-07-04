<template>
  <div>
    <v-btn
      color="primary"
      @click="dialog = true"
    >New Club</v-btn>

    <v-dialog
      v-model="dialog"
      width="500"
    >
      <v-card
        max-width="400"
        title="Enter new Club Name"
      >
        <v-card-text>
          <v-text-field
            v-model="clubName"
            label="Club Name"
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
              $emit('new-club', clubName)
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
  const clubName = ref<string>()

  const validateInput = () => {
    validInput.value = !!(clubName.value && clubName.value.trim() !== '')
  }
</script>
