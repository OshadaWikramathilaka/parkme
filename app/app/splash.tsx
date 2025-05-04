import { View, Text, Image } from 'react-native';
import { Stack } from 'expo-router';
import { useEffect } from 'react';
import { useRouter } from 'expo-router';

export default function SplashScreen() {
  const router = useRouter();

  useEffect(() => {
    // Automatically navigate to the main screen after 3 seconds
    const timer = setTimeout(() => {
      router.replace('/(tabs)/vehicles');
    }, 3000);

    return () => clearTimeout(timer);
  }, []);

  return (
    <>
      <Stack.Screen options={{ headerShown: false }} />
      <View className="flex-1 bg-blue-600 items-center justify-center">
        {/* Logo Container */}
        <View className="items-center space-y-4">
          <View className="bg-white p-6 rounded-full shadow-lg">
            <Image 
              source={require('../assets/images/logo.png')}
              className="w-24 h-24"
              style={{ width: 96, height: 96 }}
              resizeMode="contain"
            />
          </View>
          
          {/* App Name */}
          <Text className="text-white text-4xl font-bold tracking-wider">
            ParkMe
          </Text>
          
          {/* Tagline */}
          <Text className="text-blue-100 text-lg text-center px-6">
            Find your perfect parking spot
          </Text>
        </View>

        {/* Loading indicator */}
        <View className="absolute bottom-20">
          <Text className="text-blue-200 text-base">
            Loading...
          </Text>
        </View>
      </View>
    </>
  );
} 