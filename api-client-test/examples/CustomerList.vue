<template>
  <div class="customer-list">
    <h2>Customers</h2>
    
    <!-- Loading state -->
    <div v-if="loading" class="loading">
      Loading customers...
    </div>
    
    <!-- Error state -->
    <div v-if="error" class="error">
      {{ error }}
    </div>
    
    <!-- Create new customer form -->
    <div class="create-form">
      <h3>Add New Customer</h3>
      <form @submit.prevent="handleCreate">
        <input
          v-model="newCustomer.name"
          type="text"
          placeholder="Customer name"
          required
        />
        <input
          v-model="newCustomer.email"
          type="email"
          placeholder="Email"
          required
        />
        <button type="submit" :disabled="loading">
          Add Customer
        </button>
      </form>
    </div>
    
    <!-- Customer list -->
    <div v-if="items.length > 0" class="customer-grid">
      <div
        v-for="customer in items"
        :key="customer.id"
        class="customer-card"
      >
        <h4>{{ customer.name }}</h4>
        <p>{{ customer.email }}</p>
        <p class="meta">
          Created: {{ formatDate(customer.created_at) }}
        </p>
        
        <div class="actions">
          <button 
            @click="editCustomer(customer)"
            class="edit-btn"
          >
            Edit
          </button>
          <button 
            @click="handleDelete(customer.id)"
            class="delete-btn"
            :disabled="loading"
          >
            Delete
          </button>
        </div>
      </div>
    </div>
    
    <!-- Empty state -->
    <div v-else-if="!loading" class="empty-state">
      No customers found. Add one above!
    </div>
    
    <!-- Pagination -->
    <div v-if="pagination.totalPages > 1" class="pagination">
      <button 
        @click="prevPage" 
        :disabled="pagination.page <= 1"
      >
        Previous
      </button>
      
      <span class="page-info">
        Page {{ pagination.page }} of {{ pagination.totalPages }}
        ({{ pagination.total }} total)
      </span>
      
      <button 
        @click="nextPage" 
        :disabled="pagination.page >= pagination.totalPages"
      >
        Next
      </button>
    </div>
    
    <!-- Edit modal -->
    <div v-if="editingCustomer" class="modal-overlay" @click="cancelEdit">
      <div class="modal" @click.stop>
        <h3>Edit Customer</h3>
        <form @submit.prevent="handleUpdate">
          <input
            v-model="editingCustomer.name"
            type="text"
            placeholder="Customer name"
            required
          />
          <input
            v-model="editingCustomer.email"
            type="email"
            placeholder="Email"
            required
          />
          <div class="modal-actions">
            <button type="submit" :disabled="loading">
              Update
            </button>
            <button type="button" @click="cancelEdit">
              Cancel
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import { useCustomers } from '@unburdy/ae-saas-api-client';

// Use the customers composable
const {
  items,
  loading,
  error,
  pagination,
  fetchItems,
  createItem,
  updateItem,
  deleteItem,
  nextPage,
  prevPage
} = useCustomers();

// Form data for new customer
const newCustomer = reactive({
  name: '',
  email: ''
});

// Edit customer state
const editingCustomer = ref(null);

// Load customers on mount
onMounted(() => {
  fetchItems();
});

// Handle create
const handleCreate = async () => {
  try {
    await createItem({ ...newCustomer });
    // Reset form
    newCustomer.name = '';
    newCustomer.email = '';
  } catch (err) {
    console.error('Failed to create customer:', err);
  }
};

// Handle edit
const editCustomer = (customer: any) => {
  editingCustomer.value = { ...customer };
};

const cancelEdit = () => {
  editingCustomer.value = null;
};

const handleUpdate = async () => {
  if (!editingCustomer.value) return;
  
  try {
    await updateItem(editingCustomer.value.id, editingCustomer.value);
    editingCustomer.value = null;
  } catch (err) {
    console.error('Failed to update customer:', err);
  }
};

// Handle delete
const handleDelete = async (id: number) => {
  if (!confirm('Are you sure you want to delete this customer?')) return;
  
  try {
    await deleteItem(id);
  } catch (err) {
    console.error('Failed to delete customer:', err);
  }
};

// Utility functions
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString();
};
</script>

<style scoped>
.customer-list {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

.loading, .error, .empty-state {
  text-align: center;
  padding: 2rem;
  margin: 1rem 0;
}

.error {
  background-color: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
  border-radius: 4px;
}

.create-form {
  background-color: #f8f9fa;
  padding: 1.5rem;
  border-radius: 8px;
  margin-bottom: 2rem;
}

.create-form form {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.create-form input {
  flex: 1;
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.customer-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1rem;
}

.customer-card {
  background: white;
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.customer-card h4 {
  margin: 0 0 0.5rem 0;
  color: #333;
}

.customer-card p {
  margin: 0.25rem 0;
  color: #666;
}

.meta {
  font-size: 0.9rem;
  color: #999;
}

.actions {
  margin-top: 1rem;
  display: flex;
  gap: 0.5rem;
}

.edit-btn {
  background-color: #007bff;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
}

.delete-btn {
  background-color: #dc3545;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
}

.delete-btn:disabled {
  background-color: #6c757d;
  cursor: not-allowed;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1rem;
  margin-top: 2rem;
}

.pagination button {
  padding: 0.5rem 1rem;
  border: 1px solid #ddd;
  background: white;
  border-radius: 4px;
  cursor: pointer;
}

.pagination button:disabled {
  background-color: #f8f9fa;
  cursor: not-allowed;
}

.page-info {
  color: #666;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  min-width: 400px;
  max-width: 90vw;
}

.modal form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.modal input {
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.modal-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
}

.modal-actions button {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.modal-actions button[type="submit"] {
  background-color: #007bff;
  color: white;
}

.modal-actions button[type="button"] {
  background-color: #6c757d;
  color: white;
}
</style>