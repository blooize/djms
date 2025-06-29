<template>
  <HelloWorld />
</template>

<script lang="ts" setup>
import { onMounted } from 'vue'
import axios from 'axios'

// Configure axios to include credentials with all requests
axios.defaults.withCredentials = true

const checkUserAuth = async () => {
  try {
    const response = await axios.get('http://localhost:4000/me', {
      withCredentials: true
    })
    console.log('User authenticated:', response.data)
    // You can store the user data in a store or reactive variable here
  } catch (error) {
    console.log('User not authenticated:', error)
    // Handle unauthenticated user case here
  }
}

// Make the API call when the component mounts
onMounted(() => {
  checkUserAuth()
})
</script>
