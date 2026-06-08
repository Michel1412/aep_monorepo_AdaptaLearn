import { Platform } from 'react-native';
import Constants from 'expo-constants';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { createApiClient } from '@adaptalearn/api-client';

function resolveApiUrl(): string {
  if (process.env.EXPO_PUBLIC_API_URL) {
    return process.env.EXPO_PUBLIC_API_URL;
  }

  // Expo Go / dev: usa IP da máquina para celular físico; localhost no web
  const debuggerHost = Constants.expoConfig?.hostUri?.split(':')[0];
  if (debuggerHost && debuggerHost !== 'localhost' && debuggerHost !== '127.0.0.1') {
    return `http://${debuggerHost}:8080`;
  }

  if (Platform.OS === 'web') {
    return 'http://localhost:8080';
  }

  return 'http://localhost:8080';
}

export const API_URL = resolveApiUrl();

let token: string | null = null;

export async function loadToken() {
  token = await AsyncStorage.getItem('student_token');
  return token;
}

export async function setToken(t: string | null) {
  token = t;
  if (t) await AsyncStorage.setItem('student_token', t);
  else await AsyncStorage.removeItem('student_token');
}

export async function getStoredStudent() {
  const raw = await AsyncStorage.getItem('student');
  return raw ? JSON.parse(raw) : null;
}

export async function setStoredStudent(student: { ra: string; name: string } | null) {
  if (student) await AsyncStorage.setItem('student', JSON.stringify(student));
  else await AsyncStorage.removeItem('student');
}

export const api = createApiClient({
  baseUrl: API_URL,
  getToken: () => token,
});
