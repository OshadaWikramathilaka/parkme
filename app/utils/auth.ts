import { ENDPOINTS } from '@/constants/Config';
import AsyncStorage from '@react-native-async-storage/async-storage';

const TOKEN_KEY = '@ParkMe:token';

export const storeToken = async (token: string) => {
  try {
    await AsyncStorage.setItem(TOKEN_KEY, token);
  } catch (error) {
    console.error('Error storing token:', error);
  }
};

export const getToken = async () => {
  try {
    const token = await AsyncStorage.getItem(TOKEN_KEY);
    return token;
  } catch (error) {
    console.error('Error getting token:', error);
    return null;
  }
};

export const getUser = async () => {
  try {
    const token = await getToken();
    const response = await fetch(ENDPOINTS.PROFILE, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    });
    const result = await response.json();
    return result.data;
  } catch (error) {
    console.error('Error getting user:', error);
    return null;
  }
};

export const removeToken = async () => {
  try {
    await AsyncStorage.removeItem(TOKEN_KEY);
  } catch (error) {
    console.error('Error removing token:', error);
  }
}; 