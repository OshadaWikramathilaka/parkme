import React, { useState, useEffect } from 'react';
import { View, Text, TouchableOpacity, ScrollView } from 'react-native';
import { Stack, useRouter, useLocalSearchParams } from 'expo-router';
import { StatusBar } from 'expo-status-bar';
import { MaterialIcons } from '@expo/vector-icons';
import Spinner from 'react-native-loading-spinner-overlay';
import { ENDPOINTS } from '@/constants/Config';
import { getToken } from '@/utils/auth';
import { showError } from '@/utils/flashMessage';

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

export default function SelectLocationModal() {
  const router = useRouter();
  const { onSelect } = useLocalSearchParams<{ onSelect?: string }>();
  const [locations, setLocations] = useState<Location[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    fetchLocations();
  }, []);

  const fetchLocations = async () => {
    try {
      const token = await getToken();
      const response = await fetch(ENDPOINTS.LOCATIONS, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });
      const data = await response.json();
      if (response.ok) {
        setLocations(Array.isArray(data) ? data : []);
      }
    } catch (error) {
      console.error('Error fetching locations:', error);
      setLocations([]);
    } finally {
      setIsLoading(false);
    }
  };

  const getAvailableSpots = (location: Location) => {
    if (!location.slots) return 0;
    return location.slots.filter(slot => !slot.is_occupied).length;
  };

  const handleSelect = (location: Location) => {
    if (onSelect) {
      //check if location has slots
      if(location.slots.filter(slot => !slot.is_occupied).length === 0){
        showError('Error', 'No slots available');
        return;
      }
      router.back();
      setTimeout(() => {
        router.setParams({ selectedLocation: JSON.stringify(location) });
      }, 100);
    } else {
      router.back();
    }
  };

  return (
    <View className="flex-1 bg-white">
      <StatusBar style="dark" />
      <Stack.Screen 
        options={{ 
          headerShown: true,
          title: 'Select Location',
          headerStyle: {
            backgroundColor: '#fff',
          },
          headerShadowVisible: false,
        }} 
      />
      
      <Spinner
        visible={isLoading}
        textContent={'Loading...'}
        textStyle={{ color: '#FFF', fontFamily: 'Inter-Medium' }}
        overlayColor="rgba(0, 0, 0, 0.7)"
      />
      
      <ScrollView className="flex-1 p-6">
        <View className="space-y-4">
          {locations && locations.length > 0 ? (
            locations.map(location => (
              <TouchableOpacity
                key={location.id}
                onPress={() => handleSelect(location)}
                className="bg-white rounded-xl p-4 border border-gray-200"
              >
                <View>
                  <View className="flex-row items-center justify-between">
                    <Text style={{ fontFamily: 'Poppins-SemiBold' }} className="text-gray-800 text-lg">
                      {location.name}
                    </Text>
                    <View className="bg-green-50 px-3 py-1 rounded-full">
                      <Text style={{ fontFamily: 'Inter-Medium' }} className="text-green-800">
                        {getAvailableSpots(location)} spots
                      </Text>
                    </View>
                  </View>
                  <Text style={{ fontFamily: 'Inter-Regular' }} className="text-gray-600 mt-1">
                    {location.address}
                  </Text>
                </View>
              </TouchableOpacity>
            ))
          ) : (
            <View className="p-4">
              <Text style={{ fontFamily: 'Inter-Regular' }} className="text-gray-600 text-center">
                No locations available
              </Text>
            </View>
          )}
        </View>
      </ScrollView>
    </View>
  );
} 