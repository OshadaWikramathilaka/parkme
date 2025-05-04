import { View, Text } from 'react-native';

export default function TestTailwind() {
  return (
    <View className="flex-1 items-center justify-center bg-blue-500 p-4">
      <Text className="text-white text-xl font-bold">
        Tailwind CSS is working!
      </Text>
    </View>
  );
} 