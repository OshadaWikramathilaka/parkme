import React, { useEffect, useState, useCallback } from 'react';
import { View, Text, ScrollView, TouchableOpacity, RefreshControl } from 'react-native';
import { MaterialIcons } from '@expo/vector-icons';
import { useRouter, useFocusEffect } from 'expo-router';
import Spinner from 'react-native-loading-spinner-overlay';
import { SafeAreaView } from 'react-native-safe-area-context';
import { ENDPOINTS } from '@/constants/Config';
import { showError, showSuccess } from '@/utils/flashMessage';
import { getToken, getUser } from '@/utils/auth';
import Header from '@/components/ui/Header';

interface Vehicle {
  id: string;
  plate_number: string;
  brand: string;
  model: string;
  owner: string;
}

interface Booking {
  id: string;
  vehicleId: string;
  locationId: string;
  startTime: string;
  endTime?: string;
  status: 'pending' | 'active' | 'completed' | 'cancelled';
  spotNumber: string;
}

export default function VehiclesScreen() {
  const router = useRouter();
  const [vehicles, setVehicles] = useState<Vehicle[]>([]);
  const [bookings, setBookings] = useState<Record<string, Booking[]>>({});
  const [isLoading, setIsLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);

  const fetchVehicles = async () => {
    try {
      const token = await getToken();
      const User = await getUser();
      if (!token) {
        router.replace('/(auth)/login');
        return;
      }

      const response = await fetch(`${ENDPOINTS.VEHICLES}/user/${User.id}`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      const result = await response.json();
      if (response.ok && result.success) {
        setVehicles(result.data || []);
        // Fetch bookings for each vehicle
        result.data?.forEach((vehicle: Vehicle) => {
          fetchVehicleBookings(vehicle.id);
        });
      } else {
        if (response.status === 401) {
          router.replace('/(auth)/login');
        } else {
          showError('Error', result.message || 'Failed to load vehicles');
        }
      }
    } catch (error) {
      showError('Network Error', 'Failed to connect to server');
      console.error('Vehicles fetch error:', error);
    } finally {
      setIsLoading(false);
      setRefreshing(false);
    }
  };

  const fetchVehicleBookings = async (vehicleId: string) => {
    try {
      const token = await getToken();
      const response = await fetch(`${ENDPOINTS.BOOKINGS}/vehicle/${vehicleId}`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      const result = await response.json();
      if (response.ok && result.success) {
        setBookings(prev => ({
          ...prev,
          [vehicleId]: result.data || [],
        }));
      }
    } catch (error) {
      console.error('Bookings fetch error:', error);
    }
  };

  const handleDelete = async (id: string) => {
    try {
      const token = await getToken();
      if (!token) {
        router.replace('/(auth)/login');
        return;
      }

      const response = await fetch(`${ENDPOINTS.VEHICLES}/${id}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      const result = await response.json();
      if (response.ok && result.success) {
        showSuccess('Success', 'Vehicle deleted successfully');
        setVehicles(prev => prev.filter(vehicle => vehicle.id !== id));
      } else {
        showError('Error', result.message || 'Failed to delete vehicle');
      }
    } catch (error) {
      showError('Network Error', 'Failed to connect to server');
      console.error('Vehicle delete error:', error);
    }
  };

  const getActiveBooking = (vehicleId: string) => {
    return bookings[vehicleId]?.find(b => b.status === 'active' || b.status === 'pending');
  };

  useFocusEffect(
    useCallback(() => {
      fetchVehicles();
    }, [])
  );

  const onRefresh = () => {
    setRefreshing(true);
    fetchVehicles();
  };

  return (
    <SafeAreaView className="flex-1 bg-gray-50">
      <ScrollView
        className="flex-1"
        refreshControl={
          <RefreshControl refreshing={refreshing} onRefresh={onRefresh} />
        }
      >
        <Header title="Vehicles" />
        <Spinner
          visible={isLoading}
          textContent={'Loading...'}
          textStyle={{ color: '#FFF', fontFamily: 'Inter-Medium' }}
          overlayColor="rgba(0, 0, 0, 0.7)"
        />

        <View className="p-6">
          <View className="flex-row items-center justify-between mb-6">
            <Text style={{ fontFamily: 'Poppins-Bold' }} className="text-2xl text-gray-800">
              My Vehicles
            </Text>
            <View className="flex-row space-x-3">
              <TouchableOpacity
                onPress={() => router.push('/modals/create-booking')}
                className="bg-green-600 p-3 rounded-full"
              >
                <MaterialIcons name="event" size={24} color="#FFFFFF" />
              </TouchableOpacity>
              <TouchableOpacity
                onPress={() => router.push('/modals/add-vehicle')}
                className="bg-blue-600 p-3 rounded-full"
              >
                <MaterialIcons name="add" size={24} color="#FFFFFF" />
              </TouchableOpacity>
            </View>
          </View>

          {!isLoading && vehicles.length === 0 && (
            <View className="items-center justify-center py-8">
              <MaterialIcons name="directions-car" size={48} color="#9CA3AF" />
              <Text style={{ fontFamily: 'Inter-Medium' }} className="text-gray-400 mt-4">
                No vehicles added yet
              </Text>
            </View>
          )}

          <View className="space-y-4">
            {vehicles.map(vehicle => {
              const activeBooking = getActiveBooking(vehicle.id);
              return (
                <View
                  key={vehicle.id}
                  className="bg-white rounded-xl p-4 shadow-sm"
                >
                  <View className="flex-row items-center justify-between mb-2">
                    <View className="flex-row items-center">
                      <MaterialIcons name="directions-car" size={24} color="#2563eb" />
                      <Text style={{ fontFamily: 'Poppins-SemiBold' }} className="text-lg text-gray-800 ml-2">
                        {vehicle.plate_number}
                      </Text>
                    </View>
                    <View className="flex-row">
                      <TouchableOpacity
                        onPress={() => router.push({
                          pathname: '/modals/edit-vehicle',
                          params: { id: vehicle.id }
                        })}
                        className="p-2"
                      >
                        <MaterialIcons name="edit" size={20} color="#2563eb" />
                      </TouchableOpacity>
                      <TouchableOpacity
                        onPress={() => handleDelete(vehicle.id)}
                        className="p-2"
                      >
                        <MaterialIcons name="delete" size={20} color="#DC2626" />
                      </TouchableOpacity>
                    </View>
                  </View>

                  <View className="ml-8">
                    <Text style={{ fontFamily: 'Inter-Regular' }} className="text-gray-600">
                      {vehicle.brand} {vehicle.model}
                    </Text>
                    
                    {activeBooking ? (
                      <View className="mt-2 bg-blue-50 p-3 rounded-lg">
                        <Text style={{ fontFamily: 'Inter-Medium' }} className="text-blue-800">
                          {activeBooking.status === 'active' ? 'Currently Parked' : 'Booking Pending'}
                        </Text>
                        <Text style={{ fontFamily: 'Inter-Regular' }} className="text-blue-600 text-sm mt-1">
                          Spot: {activeBooking.spotNumber}
                        </Text>
                        <Text style={{ fontFamily: 'Inter-Regular' }} className="text-blue-600 text-sm">
                          Start: {new Date(activeBooking.startTime).toLocaleString()}
                        </Text>
                      </View>
                    ) : (
                      <TouchableOpacity
                        onPress={() => router.push({
                          pathname: '/modals/create-booking',
                          params: { vehicleId: vehicle.id }
                        })}
                        className="mt-2 bg-green-50 p-3 rounded-lg"
                      >
                        <Text style={{ fontFamily: 'Inter-Medium' }} className="text-green-800">
                          Available for Booking
                        </Text>
                        <Text style={{ fontFamily: 'Inter-Regular' }} className="text-green-600 text-sm">
                          Tap to create a booking
                        </Text>
                      </TouchableOpacity>
                    )}
                  </View>
                </View>
              );
            })}
          </View>
        </View>
      </ScrollView>
    </SafeAreaView>
  );
} 