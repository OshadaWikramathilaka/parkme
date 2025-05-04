import React from 'react';
import { View, Text, TouchableOpacity } from 'react-native';
import { MaterialIcons } from '@expo/vector-icons';
import { useRouter } from 'expo-router';

interface HeaderProps {
  title: string;
  showBack?: boolean;
  showClose?: boolean;
  rightAction?: () => void;
  rightIcon?: keyof typeof MaterialIcons.glyphMap;
}

export default function Header({ 
  title, 
  showBack, 
  showClose,
  rightAction,
  rightIcon
}: HeaderProps) {
  const router = useRouter();

  const handleBack = () => {
    if (showClose) {
      router.back();
    } else {
      router.back();
    }
  };

  return (
    <View className="bg-white border-b border-gray-200 px-4 py-3 flex-row items-center justify-between">
      {/* Left Action */}
      {(showBack || showClose) ? (
        <TouchableOpacity 
          onPress={handleBack}
          className="p-2 -ml-2"
        >
          <MaterialIcons 
            name={showClose ? "close" : "arrow-back"} 
            size={24} 
            color="#374151" 
          />
        </TouchableOpacity>
      ) : (
        <View style={{ width: 40 }} />
      )}

      {/* Title */}
      <Text style={{ fontFamily: 'Poppins-Bold' }} className="text-xl text-gray-800">
        {title}
      </Text>

      {/* Right Action */}
      {rightAction && rightIcon ? (
        <TouchableOpacity 
          onPress={rightAction}
          className="p-2 -mr-2"
        >
          <MaterialIcons name={rightIcon} size={24} color="#374151" />
        </TouchableOpacity>
      ) : (
        <View style={{ width: 40 }} />
      )}
    </View>
  );
} 