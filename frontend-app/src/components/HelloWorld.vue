<template>
  <v-container class="fill-height d-flex align-center justify-center" max-width="900">
    <div class="text-center">
      <v-img
        class="mb-4 mx-auto"
        height="200"
        src="@/assets/logo.png"
      />
      <div class="mb-8">
        <div class="text-body-2 font-weight-light mb-n1">Welcome to</div>
        <h1 class="text-h2 font-weight-bold">DJMS</h1>
      </div>
      <v-card class="text-center">
        <v-btn
        :ref="loginButton"
        size="large"
        prepend-icon="mdi-account-circle"
        @click="redirectToDiscord"
        >
        Login
      </v-btn>
      </v-card>

    </div>
  </v-container>
</template>

<script setup lang="ts">
const redirectToDiscord = () => {
  window.location.href = 'http://localhost:3000/api/auth/discord';
}
</script>
