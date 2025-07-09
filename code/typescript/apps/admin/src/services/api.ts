import axios from 'axios';

export const api = axios.create({
  baseURL: 'https://api.example.com/v1',
  headers: {
    'Content-Type': 'application/json',
  },
});
