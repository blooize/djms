<template>
    <div>
        <v-btn
            color="primary"
            class="ma-2"
            icon="mdi-plus"
            @click="dialog = true"
        />

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
                            $emit('new-club', clubName)
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
let clubName = ref<string>()

const validateInput = () => {
    if(clubName.value && clubName.value !== '') {
        validInput.value = true
    } else {
        validInput.value = false
    }
}
</script>