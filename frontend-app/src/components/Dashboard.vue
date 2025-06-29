<template>
  <div class="pa-4 ">
    <v-btn
        size="large"
        prepend-icon="mdi-account-circle"
        @click="logout"
        >
        Log out
    </v-btn>
  </div>
</template>

<script lang="ts" setup>
import { defineProps, inject } from 'vue'

const props = defineProps([
    'DiscordID',
    'Username',
    'AvatarURL'
])

const logout = () => {
    console.log('Logging out user:', props.Username)
    localStorage.removeItem('jwt')
    window.location.reload()
}
</script>