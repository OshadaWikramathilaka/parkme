import { useState, useEffect } from 'react';
import { View, Text, TouchableOpacity, ScrollView } from 'react-native';
import { Stack, useRouter, useLocalSearchParams } from 'expo-router';
import { StatusBar } from 'expo-status-bar';
import { MaterialIcons } from '@expo/vector-icons';
import Spinner from 'react-native-loading-spinner-overlay';
import { ENDPOINTS } from '@/constants/Config';
import { getToken, getUser } from '@/utils/auth';

interface Vehicle {
  id: string;
  brand: string;
  model: string;
  plate_number: string;
}

export default function SelectVehicleModal() {
  const router = useRouter();
  const { onSelect } = useLocalSearchParams<{ onSelect?: string }>();
  const [vehicles, setVehicles] = useState<Vehicle[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    fetchVehicles();
  }, []);

  const fetchVehicles = async () => {
    try {
      const token = await getToken();
      const user = await getUser();
      const response = await fetch(`${ENDPOINTS.VEHICLES}/user/${user.id}`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });
      const data = await response.json();
      if (response.ok) {
        setVehicles(data.data);
      }
    } catch (error) {
      console.error('Error fetching vehicles:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleSelect = (vehicle: Vehicle) => {
    if (onSelect) {
      router.back();
      setTimeout(() => {
        router.setParams({ selectedVehicle: JSON.stringify(vehicle) });
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
          title: 'Select Vehicle',
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
          {vehicles.map(vehicle => (
            <TouchableOpacity
              key={vehicle.id}
              onPress={() => handleSelect(vehicle)}
              className="bg-white rounded-xl p-4 border border-gray-200"
            >
              <View className="flex-row items-center">
                <MaterialIcons name="directions-car" size={24} color="#2563eb" />
                <View className="ml-3">
                  <Text style={{ fontFamily: 'Poppins-SemiBold' }} className="text-gray-800">
                    {vehicle.plate_number}
                  </Text>
                  <Text style={{ fontFamily: 'Inter-Regular' }} className="text-gray-600">
                    {vehicle.brand} {vehicle.model}
                  </Text>
                </View>
              </View>
            </TouchableOpacity>
          ))}
        </View>
      </ScrollView>
    </View>
  );
} 