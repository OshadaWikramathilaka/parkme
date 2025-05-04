// Use your machine's IP address when testing on physical device or emulator
export const API_URL = 'http://192.168.1.4:8080';

export const ENDPOINTS = {
  LOGIN: `${API_URL}/api/auth/login`,
  REGISTER: `${API_URL}/api/auth/register`,
  PROFILE: `${API_URL}/api/auth/profile`,
  VEHICLES: `${API_URL}/api/vehicles`,
  WALLET_BALANCE: `${API_URL}/api/wallet/balance`,
  WALLET_TRANSACTIONS: `${API_URL}/api/wallet/transactions`,
  WALLET_TOPUP: `${API_URL}/api/wallet/topup`,
  LOCATIONS: `${API_URL}/api/locations`,
  BOOKINGS: `${API_URL}/api/bookings`,
} as const;