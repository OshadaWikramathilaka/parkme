import React, { useState } from 'react';
import { View, Text, TextInput, TouchableOpacity, ScrollView } from 'react-native';
import { useRouter } from 'expo-router';
import { MaterialIcons } from '@expo/vector-icons';
import { ENDPOINTS } from '@/constants/Config';
import { showError } from '@/utils/flashMessage';
import { getToken } from '@/utils/auth';
import { SafeAreaView } from 'react-native-safe-area-context';
import Header from '@/components/ui/Header';

// Default credit card values
const DEFAULT_CARD = {
  number: '4111 1111 1111 1111',
  expiry: '12/25',
  cvv: '123',
  name: 'John Doe'
};

export default function TopUpModal() {
  const router = useRouter();
  const [amount, setAmount] = useState('100');
  const [cardNumber, setCardNumber] = useState(DEFAULT_CARD.number);
  const [cardExpiry, setCardExpiry] = useState(DEFAULT_CARD.expiry);
  const [cardCvv, setCardCvv] = useState(DEFAULT_CARD.cvv);
  const [cardName, setCardName] = useState(DEFAULT_CARD.name);
  const [isLoading, setIsLoading] = useState(false);

  const handleTopUp = async () => {
    if (isLoading) return;
    setIsLoading(true);

    try {
      const token = await getToken();
      if (!token) {
        router.replace('/(auth)/login');
        return;
      }

      const response = await fetch(ENDPOINTS.WALLET_TOPUP, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ amount: parseFloat(amount) }),
      });

      const result = await response.json();
      if (response.ok && result.success) {
        router.back();
      } else {
        showError('Error', result.message || 'Failed to top up wallet');
      }
    } catch (error) {
      showError('Network Error', 'Failed to connect to server');
      console.error('Top-up error:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const formatCardNumber = (text: string) => {
    const cleaned = text.replace(/\s/g, '');
    const groups = cleaned.match(/.{1,4}/g);
    return groups ? groups.join(' ') : cleaned;
  };

  const formatExpiry = (text: string) => {
    const cleaned = text.replace(/\D/g, '');
    if (cleaned.length >= 2) {
      return `${cleaned.slice(0, 2)}/${cleaned.slice(2, 4)}`;
    }
    return cleaned;
  };

  return (
    <SafeAreaView className="flex-1 bg-gray-50" edges={['top']}>
      <Header title="Top Up Wallet" showClose />

      <ScrollView className="flex-1">
        <View className="p-6">
          {/* Amount Input */}
          <View className="mb-8">
            <Text style={{ fontFamily: 'Inter-Medium' }} className="text-gray-600 mb-2">
              Amount to Top Up ($)
            </Text>
            <TextInput
              className="bg-white border border-gray-200 rounded-xl px-4 py-3 text-lg"
              style={{ fontFamily: 'Inter-Regular' }}
              keyboardType="numeric"
              value={amount}
              onChangeText={setAmount}
              placeholder="Enter amount"
            />
          </View>

          {/* Card Details */}
          <View className="bg-white rounded-2xl p-6 mb-6 shadow-sm">
            <Text style={{ fontFamily: 'Poppins-Bold' }} className="text-lg text-gray-800 mb-4">
              Card Details
            </Text>

            {/* Card Number */}
            <View className="mb-4">
              <Text style={{ fontFamily: 'Inter-Medium' }} className="text-gray-600 mb-2">
                Card Number
              </Text>
              <TextInput
                className="border border-gray-200 rounded-xl px-4 py-3"
                style={{ fontFamily: 'Inter-Regular' }}
                value={cardNumber}
                onChangeText={(text) => setCardNumber(formatCardNumber(text))}
                maxLength={19}
                keyboardType="numeric"
              />
            </View>

            {/* Card Holder */}
            <View className="mb-4">
              <Text style={{ fontFamily: 'Inter-Medium' }} className="text-gray-600 mb-2">
                Card Holder Name
              </Text>
              <TextInput
                className="border border-gray-200 rounded-xl px-4 py-3"
                style={{ fontFamily: 'Inter-Regular' }}
                value={cardName}
                onChangeText={setCardName}
              />
            </View>

            {/* Expiry and CVV */}
            <View className="flex-row space-x-4">
              <View className="flex-1">
                <Text style={{ fontFamily: 'Inter-Medium' }} className="text-gray-600 mb-2">
                  Expiry Date
                </Text>
                <TextInput
                  className="border border-gray-200 rounded-xl px-4 py-3"
                  style={{ fontFamily: 'Inter-Regular' }}
                  value={cardExpiry}
                  onChangeText={(text) => setCardExpiry(formatExpiry(text))}
                  maxLength={5}
                  placeholder="MM/YY"
                  keyboardType="numeric"
                />
              </View>
              <View className="flex-1">
                <Text style={{ fontFamily: 'Inter-Medium' }} className="text-gray-600 mb-2">
                  CVV
                </Text>
                <TextInput
                  className="border border-gray-200 rounded-xl px-4 py-3"
                  style={{ fontFamily: 'Inter-Regular' }}
                  value={cardCvv}
                  onChangeText={setCardCvv}
                  maxLength={3}
                  keyboardType="numeric"
                  secureTextEntry
                />
              </View>
            </View>
          </View>

          {/* Submit Button */}
          <TouchableOpacity
            className="bg-blue-600 rounded-xl py-4 items-center"
            onPress={handleTopUp}
            disabled={isLoading}
          >
            <Text style={{ fontFamily: 'Inter-Medium' }} className="text-white text-lg">
              {isLoading ? 'Processing...' : `Top Up $${amount}`}
            </Text>
          </TouchableOpacity>
        </View>
      </ScrollView>
    </SafeAreaView>
  );
} 