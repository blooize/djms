<template>
  <HelloWorld v-if="!isLoggedIn" />

  <Dashboard
    v-if="isLoggedIn"
    :avatar-u-r-l="avatarURL"
    :discord-i-d="discordID"
    :username="username"
  />
</template>

<script lang="ts" setup>
  import axios from 'axios'
  import { onMounted, ref } from 'vue'

  const isLoggedIn = ref(false)
  const discordID = ref('')
  const username = ref('')
  const avatarURL = ref('')

  axios.defaults.withCredentials = true

  const checkUserAuth = async () => {
    try {
      const jwtToken = document.cookie
        .split('; ')
        .find(row => row.startsWith('jwt='))
        ?.split('=')[1]

      if (jwtToken) {
        localStorage.setItem('jwt', jwtToken)
        document.cookie = 'jwt=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/'
      }
      const response = await axios.get('http://localhost:4000/me', {
        headers: {
          Authorization: `Bearer ${jwtToken || ''}`,
        },
      })
      discordID.value = response.data.user_id
      username.value = response.data.username
      const url = response.data.avatar
        ? `https://cdn.discordapp.com/avatars/${response.data.user_id}/${response.data.avatar}.png`
        : 'https://cdn.discordapp.com/embed/avatars/0.png'
      avatarURL.value = url
      isLoggedIn.value = true
    } catch (error) {
      console.info('User not authenticated:', error)
    }
  }

  // Make the API call when the component mounts
  onMounted(() => {
    checkUserAuth()
  })
</script>
