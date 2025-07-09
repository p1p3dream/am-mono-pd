/// <reference types="jest" />
import '@testing-library/jest-dom';
import { jest } from '@jest/globals';

// Mock localStorage
const storedItems: { [key: string]: string } = {};

const localStorageMock = {
  getItem: jest.fn((key: string) => storedItems[key] || null),
  setItem: jest.fn((key: string, value: string) => {
    storedItems[key] = value;
  }),
  clear: jest.fn(() => {
    Object.keys(storedItems).forEach((key) => delete storedItems[key]);
  }),
  removeItem: jest.fn((key: string) => {
    delete storedItems[key];
  }),
  key: jest.fn((index: number) => Object.keys(storedItems)[index] || null),
  length: 0,
};

Object.defineProperty(window, 'localStorage', {
  value: localStorageMock,
});

// Mock ResizeObserver
Object.defineProperty(window, 'ResizeObserver', {
  value: jest.fn().mockImplementation(() => ({
    observe: jest.fn(),
    unobserve: jest.fn(),
    disconnect: jest.fn(),
  })),
});

// Mock IntersectionObserver
Object.defineProperty(window, 'IntersectionObserver', {
  value: jest.fn().mockImplementation(() => ({
    observe: jest.fn(),
    unobserve: jest.fn(),
    disconnect: jest.fn(),
    root: null,
    rootMargin: '',
    thresholds: [],
    takeRecords: jest.fn(),
  })),
});
