import React, { useEffect, useState } from 'react';
import { View, Text, TouchableOpacity, ScrollView, RefreshControl } from 'react-native';
import { MaterialIcons } from '@expo/vector-icons';
import { useRouter, useLocalSearchParams, useFocusEffect } from 'expo-router';
import { SafeAreaView } from 'react-native-safe-area-context';
import Spinner from 'react-native-loading-spinner-overlay';
import { ENDPOINTS } from '@/constants/Config';
import { showError } from '@/utils/flashMessage';
import { getToken } from '@/utils/auth';
import Header from '@/components/ui/Header';

interface Transaction {
  id: string;
  type: 'top_up' | 'deduct';
  amount: number;
  description: string;
  createdAt: string;
}

interface WalletData {
  balance: number;
  transactions: Transaction[];
}

export default function WalletScreen() {
  const router = useRouter();
  const { refresh } = useLocalSearchParams<{ refresh: string }>();
  const [walletData, setWalletData] = useState<WalletData | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);

  const fetchWalletData = async () => {
    try {
      const token = await getToken();
      if (!token) {
        router.replace('/(auth)/login');
        return;
      }

      // Fetch balance
      const balanceResponse = await fetch(ENDPOINTS.WALLET_BALANCE, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      // Fetch transactions
      const transactionsResponse = await fetch(ENDPOINTS.WALLET_TRANSACTIONS, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      const balanceResult = await balanceResponse.json();
      const transactionsResult = await transactionsResponse.json();

      if (balanceResponse.ok && transactionsResponse.ok) {
        setWalletData({
          balance: balanceResult.data.balance,
          transactions: transactionsResult.data || [],
        });
      } else {
        if (balanceResponse.status === 401 || transactionsResponse.status === 401) {
          router.replace('/(auth)/login');
        } else {
          showError('Error', 'Failed to load wallet data');
        }
      }
    } catch (error) {
      showError('Network Error', 'Failed to connect to server');
      console.error('Wallet fetch error:', error);
    } finally {
      setIsLoading(false);
      setRefreshing(false);
    }
  };

  const handleTopUp = () => {
    router.push('/modals/top-up');
  };

  useEffect(() => {
    fetchWalletData();
  }, []);

  useEffect(() => {
    if (refresh === 'true') {
      fetchWalletData();
    }
  }, [refresh]);

  useFocusEffect(
    React.useCallback(() => {
      fetchWalletData();
    }, [])
  );

  const onRefresh = () => {
    setRefreshing(true);
    fetchWalletData();
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  return (
    <SafeAreaView className="flex-1 bg-gray-50" edges={['top']}>
      <Header title="Wallet" />

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

        {!isLoading && !walletData && (
          <View className="flex-1 items-center justify-center p-6">
            <Text style={{ fontFamily: 'Inter-Medium' }} className="text-gray-500">
              Failed to load wallet data
            </Text>
          </View>
        )}

        {!isLoading && walletData && (
          <View className="p-6">
            {/* Balance Card */}
            <View className="bg-blue-600 rounded-2xl p-6 mb-6">
              <Text style={{ fontFamily: 'Inter-Medium' }} className="text-blue-100 mb-2">
                Available Balance
              </Text>
              <Text style={{ fontFamily: 'Poppins-Bold' }} className="text-white text-3xl mb-4">
                ${walletData.balance.toFixed(2)}
              </Text>
              <TouchableOpacity
                className="bg-white/20 rounded-xl py-3 items-center"
                onPress={handleTopUp}
              >
                <Text style={{ fontFamily: 'Inter-Medium' }} className="text-white">
                  Top Up Wallet
                </Text>
              </TouchableOpacity>
            </View>

            {/* Transactions */}
            <View className="bg-white rounded-2xl shadow-sm p-6">
              <Text style={{ fontFamily: 'Poppins-Bold' }} className="text-xl text-gray-800 mb-4">
                Recent Transactions
              </Text>
              <View className="space-y-4">
                {walletData.transactions.map((transaction) => (
                  <View key={transaction.id} className="flex-row items-center justify-between py-2 border-b border-gray-100">
                    <View className="flex-row items-center">
                      <View className={`p-2 rounded-full mr-3 ${transaction.type === 'top_up' ? 'bg-green-100' : 'bg-red-100'}`}>
                        <MaterialIcons
                          name={transaction.type === 'top_up' ? 'add' : 'remove'}
                          size={20}
                          color={transaction.type === 'top_up' ? '#059669' : '#dc2626'}
                        />
                      </View>
                      <View>
                        <Text style={{ fontFamily: 'Inter-Medium' }} className="text-gray-800">
                          {transaction.description}
                        </Text>
                        <Text style={{ fontFamily: 'Inter-Regular' }} className="text-gray-500 text-sm">
                          {formatDate(transaction.createdAt)}
                        </Text>
                      </View>
                    </View>
                    <Text
                      style={{ fontFamily: 'Inter-Medium' }}
                      className={`${transaction.type === 'top_up' ? 'text-green-600' : 'text-red-600'}`}
                    >
                      {transaction.type === 'top_up' ? '+' : '-'}${transaction.amount.toFixed(2)}
                    </Text>
                  </View>
                ))}
              </View>
            </View>
          </View>
        )}
      </ScrollView>
    </SafeAreaView>
  );
} 