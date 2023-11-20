import '@testing-library/jest-dom';
import 'cross-fetch/polyfill';
import { setupServer } from 'msw/node';
import { handlers } from '../testing/mswHandlers';

export const server = setupServer(...handlers);

beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());
