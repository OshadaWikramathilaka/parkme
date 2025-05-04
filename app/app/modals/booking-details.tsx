import React, { useEffect, useState } from 'react';
import { View, Text, TouchableOpacity, ScrollView } from 'react-native';
import { Stack, useRouter, useLocalSearchParams } from 'expo-router';
import { StatusBar } from 'expo-status-bar';
import { MaterialIcons } from '@expo/vector-icons';
import Spinner from 'react-native-loading-spinner-overlay';
import { SafeAreaView } from 'react-native-safe-area-context';
import { ENDPOINTS } from '@/constants/Config';
import { showError } from '@/utils/flashMessage';
import { getToken } from '@/utils/auth';

interface BookingDetails {
  id: string;
  vehicleId: string;
  locationId: string;
  startTime: string;
  endTime?: string;
  status: string;
  spotNumber: string;
  totalAmount?: number;
  vehicle: {
    plate_number: string;
    brand: string;
    model: string;
  };
  location: {
    name: string;
    address: string;
  };
}

export default function BookingDetailsScreen() {
  const router = useRouter();
  const { id } = useLocalSearchParams<{ id: string }>();
  const [booking, setBooking] = useState<BookingDetails | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    fetchBookingDetails();
  }, [id]);

  const fetchBookingDetails = async () => {
    try {
      const token = await getToken();
      const response = await fetch(`${ENDPOINTS.BOOKINGS}/${id}`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      const result = await response.json();
      if (response.ok && result.success) {
        setBooking(result.data);
      } else {
        showError('Error', result.message || 'Failed to load booking details');
      }
    } catch (error) {
      showError('Network Error', 'Failed to connect to server');
      console.error('Booking details fetch error:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const formatDateTime = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  };

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'completed':
        return 'bg-green-100 text-green-800';
      case 'pending':
        return 'bg-yellow-100 text-yellow-800';
      case 'cancelled':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  return (
    <SafeAreaView className="flex-1 bg-gray-50">
      <StatusBar style="dark" />
      <Stack.Screen
        options={{
          headerShown: true,
          title: 'Booking Details',
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
        textContent={'Loading...'}
        textStyle={{ color: '#FFF', fontFamily: 'Inter-Medium' }}
        overlayColor="rgba(0, 0, 0, 0.7)"
      />

      {!isLoading && booking && (
        <ScrollView className="flex-1 p-6">
          <View className="bg-white rounded-xl p-6 shadow-sm">
            {/* Status */}
            <View className="flex-row justify-between items-center mb-6">
              <Text style={{ fontFamily: 'Poppins-SemiBold' }} className="text-xl text-gray-800">
                Booking #{booking.id.slice(-6)}
              </Text>
              <View className={`px-3 py-1 rounded-full ${getStatusColor(booking.status)}`}>
                <Text style={{ fontFamily: 'Inter-Medium' }} className="text-sm capitalize">
                  {booking.status}
                </Text>
              </View>
            </View>

            {/* Vehicle Details */}
            <View className="mb-6">
              <Text style={{ fontFamily: 'Inter-Medium' }} className="text-gray-500 mb-2">
                Vehicle
              </Text>
              <View className="flex-row items-center">
                <MaterialIcons name="directions-car" size={24} color="#2563eb" />
                <View className="ml-3">
                  <Text style={{ fontFamily: 'Poppins-SemiBold' }} className="text-gray-800">
                    {booking.vehicle.plate_number}
                  </Text>
                  <Text style={{ fontFamily: 'Inter-Regular' }} className="text-gray-600">
                    {booking.vehicle.brand} {booking.vehicle.model}
                  </Text>
                </View>
              </View>
            </View>

            {/* Location Details */}
            <View className="mb-6">
              <Text style={{ fontFamily: 'Inter-Medium' }} className="text-gray-500 mb-2">
                Location
              </Text>
              <View className="flex-row items-start">
                <MaterialIcons name="location-on" size={24} color="#2563eb" />
                <View className="ml-3">
                  <Text style={{ fontFamily: 'Poppins-SemiBold' }} className="text-gray-800">
                    {booking.location.name}
                  </Text>
                  <Text style={{ fontFamily: 'Inter-Regular' }} className="text-gray-600">
                    {booking.location.address}
                  </Text>
                  <Text style={{ fontFamily: 'Inter-Regular' }} className="text-gray-600">
                    Spot: {booking.spotNumber}
                  </Text>
                </View>
              </View>
            </View>

            {/* Time Details */}
            <View className="mb-6">
              <Text style={{ fontFamily: 'Inter-Medium' }} className="text-gray-500 mb-2">
                Time
              </Text>
              <View className="space-y-2">
                <View>
                  <Text style={{ fontFamily: 'Inter-Regular' }} className="text-gray-600">
                    Start Time
                  </Text>
                  <Text style={{ fontFamily: 'Inter-Medium' }} className="text-gray-800">
                    {formatDateTime(booking.startTime)}
                  </Text>
                </View>
                {booking.endTime && (
                  <View>
                    <Text style={{ fontFamily: 'Inter-Regular' }} className="text-gray-600">
                      End Time
                    </Text>
                    <Text style={{ fontFamily: 'Inter-Medium' }} className="text-gray-800">
                      {formatDateTime(booking.endTime)}
                    </Text>
                  </View>
                )}
              </View>
            </View>

            {/* Amount */}
            {booking.totalAmount && (
              <View>
                <Text style={{ fontFamily: 'Inter-Medium' }} className="text-gray-500 mb-2">
                  Payment
                </Text>
                <Text style={{ fontFamily: 'Poppins-SemiBold' }} className="text-2xl text-gray-800">
                  ${booking.totalAmount.toFixed(2)}
                </Text>
              </View>
            )}

            {/* Cancel Button */}
            {booking.status === 'pending' && (
              <TouchableOpacity
                onPress={() => router.push({
                  pathname: '/modals/cancel-booking',
                  params: { id: booking.id }
                })}
                className="mt-8 bg-red-50 py-3 rounded-xl"
              >
                <Text style={{ fontFamily: 'Inter-Medium' }} className="text-red-600 text-center">
                  Cancel Booking
                </Text>
              </TouchableOpacity>
            )}
          </View>
        </ScrollView>
      )}
    </SafeAreaView>
  );
} 