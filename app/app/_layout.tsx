import { DarkTheme, DefaultTheme, ThemeProvider } from '@react-navigation/native';
import { useFonts } from 'expo-font';
import { Stack, useRouter } from 'expo-router';
import * as SplashScreen from 'expo-splash-screen';
import { StatusBar } from 'expo-status-bar';
import { useEffect } from 'react';
import 'react-native-reanimated';
import FlashMessage from 'react-native-flash-message';
import { View } from 'react-native';

import { useColorScheme } from '@/hooks/useColorScheme';
import { getToken } from '@/utils/auth';

// Keep the splash screen visible while we fetch resources
SplashScreen.preventAutoHideAsync();

export default function RootLayout() {
  const router = useRouter();
  const colorScheme = useColorScheme();
  const [loaded] = useFonts({
    'Poppins-Bold': require('../assets/fonts/Poppins-Bold.ttf'),
    'Poppins-SemiBold': require('../assets/fonts/Poppins-SemiBold.ttf'),
    'Poppins-Medium': require('../assets/fonts/Poppins-Medium.ttf'),
    'Inter-Regular': require('../assets/fonts/Inter-Regular.ttf'),
    'Inter-Medium': require('../assets/fonts/Inter-Medium.ttf'),
  });

  useEffect(() => {
    if (loaded) {
      SplashScreen.hideAsync();
      checkToken();
    }
  }, [loaded]);

  const checkToken = async () => {
    const token = await getToken();
    if (token) {
      router.replace('/(tabs)/vehicles');
    } else {
      router.replace('/(auth)/login');
    }
  };

  if (!loaded) {
    return (
      <View style={{ flex: 1, backgroundColor: '#2563eb', justifyContent: 'center', alignItems: 'center' }}>
        <StatusBar style="light" />
      </View>
    );
  }

  return (
    <ThemeProvider value={colorScheme === 'dark' ? DarkTheme : DefaultTheme}>
      <Stack screenOptions={{ headerShown: false }}>
        <Stack.Screen name="(auth)" />
        <Stack.Screen name="(tabs)" />
        <Stack.Screen name="modals" options={{ presentation: 'modal' }} />
      </Stack>
      <StatusBar style="light" />
      <FlashMessage position="top" floating statusBarHeight={50} />
    </ThemeProvider>
  );
}
