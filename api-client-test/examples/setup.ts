// examples/setup.ts
import { createApp } from 'vue';
import App from './App.vue';
import { AESaasApiClientPlugin } from '@unburdy/ae-saas-api-client';

const app = createApp(App);

// Install the API client plugin
app.use(AESaasApiClientPlugin, {
  baseURL: 'http://localhost:8080/api/v1',
  timeout: 10000
});

app.mount('#app');

export default app;