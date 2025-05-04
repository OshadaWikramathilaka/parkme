import React, { useState, useEffect } from 'react';
import { View, Text, TouchableOpacity, Platform, ScrollView, Modal } from 'react-native';
import { Stack, useRouter, useLocalSearchParams } from 'expo-router';
import { StatusBar } from 'expo-status-bar';
import Spinner from 'react-native-loading-spinner-overlay';
import { MaterialIcons } from '@expo/vector-icons';
import { ENDPOINTS } from '@/constants/Config';
import { showError, showSuccess } from '@/utils/flashMessage';
import { getToken } from '@/utils/auth';

interface Vehicle {
  id: string;
  brand: string;
  model: string;
  plate_number: string;
}

interface Location {
  id: string;
  name: string;
  address: string;
  slots: Array<{
    number: string;
    is_occupied: boolean;
    type: string;
  }>;
}

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

export default function CreateBookingModal() {
  const router = useRouter();
  const params = useLocalSearchParams<{
    vehicleId?: string;
    selectedVehicle?: string;
    selectedLocation?: string;
  }>();

  const [isLoading, setIsLoading] = useState(false);
  const [selectedVehicle, setSelectedVehicle] = useState<Vehicle | null>(null);
  const [selectedLocation, setSelectedLocation] = useState<Location | null>(null);
  const [dateTime, setDateTime] = useState(new Date());
  const [showDateModal, setShowDateModal] = useState(false);
  const [showTimeModal, setShowTimeModal] = useState(false);

  useEffect(() => {
    if (params.vehicleId && !selectedVehicle) {
      fetchVehicle(params.vehicleId);
    }
  }, [params.vehicleId]);

  useEffect(() => {
    if (params.selectedVehicle) {
      try {
        const vehicle = JSON.parse(params.selectedVehicle);
        setSelectedVehicle(vehicle);
      } catch (e) {
        console.error('Error parsing vehicle:', e);
      }
    }
  }, [params.selectedVehicle]);

  useEffect(() => {
    if (params.selectedLocation) {
      try {
        const location = JSON.parse(params.selectedLocation);
        setSelectedLocation(location);
      } catch (e) {
        console.error('Error parsing location:', e);
      }
    }
  }, [params.selectedLocation]);

  const fetchVehicle = async (id: string) => {
    try {
      const token = await getToken();
      const response = await fetch(`${ENDPOINTS.VEHICLES}/${id}`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });
      const data = await response.json();
      if (response.ok && data.success) {
        setSelectedVehicle(data.data);
      }
    } catch (error) {
      console.error('Error fetching vehicle:', error);
    }
  };

  const generateTimeSlots = () => {
    const slots = [];
    for (let hour = 0; hour < 24; hour++) {
      for (let minute of [0, 30]) {
        slots.push({
          hour,
          minute,
          label: `${hour.toString().padStart(2, '0')}:${minute.toString().padStart(2, '0')}`
        });
      }
    }
    return slots;
  };

  const handleSelectTime = (hour: number, minute: number) => {
    const newDateTime = new Date(dateTime);
    newDateTime.setHours(hour);
    newDateTime.setMinutes(minute);
    setDateTime(newDateTime);
    setShowTimeModal(false);
  };

  const handleSelectDate = (date: Date) => {
    const newDateTime = new Date(dateTime);
    newDateTime.setFullYear(date.getFullYear());
    newDateTime.setMonth(date.getMonth());
    newDateTime.setDate(date.getDate());
    setDateTime(newDateTime);
    setShowDateModal(false);
  };

  const handleCreateBooking = async () => {
    if (!selectedVehicle || !selectedLocation) {
      showError('Error', 'Please fill in all fields');
      return;
    }

    setIsLoading(true);
    const spotNumber = selectedLocation.slots.find(slot => slot.is_occupied === false)?.number;
    console.log(spotNumber);
    try {
      const token = await getToken();
      const response = await fetch(ENDPOINTS.BOOKINGS, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          vehicleId: selectedVehicle.id,
          locationId: selectedLocation.id,
          startTime: dateTime.toISOString(),
          bookingType: 'pre_booked',
          spotNumber: `${spotNumber}`,
        }),
      });

      const result = await response.json();
      if (response.ok) {
        showSuccess('Success', 'Booking created successfully');
        router.back();
      } else {
        showError('Booking Failed', result.message || 'Please try again');
      }
    } catch (error) {
      showError('Network Error', 'Please check your connection and try again');
      console.error('Booking error:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const generateDateArray = () => {
    const dates = [];
    const today = new Date();
    for (let i = 0; i < 30; i++) {
      const date = new Date();
      date.setDate(today.getDate() + i);
      dates.push(date);
    }
    return dates;
  };

  return (
    <View className="flex-1 bg-white">
      <StatusBar style="dark" />
      <Stack.Screen 
        options={{ 
          headerShown: true,
          title: 'Create Booking',
          headerStyle: { backgroundColor: '#fff' },
          headerShadowVisible: false,
        }} 
      />
      
      <Spinner
        visible={isLoading}
        textContent={'Creating booking...'}
        textStyle={{ color: '#FFF', fontFamily: 'Inter-Medium' }}
        overlayColor="rgba(0, 0, 0, 0.7)"
      />
      
      <ScrollView className="flex-1 p-6">
        <View className="space-y-6">
          {/* Vehicle Selection */}
          <View>
            <Text style={{ fontFamily: 'Inter-Medium' }} className="text-gray-700 mb-2 ml-1">
              Vehicle
            </Text>
            <TouchableOpacity
              onPress={() => !params.vehicleId && router.push('/modals/select-vehicle?onSelect=true')}
              disabled={!!params.vehicleId}
              className="bg-white rounded-xl p-4 border border-gray-200"
            >
              {selectedVehicle ? (
                <View className="flex-row items-center">
                  <MaterialIcons name="directions-car" size={24} color="#2563eb" />
                  <View className="ml-3">
                    <Text style={{ fontFamily: 'Poppins-SemiBold' }} className="text-gray-800">
                      {selectedVehicle.plate_number}
                    </Text>
                    <Text style={{ fontFamily: 'Inter-Regular' }} className="text-gray-600">
                      {selectedVehicle.brand} {selectedVehicle.model}
                    </Text>
                  </View>
                </View>
              ) : (
                <Text className="text-gray-500">Select a vehicle</Text>
              )}
            </TouchableOpacity>
          </View>

          {/* Location Selection */}
          <View>
            <Text style={{ fontFamily: 'Inter-Medium' }} className="text-gray-700 mb-2 ml-1">
              Location
            </Text>
            <TouchableOpacity
              onPress={() => router.push('/modals/select-location?onSelect=true')}
              className="bg-white rounded-xl p-4 border border-gray-200"
            >
              {selectedLocation ? (
                <View>
                  <Text style={{ fontFamily: 'Poppins-SemiBold' }} className="text-gray-800">
                    {selectedLocation.name}
                  </Text>
                  <Text style={{ fontFamily: 'Inter-Regular' }} className="text-gray-600 mt-1">
                    {selectedLocation.address}
                  </Text>
                </View>
              ) : (
                <Text className="text-gray-500">Select a location</Text>
              )}
            </TouchableOpacity>
          </View>

          {/* Date Selection */}
          <View>
            <Text style={{ fontFamily: 'Inter-Medium' }} className="text-gray-700 mb-2 ml-1">
              Date
            </Text>
            <TouchableOpacity
              onPress={() => setShowDateModal(true)}
              className="bg-white rounded-xl p-4 border border-gray-200"
            >
              <Text className="text-gray-800">
                {dateTime.toLocaleDateString()}
              </Text>
            </TouchableOpacity>
          </View>

          {/* Time Selection */}
          <View>
            <Text style={{ fontFamily: 'Inter-Medium' }} className="text-gray-700 mb-2 ml-1">
              Time
            </Text>
            <TouchableOpacity
              onPress={() => setShowTimeModal(true)}
              className="bg-white rounded-xl p-4 border border-gray-200"
            >
              <Text className="text-gray-800">
                {dateTime.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
              </Text>
            </TouchableOpacity>
          </View>

          {/* Create Button */}
          <TouchableOpacity
            onPress={handleCreateBooking}
            disabled={isLoading}
            className={`rounded-xl p-4 mt-4 ${isLoading ? 'bg-blue-400' : 'bg-blue-600'}`}
          >
            <Text style={{ fontFamily: 'Inter-Medium' }} className="text-white text-center text-lg">
              Create Booking
            </Text>
          </TouchableOpacity>
        </View>

        {/* Time Selection Modal */}
        <Modal
          visible={showTimeModal}
          transparent
          animationType="slide"
          onRequestClose={() => setShowTimeModal(false)}
        >
          <View className="flex-1 justify-end bg-black/50">
            <View className="bg-white rounded-t-3xl">
              <View className="p-4 border-b border-gray-200 flex-row justify-between items-center">
                <Text style={{ fontFamily: 'Inter-Medium' }} className="text-lg">
                  Select Time
                </Text>
                <TouchableOpacity 
                  onPress={() => setShowTimeModal(false)}
                  className="px-4 py-2"
                >
                  <Text className="text-blue-600" style={{ fontFamily: 'Inter-Medium' }}>
                    Done
                  </Text>
                </TouchableOpacity>
              </View>
              <ScrollView className="max-h-96">
                {generateTimeSlots().map(({ hour, minute, label }) => (
                  <TouchableOpacity
                    key={label}
                    onPress={() => handleSelectTime(hour, minute)}
                    className="p-4 border-b border-gray-100"
                  >
                    <Text 
                      className={`text-center text-lg ${
                        dateTime.getHours() === hour && dateTime.getMinutes() === minute 
                          ? 'text-blue-600 font-bold' 
                          : 'text-gray-800'
                      }`}
                      style={{ fontFamily: 'Inter-Regular' }}
                    >
                      {label}
                    </Text>
                  </TouchableOpacity>
                ))}
              </ScrollView>
            </View>
          </View>
        </Modal>

        {/* Date Selection Modal */}
        <Modal
          visible={showDateModal}
          transparent
          animationType="slide"
          onRequestClose={() => setShowDateModal(false)}
        >
          <View className="flex-1 justify-end bg-black/50">
            <View className="bg-white rounded-t-3xl">
              <View className="p-4 border-b border-gray-200 flex-row justify-between items-center">
                <Text style={{ fontFamily: 'Inter-Medium' }} className="text-lg">
                  Select Date
                </Text>
                <TouchableOpacity 
                  onPress={() => setShowDateModal(false)}
                  className="px-4 py-2"
                >
                  <Text className="text-blue-600" style={{ fontFamily: 'Inter-Medium' }}>
                    Done
                  </Text>
                </TouchableOpacity>
              </View>
              <ScrollView className="max-h-96">
                {generateDateArray().map((date, index) => (
                  <TouchableOpacity
                    key={index}
                    onPress={() => handleSelectDate(date)}
                    className="p-4 border-b border-gray-100"
                  >
                    <Text 
                      className={`text-center text-lg ${
                        dateTime.toDateString() === date.toDateString() 
                          ? 'text-blue-600 font-bold' 
                          : 'text-gray-800'
                      }`}
                      style={{ fontFamily: 'Inter-Regular' }}
                    >
                      {date.toLocaleDateString('en-US', { 
                        weekday: 'long',
                        year: 'numeric',
                        month: 'long',
                        day: 'numeric'
                      })}
                    </Text>
                  </TouchableOpacity>
                ))}
              </ScrollView>
            </View>
          </View>
        </Modal>
      </ScrollView>
    </View>
  );
}