import { Tabs } from 'expo-router';
import { MaterialIcons } from '@expo/vector-icons';
import { SafeAreaView } from 'react-native';

export default function TabLayout() {
  return (
    <SafeAreaView style={{ flex: 1 }}>
      <Tabs
        screenOptions={{
          tabBarActiveTintColor: '#2563eb',
          tabBarInactiveTintColor: '#6B7280',
          tabBarStyle: {
            borderTopWidth: 1,
            borderTopColor: '#E5E7EB',
            height: 60,
            paddingBottom: 8,
            paddingTop: 8,
          },
          tabBarLabelStyle: {
            fontFamily: 'Inter-Medium',
            fontSize: 12,
          },
          headerShown: false,
        }}
        initialRouteName="home"
      >
        <Tabs.Screen
          name="home"
          options={{
            title: 'Home',
            tabBarIcon: ({ color, size }) => (
              <MaterialIcons name="home" size={size} color={color} />
            ),
          }}
        />
        <Tabs.Screen
          name="vehicles"
          options={{
            title: 'Vehicles',
            tabBarIcon: ({ color, size }) => (
              <MaterialIcons name="directions-car" size={size} color={color} />
            ),
          }}
        />
        <Tabs.Screen
          name="wallet"
          options={{
            title: 'Wallet',
            tabBarIcon: ({ color, size }) => (
              <MaterialIcons name="account-balance-wallet" size={size} color={color} />
            ),
          }}
        />
        <Tabs.Screen
          name="profile"
          options={{
            title: 'Profile',
            tabBarIcon: ({ color, size }) => (
              <MaterialIcons name="person" size={size} color={color} />
            ),
          }}
        />
      </Tabs>
    </SafeAreaView>
  );
}
