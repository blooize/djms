<template>
    <div>
        <v-btn
            size="default"
            color="primary"
            class=""
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
                        @input="validateInput"
                        variant="outlined"
                    />
                </v-card-text>
                <template v-slot:actions>
                    <v-btn
                        v-if="validInput"
                        color="primary"
                        append-icon="mdi-check"
                        @click="() => {
                            dialog = false
                            $emit('new-event', eventName)
                        }"
                    >
                        Accept
                    </v-btn>
                    <v-btn
                        v-if="!validInput"
                        color="primary"
                        append-icon="mdi-close"
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

let validInput = ref(false)
let dialog = ref(false)
let eventName = ref<string>()

const validateInput = () => {
    if(eventName.value && eventName.value !== '') {
        validInput.value = true
    } else {
        validInput.value = false
    }
}
</script>