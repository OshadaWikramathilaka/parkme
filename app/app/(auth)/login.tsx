import React, { useState } from 'react';
import { View, Text, TextInput, TouchableOpacity, Image, KeyboardAvoidingView, Platform, ScrollView } from 'react-native';
import { Stack, useRouter } from 'expo-router';
import { StatusBar } from 'expo-status-bar';
import Spinner from 'react-native-loading-spinner-overlay';
import { ENDPOINTS } from '@/constants/Config';
import { showError, showSuccess } from '@/utils/flashMessage';
import { getUser, storeToken } from '@/utils/auth';

const shadowStyle = Platform.select({
  ios: {
    shadowColor: '#000',
    shadowOffset: {
      width: 0,
      height: 2,
    },
    shadowOpacity: 0.25,
    shadowRadius: 3.84,
  },
  android: {
    elevation: 5,
  },
  web: {
    boxShadow: '0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1)',
  },
});

const buttonShadowStyle = Platform.select({
  ios: {
    shadowColor: '#000',
    shadowOffset: {
      width: 0,
      height: 1,
    },
    shadowOpacity: 0.2,
    shadowRadius: 1.41,
  },
  android: {
    elevation: 2,
  },
  web: {
    boxShadow: '0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1)',
  },
});

export default function LoginScreen() {
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleLogin = async () => {
    if (!email || !password) {
      showError('Error', 'Please fill in all fields');
      return;
    }

    setIsLoading(true);
    console.log(email, password);
    try {
      const response = await fetch(ENDPOINTS.LOGIN, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email,
          password,
        }),
      });

      const result = await response.json();
      console.log(result);
      if (response.ok ) {
        await storeToken(result.token);
        showSuccess('Welcome back!');
        
        router.replace('/(tabs)/home');
      } else {
        showError('Login Failed', result.message || 'Please check your credentials');
      }
    } catch (error) {
      showError('Network Error', 'Please check your connection and try again');
      console.error('Login error:', error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <KeyboardAvoidingView
      behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
      className="flex-1 bg-white"
    >
      <StatusBar style="dark" />
      <Stack.Screen options={{ headerShown: false }} />
      
      <Spinner
        visible={isLoading}
        textContent={'Signing in...'}
        textStyle={{ color: '#FFF', fontFamily: 'Inter-Medium' }}
        overlayColor="rgba(0, 0, 0, 0.7)"
      />
      
      <ScrollView contentContainerStyle={{ flexGrow: 1 }}>
        <View className="flex-1 justify-between p-8">
          {/* Top Section */}
          <View className="items-center mt-20">
            <View 
              className="bg-blue-600 p-6 rounded-full mb-6"
              style={shadowStyle}
            >
              <Image 
                source={require('../../assets/images/splash-icon.png')}
                style={{ width: 80, height: 80 }}
                resizeMode="contain"
              />
            </View>
            <Text style={{ fontFamily: 'Poppins-Bold' }} className="text-3xl text-gray-800 mb-2">
              Welcome Back
            </Text>
            <Text style={{ fontFamily: 'Inter-Regular' }} className="text-gray-500 text-center mb-8">
              Sign in to continue using ParkMe
            </Text>
          </View>

          {/* Form Section */}
          <View className="space-y-6">
            <View>
              <Text style={{ fontFamily: 'Inter-Medium' }} className="text-gray-700 mb-2 ml-1">
                Email
              </Text>
              <TextInput
                className="bg-gray-50 rounded-xl p-4 text-gray-800 border border-gray-200"
                placeholder="Enter your email"
                value={email}
                onChangeText={setEmail}
                keyboardType="email-address"
                autoCapitalize="none"
                editable={!isLoading}
              />
            </View>

            <View>
              <Text style={{ fontFamily: 'Inter-Medium' }} className="text-gray-700 mb-2 ml-1">
                Password
              </Text>
              <TextInput
                className="bg-gray-50 rounded-xl p-4 text-gray-800 border border-gray-200"
                placeholder="Enter your password"
                value={password}
                onChangeText={setPassword}
                secureTextEntry
                editable={!isLoading}
              />
            </View>

            <TouchableOpacity
              onPress={handleLogin}
              disabled={isLoading}
              className={`rounded-xl p-4 ${isLoading ? 'bg-blue-400' : 'bg-blue-600'}`}
              style={buttonShadowStyle}
            >
              <Text style={{ fontFamily: 'Inter-Medium' }} className="text-white text-center text-lg">
                Sign In
              </Text>
            </TouchableOpacity>

            <TouchableOpacity
              onPress={() => router.push('/(auth)/signup')}
              disabled={isLoading}
              className="mt-4"
            >
              <Text style={{ fontFamily: 'Inter-Regular' }} className="text-center text-gray-600">
                Don't have an account? <Text className="text-blue-600 font-semibold">Sign Up</Text>
              </Text>
            </TouchableOpacity>
          </View>

          {/* Bottom Spacing */}
          <View className="h-8" />
        </View>
      </ScrollView>
    </KeyboardAvoidingView>
  );
} 