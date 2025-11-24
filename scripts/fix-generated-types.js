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
  
  // Fix 4: Add missing booking-related methods to client.ts
  const clientPath = path.join(srcDir, 'client.ts');
  if (fs.existsSync(clientPath)) {
    let clientContent = fs.readFileSync(clientPath, 'utf8');
    
    // Check if methods are already added
    if (!clientContent.includes('async getBookingFreeSlots(')) {
      // Find the position after listBookingTemplatesByUser
      const insertAfter = 'async listBookingTemplatesByUser(params?: Record<string, any>) {\n    const response = await this.request<ApiResponse<any>>(\'GET\', `/booking/templates/by-user`, undefined, params);\n    return response;\n  }';
      
      const methodsToAdd = `

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
  
  console.log('✅ All TypeScript files have been fixed!');
}

// Run the fix
fixGeneratedTypes();