'use client';

import React, { useState } from 'react';
import Footer from 'components/footer/Footer';
import Card from 'components/card';
import Logo from '../landing/Logo';
import LandingHeader from '../landing/Header';

function Pricing() {
  const [activeButton, setActiveButton] = useState('monthly');

  return (
    <div className="flex flex-col w-full min-h-screen">
      <header className="relative z-50">
        <LandingHeader
          className="absolute top-0 w-full"
          logo={<Logo className="h-9 w-auto" />}
          logoDark={<Logo className="h-9 w-auto" />}
        />
      </header>
      <div className="relative h-full w-full px-3 font-dm bg-gradient-to-r from-[#001020] via-[#001530] to-[#000810]">
        {/* Header Content */}
        <div className="mx-auto mt-[96px] flex w-full max-w-screen-md flex-col items-center text-center md:px-3">
          <h2 className="text-[28px] font-bold text-white md:text-[44px]">
            Flexible Pricing for Real-Time Satellite Insights
          </h2>
          <p className="mt-3 px-6 text-sm text-white md:px-8 md:text-base">
            Explore our plans to access satellite tracking, threat alerts, and
            country-specific leaderboards. Start for free or upgrade to unlock
            premium features.
          </p>
          {/* Monthly / Yearly Toggle */}
          <div className="mt-8 flex h-[50px] w-[280px] items-center rounded-full bg-navy-800 p-1.5">
            <button
              className={`linear flex h-full w-1/2 cursor-pointer items-center justify-center rounded-[20px] text-xs font-bold uppercase transition duration-200 ${activeButton === 'monthly'
                ? 'bg-white text-brand-500'
                : 'bg-transparent text-white'
                }`}
              onClick={() => setActiveButton('monthly')}
            >
              Monthly
            </button>
            <button
              className={`linear flex h-full w-1/2 cursor-pointer items-center justify-center rounded-[20px] text-xs font-bold uppercase transition duration-200 ${activeButton === 'yearly'
                ? 'bg-white text-brand-500'
                : 'bg-transparent text-white'
                }`}
              onClick={() => setActiveButton('yearly')}
            >
              Yearly
            </button>
          </div>
        </div>

        {/* Pricing Section */}
        <div className="relative mx-auto mb-20 mt-12 grid h-fit w-full max-w-[375px] grid-cols-1 gap-6 px-3 xl:mt-16 xl:max-w-screen-lg xl:grid-cols-3">
          {/* Basic Plan */}
          <Card extra="w-full h-full rounded-[20px] pb-6 pt-8 px-[20px]">
            <h5 className="text-3xl font-bold text-navy-700 dark:text-white">Basic</h5>
            <p className="mt-1 text-base font-medium text-gray-600">
              Perfect for enthusiasts exploring satellite data.
            </p>
            <ul className="mt-4 list-disc pl-6 text-sm text-gray-600 dark:text-white">
              <li>Real-time satellite tracking</li>
              <li>Access to live orbital data</li>
              <li>Basic country activity leaderboard</li>
              <li>Email notifications for major events</li>
            </ul>
            <button className="mt-6 w-full rounded-xl bg-brand-50 py-2 text-base font-medium text-brand-500 dark:bg-navy-700 dark:text-white">
              Start Free
            </button>
            <div className="mt-8">
              <h5 className="text-4xl font-bold text-navy-700 dark:text-white">
                {activeButton === 'monthly' ? '$10' : '$100'}
                <span className="text-gray-600">/{activeButton === 'monthly' ? 'mo' : 'yr'}</span>
              </h5>
            </div>
          </Card>

          {/* Pro Plan */}
          <Card extra="w-full h-full rounded-[20px] pb-6 pt-8 px-[20px]">
            <h5 className="text-3xl font-bold text-navy-700 dark:text-white">Pro</h5>
            <p className="mt-1 text-base font-medium text-gray-600">
              Advanced tools for satellite professionals.
            </p>
            <ul className="mt-4 list-disc pl-6 text-sm text-gray-600 dark:text-white">
              <li>Everything in Basic</li>
              <li>Advanced threat alerts</li>
              <li>Detailed country leaderboard</li>
              <li>Customizable data dashboards</li>
            </ul>
            <button className="mt-6 w-full rounded-xl bg-brand-500 py-2 text-base font-medium text-white dark:bg-blue-800">
              Get Started
            </button>
            <div className="mt-8">
              <h5 className="text-4xl font-bold text-navy-700 dark:text-white">
                {activeButton === 'monthly' ? '$30' : '$300'}
                <span className="text-gray-600">/{activeButton === 'monthly' ? 'mo' : 'yr'}</span>
              </h5>
            </div>
          </Card>

          {/* Enterprise Plan */}
          <Card extra="w-full h-full rounded-[20px] pb-6 pt-8 px-[20px]">
            <h5 className="text-3xl font-bold text-navy-700 dark:text-white">Enterprise</h5>
            <p className="mt-1 text-base font-medium text-gray-600">
              Tailored solutions for large organizations.
            </p>
            <ul className="mt-4 list-disc pl-6 text-sm text-gray-600 dark:text-white">
              <li>Everything in Pro</li>
              <li>Custom integrations and APIs</li>
              <li>24/7 priority support</li>
              <li>Dedicated account manager</li>
            </ul>
            <button className="mt-6 w-full rounded-xl bg-brand-50 py-2 text-base font-medium text-brand-500 dark:bg-navy-700 dark:text-white">
              Contact Us
            </button>
            <div className="mt-8">
              <h5 className="text-4xl font-bold text-navy-700 dark:text-white">Custom Pricing</h5>
            </div>
          </Card>
        </div>

        {/* Frequently Asked Questions */}
        <div className="mx-auto mt-16 max-w-screen-md text-white">
          <h2 className="text-center text-[28px] font-bold">Frequently Asked Questions</h2>
          <div className="mt-8 grid gap-8">
            <div>
              <h3 className="text-lg font-semibold">What is included in the Basic plan?</h3>
              <p className="mt-2 text-sm text-gray-300">
                The Basic plan includes real-time satellite tracking, access to live orbital data, and email notifications for major events.
              </p>
            </div>
            <div>
              <h3 className="text-lg font-semibold">Can I upgrade or downgrade my plan later?</h3>
              <p className="mt-2 text-sm text-gray-300">
                Yes, you can change your plan at any time based on your needs. Contact support for assistance.
              </p>
            </div>
            <div>
              <h3 className="text-lg font-semibold">Do you offer custom solutions?</h3>
              <p className="mt-2 text-sm text-gray-300">
                Absolutely! Our Enterprise plan is designed for organizations with custom requirements.
              </p>
            </div>
          </div>
        </div>

        {/* Footer */}
        <div className="mx-auto max-w-screen-xl mt-20">
          <Footer />
        </div>
      </div>
    </div>
  );
}

export default Pricing;
