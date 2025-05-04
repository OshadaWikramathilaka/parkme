import { showMessage } from 'react-native-flash-message';

export const showSuccess = (message: string, description?: string) => {
  showMessage({
    message,
    description,
    type: 'success',
    icon: 'success',
    duration: 4000,
    style: { borderRadius: 12 },
    titleStyle: { fontFamily: 'Inter-Medium' },
    textStyle: { fontFamily: 'Inter-Regular' },
  });
};

export const showError = (message: string, description?: string) => {
  showMessage({
    message,
    description,
    type: 'danger',
    icon: 'danger',
    duration: 4000,
    style: { borderRadius: 12 },
    titleStyle: { fontFamily: 'Inter-Medium' },
    textStyle: { fontFamily: 'Inter-Regular' },
  });
}; 