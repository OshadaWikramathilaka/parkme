import React, { useEffect, useState } from 'react';
import { View, Text, TextInput, TouchableOpacity, ScrollView, KeyboardAvoidingView, Platform } from 'react-native';
import { useRouter, useLocalSearchParams } from 'expo-router';
import { SafeAreaView } from 'react-native-safe-area-context';
import Spinner from 'react-native-loading-spinner-overlay';
import { MaterialIcons } from '@expo/vector-icons';
import { ENDPOINTS } from '@/constants/Config';
import { showError, showSuccess } from '@/utils/flashMessage';
import { getToken } from '@/utils/auth';

interface VehicleForm {
  plate_number: string;
  brand: string;
  model: string;
}

export default function EditVehicleScreen() {
  const router = useRouter();
  const { id } = useLocalSearchParams();
  const [isLoading, setIsLoading] = useState(true);
  const [form, setForm] = useState<VehicleForm>({
    plate_number: '',
    brand: '',
    model: '',
  });

  const fetchVehicle = async () => {
    try {
      const token = await getToken();
      if (!token) {
        router.replace('/(auth)/login');
        return;
      }

      const response = await fetch(`${ENDPOINTS.VEHICLES}/${id}`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      const result = await response.json();
      if (response.ok && result.success) {
        const { plate_number, brand, model } = result.data;
        setForm({ plate_number, brand, model });
      } else {
        showError('Error', result.message || 'Failed to load vehicle');
        router.back();
      }
    } catch (error) {
      showError('Network Error', 'Failed to connect to server');
      console.error('Vehicle fetch error:', error);
      router.back();
    } finally {
      setIsLoading(false);
    }
  };

  const handleSubmit = async () => {
    if (!form.plate_number || !form.brand || !form.model) {
      showError('Error', 'Please fill in all fields');
      return;
    }

    setIsLoading(true);
    try {
      const token = await getToken();
      if (!token) {
        router.replace('/(auth)/login');
        return;
      }

      const response = await fetch(`${ENDPOINTS.VEHICLES}/${id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(form),
      });

      const result = await response.json();
      if (response.ok && result.success) {
        showSuccess('Success', 'Vehicle updated successfully');
        router.back();
      } else {
        showError('Error', result.message || 'Failed to update vehicle');
      }
    } catch (error) {
      showError('Network Error', 'Failed to connect to server');
      console.error('Update vehicle error:', error);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchVehicle();
  }, [id]);

  return (
    <SafeAreaView className="flex-1 bg-gray-50">
      <KeyboardAvoidingView
        behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
        className="flex-1"
      >
        <ScrollView className="flex-1">
          <Spinner
            visible={isLoading}
            textContent={form.plate_number ? 'Updating vehicle...' : 'Loading...'}
            textStyle={{ color: '#FFF', fontFamily: 'Inter-Medium' }}
            overlayColor="rgba(0, 0, 0, 0.7)"
          />

          <View className="p-6">
            <View className="flex-row items-center mb-6">
              <TouchableOpacity
                onPress={() => router.back()}
                className="mr-4"
              >
                <MaterialIcons name="arrow-back" size={24} color="#374151" />
              </TouchableOpacity>
              <Text style={{ fontFamily: 'Poppins-Bold' }} className="text-2xl text-gray-800">
                Edit Vehicle
              </Text>
            </View>

            <View className="space-y-6">
              <View>
                <Text style={{ fontFamily: 'Inter-Medium' }} className="text-gray-700 mb-2 ml-1">
                  Plate Number
                </Text>
                <TextInput
                  className="bg-white rounded-xl p-4 text-gray-800 border border-gray-200"
                  placeholder="Enter plate number"
                  value={form.plate_number}
                  onChangeText={(text) => setForm(prev => ({ ...prev, plate_number: text }))}
                  autoCapitalize="characters"
                  editable={!isLoading}
                />
              </View>

              <View>
                <Text style={{ fontFamily: 'Inter-Medium' }} className="text-gray-700 mb-2 ml-1">
                  Brand
                </Text>
                <TextInput
                  className="bg-white rounded-xl p-4 text-gray-800 border border-gray-200"
                  placeholder="Enter vehicle brand"
                  value={form.brand}
                  onChangeText={(text) => setForm(prev => ({ ...prev, brand: text }))}
                  editable={!isLoading}
                />
              </View>

              <View>
                <Text style={{ fontFamily: 'Inter-Medium' }} className="text-gray-700 mb-2 ml-1">
                  Model
                </Text>
                <TextInput
                  className="bg-white rounded-xl p-4 text-gray-800 border border-gray-200"
                  placeholder="Enter vehicle model"
                  value={form.model}
                  onChangeText={(text) => setForm(prev => ({ ...prev, model: text }))}
                  editable={!isLoading}
                />
              </View>

              <TouchableOpacity
                onPress={handleSubmit}
                disabled={isLoading}
                className={`rounded-xl p-4 shadow-sm mt-4 ${isLoading ? 'bg-blue-400' : 'bg-blue-600'}`}
              >
                <Text style={{ fontFamily: 'Inter-Medium' }} className="text-white text-center text-lg">
                  Update Vehicle
                </Text>
              </TouchableOpacity>
            </View>
          </View>
        </ScrollView>
      </KeyboardAvoidingView>
    </SafeAreaView>
  );
} 