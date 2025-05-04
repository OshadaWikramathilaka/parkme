import React, { useState } from 'react';
import { View, Text, TouchableOpacity } from 'react-native';
import { Stack, useRouter, useLocalSearchParams } from 'expo-router';
import { StatusBar } from 'expo-status-bar';
import { MaterialIcons } from '@expo/vector-icons';
import Spinner from 'react-native-loading-spinner-overlay';
import { SafeAreaView } from 'react-native-safe-area-context';
import { ENDPOINTS } from '@/constants/Config';
import { showError, showSuccess } from '@/utils/flashMessage';
import { getToken } from '@/utils/auth';

export default function CancelBookingScreen() {
  const router = useRouter();
  const { id } = useLocalSearchParams<{ id: string }>();
  const [isLoading, setIsLoading] = useState(false);

  const handleCancel = async () => {
    setIsLoading(true);
    try {
      const token = await getToken();
      const response = await fetch(`${ENDPOINTS.BOOKINGS}/${id}/cancel`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (response.ok) {
        let result;
        const text = await response.text();
        try {
          result = text ? JSON.parse(text) : { success: true };
        } catch (e) {
          // If response is empty or not JSON, consider it success
          result = { success: true };
        }

        showSuccess('Success', 'Booking cancelled successfully');
        router.back();
        // Refresh the bookings list
        router.setParams({ refresh: 'true' });
      } else {
        const errorText = await response.text();
        let errorMessage;
        try {
          const errorJson = JSON.parse(errorText);
          errorMessage = errorJson.message;
        } catch (e) {
          errorMessage = 'Failed to cancel booking';
        }
        showError('Error', errorMessage);
      }
    } catch (error) {
      showError('Network Error', 'Failed to connect to server');
      console.error('Cancel booking error:', error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <SafeAreaView className="flex-1 bg-gray-50">
      <StatusBar style="dark" />
      <Stack.Screen
        options={{
          headerShown: true,
          title: 'Cancel Booking',
          headerStyle: { backgroundColor: '#fff' },
          headerShadowVisible: false,
          headerLeft: () => (
            <TouchableOpacity
              onPress={() => router.back()}
              className="ml-4"
            >
              <MaterialIcons name="arrow-back" size={24} color="#000" />
            </TouchableOpacity>
          ),
        }}
      />

      <Spinner
        visible={isLoading}
        textContent={'Cancelling booking...'}
        textStyle={{ color: '#FFF', fontFamily: 'Inter-Medium' }}
        overlayColor="rgba(0, 0, 0, 0.7)"
      />

      <View className="flex-1 p-6">
        <View className="bg-white rounded-xl p-6 shadow-sm">
          <View className="items-center mb-6">
            <View className="bg-red-100 p-4 rounded-full mb-4">
              <MaterialIcons name="warning" size={32} color="#DC2626" />
            </View>
            <Text style={{ fontFamily: 'Poppins-SemiBold' }} className="text-xl text-gray-800 mb-2">
              Cancel Booking
            </Text>
            <Text style={{ fontFamily: 'Inter-Regular' }} className="text-gray-600 text-center">
              Are you sure you want to cancel this booking? This action cannot be undone.
            </Text>
          </View>

          <View className="space-y-4">
            <TouchableOpacity
              onPress={handleCancel}
              disabled={isLoading}
              className="bg-red-600 py-4 rounded-xl"
            >
              <Text style={{ fontFamily: 'Inter-Medium' }} className="text-white text-center">
                Yes, Cancel Booking
              </Text>
            </TouchableOpacity>

            <TouchableOpacity
              onPress={() => router.back()}
              disabled={isLoading}
              className="bg-gray-100 py-4 rounded-xl"
            >
              <Text style={{ fontFamily: 'Inter-Medium' }} className="text-gray-800 text-center">
                No, Keep Booking
              </Text>
            </TouchableOpacity>
          </View>
        </View>
      </View>
    </SafeAreaView>
  );
} 