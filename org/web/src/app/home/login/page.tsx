'use client';
import { useUser } from '@clerk/nextjs';
import LoginPage from 'components/home/login';

const PricingPage = () => {
  const { isLoaded, isSignedIn, user } = useUser();

  if (!isLoaded) {
    return <div>Loading...</div>;
  }

  return (
    <LoginPage
      user={isSignedIn ? { name: user?.fullName || 'User' } : undefined}
    />)
};

export default PricingPage;
