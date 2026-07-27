<script setup lang="ts">
import { Icon } from '@iconify/vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const router = useRouter()
const LogOut = () => {
    authStore.logout()
    router.push('/login')
}
</script>

<template>
    <button class="btn btn-ghost btn-circle" @click="LogOut" data-testid="logout-button">
        <Icon icon="radix-icons:exit" class="h-[1.2rem] w-[1.2rem]" />
    </button>
</template>