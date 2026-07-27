<script setup lang="ts">
import ResetPasswordForm from "@@/components/ResetPasswordForm.vue"
import AuthLayout from "@@/components/AuthLayout.vue"
</script>

<template>
    <AuthLayout :showLogoutButton="true">
        <ResetPasswordForm />
    </AuthLayout>
</template>