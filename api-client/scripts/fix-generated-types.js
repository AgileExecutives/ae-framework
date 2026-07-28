#!/usr/bin/env node

import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

/**
 * Simple glob function to find TypeScript files
 */
function findTsFiles(dir) {
  const results = [];
  const files = fs.readdirSync(dir);
  
  for (const file of files) {
    const filePath = path.join(dir, file);
    const stat = fs.statSync(filePath);
    
    if (stat.isDirectory()) {
      results.push(...findTsFiles(filePath));
    } else if (file.endsWith('.ts')) {
      results.push(filePath);
    }
  }
  
  return results;
}

/**
 * Post-generation script to fix TypeScript issues in generated API client
 */
function fixGeneratedTypes() {
  console.log('🔧 Fixing generated TypeScript files...');
  
  // Find all TypeScript files in src directory
  const srcDir = path.join(__dirname, '..', 'src');
  const files = findTsFiles(srcDir);
  
  files.forEach(filePath => {
    let content = fs.readFileSync(filePath, 'utf8');
    let modified = false;
    
    // Fix 1: Convert Axios imports to type-only imports
    if (content.includes("import axios, { AxiosInstance, AxiosResponse, AxiosRequestConfig }")) {
      content = content.replace(
        /import axios, \{ AxiosInstance, AxiosResponse, AxiosRequestConfig \} from 'axios';/g,
        "import axios from 'axios';\nimport type { AxiosInstance, AxiosResponse, AxiosRequestConfig } from 'axios';"
      );
      modified = true;
      console.log(`✅ Fixed Axios imports in ${filePath}`);
    }
    
    // Fix 2: Handle readonly function for reactive objects
    if (content.includes("function readonly<T>(ref: Ref<T>)")) {
      content = content.replace(
        /function readonly<T>\(ref: Ref<T>\) \{\s*return computed\(\(\) => ref\.value\);\s*\}/g,
        `function readonly<T>(source: Ref<T> | T) {
  // Check if it's a Ref by looking for the 'value' property
  if (source && typeof source === 'object' && 'value' in source) {
    // It's a Ref, return computed with .value
    return computed(() => (source as Ref<T>).value);
  } else {
    // It's a reactive object, return computed without .value
    return computed(() => source as T);
  }
}`
      );
      modified = true;
      console.log(`✅ Fixed readonly function in ${filePath}`);
    }
    
    // Fix 3: Add missing imports if needed
    if (content.includes('computed(') && !content.includes("computed") && !content.includes("from 'vue'")) {
      content = "import { computed } from 'vue';\n" + content;
      modified = true;
      console.log(`✅ Added computed import in ${filePath}`);
    }
    
    // Write back if modified
    if (modified) {
      fs.writeFileSync(filePath, content);
    }
  });

  // Fix 4: Avoid hard-referencing schemas that may not exist in every generated spec
  const indexPath = path.join(srcDir, 'index.ts');
  if (fs.existsSync(indexPath)) {
    let indexContent = fs.readFileSync(indexPath, 'utf8');
    const clientTypePattern = /\/\/ Export Client type from generated types\s*\nexport type Client = components\['schemas'\]\['entities\.ClientResponse'\];/;

    if (clientTypePattern.test(indexContent)) {
      indexContent = indexContent.replace(
        clientTypePattern,
        `// Export Client type from generated types. Some generated specs do not expose
// entities.ClientResponse, so keep this permissive for DTS generation.
export type Client = any;`
      );
      fs.writeFileSync(indexPath, indexContent);
      console.log('✅ Fixed optional Client type export in index.ts');
    }
  }
  
  // Fix 5: Add missing auth and dashboard methods
  const clientPath = path.join(srcDir, 'client.ts');
  if (fs.existsSync(clientPath)) {
    let clientContent = fs.readFileSync(clientPath, 'utf8');
    
    // Add auth methods after clearToken()
    if (!clientContent.includes('async login(')) {
      const authMethodsToAdd = `

  // Auth methods (not auto-generated, manually maintained)
  async login(credentials: { email: string; password: string }) {
    const response = await this.request<ApiResponse<{ token: string; user: any }>>('POST', '/auth/login', credentials);
    return response;
  }

  async register(credentials: { email: string; password: string; name?: string }) {
    const response = await this.request<ApiResponse<{ token: string; user: any }>>('POST', '/auth/register', credentials);
    return response;
  }

  async logout() {
    const response = await this.request<ApiResponse<any>>('POST', '/auth/logout', undefined);
    return response;
  }

  async getCurrentUser() {
    const response = await this.request<ApiResponse<any>>('GET', '/auth/me', undefined);
    return response;
  }

  async changePassword(credentials: { current_password: string; new_password: string }) {
    const response = await this.request<ApiResponse<any>>('POST', '/auth/change-password', credentials);
    return response;
  }

  async forgotPassword(email: string) {
    const response = await this.request<ApiResponse<any>>('POST', '/auth/forgot-password', { email });
    return response;
  }

  async resetPassword(token: string, newPassword: string) {
    const response = await this.request<ApiResponse<any>>('POST', '/auth/reset-password', { token, new_password: newPassword });
    return response;
  }

  async getPlans() {
    const response = await this.request<ApiResponse<any>>('GET', '/plans', undefined);
    return response;
  }

  async getCustomers(params?: Record<string, any>) {
    const response = await this.request<ApiResponse<any>>('GET', '/customers', undefined, params);
    return response;
  }

  async getEmailStats() {
    const response = await this.request<ApiResponse<any>>('GET', '/email-stats', undefined);
    return response;
  }

  async getNewsletterSubscriptions(params?: Record<string, any>) {
    const response = await this.request<ApiResponse<any>>('GET', '/newsletter-subscriptions', undefined, params);
    return response;
  }

  // Calendar methods (not in Swagger spec)
  async getCalendars(params?: Record<string, any>) {
    const response = await this.request<any>('GET', '/calendars', undefined, params);
    return response;
  }

  async createCalendarEntry(data: any) {
    const response = await this.request<ApiResponse<any>>('POST', '/calendar/entries', data);
    return response;
  }

  async updateCalendarEntry(id: number, data: any) {
    const response = await this.request<ApiResponse<any>>('PUT', \`/calendar/entries/\${id}\`, data);
    return response;
  }

  async deleteCalendarEntry(id: number) {
    const response = await this.request<ApiResponse<any>>('DELETE', \`/calendar/entries/\${id}\`, undefined);
    return response;
  }

  async getCalendarEntryById(id: number) {
    const response = await this.request<ApiResponse<any>>('GET', \`/calendar/entries/\${id}\`, undefined);
    return response;
  }

  async createCalendarSeries(data: any) {
    const response = await this.request<ApiResponse<any>>('POST', '/calendar/series', data);
    return response;
  }

  async updateCalendarSeries(id: number, data: any) {
    const response = await this.request<ApiResponse<any>>('PUT', \`/calendar/series/\${id}\`, data);
    return response;
  }

  async deleteCalendarSeries(id: number, options?: any) {
    const response = await this.request<ApiResponse<any>>('DELETE', \`/calendar/series/\${id}\`, options);
    return response;
  }

  async updateCalendar(id: number, data: any) {
    const response = await this.request<ApiResponse<any>>('PUT', \`/calendars/\${id}\`, data);
    return response;
  }

  // Additional booking methods
  async listBookingTemplatesByUser(params?: Record<string, any>) {
    const response = await this.request<ApiResponse<any>>('GET', '/booking/templates/by-user', undefined, params);
    return response;
  }

  async getStaticFile(filename: string) {
    const response = await this.request<any>('GET', \`/static/\${filename}\`, undefined);
    return response;
  }
`;
      
      // Insert after clearToken() method
      const clearTokenPos = clientContent.indexOf('  clearToken() {\n    this.token = null;\n  }');
      if (clearTokenPos !== -1) {
        const insertPos = clearTokenPos + '  clearToken() {\n    this.token = null;\n  }'.length;
        clientContent = clientContent.slice(0, insertPos) + authMethodsToAdd + clientContent.slice(insertPos);
        fs.writeFileSync(clientPath, clientContent);
        console.log('✅ Added missing auth and calendar methods to client.ts');
      }
    }
    
    // Only add booking methods if they don't exist
    if (!clientContent.includes('async getBookingFreeSlots(')) {
      // Find the position after listBookingTemplatesByUser
      const insertAfter = 'async listBookingTemplatesByUser(params?: Record<string, any>) {\n    const response = await this.request<ApiResponse<any>>(\'GET\', `/booking/templates/by-user`, undefined, params);\n    return response;\n  }';
      
      const methodsToAdd = `

  // Booking methods for endpoints without operationIds (temporary until backend adds them)
  async getBookingFreeSlots(token: string, params?: Record<string, any>) {
    if (!token) throw new Error('token is required');
    const response = await this.request<any>('GET', \`/booking/freeslots/\${token}\`, undefined, params);
    return response;
  }

  async createBookingLink(data: any) {
    const response = await this.request<ApiResponse<any>>('POST', \`/booking/link\`, data);
    return response;
  }

  async createBookingTemplate(data: any) {
    const response = await this.request<ApiResponse<any>>('POST', \`/booking/templates\`, data);
    return response;
  }

  async getBookingTemplateById(id: number) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('GET', \`/booking/templates/\${id}\`, undefined);
    return response;
  }

  async updateBookingTemplate(id: number, data: any) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('PUT', \`/booking/templates/\${id}\`, data);
    return response;
  }

  async deleteBookingTemplate(id: number) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('DELETE', \`/booking/templates/\${id}\`, undefined);
    return response || { success: true };
  }`;
      
      if (clientContent.includes(insertAfter)) {
        clientContent = clientContent.replace(insertAfter, insertAfter + methodsToAdd);
        fs.writeFileSync(clientPath, clientContent);
        console.log('✅ Added missing booking methods to client.ts');
      }
    }
  }
  
  // Fix 6: Fix composables auth methods
  const composablesPath = path.join(srcDir, 'composables', 'index.ts');
  if (fs.existsSync(composablesPath)) {
    let composablesContent = fs.readFileSync(composablesPath, 'utf8');
    let modified = false;
    
    // Fix user ref type
    if (composablesContent.includes('const user = ref(null);')) {
      composablesContent = composablesContent.replace(
        'const user = ref(null);',
        'const user = ref<any | null>(null);'
      );
      modified = true;
    }
    
    // Fix login credentials: username -> email
    if (composablesContent.includes('{ username: string; password: string }')) {
      composablesContent = composablesContent.replace(
        '{ username: string; password: string }',
        '{ email: string; password: string }'
      );
      modified = true;
    }
    
    // Fix getCurrentUser to access userData.data
    if (composablesContent.includes('user.value = userData;')) {
      composablesContent = composablesContent.replace(
        /const userData = await client\.getCurrentUser\(\);\s+user\.value = userData;/,
        'const userData = await client.getCurrentUser();\n      user.value = userData.data;'
      );
      modified = true;
    }
    
    if (modified) {
      fs.writeFileSync(composablesPath, composablesContent);
      console.log('✅ Fixed composables auth methods');
    }
  }
  
  console.log('✅ All TypeScript files have been fixed!');
}

// Run the fix
fixGeneratedTypes();