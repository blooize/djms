<template>
  <HelloWorld v-if="!isLoggedIn" />

  <Dashboard 
    v-if="isLoggedIn" 
    :DiscordID="discordID" 
    :Username="username" 
    :AvatarURL="avatarURL" 
  />
</template>

<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import axios from 'axios'


let isLoggedIn = ref(false)
let discordID = ref('')
let username = ref('')
let avatarURL = ref('')

axios.defaults.withCredentials = true

const checkUserAuth = async () => {
  try {
    const jwtToken = document.cookie
      .split('; ')
      .find(row => row.startsWith('jwt='))
      ?.split('=')[1];
    
    if (jwtToken) {
      localStorage.setItem('jwt', jwtToken);
      document.cookie = 'jwt=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/'
    }
    const response = await axios.get('http://localhost:4000/me', {
      headers: {
        'Authorization': `Bearer ${jwtToken || ''}`
      },
    })
    console.log('User authenticated:', response.data)
    discordID.value = response.data.discord_id
    username.value = response.data.username
    avatarURL.value = response.data.avatar
    isLoggedIn.value = true
  } catch (error) {
    console.log('User not authenticated:', error)
  }
}

// Make the API call when the component mounts
onMounted(() => {
  checkUserAuth()
})
</script>
