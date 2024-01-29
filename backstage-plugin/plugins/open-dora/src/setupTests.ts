import '@testing-library/jest-dom';
import 'cross-fetch/polyfill';
import { setupServer } from 'msw/node';
import { handlers } from '../testing/mswHandlers';

export const server = setupServer(...handlers);

beforeAll(() => {
  server.listen();
  Object.defineProperty(window, 'matchMedia', {
    writable: true,
    value: jest.fn().mockImplementation(query => ({
      matches: false,
      media: query,
      onchange: null,
      addListener: jest.fn(), // Deprecated
      removeListener: jest.fn(), // Deprecated
      addEventListener: jest.fn(),
      removeEventListener: jest.fn(),
      dispatchEvent: jest.fn(),
    })),
  });
});
afterEach(() => server.resetHandlers());
afterAll(() => server.close());
