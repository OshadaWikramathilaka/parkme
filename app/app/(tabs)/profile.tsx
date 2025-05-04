import React, { useEffect, useState } from 'react';
import { View, Text, TouchableOpacity, ScrollView, RefreshControl, Image } from 'react-native';
import { MaterialIcons, Ionicons } from '@expo/vector-icons';
import { useRouter, useLocalSearchParams } from 'expo-router';
import { SafeAreaView } from 'react-native-safe-area-context';
import Spinner from 'react-native-loading-spinner-overlay';
import { ENDPOINTS } from '@/constants/Config';
import { showError, showSuccess } from '@/utils/flashMessage';
import { getToken, removeToken } from '@/utils/auth';
import Header from '@/components/ui/Header';

interface UserProfile {
  name: string;
  email: string;
  role: string;
  status: string;
  profile_image_url: string;
}

export default function ProfileScreen() {
  const router = useRouter();
  const { refresh } = useLocalSearchParams<{ refresh?: string }>();
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);

  const fetchProfile = async () => {
    try {
      const token = await getToken();
      if (!token) {
        router.replace('/(auth)/login');
        return;
      }

      const response = await fetch(ENDPOINTS.PROFILE, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      const result = await response.json();
      if (response.ok && result.success) {
        setProfile(result.data);
      } else {
        if (response.status === 401) {
          await removeToken();
          router.replace('/(auth)/login');
        } else {
          showError('Error', result.message || 'Failed to load profile');
        }
      }
    } catch (error) {
      showError('Network Error', 'Failed to connect to server');
      console.error('Profile fetch error:', error);
    } finally {
      setIsLoading(false);
      setRefreshing(false);
    }
  };

  useEffect(() => {
    fetchProfile();
  }, []);

  useEffect(() => {
    if (refresh === 'true') {
      fetchProfile();
      // Clear the refresh parameter
      router.setParams({ refresh: undefined });
    }
  }, [refresh]);

  const handleLogout = async () => {
    try {
      await removeToken();
      showSuccess('Logout successful');
      router.replace('/(auth)/login');
    } catch (error) {
      showError('Error', 'Failed to logout');
      console.error('Logout error:', error);
    }
  };

  const onRefresh = () => {
    setRefreshing(true);
    fetchProfile();
  };

  return (
    <SafeAreaView className="flex-1 bg-gray-50">
        <Header title="Profile" />
      <ScrollView 
        className="flex-1"
        refreshControl={
          <RefreshControl refreshing={refreshing} onRefresh={onRefresh} />
        }
      >
        <Spinner
          visible={isLoading}
          textContent={'Loading...'}
          textStyle={{ color: '#FFF', fontFamily: 'Inter-Medium' }}
          overlayColor="rgba(0, 0, 0, 0.7)"
        />

        {!isLoading && !profile && (
          <View className="flex-1 items-center justify-center p-6">
            <Text style={{ fontFamily: 'Inter-Medium' }} className="text-gray-500">
              Failed to load profile
            </Text>
          </View>
        )}

        {!isLoading && profile && (
          <View className="p-6">
            {/* Profile Header */}
            <View className="items-center mb-8">
              <View className="w-24 h-24 bg-blue-600 rounded-full items-center justify-center mb-4">
                {/* first check user have profile image */}
                {profile.profile_image_url ? (
                  <Image source={{ uri: profile.profile_image_url }} style={{ width: '100%', height: '100%', borderRadius: 100 }} />
                ) : (
                  <Text style={{ fontFamily: 'Poppins-Bold' }} className="text-white text-3xl">
                    {profile.name.split(' ').map(word => word.charAt(0).toUpperCase()).slice(0, 2).join('')}
                  </Text>
                )}
              </View>
              <Text style={{ fontFamily: 'Poppins-Bold' }} className="text-2xl text-gray-800">
                {profile.name}
              </Text>
              <Text style={{ fontFamily: 'Inter-Regular' }} className="text-gray-500">
                {profile.email}
              </Text>
            </View>

            {/* Profile Details */}
            <View className="bg-white rounded-2xl shadow-sm p-6 mb-6">
              {/* <View className="flex-row items-center mb-6">
                <View className="bg-blue-100 p-3 rounded-full mr-4">
                  <Ionicons name="person-outline" size={24} color="#2563eb" />
                </View>
                <View>
                  <Text style={{ fontFamily: 'Inter-Medium' }} className="text-gray-500 text-sm">
                    Role
                  </Text>
                  <Text style={{ fontFamily: 'Inter-Regular' }} className="text-gray-800">
                    {profile.role}
                  </Text>
                </View>
              </View> */}

              <View className="flex-row items-center">
                <View className="bg-green-100 p-3 rounded-full mr-4">
                  <Ionicons name="checkmark-circle-outline" size={24} color="#059669" />
                </View>
                <View>
                  <Text style={{ fontFamily: 'Inter-Medium' }} className="text-gray-500 text-sm">
                    Status
                  </Text>
                  <Text style={{ fontFamily: 'Inter-Regular' }} className="text-gray-800">
                    {profile.status}
                  </Text>
                </View>
              </View>
            </View>

            {/* Actions */}
            <View className="space-y-4">
              <TouchableOpacity 
                className="flex-row items-center bg-white p-4 rounded-xl"
                onPress={() => router.push('/modals/bookings')}
                activeOpacity={0.7}
              >
                <MaterialIcons name="history" size={24} color="#2563eb" />
                <Text style={{ fontFamily: 'Inter-Medium' }} className="ml-3 text-gray-800">
                  Booking History
                </Text>
              </TouchableOpacity>
              <TouchableOpacity 
                className="flex-row items-center bg-white p-4 rounded-xl"
                onPress={() => router.push('/modals/edit-profile')}
                activeOpacity={0.7}
              >
                <MaterialIcons name="edit" size={24} color="#2563eb" />
                <Text style={{ fontFamily: 'Inter-Medium' }} className="ml-3 text-gray-800">
                  Edit Profile
                </Text>
              </TouchableOpacity>

            

              <TouchableOpacity 
                className="flex-row items-center bg-red-50 p-4 rounded-xl mt-6"
                onPress={handleLogout}
              >
                <MaterialIcons name="logout" size={24} color="#dc2626" />
                <Text style={{ fontFamily: 'Inter-Medium' }} className="ml-3 text-red-600">
                  Logout
                </Text>
              </TouchableOpacity>
            </View>
          </View>
        )}
      </ScrollView>
    </SafeAreaView>
  );
} 