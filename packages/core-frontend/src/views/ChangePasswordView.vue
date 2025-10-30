<script setup lang="ts">
import ChangePasswordForm from "@@/components/ChangePasswordForm.vue"
import AuthLayout from "@@/components/AuthLayout.vue"
</script>

<template>
    <AuthLayout :showLogoutButton="true">
        <ChangePasswordForm />
    </AuthLayout>
</template>