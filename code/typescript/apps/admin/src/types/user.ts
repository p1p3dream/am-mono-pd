export type User = {
  id: number;
  name: string;
  company: string;
  email: string;
  apiKey: string;
  status: 'active' | 'inactive';
};

export type UserResponse = {
  users: User[];
  total: number;
  page: number;
  limit: number;
};
