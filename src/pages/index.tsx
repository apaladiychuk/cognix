import { Outlet } from 'react-router-dom';
import ThemeProvider from '@/context/ThemeContext';
import AuthProvider from '@/context/AuthContext';
import { Toaster } from '@/components/ui/sonner';
import NotificationProvider from '@/context/NotificationContext';
import { useMixpanelTracking } from '@/lib/hooks/use-mixpanel-tracking';

export function Providers() {
  useMixpanelTracking();

  return (
    <NotificationProvider>
      <AuthProvider>
        <ThemeProvider>
          <Outlet />
          <Toaster />
        </ThemeProvider>
      </AuthProvider>
    </NotificationProvider>
  );
}

export { Providers as Component };
