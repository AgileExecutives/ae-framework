<template>
  <div class="search-component">
    <h2>Search</h2>
    
    <div class="search-form">
      <input
        v-model="searchQuery"
        type="text"
        placeholder="Search customers, contacts, etc..."
        @keyup.enter="performSearch"
        :disabled="loading"
      />
      <button 
        @click="performSearch" 
        :disabled="loading || !searchQuery.trim()"
      >
        {{ loading ? 'Searching...' : 'Search' }}
      </button>
      <button 
        @click="performQuickSearch" 
        :disabled="loading || !searchQuery.trim()"
      >
        Quick Search
      </button>
    </div>
    
    <div v-if="error" class="error">
      {{ error }}
    </div>
    
    <div v-if="results.length > 0" class="results">
      <h3>Search Results ({{ results.length }})</h3>
      
      <div class="result-list">
        <div
          v-for="(result, index) in results"
          :key="index"
          class="result-item"
        >
          <h4>{{ result.title || result.name || 'Untitled' }}</h4>
          <p v-if="result.description">{{ result.description }}</p>
          <p v-if="result.email" class="email">{{ result.email }}</p>
          <div class="result-meta">
            <span v-if="result.type" class="type-badge">{{ result.type }}</span>
            <span v-if="result.score" class="score">Score: {{ result.score }}</span>
          </div>
        </div>
      </div>
    </div>
    
    <div v-else-if="query && !loading" class="no-results">
      No results found for "{{ query }}"
    </div>
    
    <div class="search-actions">
      <button 
        v-if="results.length > 0" 
        @click="clearSearch"
        class="clear-btn"
      >
        Clear Results
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useSearch } from '@unburdy/ae-saas-api-client';

// Use the search composable
const {
  results,
  loading,
  error,
  query,
  search,
  quickSearch,
  clearResults
} = useSearch();

// Local search query
const searchQuery = ref('');

// Perform full search
const performSearch = async () => {
  if (!searchQuery.value.trim()) return;
  
  try {
    await search(searchQuery.value, {
      limit: 20,
      include_score: true
    });
  } catch (err) {
    console.error('Search failed:', err);
  }
};

// Perform quick search
const performQuickSearch = async () => {
  if (!searchQuery.value.trim()) return;
  
  try {
    await quickSearch(searchQuery.value, 10);
  } catch (err) {
    console.error('Quick search failed:', err);
  }
};

// Clear search
const clearSearch = () => {
  searchQuery.value = '';
  clearResults();
};
</script>

<style scoped>
.search-component {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;
}

.search-form {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
}

.search-form input {
  flex: 1;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.search-form button {
  padding: 0.75rem 1.5rem;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  white-space: nowrap;
}

.search-form button:disabled {
  background-color: #6c757d;
  cursor: not-allowed;
}

.search-form button:hover:not(:disabled) {
  background-color: #0056b3;
}

.error {
  background-color: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
  border-radius: 4px;
  padding: 1rem;
  margin-bottom: 1rem;
}

.results {
  margin-top: 2rem;
}

.results h3 {
  margin-bottom: 1rem;
  color: #333;
}

.result-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.result-item {
  background: white;
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.result-item h4 {
  margin: 0 0 0.5rem 0;
  color: #333;
}

.result-item p {
  margin: 0.25rem 0;
  color: #666;
}

.email {
  color: #007bff;
  font-family: monospace;
}

.result-meta {
  display: flex;
  gap: 1rem;
  margin-top: 1rem;
  align-items: center;
}

.type-badge {
  background-color: #e9ecef;
  color: #495057;
  padding: 0.25rem 0.5rem;
  border-radius: 12px;
  font-size: 0.8rem;
  font-weight: 500;
}

.score {
  font-size: 0.9rem;
  color: #6c757d;
}

.no-results {
  text-align: center;
  padding: 2rem;
  color: #666;
  background-color: #f8f9fa;
  border-radius: 8px;
}

.search-actions {
  margin-top: 2rem;
  text-align: center;
}

.clear-btn {
  background-color: #6c757d;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
}

.clear-btn:hover {
  background-color: #5a6268;
}
</style>