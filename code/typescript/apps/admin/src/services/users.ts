import type { User, UserResponse } from '@/types/user';

// Mock data
const mockUsers: User[] = [
  {
    id: 1,
    name: 'John Doe',
    company: 'Tech Corp',
    email: 'john@techcorp.com',
    apiKey: 'tk_1234567890abcdef',
    status: 'active',
  },
  {
    id: 2,
    name: 'Jane Smith',
    company: 'Digital Solutions',
    email: 'jane@digitalsolutions.com',
    apiKey: 'tk_abcdef1234567890',
    status: 'active',
  },
  {
    id: 3,
    name: 'Bob Wilson',
    company: 'Cloud Systems',
    email: 'bob@cloudsystems.com',
    apiKey: 'tk_9876543210abcdef',
    status: 'inactive',
  },
];

export const userService = {
  list: async (page = 1, limit = 10): Promise<UserResponse> => {
    try {
      await new Promise((resolve) => setTimeout(resolve, 800)); // Simula delay da rede
      return {
        users: mockUsers,
        total: mockUsers.length,
        page,
        limit,
      };
    } catch (error) {
      console.error('Error fetching users:', error);
      throw error;
    }
  },

  create: async (user: Omit<User, 'id' | 'apiKey'>): Promise<User> => {
    try {
      await new Promise((resolve) => setTimeout(resolve, 800));
      return {
        ...user,
        id: Math.floor(Math.random() * 1000),
        apiKey: `tk_${Math.random().toString(36).substring(7)}`,
      };
    } catch (error) {
      console.error('Error creating user:', error);
      throw error;
    }
  },

  update: async (id: number, user: Partial<User>): Promise<User> => {
    try {
      await new Promise((resolve) => setTimeout(resolve, 800));
      return {
        ...mockUsers.find((u) => u.id === id)!,
        ...user,
      };
    } catch (error) {
      console.error('Error updating user:', error);
      throw error;
    }
  },

  delete: async (id: number): Promise<void> => {
    console.log('Deleting user:', id);
    try {
      await new Promise((resolve) => setTimeout(resolve, 800));
    } catch (error) {
      console.error('Error deleting user:', error);
      throw error;
    }
  },
};
