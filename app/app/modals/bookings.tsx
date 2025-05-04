import React, { useEffect, useState } from 'react';
import { View, Text, TouchableOpacity, ScrollView, RefreshControl } from 'react-native';
import { Stack, useRouter, useLocalSearchParams } from 'expo-router';
import { StatusBar } from 'expo-status-bar';
import Spinner from 'react-native-loading-spinner-overlay';
import { MaterialIcons } from '@expo/vector-icons';
import { SafeAreaView } from 'react-native-safe-area-context';
import { ENDPOINTS } from '@/constants/Config';
import { showError } from '@/utils/flashMessage';
import { getToken } from '@/utils/auth';

interface Vehicle {
  id: string;
  plate_number: string;
  brand: string;
  model: string;
  owner: string;
}

interface User {
  id: string;
  name: string;
  email: string;
  role: string;
  status: string;
  profile_image_url: string;
}

interface Booking {
  id: string;
  vehicleId: string;
  userId: string;
  locationId: string;
  startTime: string;
  endTime?: string;
  status: string;
  spotNumber: string;
  totalAmount?: number;
  bookingType: string;
  createdAt: string;
  updatedAt: string;
  vehicle: Vehicle;
  user: User;
  location: Location;
}


interface Location {
  id: string;
  name: string;
  address: string;
  slots: Slot[];
}

interface Slot {
    id: string;
    is_occupied: boolean;
    slot_number: string;
}

export default function BookingHistoryScreen() {
  const router = useRouter();
  const { refresh } = useLocalSearchParams<{ refresh?: string }>();
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);

  useEffect(() => {
    if (refresh === 'true') {
      fetchBookings();
      router.setParams({ refresh: undefined });
    }
  }, [refresh]);

  const fetchBookings = async () => {
    try {
      const token = await getToken();
      if (!token) {
        router.replace('/(auth)/login');
        return;
      }

      const response = await fetch(`${ENDPOINTS.BOOKINGS}/user`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      const result = await response.json();
      if (response.ok && result.success) {
        setBookings(result.data);
      }
    } catch (error) {
      showError('Network Error', 'Failed to connect to server');
      console.error('Bookings fetch error:', error);
    } finally {
      setIsLoading(false);
      setRefreshing(false);
    }
  };

  useEffect(() => {
    fetchBookings();
  }, []);

  const onRefresh = () => {
    setRefreshing(true);
    fetchBookings();
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

  const formatDateTime = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  };

  return (
    <SafeAreaView className="flex-1 bg-gray-50">
      <StatusBar style="dark" />
      <Stack.Screen
        options={{
          headerShown: true,
          title: 'Booking History',
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

      <ScrollView
        className="flex-1"
        refreshControl={
          <RefreshControl refreshing={refreshing} onRefresh={onRefresh} />
        }
      >
        {!isLoading && bookings.length === 0 ? (
          <View className="flex-1 items-center justify-center p-6">
            <Text style={{ fontFamily: 'Inter-Medium' }} className="text-gray-500">
              No bookings found
            </Text>
          </View>
        ) : (
          <View className="p-6 space-y-4">
            {bookings.map((booking) => (
              <TouchableOpacity
                key={booking.id}
                // onPress={() => router.push({
                //   pathname: '/modals/booking-details',
                //   params: { id: booking.id }
                // })}
                className="bg-white rounded-xl p-4 shadow-sm border border-gray-100"
              >
                <View className="flex-row justify-between items-start mb-3">
                  <View className="flex-1">
                    <Text style={{ fontFamily: 'Poppins-SemiBold' }} className="text-lg text-gray-800">
                      {booking.vehicle.brand} {booking.vehicle.model}
                    </Text>
                    <Text style={{ fontFamily: 'Inter-Regular' }} className="text-gray-600">
                      {booking.vehicle.plate_number}
                    </Text>
                  </View>
                  <View className={`px-3 py-1 rounded-full ${getStatusColor(booking.status)}`}>
                    <Text style={{ fontFamily: 'Inter-Medium' }} className="text-sm capitalize">
                      {booking.status}
                    </Text>
                  </View>
                </View>

                <View className="space-y-2">
                  <View className="flex-row items-center">
                    <MaterialIcons name="location-on" size={16} color="#6B7280" />
                    <Text style={{ fontFamily: 'Inter-Regular' }} className="text-gray-600 ml-1">
                    {booking.spotNumber}, {booking.location.name}, {booking.location.address}
                    </Text>
                  </View>

                  <View>
                    <Text style={{ fontFamily: 'Inter-Regular' }} className="text-gray-500 text-sm">
                      Start Time
                    </Text>
                    <Text style={{ fontFamily: 'Inter-Medium' }} className="text-gray-700">
                      {formatDateTime(booking.startTime)}
                    </Text>
                  </View>

                  {booking.endTime && (
                    <View>
                      <Text style={{ fontFamily: 'Inter-Regular' }} className="text-gray-500 text-sm">
                        End Time
                      </Text>
                      <Text style={{ fontFamily: 'Inter-Medium' }} className="text-gray-700">
                        {formatDateTime(booking.endTime)}
                      </Text>
                    </View>
                  )}

                  {booking.totalAmount && (
                    <View>
                      <Text style={{ fontFamily: 'Inter-Regular' }} className="text-gray-500 text-sm">
                        Total Amount
                      </Text>
                      <Text style={{ fontFamily: 'Inter-Medium' }} className="text-gray-700">
                        ${booking.totalAmount}
                      </Text>
                    </View>
                  )}

                  <View>
                    <Text style={{ fontFamily: 'Inter-Regular' }} className="text-gray-500 text-sm">
                      Booking Type
                    </Text>
                    <Text style={{ fontFamily: 'Inter-Medium' }} className="text-gray-700 capitalize">
                      {booking.bookingType.replace('_', ' ')}
                    </Text>
                  </View>
                </View>

                {booking.status === 'pending' && (
                  <TouchableOpacity
                    onPress={() => router.push({
                      pathname: '/modals/cancel-booking',
                      params: { id: booking.id }
                    })}
                    className="mt-4 bg-red-50 py-2 rounded-lg"
                  >
                    <Text style={{ fontFamily: 'Inter-Medium' }} className="text-red-600 text-center">
                      Cancel Booking
                    </Text>
                  </TouchableOpacity>
                )}
              </TouchableOpacity>
            ))}
          </View>
        )}
      </ScrollView>
    </SafeAreaView>
  );
}