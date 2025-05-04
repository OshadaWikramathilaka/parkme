import { Stack } from 'expo-router';

export default function ModalLayout() {
  return (
    <Stack screenOptions={{ headerShown: false }}>
      <Stack.Screen 
        name="add-vehicle" 
        options={{ 
          presentation: 'modal',
          headerShown: false,
        }} 
      />
      <Stack.Screen 
        name="edit-vehicle" 
        options={{ 
          presentation: 'modal',
          headerShown: false,
        }} 
      />
      <Stack.Screen 
        name="edit-profile" 
        options={{ 
          presentation: 'modal',
          headerShown: false,
        }} 
      />
    </Stack>
  );
} 