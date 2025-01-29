'use client';
// import { useUser } from '@clerk/nextjs';
import Pricing from 'components/home/pricing';


const PricingPage = () => {
  // const { isLoaded, isSignedIn, user } = useUser();

  // if (!isLoaded) {
  //   return <div>Loading...</div>;
  // }

  return (
    <Pricing></Pricing>)
};

export default PricingPage;
