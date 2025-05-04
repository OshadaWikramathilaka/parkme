import React, { useEffect, useState } from 'react';
import { View, ScrollView, StyleSheet, ActivityIndicator, Image } from 'react-native';
import { SafeAreaView, useSafeAreaInsets } from 'react-native-safe-area-context';
import { ThemedText } from '@/components/ThemedText';
import { MaterialCommunityIcons } from '@expo/vector-icons';
import { format } from 'date-fns';
import axios from 'axios';
import { Colors } from '@/constants/Colors';
import { useColorScheme } from 'react-native';
import { getToken, getUser } from '@/utils/auth';
import { API_URL } from '@/constants/Config';

interface User {
  id: string;
  name: string;
  email: string;
  profile_image_url?: string;
  role: string;
  status: string;
}

interface UserStats {
  total_bookings: number;
  active_bookings: number;
  completed_bookings: number;
  cancelled_bookings: number;
  total_spent_amount: number;
  average_booking_duration: number;
  last_booking_date: string;
  upcoming_booking_date: string;
}

function HomeScreen() {
  const [user, setUser] = useState<User | null>(null);
  const [stats, setStats] = useState<UserStats | null>(null);
  const [loading, setLoading] = useState(true);
  const colorScheme = useColorScheme();
  const colors = Colors[colorScheme ?? 'light'];
  const insets = useSafeAreaInsets();

  useEffect(() => {
    loadUserAndStats();
  }, []);

  const loadUserAndStats = async () => {
    try {
      const [userData, token] = await Promise.all([getUser(), getToken()]);
      setUser(userData);
      
      if (token) {
        const response = await axios.get(`${API_URL}/api/user/stats`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        setStats(response.data.data);
      }
    } catch (error) {
      console.error('Error loading data:', error);
    } finally {
      setLoading(false);
    }
  };

  const getGreeting = () => {
    const hour = new Date().getHours();
    if (hour < 12) return 'Good Morning';
    if (hour < 17) return 'Good Afternoon';
    return 'Good Evening';
  };

  if (loading) {
    return (
      <View style={[styles.container, { backgroundColor: colors.background }]}>
        <View style={styles.loadingContainer}>
          <ActivityIndicator size="large" color={colors.tint} />
        </View>
      </View>
    );
  }

  return (
    <View style={[styles.container, { backgroundColor: colors.background }]}>
      <View style={[styles.headerBackground, { backgroundColor: colors.tint }]} />
      <SafeAreaView edges={['top']} style={styles.safeArea}>
        <ScrollView style={styles.scrollView} contentContainerStyle={styles.scrollContent}>
          {/* User Greeting Section */}
          <View style={[styles.greetingContainer, { backgroundColor: colors.tint, paddingTop: insets.top + 20 }]}>
            <View style={styles.userInfo}>
              <ThemedText style={styles.greeting}>{getGreeting()},</ThemedText>
              <ThemedText style={styles.userName}>{user?.name}</ThemedText>
            </View>
            {user?.profile_image_url && (
              <Image
                source={{ uri: user.profile_image_url }}
                style={styles.profileImage}
                resizeMode="cover"
              />
            )}
          </View>

          {/* Stats Grid */}
          <View style={styles.statsGrid}>
            <View style={[styles.statCard, styles.elevation, { backgroundColor: colors.background }]}>
              <MaterialCommunityIcons name="car-clock" size={24} color={colors.tint} />
              <ThemedText style={[styles.statNumber, { color: colors.text }]}>{stats?.active_bookings || 0}</ThemedText>
              <ThemedText style={[styles.statLabel, { color: colors.icon }]}>Active Bookings</ThemedText>
            </View>

            <View style={[styles.statCard, styles.elevation, { backgroundColor: colors.background }]}>
              <MaterialCommunityIcons name="check-circle" size={24} color={colors.tint} />
              <ThemedText style={[styles.statNumber, { color: colors.text }]}>{stats?.completed_bookings || 0}</ThemedText>
              <ThemedText style={[styles.statLabel, { color: colors.icon }]}>Completed</ThemedText>
            </View>

            <View style={[styles.statCard, styles.elevation, { backgroundColor: colors.background }]}>
              <MaterialCommunityIcons name="currency-usd" size={24} color={colors.tint} />
              <ThemedText style={[styles.statNumber, { color: colors.text }]}>${stats?.total_spent_amount || 0}</ThemedText>
              <ThemedText style={[styles.statLabel, { color: colors.icon }]}>Total Spent</ThemedText>
            </View>

            <View style={[styles.statCard, styles.elevation, { backgroundColor: colors.background }]}>
              <MaterialCommunityIcons name="car-multiple" size={24} color={colors.tint} />
              <ThemedText style={[styles.statNumber, { color: colors.text }]}>{stats?.total_bookings || 0}</ThemedText>
              <ThemedText style={[styles.statLabel, { color: colors.icon }]}>Total Bookings</ThemedText>
            </View>
          </View>

          {/* Upcoming Booking */}
          {stats?.upcoming_booking_date && (
            <View style={[styles.upcomingCard, styles.elevation, { backgroundColor: colors.background }]}>
              <View style={styles.upcomingHeader}>
                <MaterialCommunityIcons name="calendar-clock" size={24} color={colors.tint} />
                <ThemedText style={[styles.upcomingTitle, { color: colors.text }]}>Upcoming Booking</ThemedText>
              </View>
              <ThemedText style={[styles.upcomingDate, { color: colors.text }]}>
                {format(new Date(stats.upcoming_booking_date), 'PPP')}
              </ThemedText>
              <ThemedText style={[styles.upcomingTime, { color: colors.icon }]}>
                {format(new Date(stats.upcoming_booking_date), 'p')}
              </ThemedText>
            </View>
          )}

          {/* Last Booking */}
          {stats?.last_booking_date && (
            <View style={[styles.upcomingCard, styles.elevation, { backgroundColor: colors.background }]}>
              <View style={styles.upcomingHeader}>
                <MaterialCommunityIcons name="calendar-clock" size={24} color={colors.tint} />
                <ThemedText style={[styles.upcomingTitle, { color: colors.text }]}>Last Booking</ThemedText>
              </View>
              <ThemedText style={[styles.upcomingDate, { color: colors.text }]}>
                {format(new Date(stats.last_booking_date), 'PPP')}
              </ThemedText>
              <ThemedText style={[styles.upcomingTime, { color: colors.icon }]}>
                {format(new Date(stats.last_booking_date), 'p')}
              </ThemedText>
            </View>
          )}
        </ScrollView>
      </SafeAreaView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  headerBackground: {
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    height: 200, // Adjust this value based on your needs
  },
  safeArea: {
    flex: 1,
  },
  scrollView: {
    flex: 1,
  },
  scrollContent: {
    flexGrow: 1,
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  greetingContainer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 20,
    paddingTop: 0, // This will be set dynamically
    borderBottomLeftRadius: 30,
    borderBottomRightRadius: 30,
    marginTop: -20, // To compensate for the SafeArea padding
  },
  userInfo: {
    flex: 1,
  },
  greeting: {
    fontSize: 16,
    color: '#fff',
    opacity: 0.9,
  },
  userName: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#fff',
    marginTop: 4,
  },
  profileImage: {
    width: 50,
    height: 50,
    borderRadius: 25,
    marginLeft: 15,
    backgroundColor: '#fff',
  },
  statsGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    padding: 15,
    justifyContent: 'space-between',
  },
  statCard: {
    borderRadius: 15,
    padding: 15,
    alignItems: 'center',
    width: '47%',
    marginBottom: 15,
  },
  elevation: {
    elevation: 3,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
  },
  statNumber: {
    fontSize: 24,
    fontWeight: 'bold',
    marginVertical: 8,
  },
  statLabel: {
    fontSize: 14,
    textAlign: 'center',
  },
  upcomingCard: {
    margin: 15,
    padding: 20,
    borderRadius: 15,
  },
  upcomingHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 15,
  },
  upcomingTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    marginLeft: 10,
  },
  upcomingDate: {
    fontSize: 16,
  },
  upcomingTime: {
    fontSize: 14,
    marginTop: 5,
  },
});

export default HomeScreen; 