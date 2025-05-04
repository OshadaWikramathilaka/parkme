import React, { useState, useEffect } from 'react';
import { View, TextInput, TouchableOpacity, StyleSheet, Image, Platform } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useRouter } from 'expo-router';
import { MaterialIcons } from '@expo/vector-icons';
import * as ImagePicker from 'expo-image-picker';
import { ThemedText } from '@/components/ThemedText';
import { Colors } from '@/constants/Colors';
import { useColorScheme } from 'react-native';
import { getToken } from '@/utils/auth';
import { ENDPOINTS } from '@/constants/Config';
import { showError, showSuccess } from '@/utils/flashMessage';
import Header from '@/components/ui/Header';
import Spinner from 'react-native-loading-spinner-overlay';

export default function EditProfileModal() {
  const router = useRouter();
  const colorScheme = useColorScheme();
  const colors = Colors[colorScheme ?? 'light'];
  const [name, setName] = useState('');
  const [currentImage, setCurrentImage] = useState<string | null>(null);
  const [selectedImage, setSelectedImage] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    loadProfile();
    (async () => {
      if (Platform.OS !== 'web') {
        const { status } = await ImagePicker.requestMediaLibraryPermissionsAsync();
        if (status !== 'granted') {
          showError('Permission Required', 'Sorry, we need camera roll permissions to upload images!');
        }
      }
    })();
  }, []);

  const loadProfile = async () => {
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
        setName(result.data.name);
        setCurrentImage(result.data.profile_image_url);
      } else {
        showError('Error', result.message || 'Failed to load profile');
      }
    } catch (error) {
      showError('Error', 'Failed to load profile');
      console.error('Profile load error:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const pickImage = async () => {
    try {
      // No permissions request is necessary for launching the image library
      const result = await ImagePicker.launchImageLibraryAsync({
        mediaTypes: 'images',
        allowsEditing: true,
        aspect: [1, 1],
        quality: 1,
      });

      if (!result.canceled) {
        setSelectedImage(result.assets[0].uri);
      }
    } catch (error) {
      showError('Error', 'Failed to pick image');
      console.error('Image picker error:', error);
    }
  };

  const handleSubmit = async () => {
    if (!name.trim()) {
      showError('Validation Error', 'Name is required');
      return;
    }

    try {
      setIsSubmitting(true);
      const token = await getToken();
      if (!token) {
        router.replace('/(auth)/login');
        return;
      }

      const formData = new FormData();
      formData.append('name', name.trim());

      if (selectedImage) {
        const filename = selectedImage.split('/').pop() || 'image.jpg';
        const match = /\.(\w+)$/.exec(filename);
        const type = match ? `image/${match[1]}` : 'image/jpeg';

        formData.append('image', {
          uri: selectedImage,
          type,
          name: filename,
        } as any);
      }

      const response = await fetch(ENDPOINTS.PROFILE, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'multipart/form-data',
        },
        body: formData,
      });

      const result = await response.json();
      
      if (response.ok && result.success) {
        showSuccess('Profile Updated', 'Your profile has been updated successfully');
        router.back();
        setTimeout(() => {
          router.setParams({ refresh: 'true' });
        }, 100);
      } else {
        showError('Update Failed', result.message || 'Failed to update profile');
      }
    } catch (error) {
      showError('Error', 'An error occurred while updating profile');
      console.error('Profile update error:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <SafeAreaView style={[styles.container, { backgroundColor: colors.background }]} edges={['top']}>
      <Header title="Edit Profile" />
      
      <Spinner
        visible={isLoading}
        textContent={'Loading...'}
        textStyle={{ color: '#FFF' }}
        overlayColor="rgba(0, 0, 0, 0.7)"
      />

      <View style={styles.content}>
        {/* Profile Image Section */}
        <TouchableOpacity onPress={pickImage} style={styles.imageContainer}>
          {selectedImage ? (
            <Image source={{ uri: selectedImage }} style={styles.profileImage} />
          ) : currentImage ? (
            <Image source={{ uri: currentImage }} style={styles.profileImage} />
          ) : (
            <View style={[styles.imagePlaceholder, { backgroundColor: colors.tint }]}>
              <MaterialIcons name="person" size={32} color="#fff" />
            </View>
          )}
          <ThemedText style={styles.changePhotoText}>Change Photo</ThemedText>
        </TouchableOpacity>

        {/* Name Input */}
        <View style={styles.inputContainer}>
          <ThemedText style={styles.label}>Name</ThemedText>
          <TextInput
            style={[styles.input, { color: colors.text, borderColor: colors.icon }]}
            placeholder="Enter your name"
            placeholderTextColor={colors.icon}
            value={name}
            onChangeText={setName}
          />
        </View>

        {/* Submit Button */}
        <TouchableOpacity
          style={[styles.submitButton, { backgroundColor: colors.tint }]}
          onPress={handleSubmit}
          disabled={isSubmitting}
        >
          <ThemedText style={styles.submitButtonText}>
            {isSubmitting ? 'Updating...' : 'Update Profile'}
          </ThemedText>
        </TouchableOpacity>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  content: {
    padding: 20,
  },
  imageContainer: {
    alignItems: 'center',
    marginBottom: 30,
  },
  profileImage: {
    width: 120,
    height: 120,
    borderRadius: 60,
    marginBottom: 10,
  },
  imagePlaceholder: {
    width: 120,
    height: 120,
    borderRadius: 60,
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: 10,
  },
  changePhotoText: {
    fontSize: 14,
    color: '#666',
    marginTop: 8,
  },
  inputContainer: {
    marginBottom: 20,
  },
  label: {
    fontSize: 16,
    marginBottom: 8,
  },
  input: {
    height: 50,
    borderWidth: 1,
    borderRadius: 8,
    paddingHorizontal: 15,
    fontSize: 16,
  },
  submitButton: {
    height: 50,
    borderRadius: 25,
    justifyContent: 'center',
    alignItems: 'center',
    marginTop: 20,
  },
  submitButtonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: 'bold',
  },
}); 